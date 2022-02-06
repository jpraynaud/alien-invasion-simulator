package simulator

import (
	"context"
	"testing"

	"github.com/jpraynaud/alien-invasion-simulator/pkg/simulator/entity"
	"github.com/stretchr/testify/require"
)

func Test_World_CityScenario(t *testing.T) {
	ctx := context.Background()
	world := NewWorld()

	cityNameA := "CityA"
	cityNameB := "CityB"

	// No empty city name allowed
	cityEmpty, err := world.AddCity(ctx, "")
	require.ErrorIs(t, err, entity.ErrEmptyCityName)
	require.Nil(t, cityEmpty)

	// CityA does not exist yet
	cityA, err := world.GetCity(ctx, cityNameA)
	require.NoError(t, err)
	require.Nil(t, cityA)

	// CityB does not exist yet
	cityB, err := world.GetCity(ctx, cityNameB)
	require.NoError(t, err)
	require.Nil(t, cityB)

	// No alive city
	aliveCities, err := world.GetAliveCities(ctx)
	require.NoError(t, err)
	require.Equal(t, []*entity.City(nil), aliveCities)

	// CityA is added
	cityNewA, err := world.AddCity(ctx, cityNameA)
	require.NoError(t, err)
	require.NotNil(t, cityNewA)

	// CityA exists now
	cityA, err = world.GetCity(ctx, cityNameA)
	require.NoError(t, err)
	require.Equal(t, cityNewA, cityA)

	// CityA already exists and can't be added again
	cityDuplicateA, err := world.AddCity(ctx, cityNameA)
	require.ErrorIs(t, err, entity.ErrDuplicateCity)
	require.Nil(t, cityDuplicateA)

	// CityB does not exist yet
	cityB, err = world.GetCity(ctx, cityNameB)
	require.NoError(t, err)
	require.Nil(t, cityB)

	// CityA is an alive city
	aliveCities, err = world.GetAliveCities(ctx)
	require.NoError(t, err)
	require.Equal(t, []*entity.City{cityA}, aliveCities)

	// CityB is added
	cityNewB, err := world.AddCity(ctx, cityNameB)
	require.NoError(t, err)
	require.NotNil(t, cityNewB)

	// CityB exists now
	cityB, err = world.GetCity(ctx, cityNameB)
	require.NoError(t, err)
	require.Equal(t, cityNewB, cityB)

	// CityB already exists and can't be added again
	cityDuplicateB, err := world.AddCity(ctx, cityNameB)
	require.ErrorIs(t, err, entity.ErrDuplicateCity)
	require.Nil(t, cityDuplicateB)

	// CityA and CityB are alive cities
	aliveCities, err = world.GetAliveCities(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t, []*entity.City{cityA, cityB}, aliveCities)

	// CityA is destroyed
	err = world.DestroyCity(ctx, cityA)
	require.NoError(t, err)

	// CityA does not exist any more
	cityA, err = world.GetCity(ctx, cityNameA)
	require.NoError(t, err)
	require.Nil(t, cityA)

	// CityB still exists
	cityB, err = world.GetCity(ctx, cityNameB)
	require.NoError(t, err)
	require.Equal(t, cityNewB, cityB)

	// CityB is the only alive city
	aliveCities, err = world.GetAliveCities(ctx)
	require.NoError(t, err)
	require.Equal(t, []*entity.City{cityB}, aliveCities)

	// CityB is destroyed
	err = world.DestroyCity(ctx, cityB)
	require.NoError(t, err)

	// CityB does not exist any more
	cityB, err = world.GetCity(ctx, cityNameB)
	require.NoError(t, err)
	require.Nil(t, cityB)

	// No more alive city
	aliveCities, err = world.GetAliveCities(ctx)
	require.NoError(t, err)
	require.Equal(t, []*entity.City(nil), aliveCities)
}

