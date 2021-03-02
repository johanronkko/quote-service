package handler

import "net/http"

func (s *Handler) routes() {
	s.router.HandleFunc(http.MethodGet, "/api.v1/healthcheck", s.handleHealthCheck())
	s.router.HandleFunc(http.MethodGet, "/api.v1/quotes", s.handleListQuotes())
}
