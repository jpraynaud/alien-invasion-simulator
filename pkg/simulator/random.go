package simulator

import (
	"math/rand"
	"time"

	"github.com/jpraynaud/alien-invasion-simulator/pkg/simulator/entity"
	log "github.com/sirupsen/logrus"
)

// RandomSimple is a simple random integer generator
// Important: it does not rely on crypto safe random generator
type RandomSimple struct {
}

// NewRandomSimple is a simple random generator
func NewRandomSimple() *RandomSimple {
	return &RandomSimple{}
}

// GetRandomInt retrieves a random integer between 0 and n-1 given n
// Returns an error if n <= 0
func (rs *RandomSimple) GetRandomInt(n int) (int, error) {
	rand.Seed(time.Now().UnixNano())
	r := 0
	if n <= 0 {
		return r, entity.ErrRandomOutOfBounds
	}
	r = rand.Intn(n)
	log.WithFields(log.Fields{
		"n":      n,
		"random": r,
	}).Debug("GetRandomInt")

	return r, nil
}