func Test_World_AlienScenario(t *testing.T) {
	ctx := context.Background()
	world := NewWorld()

	alienID1 := 1
	alienID2 := 2

	// Alien1 does not exist yet
	alien1, err := world.GetAlien(ctx, alienID1)
	require.NoError(t, err)
	require.Nil(t, alien1)

	// Alien2 does not exist yet
	alien2, err := world.GetAlien(ctx, alienID2)
	require.NoError(t, err)
	require.Nil(t, alien2)

	// No alive alien
	untrappedAliens, err := world.GetUntrappedAliens(ctx)
	require.NoError(t, err)
	require.Equal(t, []*entity.Alien(nil), untrappedAliens)

	// Alien1 is added
	alienNew1, err := world.AddAlien(ctx, alienID1)
	require.NoError(t, err)
	require.NotNil(t, alienNew1)

	// Alien1 exists now
	alien1, err = world.GetAlien(ctx, alienID1)
	require.NoError(t, err)
	require.Equal(t, alienNew1, alien1)

	// Alien1 already exists and can't be added again
	alienDuplicate1, err := world.AddAlien(ctx, alienID1)
	require.ErrorIs(t, err, entity.ErrDuplicateAlien)
	require.Nil(t, alienDuplicate1)

	// Alien1 is not trapped
	trapped1, err := world.IsTrappedAlien(ctx, alien1)
	require.NoError(t, err)
	require.False(t, trapped1)

	// Alien2 does not exist yet
	alien2, err = world.GetAlien(ctx, alienID2)
	require.NoError(t, err)
	require.Nil(t, alien2)

	// Alien1 is an untrapped alien
	untrappedAliens, err = world.GetUntrappedAliens(ctx)
	require.NoError(t, err)
	require.Equal(t, []*entity.Alien{alien1}, untrappedAliens)

	// Alien2 is added
	alienNew2, err := world.AddAlien(ctx, alienID2)
	require.NoError(t, err)
	require.NotNil(t, alienNew2)

	// Alien2 exists now
	alien2, err = world.GetAlien(ctx, alienID2)
	require.NoError(t, err)
	require.Equal(t, alienNew2, alien2)

	// Alien2 already exists and can't be added again
	alienDuplicate2, err := world.AddAlien(ctx, alienID2)
	require.ErrorIs(t, err, entity.ErrDuplicateAlien)
	require.Nil(t, alienDuplicate2)

	// Alien2 is not trapped
	trapped2, err := world.IsTrappedAlien(ctx, alien2)
	require.NoError(t, err)
	require.False(t, trapped2)

	// Alien1 and Alien2 are untrapped aliens
	untrappedAliens, err = world.GetUntrappedAliens(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t, []*entity.Alien{alien1, alien2}, untrappedAliens)

	// Alien1 gets trapped
	err = world.TrapAlien(ctx, alien1)
	require.NoError(t, err)

	// Alien1 is trapped
	trapped1, err = world.IsTrappedAlien(ctx, alien1)
	require.NoError(t, err)
	require.True(t, trapped1)

	// Alien2 is the only untrapped alien
	untrappedAliens, err = world.GetUntrappedAliens(ctx)
	require.NoError(t, err)
	require.Equal(t, []*entity.Alien{alien2}, untrappedAliens)

	// Alien2 gets trapped
	err = world.TrapAlien(ctx, alien2)
	require.NoError(t, err)

	// Alien2 is not trapped
	trapped2, err = world.IsTrappedAlien(ctx, alien2)
	require.NoError(t, err)
	require.True(t, trapped2)

	// No more untrapped alien
	untrappedAliens, err = world.GetUntrappedAliens(ctx)
	require.NoError(t, err)
	require.Equal(t, []*entity.Alien(nil), untrappedAliens)
}

