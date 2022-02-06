package simulator

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jpraynaud/alien-invasion-simulator/pkg/simulator/entity"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_SimulationEngine_Prepare(t *testing.T) {
	var alienNil *entity.Alien
	alien1 := entity.NewAlien(1)
	alien2 := entity.NewAlien(2)
	city1 := entity.NewCity("City1")
	city2 := entity.NewCity("City2")

	t.Run("Case 1: OK", func(t *testing.T) {
		ctx := context.Background()

		worldStorerMock := &WorldStorerMock{}
		// Add Alien1 to unoccupied city
		worldStorerMock.On("AddAlien", ctx, alien1.AlienID).Return(alien1, nil).Once()
		worldStorerMock.On("GetAliveCities", ctx).Return([]*entity.City{city1, city2}, nil).Once()
		worldStorerMock.On("GetAlienAtCity", ctx, city1).Return(alienNil, nil).Once()
		worldStorerMock.On("MoveAlien", ctx, alien1, city1).Return(nil, nil).Once()
		// Add Alien2
		worldStorerMock.On("AddAlien", ctx, alien2.AlienID).Return(alien2, nil).Once()
		worldStorerMock.On("GetAliveCities", ctx).Return([]*entity.City{city1, city2}, nil).Once()
		worldStorerMock.On("GetAlienAtCity", ctx, city1).Return(alien1, nil).Once()
		worldStorerMock.On("TrapAlien", ctx, alien1).Return(nil, nil).Once()
		worldStorerMock.On("TrapAlien", ctx, alien2).Return(nil, nil).Once()
		worldStorerMock.On("DestroyCity", ctx, city1).Return(nil, nil).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		randomerMock.On("GetRandomInt", mock.Anything).Return(0, nil).Times(2)
		defer randomerMock.AssertExpectations(t)

		out := &bytes.Buffer{}

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          &bytes.Buffer{},
			out:         out,
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 2,
		}

		err := s.Prepare(ctx)
		require.NoError(t, err)
		require.Equal(t, "City1 has been destroyed by Alien #2 and Alien #1\n", out.String())
	})

	t.Run("Case 2: Error occurs", func(t *testing.T) {
		ctx := context.Background()

		error1 := fmt.Errorf("error 1")

		worldStorerMock := &WorldStorerMock{}
		// Add Alien1 to unoccupied city
		worldStorerMock.On("AddAlien", ctx, alien1.AlienID).Return(alien1, nil).Once()
		worldStorerMock.On("GetAliveCities", ctx).Return([]*entity.City{city1, city2}, nil).Once()
		worldStorerMock.On("GetAlienAtCity", ctx, city1).Return(alienNil, nil).Once()
		worldStorerMock.On("MoveAlien", ctx, alien1, city1).Return(nil, nil).Once()
		// Add Alien2
		worldStorerMock.On("AddAlien", ctx, alien2.AlienID).Return(alien2, nil).Once()
		worldStorerMock.On("GetAliveCities", ctx).Return([]*entity.City{nil}, error1).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		randomerMock.On("GetRandomInt", mock.Anything).Return(0, nil).Once()
		defer randomerMock.AssertExpectations(t)

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          &bytes.Buffer{},
			out:         &bytes.Buffer{},
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 2,
		}

		err := s.Prepare(ctx)
		require.ErrorIs(t, err, error1)
	})
}

