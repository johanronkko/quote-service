package quote_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/johanronkko/quote-service/internal/business/data/quote"
	"github.com/johanronkko/quote-service/internal/business/data/schema"
	"github.com/johanronkko/quote-service/internal/business/tests"
	"github.com/matryer/is"
)

// TODO: test validation of ID and NewQuote + ShipmentCost error case.

type TestShipmentCostCalculator struct {
	price float64
	err   error
}

func (c TestShipmentCostCalculator) ShipmentCost(weight int, ccode string) (float64, error) {
	return c.price, c.err
}

func TestQuote(t *testing.T) {
	is := is.New(t)

	db := tests.NewUnit(t)
	calc := &TestShipmentCostCalculator{
		price: 1250,
	}
	q := quote.New(db, calc)

	ctx := context.Background()

	nq := quote.NewQuote{
		To: quote.Customer{
			Name:        "Sven Svensson",
			Email:       "sven.svensson@test.com",
			Address:     "Testgatan 42B, Göteborg 12345",
			CountryCode: "sv",
		},
		From: quote.Customer{
			Name:        "John Doe",
			Email:       "john.doe@test.com",
			Address:     "Teststreet 4242, Blaine 55434",
			CountryCode: "us",
		},
		Weight: 500,
	}

	// Query empty database.
	quotes, err := q.Query(ctx)
	is.NoErr(err)
	is.Equal(len(quotes), 0)

	// Query database with 3 quotes.
	err = schema.Seed(db)
	is.NoErr(err)
	quotes, err = q.Query(ctx)
	is.NoErr(err)
	is.Equal(len(quotes), 3)

	// Create quote.
	quote, err := q.Create(ctx, nq)
	is.NoErr(err)
	is.True(cmp.Diff(quote.To, nq.To) == "")
	is.True(cmp.Diff(quote.From, nq.From) == "")
	is.True(cmp.Diff(quote.Weight, nq.Weight) == "")

	// Query by ID returns correct quote.
	saved, err := q.QueryByID(ctx, quote.ID)
	is.NoErr(err)
	is.True(cmp.Diff(quote, saved) == "")
}
