package simulator

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/jpraynaud/alien-invasion-simulator/pkg/simulator/entity"
)

// WorldStorerMock mocks a WorldStorer
type WorldStorerMock struct {
	mock.Mock
}

var _ WorldStorer = (*WorldStorerMock)(nil)

// GetCity retrieves a city
func (w *WorldStorerMock) GetCity(ctx context.Context, cityName string) (*entity.City, error) {
	args := w.Called(ctx, cityName)
	return args.Get(0).(*entity.City), args.Error(1)
}

// GetAliveCities retrieves the list of non destroyed cities
func (w *WorldStorerMock) GetAliveCities(ctx context.Context) ([]*entity.City, error) {
	args := w.Called(ctx)
	return args.Get(0).([]*entity.City), args.Error(1)
}

// AddCity adds a city
func (w *WorldStorerMock) AddCity(ctx context.Context, cityName string) (*entity.City, error) {
	args := w.Called(ctx, cityName)
	return args.Get(0).(*entity.City), args.Error(1)
}

// DestroyCity destroys a city
func (w *WorldStorerMock) DestroyCity(ctx context.Context, city *entity.City) error {
	args := w.Called(ctx, city)
	return args.Error(0)
}

// AddLink adds a link from a city to another city given a direction
func (w *WorldStorerMock) AddLink(ctx context.Context, cityFrom, cityTo *entity.City, direction entity.Direction) error {
	args := w.Called(ctx, cityFrom, cityTo, direction)
	return args.Error(0)
}

// GetAlien retrieves an alien
func (w *WorldStorerMock) GetAlien(ctx context.Context, alienID int) (*entity.Alien, error) {
	args := w.Called(ctx, alienID)
	return args.Get(0).(*entity.Alien), args.Error(1)
}

// AddAlien adds an alien
func (w *WorldStorerMock) AddAlien(ctx context.Context, alienID int) (*entity.Alien, error) {
	args := w.Called(ctx, alienID)
	return args.Get(0).(*entity.Alien), args.Error(1)
}

// MoveAlien moves an alien to a city
func (w *WorldStorerMock) MoveAlien(ctx context.Context, alien *entity.Alien, city *entity.City) error {
	args := w.Called(ctx, alien, city)
	return args.Error(0)
}

// IsTrappedAlien checks if an alien is trapped
func (w *WorldStorerMock) IsTrappedAlien(ctx context.Context, alien *entity.Alien) (bool, error) {
	args := w.Called(ctx, alien)
	return args.Bool(0), args.Error(1)
}

// TrapAlien traps an alien
func (w *WorldStorerMock) TrapAlien(ctx context.Context, alien *entity.Alien) error {
	args := w.Called(ctx, alien)
	return args.Error(0)
}

// GetAlienAtCity retrieves the alien at a given city
func (w *WorldStorerMock) GetAlienAtCity(ctx context.Context, city *entity.City) (*entity.Alien, error) {
	args := w.Called(ctx, city)
	return args.Get(0).(*entity.Alien), args.Error(1)
}

// GetUntrappedAliens retrieves the list of untrapped aliens
func (w *WorldStorerMock) GetUntrappedAliens(ctx context.Context) ([]*entity.Alien, error) {
	args := w.Called(ctx)
	return args.Get(0).([]*entity.Alien), args.Error(1)
}

// SimulatorMock mocks a Simulator
type SimulatorMock struct {
	mock.Mock
}

var _ Simulator = (*SimulatorMock)(nil)

// Prepare prepares the simulation
func (s *SimulatorMock) Prepare(ctx context.Context) error {
	args := s.Called(ctx)
	return args.Error(0)
}

// HasNextStep computes if a next step of the simulation exists
func (s *SimulatorMock) HasNextStep(ctx context.Context) (bool, error) {
	args := s.Called(ctx)
	return args.Bool(0), args.Error(1)
}

// SimulateNextStep simulates the next step of the simulation
func (s *SimulatorMock) SimulateNextStep(ctx context.Context) error {
	args := s.Called(ctx)
	return args.Error(0)
}

// Run simulates an alien invasion
func (s *SimulatorMock) Run(ctx context.Context) error {
	args := s.Called(ctx)
	return args.Error(0)
}

// Finalize finalizes the simulation
func (s *SimulatorMock) Finalize(ctx context.Context) error {
	args := s.Called(ctx)
	return args.Error(0)
}

// RandomerMock mocks a Randomer
type RandomerMock struct {
	mock.Mock
}

var _ Randomer = (*RandomerMock)(nil)

// GetRandomInt retrieves a random integer between 0 and n-1 given n
func (r *RandomerMock) GetRandomInt(n int) (int, error) {
	args := r.Called(n)
	return args.Int(0), args.Error(1)
}