func Test_SimulationEngine_HasNextStep(t *testing.T) {
	alien1 := entity.NewAlien(1)
	alien2 := entity.NewAlien(2)
	city1 := entity.NewCity("City1")
	city2 := entity.NewCity("City2")

	aliensEmpty := []*entity.Alien{}
	aliensFilled := []*entity.Alien{alien1, alien2}
	citiesEmpty := []*entity.City{}
	citiesFilled := []*entity.City{city1, city2}

	error1 := fmt.Errorf("error 1")
	error2 := fmt.Errorf("error 2")

	tests := []struct {
		testName                     string
		giveTotalSteps, giveMaxSteps uint
		giveUntrappedAliens          []*entity.Alien
		giveUntrappedAliensError     error
		giveAliveCities              []*entity.City
		giveAliveCitiesError         error
		wantGetUntrappedAliensCalls  int
		wantGetAliveCitiesCalls      int
		wantResult                   bool
		wantError                    error
	}{
		{
			testName:                    "Too many steps",
			giveTotalSteps:              10,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensEmpty,
			giveUntrappedAliensError:    nil,
			giveAliveCities:             citiesEmpty,
			giveAliveCitiesError:        nil,
			wantGetUntrappedAliensCalls: 0,
			wantGetAliveCitiesCalls:     0,
			wantResult:                  false,
			wantError:                   nil,
		},
		{
			testName:                    "All aliens trapped",
			giveTotalSteps:              2,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensEmpty,
			giveUntrappedAliensError:    nil,
			giveAliveCities:             citiesFilled,
			giveAliveCitiesError:        nil,
			wantGetUntrappedAliensCalls: 1,
			wantGetAliveCitiesCalls:     0,
			wantResult:                  false,
			wantError:                   nil,
		},
		{
			testName:                    "GetUntrappedAliens returns error",
			giveTotalSteps:              2,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensEmpty,
			giveUntrappedAliensError:    error1,
			giveAliveCities:             citiesFilled,
			giveAliveCitiesError:        nil,
			wantGetUntrappedAliensCalls: 1,
			wantGetAliveCitiesCalls:     0,
			wantResult:                  false,
			wantError:                   error1,
		},
		{
			testName:                    "All cities destroyed",
			giveTotalSteps:              2,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensFilled,
			giveUntrappedAliensError:    nil,
			giveAliveCities:             citiesEmpty,
			giveAliveCitiesError:        nil,
			wantGetUntrappedAliensCalls: 1,
			wantGetAliveCitiesCalls:     1,
			wantResult:                  false,
			wantError:                   nil,
		},
		{
			testName:                    "GetAliveCities returns error",
			giveTotalSteps:              2,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensFilled,
			giveUntrappedAliensError:    nil,
			giveAliveCities:             citiesEmpty,
			giveAliveCitiesError:        error2,
			wantGetUntrappedAliensCalls: 1,
			wantGetAliveCitiesCalls:     1,
			wantResult:                  false,
			wantError:                   error2,
		},
		{
			testName:                    "Next step exists",
			giveTotalSteps:              2,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensFilled,
			giveUntrappedAliensError:    nil,
			giveAliveCities:             citiesFilled,
			giveAliveCitiesError:        nil,
			wantGetUntrappedAliensCalls: 1,
			wantGetAliveCitiesCalls:     1,
			wantResult:                  true,
			wantError:                   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctx := context.Background()

			worldStorerMock := &WorldStorerMock{}
			if tt.wantGetUntrappedAliensCalls > 0 {
				worldStorerMock.On("GetUntrappedAliens", ctx).Return(tt.giveUntrappedAliens, tt.giveUntrappedAliensError).Times(tt.wantGetUntrappedAliensCalls)

			}
			if tt.wantGetAliveCitiesCalls > 0 {
				worldStorerMock.On("GetAliveCities", ctx).Return(tt.giveAliveCities, tt.giveAliveCitiesError).Times(tt.wantGetAliveCitiesCalls)
			}
			defer worldStorerMock.AssertExpectations(t)

			randomerMock := &RandomerMock{}
			defer randomerMock.AssertExpectations(t)

			s := SimulationEngine{
				world:       worldStorerMock,
				random:      randomerMock,
				in:          &bytes.Buffer{},
				out:         &bytes.Buffer{},
				totalSteps:  tt.giveTotalSteps,
				maxSteps:    tt.giveMaxSteps,
				startAliens: 0,
			}

			result, err := s.HasNextStep(ctx)
			require.ErrorIs(t, tt.wantError, err)
			require.Equal(t, tt.wantResult, result)
		})
	}
}

