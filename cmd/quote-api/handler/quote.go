package handler

import (
	"context"
	"net/http"

	"github.com/johanronkko/quote-service/internal/business/data/quote"
)

// Quote manages the set of API's for quote access.
type Quote interface {
	// Query retrieves a list of existing quotes.
	Query(context.Context) ([]quote.Info, error)
	// QueryByID retrieves the quote with with id.
	QueryByID(context.Context, quote.ID) (quote.Info, error)
	// Create adds a quote to the system.
	Create(context.Context, quote.NewQuote) (quote.Info, error)
}

func (h *Handler) handleGetQuote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	}
}

func (h *Handler) handleListQuotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	}
}

func (h *Handler) handleAddQuote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	}
}
