// Copyright (c) 2022 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package apierror

import (
	"encoding/json"
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

type APIError struct {
	Err    Err    `json:"error"`
	Status int    `json:"status"`
	Path   string `json:"path"`
}

type Err struct {
	Code    string                `json:"code"`
	Message string                `json:"message"`
	Fields  map[string]FieldError `json:"fields,omitempty"`
}

type FieldError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Message returns a message for humans.
func (e *APIError) Message() string {
	return e.Err.Message
}

func (e *APIError) Error() string {
	return fmt.Sprintf("server responded with an error: %v (%s) [%s] %s", e.Status, e.Path, e.Err.Code, e.Err.Message)
}

func Parse(response *http.Response) *APIError {
	var apiErr APIError
	blob, _ := ioutil.ReadAll(response.Body)

	switch response.Header.Get("Content-Type") {
	case ContentTypeV1:
		if err := json.Unmarshal(blob, &apiErr); err != nil {
			apiErr = makeDefault(response.StatusCode, blob)
		}
	default:
		apiErr = makeDefault(response.StatusCode, blob)
	}
	apiErr.Path = response.Request.URL.String()
	return &apiErr
}

func makeDefault(statusCode int, payload []byte) APIError {
	return APIError{
		Err:    Err{Code: ErrCodeInternal, Message: string(payload)},
		Status: statusCode,
	}
}