func Test_SimulationEngine_SimulateNextStep(t *testing.T) {
	var alienNil *entity.Alien
	alien1 := entity.NewAlien(1)
	alien2 := entity.NewAlien(2)
	alien3 := entity.NewAlien(3)
	city1 := entity.NewCity("City1")
	city2 := entity.NewCity("City2")

	t.Run("Case 1: Alien1 is trapped / Alien2 is moving to same city", func(t *testing.T) {
		ctx := context.Background()

		err := city1.SetCityTo(city2, entity.South)
		require.NoError(t, err)
		alien2.City = city1

		worldStorerMock := &WorldStorerMock{}
		worldStorerMock.On("GetUntrappedAliens", ctx).Return([]*entity.Alien{alien1, alien2}, nil).Once()
		// Alien1 is already trapped
		worldStorerMock.On("IsTrappedAlien", ctx, alien1).Return(true, nil).Once()
		// Alien2 is moved to its current city
		worldStorerMock.On("IsTrappedAlien", ctx, alien2).Return(false, nil).Once()
		worldStorerMock.On("GetAlienAtCity", ctx, city2).Return(alien2, nil).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		randomerMock.On("GetRandomInt", mock.Anything).Return(0, nil).Once()
		defer randomerMock.AssertExpectations(t)

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          &bytes.Buffer{},
			out:         &bytes.Buffer{},
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 0,
		}

		err = s.SimulateNextStep(ctx)
		require.NoError(t, err)
		require.Equal(t, uint(1), s.totalSteps)
	})

	t.Run("Case 2: Alien1 is trapped / Alien2 is moving to an unoccupied city", func(t *testing.T) {
		ctx := context.Background()

		err := city1.SetCityTo(city2, entity.South)
		require.NoError(t, err)
		alien2.City = city1

		worldStorerMock := &WorldStorerMock{}
		worldStorerMock.On("GetUntrappedAliens", ctx).Return([]*entity.Alien{alien1, alien2}, nil).Once()
		// Alien1 is already trapped
		worldStorerMock.On("IsTrappedAlien", ctx, alien1).Return(true, nil).Once()
		// Alien2 is moved to an unoccupied city
		worldStorerMock.On("IsTrappedAlien", ctx, alien2).Return(false, nil).Once()
		worldStorerMock.On("GetAlienAtCity", ctx, city2).Return(alienNil, nil).Once()
		worldStorerMock.On("MoveAlien", ctx, alien2, city2).Return(nil).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		randomerMock.On("GetRandomInt", mock.Anything).Return(0, nil).Once()
		defer randomerMock.AssertExpectations(t)

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          &bytes.Buffer{},
			out:         &bytes.Buffer{},
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 0,
		}

		err = s.SimulateNextStep(ctx)
		require.NoError(t, err)
		require.Equal(t, uint(1), s.totalSteps)
	})

	t.Run("Case 3:  Alien1 is trapped / Alien2 is moving to an occupied city", func(t *testing.T) {
		ctx := context.Background()

		err := city1.SetCityTo(city2, entity.South)
		require.NoError(t, err)
		alien2.City = city1

		worldStorerMock := &WorldStorerMock{}
		worldStorerMock.On("GetUntrappedAliens", ctx).Return([]*entity.Alien{alien1, alien2}, nil).Once()
		// Alien1 is already trapped
		worldStorerMock.On("IsTrappedAlien", ctx, alien1).Return(true, nil).Once()
		// Alien2 is moved to an occupied city
		worldStorerMock.On("IsTrappedAlien", ctx, alien2).Return(false, nil).Once()
		worldStorerMock.On("GetAlienAtCity", ctx, city2).Return(alien3, nil).Once()
		worldStorerMock.On("TrapAlien", ctx, alien2).Return(nil).Once()
		worldStorerMock.On("TrapAlien", ctx, alien3).Return(nil).Once()
		worldStorerMock.On("DestroyCity", ctx, city2).Return(nil).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		randomerMock.On("GetRandomInt", mock.Anything).Return(0, nil).Once()
		defer randomerMock.AssertExpectations(t)

		out := &bytes.Buffer{}

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          &bytes.Buffer{},
			out:         out,
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 0,
		}

		err = s.SimulateNextStep(ctx)
		require.NoError(t, err)
		require.Equal(t, uint(1), s.totalSteps)
		require.Equal(t, "City2 has been destroyed by Alien #2 and Alien #3\n", out.String())
	})

	t.Run("Case 2: Error", func(t *testing.T) {
		ctx := context.Background()

		error1 := fmt.Errorf("error 1")

		worldStorerMock := &WorldStorerMock{}
		worldStorerMock.On("GetUntrappedAliens", ctx).Return([]*entity.Alien{nil}, error1).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		defer randomerMock.AssertExpectations(t)

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          &bytes.Buffer{},
			out:         &bytes.Buffer{},
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 0,
		}

		err := s.SimulateNextStep(ctx)
		require.ErrorIs(t, err, error1)
	})
}

