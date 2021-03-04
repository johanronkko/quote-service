package region_test

import (
	"testing"

	. "github.com/johanronkko/quote-service/internal/business/region"
	"github.com/matryer/is"
)

func TestFrom(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cases := []struct {
			Name string

			CountryCode string
			Want        Region
		}{
			{"nordic", "sv", Nordic},
			{"within EU", "fr", WithinEU},
			{"outside EU", "jp", OutsideEU},
			{"uppercase", "CA", OutsideEU},
		}
		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				is := is.New(t)

				got, err := From(tc.CountryCode)
				is.NoErr(err)

				is.Equal(got, tc.Want)
			})
		}
	})

	t.Run("invalid", func(t *testing.T) {
		cases := []struct {
			Name string

			CountryCode string
		}{
			{"bad form", "banana"},
			{"unsupported country", "nn"},
		}
		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				is := is.New(t)

				_, err := From(tc.CountryCode)
				is.Equal(err.Error(), ErrInvalidCountryCode.Error())
			})
		}
	})
}
