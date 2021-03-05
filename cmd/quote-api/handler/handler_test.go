package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/johanronkko/quote-service/cmd/quote-api/handler"
	"github.com/johanronkko/quote-service/internal/business/data/quote"
	"github.com/johanronkko/quote-service/internal/business/mock"
	"github.com/johanronkko/quote-service/internal/business/region"
	"github.com/johanronkko/quote-service/internal/business/validate"
	"github.com/matryer/is"
)

type NoDataResponse struct {
	Code    int     `json:"code"`
	Error   *string `json:"error"`
	Success bool    `json:"success"`
}

func decodePayload(is *is.I, r io.Reader, v interface{}) {
	data, err := ioutil.ReadAll(r)
	is.NoErr(err)
	err = json.Unmarshal(data, v)
	is.NoErr(err)
}

func TestHandleHealthCheck(t *testing.T) {
	is := is.New(t)

	// Setup handler.
	h := New()

	// Make request.
	r := httptest.NewRequest(http.MethodGet, "/api.v1/healthcheck", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	// Assert response HTTP headers.
	is.Equal(w.Code, http.StatusOK)

	// Assert response payload.
	var resp NoDataResponse
	decodePayload(is, w.Body, &resp)
	is.Equal(resp.Code, http.StatusOK)
	is.True(resp.Success)
	is.Equal(resp.Error, nil)
}

type QuoteResponse struct {
	NoDataResponse
	Data struct {
		Quote quote.Info `json:"quote"`
	} `json:"data"`
}

type QuotesResponse struct {
	NoDataResponse
	Data struct {
		Quotes []quote.Info `json:"quotes"`
	} `json:"data"`
}

func TestHandleListQuotes(t *testing.T) {
	t.Run("ok", func(t *testing.T) {

		cases := []struct {
			Name string

			NumQuotes int
		}{
			{"0 quotes", 0},
			{"3 quotes", 3},
			{"10 quotes", 10},
		}
		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				is := is.New(t)

				// Mock.
				q := &mock.Quote{}
				q.QueryCall.Returns.Quotes = createTestQuotes(tc.NumQuotes)

				// Setup handler.
				h := New()
				h.Quote = q

				// Make request.
				r := httptest.NewRequest(http.MethodGet, "/api.v1/quotes/", nil)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)

				// Assert response HTTP headers.
				is.Equal(w.Code, http.StatusOK)

				// Assert response payload.
				var resp QuotesResponse
				decodePayload(is, w.Body, &resp)
				is.Equal(resp.Code, http.StatusOK)
				is.True(resp.Success)
				is.Equal(resp.Error, nil)
				is.Equal(len(resp.Data.Quotes), tc.NumQuotes)
			})
		}
	})

	t.Run("service error", func(t *testing.T) {
		is := is.New(t)

		// Mock services.
		q := &mock.Quote{}
		q.QueryCall.Returns.Err = fmt.Errorf("some error")

		// Setup handler.
		h := New()
		h.Quote = q

		// Make request.
		r := httptest.NewRequest(http.MethodGet, "/api.v1/quotes/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		// Assert response HTTP headers.
		is.Equal(w.Code, http.StatusInternalServerError)

		// Assert response payload.
		var resp NoDataResponse
		decodePayload(is, w.Body, &resp)
		is.Equal(resp.Code, http.StatusInternalServerError)
		is.True(!resp.Success)
		is.Equal(*resp.Error, "internal server error")
	})

	t.Run("bad id format", func(t *testing.T) {
		is := is.New(t)

		quoteID := "badFormat"

		// Mock services.
		q := &mock.Quote{}

		// Setup handler.
		h := New()
		h.Quote = q

		// Make request.
		r := httptest.NewRequest(http.MethodGet, "/api.v1/quotes/"+quoteID, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		// Assert response HTTP headers.
		is.Equal(w.Code, http.StatusBadRequest)

		// Assert response payload.
		var resp NoDataResponse
		decodePayload(is, w.Body, &resp)
		is.Equal(resp.Code, http.StatusBadRequest)
		is.True(!resp.Success)
		is.Equal(*resp.Error, validate.ErrInvalidID.Error())
	})
}

func TestHandleGetQuote(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		is := is.New(t)

		quoteID := validate.GenerateID()

		// Mock.
		q := &mock.Quote{}
		q.QueryByIDCall.Returns.Info = createTestQuote(quoteID)

		// Setup handler.
		h := New()
		h.Quote = q

		// Make request.
		r := httptest.NewRequest(http.MethodGet, "/api.v1/quotes/"+quoteID, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		// Assert response HTTP headers.
		is.Equal(w.Code, http.StatusOK)

		// Assert response payload.
		var resp QuoteResponse
		decodePayload(is, w.Body, &resp)
		is.Equal(resp.Code, http.StatusOK)
		is.True(resp.Success)
		is.Equal(resp.Error, nil)
		is.Equal(resp.Data.Quote, q.QueryByIDCall.Returns.Info)
	})

	t.Run("service error", func(t *testing.T) {
		cases := []struct {
			Name string

			QuoteID    string
			ServiceErr error
			ErrMsg     string

			StatusCode int
		}{
			{"unknown error", validate.GenerateID(), errors.New("some error"), "internal server error", http.StatusInternalServerError},
			{"quote not found", validate.GenerateID(), quote.ErrNotFound, quote.ErrNotFound.Error(), http.StatusBadRequest},
		}
		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				is := is.New(t)

				// Mock services.
				q := &mock.Quote{}
				q.QueryByIDCall.Returns.Err = tc.ServiceErr

				// Setup handler.
				h := New()
				h.Quote = q

				// Make request.
				r := httptest.NewRequest(http.MethodGet, "/api.v1/quotes/"+tc.QuoteID, nil)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)

				// Assert response HTTP headers.
				is.Equal(w.Code, tc.StatusCode)

				// Assert response payload.
				var resp NoDataResponse
				decodePayload(is, w.Body, &resp)
				is.Equal(resp.Code, tc.StatusCode)
				is.True(!resp.Success)
				is.Equal(*resp.Error, tc.ErrMsg)
			})
		}
	})
}

