package simulator

import (
	"context"

	"github.com/jpraynaud/alien-invasion-simulator/pkg/simulator/entity"
)

// WorldStorer is a world store interface
type WorldStorer interface {
	// GetCity retrieves a city
	GetCity(ctx context.Context, cityName string) (*entity.City, error)
	// GetAliveCities retrieves the list of non destroyed cities
	GetAliveCities(ctx context.Context) ([]*entity.City, error)
	// AddCity adds a city
	AddCity(ctx context.Context, cityName string) (*entity.City, error)
	// DestroyCity destroys a city
	DestroyCity(ctx context.Context, city *entity.City) error
	// AddLink adds a link from a city to another city given a direction
	AddLink(ctx context.Context, cityFrom, cityTo *entity.City, direction entity.Direction) error
	// GetAlien retrieves an alien
	GetAlien(ctx context.Context, alienID int) (*entity.Alien, error)
	// AddAlien adds an alien
	AddAlien(ctx context.Context, alienID int) (*entity.Alien, error)
	// MoveAlien moves an alien to a city
	MoveAlien(ctx context.Context, alien *entity.Alien, city *entity.City) error
	// IsTrappedAlien checks if an alien is trapped
	IsTrappedAlien(ctx context.Context, alien *entity.Alien) (bool, error)
	// TrapAlien traps an alien
	TrapAlien(ctx context.Context, alien *entity.Alien) error
	// GetAlienAtCity retrieves the alien at a given city
	GetAlienAtCity(ctx context.Context, city *entity.City) (*entity.Alien, error)
	// GetUntrappedAliens retrieves the list of untrapped aliens
	GetUntrappedAliens(ctx context.Context) ([]*entity.Alien, error)
}

// Simulator is an alien invasion simulator interface
type Simulator interface {
	// Prepare prepares the simulation
	Prepare(ctx context.Context) error
	// HasNextStep computes if a next step of the simulation exists
	HasNextStep(ctx context.Context) (bool, error)
	// SimulateNextStep simulates the next step of the simulation
	SimulateNextStep(ctx context.Context) error
	// Run simulates an alien invasion
	Run(ctx context.Context) error
	// Finalize finalizes the simulation
	Finalize(ctx context.Context) error
}

// Randomer is a random generator
type Randomer interface {
	// GetRandomInt retrieves a random integer between 0 and n-1 given n
	GetRandomInt(n int) (int, error)
}
