package quote

import (
	"errors"
)

var (
	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")
	// ErrNotFound is used when a specific Quote is requested but does not exist.
	ErrNotFound = errors.New("not found")
)