func Test_SimulationEngine_Run(t *testing.T) {
	totalSteps := 10
	error1 := fmt.Errorf("error 1")

	t.Run("Case 1: OK", func(t *testing.T) {
		ctx := context.Background()

		simulatorMock := &SimulatorMock{}
		simulatorMock.On("Prepare", ctx).Return(nil).Once()
		simulatorMock.On("HasNextStep", ctx).Return(true, nil).Times(totalSteps)
		simulatorMock.On("HasNextStep", ctx).Return(false, nil).Once()
		simulatorMock.On("SimulateNextStep", ctx).Return(nil).Times(totalSteps)
		simulatorMock.On("Finalize", ctx).Return(nil).Once()
		defer simulatorMock.AssertExpectations(t)

		err := run(ctx, simulatorMock)
		require.NoError(t, err)
	})

	t.Run("Case 2: Error on Prepare", func(t *testing.T) {
		ctx := context.Background()

		simulatorMock := &SimulatorMock{}
		simulatorMock.On("Prepare", ctx).Return(error1).Once()
		defer simulatorMock.AssertExpectations(t)

		err := run(ctx, simulatorMock)
		require.ErrorIs(t, err, error1)
	})

	t.Run("Case 3: Error on HasNextStep", func(t *testing.T) {
		ctx := context.Background()

		simulatorMock := &SimulatorMock{}
		simulatorMock.On("Prepare", ctx).Return(nil).Once()
		simulatorMock.On("HasNextStep", ctx).Return(true, nil).Times(totalSteps)
		simulatorMock.On("HasNextStep", ctx).Return(false, error1).Once()
		simulatorMock.On("SimulateNextStep", ctx).Return(nil).Times(totalSteps)
		defer simulatorMock.AssertExpectations(t)

		err := run(ctx, simulatorMock)
		require.ErrorIs(t, err, error1)
	})

	t.Run("Case 4: Error on SimulateNextStep", func(t *testing.T) {
		ctx := context.Background()

		simulatorMock := &SimulatorMock{}
		simulatorMock.On("Prepare", ctx).Return(nil).Once()
		simulatorMock.On("HasNextStep", ctx).Return(true, nil).Times(totalSteps)
		simulatorMock.On("SimulateNextStep", ctx).Return(nil).Times(totalSteps - 1)
		simulatorMock.On("SimulateNextStep", ctx).Return(error1).Once()
		defer simulatorMock.AssertExpectations(t)

		err := run(ctx, simulatorMock)
		require.ErrorIs(t, err, error1)
	})

	t.Run("Case 5: Error on Finalize", func(t *testing.T) {
		ctx := context.Background()

		simulatorMock := &SimulatorMock{}
		simulatorMock.On("Prepare", ctx).Return(nil).Once()
		simulatorMock.On("HasNextStep", ctx).Return(true, nil).Times(totalSteps)
		simulatorMock.On("HasNextStep", ctx).Return(false, nil).Once()
		simulatorMock.On("SimulateNextStep", ctx).Return(nil).Times(totalSteps)
		simulatorMock.On("Finalize", ctx).Return(error1).Once()
		defer simulatorMock.AssertExpectations(t)

		err := run(ctx, simulatorMock)
		require.ErrorIs(t, err, error1)
	})
}

