package simulator

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/jpraynaud/alien-invasion-simulator/pkg/simulator/entity"
)

// World represents the world with cities, aliens, and links
type World struct {
	// Map cities to their names
	cityMap map[string]*entity.City

	// Map aliens to their ids
	alienMap map[int]*entity.Alien

	// Map trapped aliens to their ids
	trappedAlienMap map[int]*entity.Alien

	// Map cities to aliens
	cityAlienMap map[*entity.City]*entity.Alien

	// Map links from cities to other cities
	linksFromCityMap map[*entity.City][]*entity.City
}

var _ WorldStorer = (*World)(nil)

// NewWorld is a world constructor
func NewWorld() *World {
	var (
		cityMap          = make(map[string]*entity.City)
		alienMap         = make(map[int]*entity.Alien)
		trappedAlienMap  = make(map[int]*entity.Alien)
		cityAlienMap     = make(map[*entity.City]*entity.Alien)
		linksFromCityMap = make(map[*entity.City][]*entity.City)
	)
	return &World{
		cityMap:          cityMap,
		alienMap:         alienMap,
		trappedAlienMap:  trappedAlienMap,
		cityAlienMap:     cityAlienMap,
		linksFromCityMap: linksFromCityMap,
	}
}

// GetCity retrieves a city
func (w *World) GetCity(ctx context.Context, cityName string) (*entity.City, error) {
	log.WithFields(log.Fields{
		"cityName": cityName,
	}).Debug("GetCity")

	// If city is already registered, return it
	var city *entity.City
	if cityFound, found := w.cityMap[cityName]; found {
		return cityFound, nil
	}

	return city, nil
}

// GetAliveCities retrieves the list of non destroyed cities
func (w *World) GetAliveCities(ctx context.Context) ([]*entity.City, error) {
	log.Debug("GetAliveCities")

	var cities []*entity.City
	for _, city := range w.cityMap {
		cities = append(cities, city)
	}

	return cities, nil
}

// AddCity adds a city to the world
func (w *World) AddCity(ctx context.Context, cityName string) (*entity.City, error) {
	log.WithFields(log.Fields{
		"cityName": cityName,
	}).Debug("AddCity")

	// Check non empty city name
	var city *entity.City
	if cityName == "" {
		return city, entity.ErrEmptyCityName
	}

	// Can't add twice the same city
	if _, found := w.cityMap[cityName]; found {
		return city, entity.ErrDuplicateCity
	}

	// Create a new city and register it
	newCity := entity.NewCity(cityName)
	w.cityMap[newCity.Name] = newCity

	return newCity, nil
}

// DestroyCity destroys a city
func (w *World) DestroyCity(ctx context.Context, city *entity.City) error {
	log.WithFields(log.Fields{
		"city": city,
	}).Debug("DestroyCity")

	if citiesFrom, found := w.linksFromCityMap[city]; found {
		for _, cityFrom := range citiesFrom {
			err := cityFrom.RemoveCityTo(city)
			if err != nil {
				return err
			}
		}
	}

	delete(w.cityMap, city.Name)
	delete(w.cityAlienMap, city)
	delete(w.linksFromCityMap, city)

	return nil
}

// AddLink adds a link from a city to another city given a direction
func (w *World) AddLink(ctx context.Context, cityFrom, cityTo *entity.City, direction entity.Direction) error {
	log.WithFields(log.Fields{
		"cityFrom":  cityFrom,
		"cityTo":    cityTo,
		"direction": direction,
	}).Debug("AddLink")

	// Check that cityFrom and cityTo are not null
	if cityFrom == nil {
		return entity.ErrMissingCity
	}
	if cityTo == nil {
		return entity.ErrMissingCity
	}

	// Check that cityFrom and cityTo are not the same
	if cityFrom.Name == cityTo.Name {
		return entity.ErrLinkSameCity
	}

	// Check that cityFrom and cityTo are already added
	cityFromFound, err := w.GetCity(ctx, cityFrom.Name)
	if err != nil {
		return err
	}
	if cityFromFound == nil {
		return entity.ErrUnknownCity
	}
	cityToFound, err := w.GetCity(ctx, cityTo.Name)
	if err != nil {
		return err
	}
	if cityToFound == nil {
		return entity.ErrUnknownCity
	}

	// Check that a link is not already registered
	cityToRegistered, err := cityFrom.GetCityTo(direction)
	if err != nil {
		return err
	}
	if cityToRegistered != nil && cityToRegistered != cityTo {
		return entity.ErrAlreadyExistsLink
	}

	// Save `city to` and direction in `city from`
	err = cityFrom.SetCityTo(cityTo, direction)
	if err != nil {
		return err
	}

	// Register mapping from new city
	citiesFrom, found := w.linksFromCityMap[cityTo]
	if !found {
		citiesFrom = make([]*entity.City, 0)
	}
	citiesFrom = append(citiesFrom, cityFrom)
	w.linksFromCityMap[cityTo] = citiesFrom

	return nil
}

