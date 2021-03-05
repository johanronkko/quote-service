package quote

import (
	"context"
	"testing"

	"github.com/johanronkko/quote-service/internal/business/data/schema"
	"github.com/johanronkko/quote-service/internal/business/tests"
	"github.com/matryer/is"
)

func TestQuote(t *testing.T) {
	is := is.New(t)

	db := tests.NewUnit(t)
	q := New(db)

	ctx := context.Background()

	// Query empty database.
	quotes, err := q.Query(ctx)
	is.NoErr(err)
	is.Equal(len(quotes), 0)

	// Create quote.
	nq := NewQuote{
		To: Customer{
			Name:        "Sven Svensson",
			Email:       "sven.svensson@test.com",
			Address:     "Testgatan 42B, Göteborg 12345",
			CountryCode: "SV",
		},
		From: Customer{
			Name:        "John Doe",
			Email:       "john.doe@test.com",
			Address:     "Teststreet 4242, Blaine 55434",
			CountryCode: "US",
		},
		Weight: 500,
	}
	quote, err := q.Create(ctx, nq)
	is.NoErr(err)

	// Query by ID returns correct quote.
	saved, err := q.QueryByID(ctx, quote.ID)
	is.NoErr(err)
	is.Equal(quote, saved)

	// Query database with 1 newly added quote and 3 seeded quotes.
	err = schema.Seed(db)
	is.NoErr(err)
	quotes, err = q.Query(ctx)
	is.NoErr(err)
	is.Equal(len(quotes), 1+3)
}

func TestCalcShipmentCost(t *testing.T) {
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
				got, err := calcShipmentCost(tc.Weight, tc.CountryCode)
				is.NoErr(err)
				is.Equal(got, tc.Want)
			})
		}
	})

	t.Run("invalid country code", func(t *testing.T) {
		cases := []struct {
			Name string

			CountryCode string
		}{
			{"not supported", "nn"},
			{"bad format", "banana"},
		}
		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				is := is.New(t)
				_, err := calcShipmentCost(42, tc.CountryCode)
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
				_, err := calcShipmentCost(tc.Weight, "sv")
				is.True(err != nil)
			})
		}
	})
}
