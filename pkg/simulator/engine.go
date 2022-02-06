package simulator

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jpraynaud/alien-invasion-simulator/pkg/simulator/entity"
)

// SimulationEngine is a simple implementation of an alien invasion simulator
type SimulationEngine struct {
	// World store
	world WorldStorer

	// Random generator
	random Randomer

	// Input reader
	in io.Reader

	// Output writer
	out io.Writer

	// Number of aliens that are spawned during initialization
	startAliens uint

	// Maximum number of iterations to simulate
	maxSteps uint

	// Number of steps already simulated
	totalSteps uint
}

var _ Simulator = (*SimulationEngine)(nil)

// NewSimulationEngine is a simulation engine constructor
func NewSimulationEngine(startAliens, maxSteps uint, world WorldStorer, random Randomer, in io.Reader, out io.Writer) *SimulationEngine {
	return &SimulationEngine{
		world:       world,
		random:      random,
		in:          in,
		out:         out,
		maxSteps:    maxSteps,
		startAliens: startAliens,
	}
}

// Prepare prepares the simulation
func (s *SimulationEngine) Prepare(ctx context.Context) error {
	log.Info("Prepare")

	// Read and parse input
	err := s.loadInputToWorld(ctx)
	if err != nil {
		return err
	}

	// Unleash all the aliens
	for i := 0; i < int(s.startAliens); i++ {
		// Create alien
		alienID := i + 1
		alien, err := s.world.AddAlien(ctx, alienID)
		if err != nil {
			return err
		}

		// Move the alien to its original city
		var nextCity *entity.City
		aliveCities, err := s.world.GetAliveCities(ctx)
		if err != nil {
			return err
		}
		if len(aliveCities) == 0 {
			return nil
		}
		r, err := s.random.GetRandomInt(len(aliveCities))
		if err != nil {
			return err
		}
		nextCity = aliveCities[r]
		_, err = s.moveAlienToCity(ctx, alien, nextCity)
		if err != nil {
			return err
		}
	}
	return nil
}

// HasNextStep computes if a next step of the simulation exists
func (s *SimulationEngine) HasNextStep(ctx context.Context) (bool, error) {
	log.WithFields(log.Fields{
		"step": s.totalSteps,
	}).Info("HasNextStep")

	// If max steps is reached, there are no more step
	if s.totalSteps >= s.maxSteps {
		return false, nil
	}

	// If all aliens have been trapped, there are no more step
	untrappedAliens, err := s.world.GetUntrappedAliens(ctx)
	if err != nil {
		return false, err
	}
	if len(untrappedAliens) == 0 {
		return false, nil
	}

	// If all cities have been destroyed, there are no more step
	aliveCities, err := s.world.GetAliveCities(ctx)
	if err != nil {
		return false, err
	}
	if len(aliveCities) == 0 {
		return false, nil
	}

	return true, nil
}

