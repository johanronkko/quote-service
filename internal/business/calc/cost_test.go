package calc_test

import (
	"testing"

	. "github.com/johanronkko/quote-service/internal/business/calc"
	"github.com/matryer/is"
)

func TestShipmentCost(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cases := []struct {
			Name string

			CountryCode string
			Weight      int

			Want float64
		}{
			{"small nordic", "no", 0, 100},
			{"medium nordic", "sv", 11, 300},
			{"large nordic", "dk", 26, 500},
			{"huge nordic", "fi", 51, 2000},

			{"small within EU", "fr", 10, 150},
			{"medium within EU", "de", 25, 450},
			{"large within EU", "lt", 50, 750},
			{"huge within EU", "cz", 1000, 3000},

			{"small outide EU", "us", 7, 250},
			{"medium outide EU", "ca", 18, 750},
			{"large outide EU", "br", 29, 1250},
			{"huge outide EU", "jp", 777, 5000},
		}
		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				is := is.New(t)
				calc := NewShipmentCost()
				got, err := calc.ShipmentCost(tc.Weight, tc.CountryCode)
				is.NoErr(err)
				is.Equal(got, tc.Want)
			})
		}
	})

	t.Run("unsupported country code", func(t *testing.T) {
		cases := []struct {
			Name string

			CountryCode string
		}{
			{"Not supported", "nn"},
			{"bad format", "banana"},
		}
		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				is := is.New(t)
				calc := NewShipmentCost()
				_, err := calc.ShipmentCost(42, tc.CountryCode)
				is.True(err != nil)
			})
		}
	})

	t.Run("invalid weight", func(t *testing.T) {
		cases := []struct {
			Name string

			Weight int
		}{
			{"below 0", -1},
			{"above 1000", 1001},
		}
		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				is := is.New(t)
				calc := NewShipmentCost()
				_, err := calc.ShipmentCost(tc.Weight, "sv")
				is.True(err != nil)
			})
		}
	})
}
