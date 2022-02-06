package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewAlien(t *testing.T) {
	a := NewAlien(123)
	require.Equal(t, 123, a.AlienID)
	require.Equal(t, "Alien #123", a.String())
}