// SimulateNextStep simulates the next step of the simulation
func (s *SimulationEngine) SimulateNextStep(ctx context.Context) error {
	log.WithFields(log.Fields{
		"step": s.totalSteps,
	}).Debug("SimulateStep")

	// Move randomly each remaining alien
	s.totalSteps++
	untrappedAliens, err := s.world.GetUntrappedAliens(ctx)
	if err != nil {
		return err
	}
	for _, alien := range untrappedAliens {
		// Check if alien is trapped
		// This occurs if it was trapped previously in this loop
		isTrapped, err := s.world.IsTrappedAlien(ctx, alien)
		if err != nil {
			return err
		}
		if isTrapped {
			continue
		}

		// Move randomly alien to next available city
		var nextCity *entity.City
		currentCity := alien.City
		availableLinks := currentCity.GetAvailableLinks()
		if len(availableLinks) > 0 {
			r, err := s.random.GetRandomInt(len(availableLinks))
			if err != nil {
				return err
			}
			i := 0
			for _, city := range availableLinks {
				if i == r {
					nextCity = city
					continue
				}
				i++
			}
			_, err = s.moveAlienToCity(ctx, alien, nextCity)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// run is a run helper function
func run(ctx context.Context, s Simulator) error {
	// Prepare
	err := s.Prepare(ctx)
	if err != nil {
		return err
	}

	// Simulate steps
	for {
		select {
		case <-ctx.Done():
			log.Warn("Simulation was cancelled")
			return entity.ErrContextCancelled
		default:
			hasNextStep, err := s.HasNextStep(ctx)
			if err != nil {
				return err
			}
			if !hasNextStep {
				return s.Finalize(ctx)
			}
			err = s.SimulateNextStep(ctx)
			if err != nil {
				return err
			}
		}
	}
}

// Run simulates an alien invasion
func (s *SimulationEngine) Run(ctx context.Context) error {
	log.Info("Run")
	return run(ctx, s)
}

// Finalize finalizes the simulation
func (s *SimulationEngine) Finalize(ctx context.Context) error {
	log.WithFields(log.Fields{
		"steps": s.totalSteps,
	}).Info("Finalize")

	cities, err := s.world.GetAliveCities(ctx)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(s.out, "")
	if err != nil {
		return err
	}
	for _, city := range cities {
		_, err = fmt.Fprintln(s.out, city)
		if err != nil {
			return err
		}
	}

	return nil
}

// loadInputToWorld parses input and saves it to the world
func (s *SimulationEngine) loadInputToWorld(ctx context.Context) error {
	// Helper function
	registerCity := func(cityName string) (*entity.City, error) {
		city, err := s.world.GetCity(ctx, cityName)
		if err != nil {
			return city, err
		}
		if city == nil {
			city, err = s.world.AddCity(ctx, cityName)
			if err != nil {
				return city, err
			}
		}
		return city, nil
	}

	// Scan input
	scanner := bufio.NewScanner(s.in)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines
		if len(line) == 0 {
			continue
		}
		// Assume that a city does not contain any space
		lineChunks := strings.Split(line, " ")
		if len(lineChunks) == 0 {
			return entity.ErrParseCityDefinition
		}
		cityFromName := lineChunks[0]
		cityFrom, err := registerCity(cityFromName)
		if err != nil {
			return err
		}
		// Parse all direction/city couples
		for _, lineChunk := range lineChunks[1:] {
			linkChunks := strings.Split(strings.TrimSpace(lineChunk), "=")
			if len(linkChunks) != 2 {
				return entity.ErrParseCityDefinition
			}
			directionName := linkChunks[0]
			cityToName := linkChunks[1]
			cityTo, err := registerCity(cityToName)
			if err != nil {
				return err
			}
			var direction entity.Direction
			switch directionName {
			case "north":
				direction = entity.North
			case "east":
				direction = entity.East
			case "south":
				direction = entity.South
			case "west":
				direction = entity.West
			default:
				return entity.ErrParseCityDefinition
			}
			err = s.world.AddLink(ctx, cityFrom, cityTo, direction)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// moveAlienToCity applies the move of an alien to a city
func (s *SimulationEngine) moveAlienToCity(ctx context.Context, alien *entity.Alien, city *entity.City) (bool, error) {
	log.WithFields(log.Fields{
		"alien": alien,
		"city":  city,
	}).Debug("moveAlienToCity")

	// Retrieve alien at city
	destroyedCity := false
	alienAlreadyInCity, err := s.world.GetAlienAtCity(ctx, city)
	if err != nil {
		return destroyedCity, err
	}

	// Decide what to do next
	switch {
	case alienAlreadyInCity == alien:
		// Same alien, nothing to do
		return destroyedCity, nil
	case alienAlreadyInCity == nil:
		// Record alien current city
		err := s.world.MoveAlien(ctx, alien, city)
		if err != nil {
			return destroyedCity, err
		}
	default:
		// Aliens fight!
		err = s.world.TrapAlien(ctx, alien)
		if err != nil {
			return destroyedCity, err
		}
		err = s.world.TrapAlien(ctx, alienAlreadyInCity)
		if err != nil {
			return destroyedCity, err
		}

		// Destroy city
		err = s.world.DestroyCity(ctx, city)
		if err != nil {
			return destroyedCity, err
		}
		// Print message
		destroyedCity = true
		_, err := fmt.Fprintf(s.out, "%s has been destroyed by %s and %s\n", city.Name, alien, alienAlreadyInCity)
		if err != nil {
			return destroyedCity, err
		}
	}

	return destroyedCity, nil
}
