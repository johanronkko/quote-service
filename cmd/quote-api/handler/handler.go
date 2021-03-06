package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/johanronkko/quote-service/internal/business/validate"
	"github.com/matryer/way"
)

// New returns a new Handler with configured routes. Does not initialize
// exported fields.
func New() *Handler {
	s := &Handler{
		router: way.NewRouter(),
	}
	s.routes()
	return s
}

// Handler handles HTTP requests.
type Handler struct {
	router *way.Router
	Quote
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// respond provides a uniform way of responding to HTTP requests.
//
// Currently responds in a JSON format, but can be extended to support different
// formats depending on the Content-Type HTTP header.
func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	dataenvelope := map[string]interface{}{"code": status}
	if err, ok := data.(error); ok {
		var ferrors validate.FieldErrors
		if errors.As(err, &ferrors) {
			dataenvelope["error"] = ferrors
		} else {
			dataenvelope["error"] = err.Error()
		}
		dataenvelope["success"] = false
	} else {
		if data != nil {
			dataenvelope["data"] = data
		}
		dataenvelope["success"] = true
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(dataenvelope); err != nil {
		log.Printf("respond: %s", err.Error())
	}
}

// decode decodes HTTP request payloads.

// Currently only support JSON unmarshalling, but can be extended to support other
// serialization formats (e.g. XML) depending on the Content-Type in the request
// header. We future proof the function by also taking the ResponseWriter, even though
// it is not used at the moment.
func decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}
