package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/johanronkko/quote-service/cmd/quote-api/handler"
	"github.com/johanronkko/quote-service/internal/business/data/quote"
	"github.com/johanronkko/quote-service/internal/business/tests"
	"github.com/johanronkko/quote-service/internal/business/validate"
	"github.com/matryer/is"
)

type NoDataResponse struct {
	Code    int     `json:"code"`
	Error   *string `json:"error"`
	Success bool    `json:"success"`
}

type FieldErrorResponse struct {
	Code        int                  `json:"code"`
	FieldErrors validate.FieldErrors `json:"error"`
	Success     bool                 `json:"success"`
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

func decodePayload(is *is.I, r io.Reader, v interface{}) {
	data, err := ioutil.ReadAll(r)
	is.NoErr(err)
	err = json.Unmarshal(data, v)
	is.NoErr(err)
}

// TestQuotes tests the quote API happy path.
func TestQuotes(t *testing.T) {
	is := is.New(t)

	db := tests.NewIntegration(t)
	numSeededQuotes := 3

	handler := handler.New()
	handler.Quote = quote.New(db)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Is able to retrieve a list of quotes.
	resp, err := http.Get(ts.URL + "/api.v1/quotes/")
	is.NoErr(err)
	var quotesResponse QuotesResponse
	json.NewDecoder(resp.Body).Decode(&quotesResponse)
	is.Equal(quotesResponse.Code, http.StatusOK)
	is.True(quotesResponse.Success)
	is.Equal(len(quotesResponse.Data.Quotes), numSeededQuotes)

	// Is able to add a new quote.
	nq := quote.NewQuote{
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
	nqReqBody, err := json.Marshal(&nq)
	is.NoErr(err)
	resp, err = http.Post(ts.URL+"/api.v1/quotes/", "application/json", bytes.NewBuffer(nqReqBody))
	is.NoErr(err)
	var newQuoteResponse QuoteResponse
	err = json.NewDecoder(resp.Body).Decode(&newQuoteResponse)
	is.NoErr(err)
	is.Equal(newQuoteResponse.Code, http.StatusCreated)
	is.True(newQuoteResponse.Success)
	is.Equal(newQuoteResponse.Data.Quote.From, nq.From)
	is.Equal(newQuoteResponse.Data.Quote.To, nq.To)
	is.Equal(newQuoteResponse.Data.Quote.Weight, nq.Weight)
	is.Equal(newQuoteResponse.Data.Quote.ShipmentCost, 2.5*2000) // From outside EU * huge package.

	// Is able to retrieve newly added quote.
	resp, err = http.Get(ts.URL + "/api.v1/quotes/" + newQuoteResponse.Data.Quote.ID)
	is.NoErr(err)
	var quoteByIDResponse QuoteResponse
	err = json.NewDecoder(resp.Body).Decode(&quoteByIDResponse)
	is.NoErr(err)
	is.Equal(quoteByIDResponse.Data.Quote, newQuoteResponse.Data.Quote)
}
