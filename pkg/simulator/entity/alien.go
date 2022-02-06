package entity

import "fmt"

// Alien represents an alien
type Alien struct {
	// Alien identifier
	AlienID int

	// Current city where alien is
	City *City
}

// NewAlien is an alien constructor
func NewAlien(alienID int) *Alien {
	return &Alien{
		AlienID: alienID,
	}
}

// String implements Stringer interface for an alien
func (a *Alien) String() string {
	return fmt.Sprintf("Alien #%d", a.AlienID)
}