// GetAlien retrieves an alien
func (w *World) GetAlien(ctx context.Context, alienID int) (*entity.Alien, error) {
	log.WithFields(log.Fields{
		"alienID": alienID,
	}).Debug("GetAlien")

	var alien *entity.Alien
	if alienFound, found := w.alienMap[alienID]; found {
		return alienFound, nil
	}

	return alien, nil
}

// AddAlien adds an alien to the world
func (w *World) AddAlien(ctx context.Context, alienID int) (*entity.Alien, error) {
	log.WithFields(log.Fields{
		"alienID": alienID,
	}).Debug("AddAlien")

	// Can't add twice the same alien
	if _, found := w.alienMap[alienID]; found {
		var alien *entity.Alien
		return alien, entity.ErrDuplicateAlien
	}

	// Create a new alien and register it
	newAlien := entity.NewAlien(alienID)
	w.alienMap[newAlien.AlienID] = newAlien

	return newAlien, nil
}

// MoveAlien moves an alien to a city
func (w *World) MoveAlien(ctx context.Context, alien *entity.Alien, city *entity.City) error {
	log.WithFields(log.Fields{
		"alien":  alien,
		"cityTo": city,
	}).Debug("MoveAlien")

	if alien == nil {
		return entity.ErrMissingAlien
	}

	if city == nil {
		return entity.ErrMissingCity
	}

	alienFound, err := w.GetAlien(ctx, alien.AlienID)
	if err != nil {
		return err
	}
	if alienFound == nil {
		return entity.ErrUnknownAlien
	}

	cityFound, err := w.GetCity(ctx, city.Name)
	if err != nil {
		return err
	}
	if cityFound == nil {
		return entity.ErrUnknownCity
	}

	if alien.City != nil {
		delete(w.cityAlienMap, alien.City)
	}

	alien.City = city
	w.cityAlienMap[alien.City] = alien

	return nil
}

// IsTrappedAlien checks if an alien is trapped
func (w *World) IsTrappedAlien(ctx context.Context, alien *entity.Alien) (bool, error) {
	log.WithFields(log.Fields{
		"alien": alien,
	}).Debug("IsTrappedAlien")

	if alien == nil {
		return false, entity.ErrMissingAlien
	}

	alienFound, err := w.GetAlien(ctx, alien.AlienID)
	if err != nil {
		return false, err
	}
	if alienFound != nil {
		_, isTrapped := w.trappedAlienMap[alien.AlienID]
		return isTrapped, nil
	}

	return false, nil
}

// TrapAlien traps an alien
func (w *World) TrapAlien(ctx context.Context, alien *entity.Alien) error {
	log.WithFields(log.Fields{
		"alien": alien,
	}).Debug("TrapAlien")

	if alien == nil {
		return entity.ErrMissingAlien
	}

	alienFound, err := w.GetAlien(ctx, alien.AlienID)
	if err != nil {
		return err
	}
	if alienFound != nil {
		delete(w.cityAlienMap, alienFound.City)
		w.trappedAlienMap[alienFound.AlienID] = alienFound
		return nil
	}

	return entity.ErrMissingAlien
}

// GetAlienAtCity retrieves the alien at a city
func (w *World) GetAlienAtCity(ctx context.Context, city *entity.City) (*entity.Alien, error) {
	log.WithFields(log.Fields{
		"city": city,
	}).Debug("GetAlienAtCity")

	var alien *entity.Alien
	if city == nil {
		return alien, entity.ErrMissingCity
	}

	cityFound, err := w.GetCity(ctx, city.Name)
	if err != nil {
		return alien, err
	}
	if cityFound == nil {
		return alien, entity.ErrUnknownCity
	}

	if alienAtCity, found := w.cityAlienMap[city]; found {
		return alienAtCity, nil
	}

	return alien, nil
}

// GetUntrappedAliens retrieves the list of untrapped aliens
func (w *World) GetUntrappedAliens(ctx context.Context) ([]*entity.Alien, error) {
	log.Debug("GetUntrappedAliens")

	var aliens []*entity.Alien
	for _, alien := range w.alienMap {
		if _, found := w.trappedAlienMap[alien.AlienID]; !found {
			aliens = append(aliens, alien)
		}
	}

	return aliens, nil
}