func Test_World_CityAlienScenario(t *testing.T) {
	ctx := context.Background()
	world := NewWorld()

	cityNameA := "CityA"
	cityNameB := "CityB"
	alienID1 := 1

	cityZ := entity.NewCity("CityZ")
	alienZ := entity.NewAlien(1000)

	var cityNull *entity.City
	var alienNull *entity.Alien

	// CityA is added
	cityA, err := world.AddCity(ctx, cityNameA)
	require.NoError(t, err)
	require.NotNil(t, cityA)

	// CityB is added
	cityB, err := world.AddCity(ctx, cityNameB)
	require.NoError(t, err)
	require.NotNil(t, cityB)

	// Alien1 is added
	alien1, err := world.AddAlien(ctx, alienID1)
	require.NoError(t, err)
	require.NotNil(t, alien1)

	// Get alien at null city
	alienFound, err := world.GetAlienAtCity(ctx, cityNull)
	require.ErrorIs(t, err, entity.ErrMissingCity)
	require.Nil(t, alienFound)

	// Get alien at unknown city
	alienFound, err = world.GetAlienAtCity(ctx, cityZ)
	require.ErrorIs(t, err, entity.ErrUnknownCity)
	require.Nil(t, alienFound)

	// Get alien at cityA
	alienFound, err = world.GetAlienAtCity(ctx, cityA)
	require.NoError(t, err)
	require.Nil(t, alienFound)

	// Move null alien to cityA
	err = world.MoveAlien(ctx, alienNull, cityA)
	require.ErrorIs(t, err, entity.ErrMissingAlien)

	// Move Alien1 to null city
	err = world.MoveAlien(ctx, alien1, cityNull)
	require.ErrorIs(t, err, entity.ErrMissingCity)

	// Move unknown alien to cityA
	err = world.MoveAlien(ctx, alienZ, cityA)
	require.ErrorIs(t, err, entity.ErrUnknownAlien)

	// Move Alien1 to unknow city
	err = world.MoveAlien(ctx, alien1, cityZ)
	require.ErrorIs(t, err, entity.ErrUnknownCity)

	// Move Alien1 to CityA
	err = world.MoveAlien(ctx, alien1, cityA)
	require.NoError(t, err)

	// Get alien at CityA
	alienFound, err = world.GetAlienAtCity(ctx, cityA)
	require.NoError(t, err)
	require.Equal(t, alien1, alienFound)

	// Move Alien1 to CityB
	err = world.MoveAlien(ctx, alien1, cityB)
	require.NoError(t, err)

	// Get alien at CityB
	alienFound, err = world.GetAlienAtCity(ctx, cityB)
	require.NoError(t, err)
	require.Equal(t, alien1, alienFound)

	// Get alien at cityA
	alienFound, err = world.GetAlienAtCity(ctx, cityA)
	require.NoError(t, err)
	require.Nil(t, alienFound)
}

func Test_World_LinkScenario(t *testing.T) {
	ctx := context.Background()
	world := NewWorld()

	cityNameA := "CityA"
	cityNameB := "CityB"
	cityNameC := "CityC"
	alienID1 := 1

	cityZ := entity.NewCity("CityZ")

	// Alien is added
	alien1, err := world.AddAlien(ctx, alienID1)
	require.NoError(t, err)
	require.NotNil(t, alien1)

	// CityA does not exist yet
	cityA, err := world.GetCity(ctx, cityNameA)
	require.NoError(t, err)
	require.Nil(t, cityA)

	// CityB does not exist yet
	cityB, err := world.GetCity(ctx, cityNameB)
	require.NoError(t, err)
	require.Nil(t, cityB)

	// AddLink between CityA and CityB not allowed
	err = world.AddLink(ctx, cityA, cityB, entity.North)
	require.ErrorIs(t, err, entity.ErrMissingCity)

	// CityA is added
	cityA, err = world.AddCity(ctx, cityNameA)
	require.NoError(t, err)
	require.NotNil(t, cityA)

	// CityB is added
	cityB, err = world.AddCity(ctx, cityNameB)
	require.NoError(t, err)
	require.NotNil(t, cityB)

	// CityC is added
	cityC, err := world.AddCity(ctx, cityNameC)
	require.NoError(t, err)
	require.NotNil(t, cityC)

	// AddLink between CityA and CityZ not allowed
	err = world.AddLink(ctx, cityA, cityZ, entity.South)
	require.ErrorIs(t, err, entity.ErrUnknownCity)

	// AddLink between CityZ and CityB not allowed
	err = world.AddLink(ctx, cityZ, cityB, entity.East)
	require.ErrorIs(t, err, entity.ErrUnknownCity)

	// AddLink between CityA and CityB with unkown direction
	err = world.AddLink(ctx, cityA, cityB, entity.Direction(0))
	require.ErrorIs(t, err, entity.ErrUnknownDirection)

	// AddLink between CityA and CityB for a direction works
	err = world.AddLink(ctx, cityA, cityB, entity.East)
	require.NoError(t, err)

	// AddLink between CityA and CityC for the same direction does not work
	err = world.AddLink(ctx, cityA, cityC, entity.East)
	require.ErrorIs(t, err, entity.ErrAlreadyExistsLink)

	// AddLink between CityA and CityC for a different direction works
	err = world.AddLink(ctx, cityA, cityC, entity.West)
	require.NoError(t, err)
}
