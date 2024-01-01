package responses

import (
	"encoding/json"
	"errors"
	errors2 "github.com/sdq-codes/claimc/gateway/cmd/errors"
	"net/http"
	"strings"
)

type response struct {
	Body       *responseBody
	StatusCode int
}

type responseBody struct {
	Success              bool        `json:"success"`
	Title                string      `json:"title,omitempty"`
	Message              string      `json:"message,omitempty"`
	Data                 interface{} `json:"data,omitempty"`
	InternalResponseCode int         `json:"internal_response_code,omitempty"`
}

// ToJSON writes the response to the given http.ResponseWriter
// with an application/json Content-Type header.
func (r response) ToJSON(w http.ResponseWriter) error {
	w.WriteHeader(r.StatusCode)
	return json.NewEncoder(w).Encode(r.Body)
}

// Created returns response when a models is created
func Created(title string, message string, data interface{}) *response {
	return newResponse(
		true,
		title,
		message,
		data,
		http.StatusCreated,
		CustomResponseCodes["create_successful"],
	)
}

// Updated returns response when a models is updated
func Updated(title string, message string, data interface{}) *response {
	return newResponse(
		true,
		title,
		message,
		data,
		http.StatusOK,
		CustomResponseCodes["update_successful"],
	)
}

// Fetched returns response when models(s) is fetched
func Fetched(title string, message string, data interface{}) *response {
	return newResponse(
		true,
		title,
		message,
		data,
		http.StatusOK,
		CustomResponseCodes["fetch_successful"],
	)
}

// Deleted returns response when models(s) is fetched
func Deleted(title string, message string, data interface{}) *response {
	return newResponse(
		true,
		title,
		message,
		data,
		http.StatusOK,
		CustomResponseCodes["delete"],
	)
}

// Fail returns a new failed response.
func Fail(title string, message string, statusCode int, data interface{}) *response {
	return newResponse(
		false,
		title,
		message,
		data,
		statusCode,
		CustomResponseCodes["custom_error"],
	)
}

// Fail returns a new failed response.
func CustomError(title string, message string, statusCode int, data interface{}) *response {
	return newResponse(
		false,
		title,
		message,
		data,
		statusCode,
		CustomResponseCodes["custom_error"],
	)
}

// CustomFail returns a new failed response.
func CustomFail(err error, statusCode int) *response {
	appError := new(errors2.AppError)
	if errors.As(err, &appError) { // client error
		return newResponse(
			false,
			toSentenceCase(err.Error()),
			string(errTypeStatusCode(appError.Type())),
			err,
			statusCode,
			CustomResponseCodes["fetch_successful"],
		)
	}
	return newResponse(
		false,
		toSentenceCase(err.Error()),
		string(errTypeStatusCode(appError.Type())),
		nil,
		errTypeStatusCode(appError.Type()),
		CustomResponseCodes["custom_error"],
	)
}

// ValidationFailed returns a validation failed expression
func ValidationFailed(message string) *response {
	return newResponse(
		false,
		"Validation failed",
		message,
		nil,
		400,
		CustomResponseCodes["validation_error"],
	)
}

func newResponse(
	success bool,
	title string,
	message string,
	data interface{},
	statusCode int,
	internalResponseCode int,
) *response {
	return &response{
		Body: &responseBody{
			Success:              success,
			Message:              message,
			Data:                 data,
			Title:                title,
			InternalResponseCode: internalResponseCode,
		},
		StatusCode: statusCode,
	}
}

func errTypeStatusCode(errType errors2.Type) int {
	switch errType {
	case errors2.TypeBadRequest:
		return http.StatusBadRequest
	case errors2.TypeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}

func toSentenceCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}
