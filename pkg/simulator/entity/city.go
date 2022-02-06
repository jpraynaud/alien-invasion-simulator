package entity

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// City represents a City
type City struct {
	// City name
	Name string

	// Cities where to go from this city
	North, East, South, West *City
}

// NewCity is a city constructor
func NewCity(Name string) *City {
	return &City{
		Name: Name,
	}
}

// GetCityTo retrieves the destination city given a direction
func (c *City) GetCityTo(direction Direction) (*City, error) {
	log.WithFields(log.Fields{
		"city":      c,
		"direction": direction,
	}).Debug("GetCityTo")

	switch direction {
	case North:
		return c.North, nil
	case East:
		return c.East, nil
	case South:
		return c.South, nil
	case West:
		return c.West, nil
	default:
		var city *City
		return city, ErrUnknownDirection
	}
}

// SetCityTo sets the destination city given a direction
func (c *City) SetCityTo(cityTo *City, direction Direction) error {
	log.WithFields(log.Fields{
		"city":      c,
		"cityTo":    cityTo,
		"direction": direction,
	}).Debug("SetCityTo")

	switch direction {
	case North:
		c.North = cityTo
	case East:
		c.East = cityTo
	case South:
		c.South = cityTo
	case West:
		c.West = cityTo
	default:
		return ErrUnknownDirection
	}
	return nil
}

// RemoveCityTo removes the destination city in any direction if it exists
func (c *City) RemoveCityTo(city *City) error {
	log.WithFields(log.Fields{
		"city":   c,
		"cityTo": city,
	}).Debug("RemoveCityTo")

	switch {
	case c.North == city:
		c.North = nil
	case c.East == city:
		c.East = nil
	case c.South == city:
		c.South = nil
	case c.West == city:
		c.West = nil
	default:
		return ErrUnknownCity
	}

	return nil
}

// GetAvailableLinks retrieves the available links from this city
func (c *City) GetAvailableLinks() map[Direction]*City {
	log.WithFields(log.Fields{
		"city": c,
	}).Debug("GetAvailableLinks")

	links := make(map[Direction]*City)
	if c.North != nil {
		links[North] = c.North
	}
	if c.East != nil {
		links[East] = c.East
	}
	if c.South != nil {
		links[South] = c.South
	}
	if c.West != nil {
		links[West] = c.West
	}
	log.WithFields(log.Fields{
		"links": links,
	}).Debug("GetAvailableLinks result")
	return links
}

// String implementats Stringer interface for a city
func (c *City) String() string {
	chunks := []string{c.Name}
	if c.North != nil {
		chunks = append(chunks, fmt.Sprintf("north=%s", c.North.Name))
	}
	if c.East != nil {
		chunks = append(chunks, fmt.Sprintf("east=%s", c.East.Name))
	}
	if c.South != nil {
		chunks = append(chunks, fmt.Sprintf("south=%s", c.South.Name))
	}
	if c.West != nil {
		chunks = append(chunks, fmt.Sprintf("west=%s", c.West.Name))
	}
	return strings.Join(chunks, " ")
}
