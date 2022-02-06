package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_City_NewCity(t *testing.T) {
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

func Test_City_Flow(t *testing.T) {

	tests := []struct {
		name          string
		giveDirection Direction
		wantError     error
	}{
		{
			name:          "North",
			giveDirection: North,
			wantError:     nil,
		},
		{
			name:          "East",
			giveDirection: East,
			wantError:     nil,
		},
		{
			name:          "South",
			giveDirection: South,
			wantError:     nil,
		},
		{
			name:          "West",
			giveDirection: West,
			wantError:     nil,
		},
		{
			name:          "UnknownDirection",
			giveDirection: Direction(100),
			wantError:     ErrUnknownDirection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			city1 := NewCity("City1")
			city2 := NewCity("City2")

			cities := city1.GetAvailableLinks()
			require.Equal(t, map[Direction]*City{}, cities)

			cityTo, err := city1.GetCityTo(tt.giveDirection)
			require.Equal(t, err, tt.wantError)
			require.Nil(t, cityTo)

			err = city1.SetCityTo(city2, tt.giveDirection)
			require.Equal(t, err, tt.wantError)

			cityTo, err = city1.GetCityTo(tt.giveDirection)
			require.Equal(t, err, tt.wantError)
			switch tt.wantError {
			case nil:
				require.Equal(t, city2, cityTo)
			default:
				require.Nil(t, cityTo)
			}

			cities = city1.GetAvailableLinks()
			switch tt.wantError {
			case nil:
				require.Equal(t, map[Direction]*City{tt.giveDirection: city2}, cities)
			default:
				require.Equal(t, map[Direction]*City{}, cities)
			}

			err = city1.RemoveCityTo(city2)
			switch tt.wantError {
			case nil:
				require.Equal(t, err, tt.wantError)
			default:
				require.Equal(t, err, ErrUnknownCity)
			}

			cityTo, err = city1.GetCityTo(tt.giveDirection)
			require.Equal(t, err, tt.wantError)
			require.Nil(t, cityTo)

			cities = city1.GetAvailableLinks()
			require.Equal(t, map[Direction]*City{}, cities)
		})
	}
}
