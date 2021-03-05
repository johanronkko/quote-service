package mock

import (
	"github.com/johanronkko/quote-service/internal/business/data/quote"
	"golang.org/x/net/context"
)

// Quote is a mock implementation of quote.Quote.
type Quote struct {
	QueryCall struct {
		Recieves struct {
			Ctx context.Context
		}
		Returns struct {
			Quotes []quote.Info
			Err    error
		}
	}
	QueryByIDCall struct {
		Recieves struct {
			Ctx context.Context
			ID  string
		}
		Returns struct {
			Info quote.Info
			Err  error
		}
	}
	CreateCall struct {
		Recieves struct {
			Ctx context.Context
			Nq  quote.NewQuote
		}
		Returns struct {
			Info quote.Info
			Err  error
		}
	}
}

// Query mocks the Query func of quote.Quote.
func (q *Quote) Query(ctx context.Context) ([]quote.Info, error) {
	return q.QueryCall.Returns.Quotes, q.QueryCall.Returns.Err
}

// QueryByID mocks the QueryByID func of quote.Quote.
func (q *Quote) QueryByID(ctx context.Context, id string) (quote.Info, error) {
	q.QueryByIDCall.Recieves.Ctx = ctx
	q.QueryByIDCall.Recieves.ID = id
	return q.QueryByIDCall.Returns.Info, q.QueryByIDCall.Returns.Err
}

// Create mocks the Create func of quote.Quote.
func (q *Quote) Create(ctx context.Context, nq quote.NewQuote) (quote.Info, error) {
	q.QueryByIDCall.Recieves.Ctx = ctx
	q.CreateCall.Recieves.Nq = nq
	return q.CreateCall.Returns.Info, q.CreateCall.Returns.Err
}
