package simulator

import (
	"fmt"
	"testing"

	"github.com/jpraynaud/alien-invasion-simulator/pkg/simulator/entity"
	"github.com/stretchr/testify/require"
)

func Test_RandomSimple(t *testing.T) {
	tests := []struct {
		giveN     int
		wantError error
	}{
		{
			giveN:     -10,
			wantError: entity.ErrRandomOutOfBounds,
		},
		{
			giveN:     -1,
			wantError: entity.ErrRandomOutOfBounds,
		},
		{
			giveN:     0,
			wantError: entity.ErrRandomOutOfBounds,
		},
		{
			giveN:     1,
			wantError: nil,
		},
		{
			giveN:     10,
			wantError: nil,
		},
		{
			giveN:     100,
			wantError: nil,
		},
		{
			giveN:     279,
			wantError: nil,
		},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("Input: %d", tt.giveN)
		t.Run(testName, func(t *testing.T) {
			rs := NewRandomSimple()
			r, err := rs.GetRandomInt(tt.giveN)
			require.Equal(t, tt.wantError, err)
			switch err {
			case nil:
				require.GreaterOrEqual(t, r, 0)
				require.Less(t, r, tt.giveN)
			default:
				require.Equal(t, 0, r)
			}
		})
	}
}