func Test_SimulationEngine_Finalize(t *testing.T) {
	city1 := entity.NewCity("City1")
	city2 := entity.NewCity("City2")

	t.Run("Case 1: OK", func(t *testing.T) {
		ctx := context.Background()

		worldStorerMock := &WorldStorerMock{}
		worldStorerMock.On("GetAliveCities", ctx).Return([]*entity.City{city1, city2}, nil).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		defer randomerMock.AssertExpectations(t)

		out := &bytes.Buffer{}

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          &bytes.Buffer{},
			out:         out,
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 0,
		}

		err := s.Finalize(ctx)
		require.NoError(t, err)
		require.Equal(t, "\nCity1\nCity2\n", out.String())
	})

	t.Run("Case 2: Error", func(t *testing.T) {
		ctx := context.Background()

		error1 := fmt.Errorf("error 1")

		worldStorerMock := &WorldStorerMock{}
		worldStorerMock.On("GetAliveCities", ctx).Return([]*entity.City{nil}, error1).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		defer randomerMock.AssertExpectations(t)

		out := &bytes.Buffer{}

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          &bytes.Buffer{},
			out:         out,
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 0,
		}

		err := s.Finalize(ctx)
		require.ErrorIs(t, err, error1)
		require.Equal(t, "", out.String())
	})
}

