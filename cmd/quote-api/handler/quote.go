package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/johanronkko/quote-service/internal/business/data/quote"
	"github.com/johanronkko/quote-service/internal/business/validate"
	"github.com/matryer/way"
)

// Quote manages the set of API's for quote access.
type Quote interface {
	// Query retrieves a list of existing quotes.
	Query(ctx context.Context) ([]quote.Info, error)
	// QueryByID retrieves the quote with id. Returns quote.ErrNotFound if
	// quote not found.
	QueryByID(ctx context.Context, id string) (quote.Info, error)
	// Create adds a quote to the system.
	Create(ctx context.Context, nq quote.NewQuote) (quote.Info, error)
}

func (h *Handler) handleGetQuote() http.HandlerFunc {
	type response struct {
		Quote quote.Info `json:"quote"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := way.Param(r.Context(), "id")
		if err := validate.CheckID(id); err != nil {
			respond(w, r, http.StatusBadRequest, err)
			return
		}
		q, err := h.Quote.QueryByID(r.Context(), id)
		if err == quote.ErrNotFound {
			respond(w, r, http.StatusBadRequest, err)
			return
		} else if err != nil {
			respond(w, r, http.StatusInternalServerError, fmt.Errorf("internal server error"))
			return
		}
		respond(w, r, http.StatusOK, &response{q})
	}
}

func (h *Handler) handleListQuotes() http.HandlerFunc {
	type response struct {
		Quotes []quote.Info `json:"quotes"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		qs, err := h.Quote.Query(r.Context())
		if err != nil {
			respond(w, r, http.StatusInternalServerError, fmt.Errorf("internal server error"))
			return
		}
		respond(w, r, http.StatusOK, &response{qs})
	}

}

func (h *Handler) handleAddQuote() http.HandlerFunc {
	type response struct {
		Quote quote.Info `json:"quote"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var nq quote.NewQuote
		decode(w, r, &nq)
		q, _ := h.Quote.Create(r.Context(), nq)
		respond(w, r, http.StatusCreated, &response{q})
	}
}
