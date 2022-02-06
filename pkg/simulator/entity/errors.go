package entity

import (
	"fmt"
)

var (
	// ErrParseCityDefinition is triggered when a city definition os unparsable
	ErrParseCityDefinition error = fmt.Errorf("impossible to parse the city definition")

	// ErrEmptyCityName is triggered in case of empty city name
	ErrEmptyCityName error = fmt.Errorf("empty city name not allowed")

	// ErrDuplicateCity is triggered in case of duplicate city
	ErrDuplicateCity error = fmt.Errorf("city is already registered")

	// ErrMissingCity is triggered in case of missing city
	ErrMissingCity error = fmt.Errorf("city is missing")

	// ErrUnknownCity is triggered in case of unknown city
	ErrUnknownCity error = fmt.Errorf("city is unknown")

	// ErrLinkSameCity is triggered in case of adding a link between the same city
	ErrLinkSameCity error = fmt.Errorf("no possible link between same city")

	// ErrDuplicateAlien is triggered in case of duplicate alien
	ErrDuplicateAlien error = fmt.Errorf("duplicate alien not allowed")

	// ErrMissingAlien is triggered in case of missing alien
	ErrMissingAlien error = fmt.Errorf("alien is missing")

	// ErrUnknownAlien is triggered in case of unknown alien
	ErrUnknownAlien error = fmt.Errorf("alien is unknown")

	// ErrUnknownDirection is triggered when an unknown Direction is provided
	ErrUnknownDirection error = fmt.Errorf("unknown direction provided")

	// ErrAlreadyExistsLink is triggered when a link between two cities already exists
	ErrAlreadyExistsLink error = fmt.Errorf("a link already exists between the two cities")

	// ErrRandomOutOfBounds is trigerred when the random number generation is not possible
	ErrRandomOutOfBounds error = fmt.Errorf("random input out of bounds")

	// ErrContextCancelled is trigerred when the context is cancelled
	ErrContextCancelled error = fmt.Errorf("the context was cancelled")
)