func TestHandleAddQuote(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		is := is.New(t)

		nq := createTestNewQuote()
		nq.From.CountryCode = "US"
		nq.Weight = 500

		// Mock services.
		q := &mock.Quote{}
		q.CreateCall.Returns.Info = quote.Info{
			ID:           validate.GenerateID(),
			To:           nq.To,
			From:         nq.From,
			Weight:       nq.Weight,
			ShipmentCost: 2.5 * 2000, // Outside EU * huge package
		}

		// Setup handler.
		h := New()
		h.Quote = q

		// Make request.
		reqBody, err := json.Marshal(&nq)
		is.NoErr(err)
		r := httptest.NewRequest(http.MethodPost, "/api.v1/quotes/", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		// Assert response HTTP headers.
		is.Equal(w.Code, http.StatusCreated)

		// Assert response payload.
		var resp QuoteResponse
		decodePayload(is, w.Body, &resp)
		is.Equal(resp.Code, http.StatusCreated)
		is.Equal(resp.Error, nil)
		is.Equal(resp.Data.Quote, q.CreateCall.Returns.Info)
	})

	t.Run("service error", func(t *testing.T) {
		cases := []struct {
			Name string

			ServiceErr error
			ErrMsg     string

			StatusCode int
		}{
			{"unknown error", errors.New("some error"), "internal server error", http.StatusInternalServerError},
			{"unsupported country code", region.ErrUnsupportedCountryCode, region.ErrUnsupportedCountryCode.Error(), http.StatusBadRequest},
		}
		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				is := is.New(t)

				// Mock services.
				q := &mock.Quote{}
				q.CreateCall.Returns.Err = tc.ServiceErr

				// Setup handler.
				h := New()
				h.Quote = q

				// Make request.
				nq := createTestNewQuote()
				reqBody, err := json.Marshal(&nq)
				is.NoErr(err)
				r := httptest.NewRequest(http.MethodPost, "/api.v1/quotes/", bytes.NewBuffer(reqBody))
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)

				// Assert response HTTP headers.
				is.Equal(w.Code, tc.StatusCode)

				// Assert response payload.
				var resp NoDataResponse
				decodePayload(is, w.Body, &resp)
				is.Equal(resp.Code, tc.StatusCode)
				is.True(!resp.Success)
				is.Equal(*resp.Error, tc.ErrMsg)
			})
		}
	})

	// TODO: bad request body (decode error)

	// TODO: field validation
}

func createTestQuotes(numQuotes int) []quote.Info {
	qs := []quote.Info{}
	for i := 0; i < numQuotes; i++ {
		qs = append(qs, createTestQuote(validate.GenerateID()))
	}
	return qs
}

func createTestQuote(id string) quote.Info {
	return quote.Info{
		ID:           id,
		From:         createTestCustomer("John Doe", "US"),
		To:           createTestCustomer("Sven Svensson", "SE"),
		Weight:       500,
		ShipmentCost: 1250,
	}
}

func createTestNewQuote() quote.NewQuote {
	return quote.NewQuote{
		To: quote.Customer{
			Name:        "Sven Svensson",
			Email:       "sven.svensson@test.com",
			Address:     "Testgatan 42B, Göteborg 12345",
			CountryCode: "SV",
		},
		From: quote.Customer{
			Name:        "John Doe",
			Email:       "john.doe@test.com",
			Address:     "Teststreet 4242, Blaine 55434",
			CountryCode: "US",
		},
		Weight: 500,
	}
}

func createTestCustomer(name string, ccode string) quote.Customer {
	return quote.Customer{
		Name:        name,
		Email:       "example@test.com",
		Address:     "Vasagatan 5B, Göteborg 41124",
		CountryCode: ccode,
	}
}
