// Copyright (c) 2022 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package apierror

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ContentTypeV1           = "application/vnd.mysterium.error+json"
	ErrCodeNotFound         = "not_found"
	ErrCodeInternal         = "internal"
	ErrCodeParseFailed      = "parse_failed"
	ErrCodeValidationFailed = "validation_failed"
	ErrCodeUnavailable      = "unavailable"
	ErrCodeUnauthorized     = "shall_not_pass"
)

// APIError represents an error response from REST API service.
type APIError struct {
	Err    Err    `json:"error"`
	Status int    `json:"status"`
	Path   string `json:"path"`
}

// Err contains the error details.
type Err struct {
	Code    string                `json:"code"`
	Message string                `json:"message"`
	Detail  string                `json:"detail,omitempty"`
	Fields  map[string]FieldError `json:"fields,omitempty"`
}

// FieldError contains the reason why a field failed validation.
type FieldError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Message returns a message for humans.
func (e *APIError) Message() string {
	return e.Err.Message
}

// Detail returns a detailed message for humans.
func (e *APIError) Detail() string {
	if e.Err.Detail == "" {
		return e.Err.Message
	}
	return e.Err.Detail
}

func (e *APIError) Error() string {
	return fmt.Sprintf("server responded with an error: %v (%s) [%s] %s", e.Status, e.Path, e.Err.Code, e.Err.Message)
}

// Parse parses http.Response into an APIError closing (consuming) http.Response.Body.
func Parse(response *http.Response) *APIError {
	var apiErr APIError
	blob, _ := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()

	switch response.Header.Get("Content-Type") {
	case ContentTypeV1:
		if err := json.Unmarshal(blob, &apiErr); err != nil {
			apiErr = makeDefault(response.StatusCode, blob)
		}
	default:
		apiErr = makeDefault(response.StatusCode, blob)
	}
	if response.Request != nil && response.Request.URL != nil {
		apiErr.Path = response.Request.URL.String()
	}
	return &apiErr
}

// TryParse parses http.Response into an APIError closing (consuming) http.Response.Body.
// If header matches and error can be parsed in to APIError it will return that.
// If No such header exists a regular error is returned with it's contents being a response body.
func TryParse(response *http.Response) error {
	blob, _ := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()

	switch response.Header.Get("Content-Type") {
	case ContentTypeV1:
		var apiErr APIError
		if err := json.Unmarshal(blob, &apiErr); err != nil {
			apiErr = makeDefault(response.StatusCode, blob)
		}
		if response.Request != nil && response.Request.URL != nil {
			apiErr.Path = response.Request.URL.String()
		}

		return &apiErr
	default:
		return errors.New(string(blob))
	}
}

func makeDefault(statusCode int, payload []byte) APIError {
	return APIError{
		Err:    Err{Code: ErrCodeInternal, Message: string(payload)},
		Status: statusCode,
	}
}
