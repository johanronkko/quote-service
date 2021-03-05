package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/johanronkko/quote-service/cmd/quote-api/handler"
	"github.com/johanronkko/quote-service/internal/business/data/quote"
	"github.com/johanronkko/quote-service/internal/business/mock"
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
				r := httptest.NewRequest(http.MethodGet, "/api.v1/quotes", nil)
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

	t.Run("unknown error", func(t *testing.T) {
		is := is.New(t)

		// Mock services.
		q := &mock.Quote{}
		q.QueryCall.Returns.Err = fmt.Errorf("some error")

		// Setup handler.
		h := New()
		h.Quote = q

		// Make request.
		r := httptest.NewRequest(http.MethodGet, "/api.v1/quotes", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		// Assert response HTTP headers.
		is.Equal(w.Code, http.StatusInternalServerError)

		// Assert response payload.
		var resp NoDataResponse
		decodePayload(is, w.Body, &resp)
		is.Equal(resp.Code, http.StatusInternalServerError)
		is.True(!resp.Success)
		is.True(resp.Error != nil)
	})

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

func createTestCustomer(name string, ccode string) quote.Customer {
	return quote.Customer{
		Name:        name,
		Email:       "example@test.com",
		Address:     "Vasagatan 5B, Göteborg 41124",
		CountryCode: ccode,
	}
}
