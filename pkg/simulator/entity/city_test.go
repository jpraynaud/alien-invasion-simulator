package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewCity(t *testing.T) {
	tests := []struct {
		cityName                                 string
		cityNorth, cityEast, citySouth, cityWest *City
		want                                     string
	}{
		{"City1", &City{Name: "CityN"}, nil, nil, nil, "City1 north=CityN"},
		{"City2", &City{Name: "CityN"}, &City{Name: "CityE"}, nil, nil, "City2 north=CityN east=CityE"},
		{"City3", &City{Name: "CityN"}, &City{Name: "CityE"}, &City{Name: "CityS"}, nil, "City3 north=CityN east=CityE south=CityS"},
		{"City4", &City{Name: "CityN"}, &City{Name: "CityE"}, &City{Name: "CityS"}, &City{Name: "CityW"}, "City4 north=CityN east=CityE south=CityS west=CityW"},
		{"City5", nil, &City{Name: "CityE"}, nil, &City{Name: "CityW"}, "City5 east=CityE west=CityW"},
	}

	for _, tt := range tests {
		t.Run(tt.cityName, func(t *testing.T) {
			c := NewCity(tt.cityName)
			c.North = tt.cityNorth
			c.East = tt.cityEast
			c.South = tt.citySouth
			c.West = tt.cityWest
			require.Equal(t, tt.want, c.String())
		})
	}
}
