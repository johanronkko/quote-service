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

	quotes, err := q.Query(ctx)
	is.NoErr(err)
	is.Equal(len(quotes), 0)

	quote, err := q.Create(ctx, nq)
	is.NoErr(err)

	saved, err := q.QueryByID(ctx, quote.ID)
	is.NoErr(err)

	diff := cmp.Diff(quote, saved)
	is.True(diff != "")

	err = schema.Seed(db)
	is.NoErr(err)

	quotes, err = q.Query(ctx)
	is.NoErr(err)
	is.Equal(len(quotes), 1+3) // 1 newly added quote and 3 seeded.
}
