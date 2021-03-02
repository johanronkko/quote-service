package handler

import "net/http"

func (h *Handler) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, http.StatusOK, nil)
	}
}
