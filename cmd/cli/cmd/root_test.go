package cmd

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/jpraynaud/alien-invasion-simulator/pkg/simulator/entity"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_runSimulator(t *testing.T) {
	log.SetLevel(log.WarnLevel)

	inputOK := `
City1 north=City2 east=City3 south=City4 west=City5
City2 east=City1 south=City4
	City3 west=City5 east=City7 south=City5
City4 east=City3 north=City5
City6
`
	inputKO := `
City1 test=City2
	`

	tests := []struct {
		name, input                   string
		giveTotalAliens, giveMaxSteps uint
		wantError                     error
	}{
		{
			name:            "Case 1: input KO",
			input:           inputKO,
			giveTotalAliens: 3,
			giveMaxSteps:    10000,
			wantError:       entity.ErrParseCityDefinition,
		},
		{
			name:            "Case 2: input OK + aliens < cities",
			input:           inputOK,
			giveTotalAliens: 3,
			giveMaxSteps:    10000,
			wantError:       nil,
		},
		{
			name:            "Case 3: input OK + aliens >= cities",
			input:           inputOK,
			giveTotalAliens: 100,
			giveMaxSteps:    10,
			wantError:       nil,
		},
		{
			name:            "Case 4: input OK + zero steps",
			input:           inputOK,
			giveTotalAliens: 2,
			giveMaxSteps:    0,
			wantError:       nil,
		},
		{
			name:            "Case 5: input OK + zero aliens",
			input:           inputOK,
			giveTotalAliens: 0,
			giveMaxSteps:    10000,
			wantError:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			in := io.NopCloser(strings.NewReader(tt.input))
			out := &bytes.Buffer{}

			c := &config{
				totalAliens: tt.giveTotalAliens,
				maxSteps:    tt.giveMaxSteps,
				in:          in,
				out:         out,
			}
			err := runSimulator(ctx, c)
			require.Equal(t, tt.wantError, err)
		})
	}

}
