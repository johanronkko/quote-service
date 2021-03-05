package validate

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator *ut.UniversalTranslator

// TODO: is first letter uppercase assumption too strong? Should it just be
// [[:alpha:]]?
var personNameRegex = regexp.MustCompile("^[[:upper:]]([[:alpha:]]|[[:punct:]]|[[:space:]]){0,29}$")

func init() {

	// Instantiate the validator for use.
	validate = validator.New()

	_ = validate.RegisterValidation("personname", func(fl validator.FieldLevel) bool {
		return personNameRegex.MatchString(fl.Field().String())
	})

	// Instantiate the english locale for the validator library.
	enLocale := en.New()

	// Create a value using English as the fallback locale (first argument).
	// Provide one or more arguments for additional supported locales.
	translator = ut.New(enLocale, enLocale)

	// Register the english error messages for validation errors.
	lang, _ := translator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, lang)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

}

// ErrInvalidID occurs when an ID is not in a valid form.
var ErrInvalidID = errors.New("ID is not in its proper form")

// GenerateID generate a unique id for entities.
func GenerateID() string {
	return uuid.New().String()
}

// CheckID validates that the format of an id is valid.
func CheckID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}
	return nil
}

// Check validates the provided model against it's declared tags.
func Check(val interface{}) error {

	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// lang controls the language of the error messages. You could look at the
		// Accept-Language header if you intend to support multiple languages.
		lang, _ := translator.GetTranslator("en")

		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(lang),
			}

			fields = append(fields, field)
		}

		return fields
	}

	return nil
}

// FieldError is used to indicate an error with a specific field on a struct.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// FieldErrors represents a collection of field errors.
type FieldErrors []FieldError

// Error implments the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}
