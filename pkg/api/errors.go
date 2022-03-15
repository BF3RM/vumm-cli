package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrModVersionNotFound = errors.New("mod version was not found")
)

type GenericError struct {
	statusCode int
	message    string
}

func (e GenericError) StatusCode() int {
	return e.statusCode
}

func (e GenericError) Status() string {
	return http.StatusText(e.statusCode)
}

func (e GenericError) Error() string {
	return fmt.Sprintf("%s: %s", e.message, e.Status())
}

type ValidationError struct {
	GenericError
	errors map[string][]string
}

func (e ValidationError) GetKeyValidationErrors(key string) ([]string, bool) {
	val, ok := e.errors[key]
	return val, ok
}

func (e ValidationError) GetValidationErrors() []string {
	var errs []string
	for _, valErrs := range e.errors {
		errs = append(errs, valErrs...)
	}

	return errs
}

func (e *ValidationError) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	e.errors = map[string][]string{}

	if errors, ok := raw["errors"]; ok {
		for key, val := range errors.(map[string]interface{}) {
			var errs []string
			for _, valErr := range val.([]interface{}) {
				errs = append(errs, valErr.(string))
			}
			e.errors[key] = errs
		}
	}

	return nil
}

func (e ValidationError) Error() string {
	if len(e.errors) == 0 {
		return e.GenericError.Error()
	}

	builder := strings.Builder{}
	builder.WriteString(e.message)
	for _, err := range e.GetValidationErrors() {
		builder.WriteString(fmt.Sprintf("\n\t- %s", err))
	}

	return builder.String()
}