func Test_SimulationEngine_loadInputToWorld(t *testing.T) {
	var cityNil *entity.City
	city1 := entity.NewCity("City1")
	city2 := entity.NewCity("City2")
	city3 := entity.NewCity("City3")
	city4 := entity.NewCity("City4")
	city5 := entity.NewCity("City5")
	city6 := entity.NewCity("City6")
	city7 := entity.NewCity("City7")

	t.Run("Case 1: Empty file", func(t *testing.T) {
		ctx := context.Background()

		worldStorerMock := &WorldStorerMock{}
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		defer randomerMock.AssertExpectations(t)

		input := ""
		in := strings.NewReader(input)

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          in,
			out:         &bytes.Buffer{},
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 0,
		}

		err := s.loadInputToWorld(ctx)
		require.NoError(t, err)
	})

	t.Run("Case 1: Correctly formed file", func(t *testing.T) {
		ctx := context.Background()

		worldStorerMock := &WorldStorerMock{}
		// City1
		worldStorerMock.On("GetCity", ctx, "City1").Return(cityNil, nil).Once()
		worldStorerMock.On("GetCity", ctx, "City1").Return(city1, nil).Once()
		worldStorerMock.On("AddCity", ctx, "City1").Return(city1, nil).Once()
		// City2
		worldStorerMock.On("GetCity", ctx, "City2").Return(cityNil, nil).Once()
		worldStorerMock.On("GetCity", ctx, "City2").Return(city2, nil).Once()
		worldStorerMock.On("AddCity", ctx, "City2").Return(city2, nil).Once()
		// City3
		worldStorerMock.On("GetCity", ctx, "City3").Return(cityNil, nil).Once()
		worldStorerMock.On("GetCity", ctx, "City3").Return(city3, nil).Times(2)
		worldStorerMock.On("AddCity", ctx, "City3").Return(city3, nil).Once()
		// City4
		worldStorerMock.On("GetCity", ctx, "City4").Return(cityNil, nil).Once()
		worldStorerMock.On("GetCity", ctx, "City4").Return(city4, nil).Times(2)
		worldStorerMock.On("AddCity", ctx, "City4").Return(city4, nil).Once()
		// City5
		worldStorerMock.On("GetCity", ctx, "City5").Return(cityNil, nil).Once()
		worldStorerMock.On("GetCity", ctx, "City5").Return(city5, nil).Times(3)
		worldStorerMock.On("AddCity", ctx, "City5").Return(city5, nil).Once()
		// City6
		worldStorerMock.On("GetCity", ctx, "City6").Return(cityNil, nil).Once()
		worldStorerMock.On("AddCity", ctx, "City6").Return(city6, nil).Once()
		// City7
		worldStorerMock.On("GetCity", ctx, "City7").Return(cityNil, nil).Once()
		worldStorerMock.On("AddCity", ctx, "City7").Return(city7, nil).Once()
		// Links from City1
		worldStorerMock.On("AddLink", ctx, city1, city2, entity.North).Return(nil).Once()
		worldStorerMock.On("AddLink", ctx, city1, city3, entity.East).Return(nil).Once()
		worldStorerMock.On("AddLink", ctx, city1, city4, entity.South).Return(nil).Once()
		worldStorerMock.On("AddLink", ctx, city1, city5, entity.West).Return(nil).Once()
		// Links from City2
		worldStorerMock.On("AddLink", ctx, city2, city1, entity.East).Return(nil).Once()
		worldStorerMock.On("AddLink", ctx, city2, city4, entity.South).Return(nil).Once()
		// Links from City3
		worldStorerMock.On("AddLink", ctx, city3, city5, entity.West).Return(nil).Once()
		worldStorerMock.On("AddLink", ctx, city3, city7, entity.East).Return(nil).Once()
		worldStorerMock.On("AddLink", ctx, city3, city5, entity.South).Return(nil).Once()
		// Links from City4
		worldStorerMock.On("AddLink", ctx, city4, city3, entity.East).Return(nil).Once()
		worldStorerMock.On("AddLink", ctx, city4, city5, entity.North).Return(nil).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		defer randomerMock.AssertExpectations(t)

		input := `
City1 north=City2 east=City3 south=City4 west=City5
City2 east=City1 south=City4
	City3 west=City5 east=City7 south=City5
City4 east=City3 north=City5
City6

		`
		in := strings.NewReader(input)

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          in,
			out:         &bytes.Buffer{},
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 0,
		}

		err := s.loadInputToWorld(ctx)
		require.NoError(t, err)
	})

	t.Run("Case 3: Incorrect direction", func(t *testing.T) {
		ctx := context.Background()

		worldStorerMock := &WorldStorerMock{}
		// City1
		worldStorerMock.On("GetCity", ctx, "City1").Return(cityNil, nil).Once()
		worldStorerMock.On("AddCity", ctx, "City1").Return(city1, nil).Once()
		// City2
		worldStorerMock.On("GetCity", ctx, "City2").Return(cityNil, nil).Once()
		worldStorerMock.On("AddCity", ctx, "City2").Return(city2, nil).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		defer randomerMock.AssertExpectations(t)

		input := `
City1 test=City2	
City2 east=City1 south=City4
		`
		in := strings.NewReader(input)

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          in,
			out:         &bytes.Buffer{},
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 0,
		}

		err := s.loadInputToWorld(ctx)
		require.ErrorIs(t, err, entity.ErrParseCityDefinition)
	})

	t.Run("Case 3: Incorrect format", func(t *testing.T) {
		ctx := context.Background()

		worldStorerMock := &WorldStorerMock{}
		// City1
		worldStorerMock.On("GetCity", ctx, "City1").Return(cityNil, nil).Once()
		worldStorerMock.On("AddCity", ctx, "City1").Return(city1, nil).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		defer randomerMock.AssertExpectations(t)

		input := `
City1 test=City2=City3	
City2 east=City1 south=City4
		`
		in := strings.NewReader(input)

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          in,
			out:         &bytes.Buffer{},
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 0,
		}

		err := s.loadInputToWorld(ctx)
		require.ErrorIs(t, err, entity.ErrParseCityDefinition)
	})

	t.Run("Case 4: Error in AddCity", func(t *testing.T) {
		ctx := context.Background()

		error1 := fmt.Errorf("error 1")

		worldStorerMock := &WorldStorerMock{}
		// City1
		worldStorerMock.On("GetCity", ctx, "City1").Return(cityNil, nil).Once()
		worldStorerMock.On("AddCity", ctx, "City1").Return(city1, nil).Once()
		// City2
		worldStorerMock.On("GetCity", ctx, "City2").Return(cityNil, nil).Once()
		worldStorerMock.On("AddCity", ctx, "City2").Return(city2, error1).Once()
		defer worldStorerMock.AssertExpectations(t)

		randomerMock := &RandomerMock{}
		defer randomerMock.AssertExpectations(t)

		input := `
City1 north=City2	
City2 east=City1 south=City4
		`
		in := strings.NewReader(input)

		s := SimulationEngine{
			world:       worldStorerMock,
			random:      randomerMock,
			in:          in,
			out:         &bytes.Buffer{},
			totalSteps:  0,
			maxSteps:    10,
			startAliens: 0,
		}

		err := s.loadInputToWorld(ctx)
		require.ErrorIs(t, err, error1)
	})
}
