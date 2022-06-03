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
	Response *http.Response
	Message  string              `json:"message"`
	Errors   map[string][]string `json:"errors"`
}

func (e GenericError) Error() string {
	return fmt.Sprintf("%v %v: %d, %v %+v", e.Response.Request.Method, e.Response.Request.URL, e.Response.StatusCode, e.Message, e.Errors)
}

type UnauthorizedError GenericError

func (e *UnauthorizedError) Error() string {
	return (*GenericError)(e).Error()
}

type BadRequestError GenericError

func (e *BadRequestError) Error() string {
	return (*GenericError)(e).Error()
}

type ConflictError GenericError

func (e *ConflictError) Error() string {
	return (*GenericError)(e).Error()
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
	builder.WriteString(e.Message)
	for _, err := range e.GetValidationErrors() {
		builder.WriteString(fmt.Sprintf("\n\t- %s", err))
	}

	return builder.String()
}
