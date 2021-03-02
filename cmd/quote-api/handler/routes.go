package handler

import "net/http"

func (s *Handler) routes() {
	s.router.HandleFunc(http.MethodGet, "/api.v1/healthcheck", s.handleHealthCheck())
	s.router.HandleFunc(http.MethodGet, "/api.v1/quotes/:id", s.handleGetQuote())
	s.router.HandleFunc(http.MethodGet, "/api.v1/quotes", s.handleListQuotes())
	s.router.HandleFunc(http.MethodPost, "/api.v1/quotes", s.handleAddQuote())
}
