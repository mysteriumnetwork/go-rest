// Copyright (c) 2022 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package apierror

import "net/http"

func Error(status int, message string, code string) *APIError {
	return &APIError{
		Err:    Err{Code: code, Message: message},
		Status: status,
	}
}

func NotFound(message string) *APIError {
	return &APIError{
		Err:    Err{Code: ErrCodeNotFound, Message: message},
		Status: http.StatusNotFound,
	}
}

func ParseFailed() *APIError {
	return &APIError{
		Err:    Err{Code: ErrCodeParseFailed, Message: "Could not parse request"},
		Status: http.StatusBadRequest,
	}
}

func BadRequestField(message string, code string, field string) *APIError {
	return &APIError{
		Err: Err{Code: ErrCodeValidationFailed, Message: "Request validation failed", Fields: map[string]FieldError{
			field: {
				Error:   code,
				Message: message,
			},
		}},
		Status: http.StatusBadRequest,
	}
}

func BadRequestFields(fields map[string]FieldError) *APIError {
	return &APIError{
		Err:    Err{Code: ErrCodeValidationFailed, Message: "Request validation failed", Fields: fields},
		Status: http.StatusBadRequest,
	}
}

func BadRequest(message string, code string) *APIError {
	return &APIError{
		Err:    Err{Code: code, Message: message},
		Status: http.StatusBadRequest,
	}
}

func Conflict(message string, code string, field string) *APIError {
	return &APIError{
		Err: Err{Code: code, Message: message, Fields: map[string]FieldError{
			field: {
				Error:   code,
				Message: message,
			},
		}},
		Status: http.StatusConflict,
	}
}

func Forbidden(message string, code string) *APIError {
	return &APIError{
		Err:    Err{Code: code, Message: message},
		Status: http.StatusForbidden,
	}
}

func Internal(message string, code string) *APIError {
	return &APIError{
		Err:    Err{Code: code, Message: message},
		Status: http.StatusInternalServerError,
	}
}

func InternalDefault() *APIError {
	return Internal(http.StatusText(http.StatusInternalServerError), ErrCodeInternal)
}

func Unauthorized() *APIError {
	return &APIError{
		Err:    Err{Code: ErrCodeUnauthorized, Message: http.StatusText(http.StatusUnauthorized)},
		Status: http.StatusUnauthorized,
	}
}

func Unprocessable(message string, code string) *APIError {
	return &APIError{
		Err:    Err{Code: code, Message: message},
		Status: http.StatusUnprocessableEntity,
	}
}

func ServiceUnavailable() *APIError {
	return &APIError{
		Err:    Err{Code: ErrCodeUnavailable, Message: http.StatusText(http.StatusServiceUnavailable)},
		Status: http.StatusServiceUnavailable,
	}
}
