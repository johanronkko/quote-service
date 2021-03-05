package validate_test

import (
	"strings"
	"testing"

	"github.com/johanronkko/quote-service/internal/business/tests"
	"github.com/johanronkko/quote-service/internal/business/validate"
	"github.com/matryer/is"
)

type testStruct struct {
	Name string `json:"name" validate:"personname"`
}

func TestPersonNameTag(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cases := []struct {
			Name string

			PersonName string
		}{
			{"only letters", "John"},
			{"space", "John doe"},
			{"punct", "John-doe"},
			{"min length", "A"},
			{"max length", strings.Title(tests.GenRandomAlpha(30))},
		}
		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				is := is.New(t)

				err := validate.Check(testStruct{tc.PersonName})
				is.NoErr(err)
			})
		}
	})
	t.Run("invalid", func(t *testing.T) {
		cases := []struct {
			Name string

			PersonName string
		}{
			{"first lowercase", "john"},
			{"number", "John1"},
			{"empty", ""},
			{"too long", strings.Title(tests.GenRandomAlpha(31))},
		}
		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				is := is.New(t)

				err := validate.Check(testStruct{tc.PersonName})
				is.True(err != nil)
			})
		}
	})
}
