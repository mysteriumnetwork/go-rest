// Copyright (c) 2022 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package apierror

import (
	"fmt"
)

const (
	ValidateErrRequired   = "required"
	ValidateErrInvalidVal = "invalid_value"
)

type Validator struct {
	fields map[string]FieldError
}

func NewValidator() *Validator {
	return &Validator{
		fields: map[string]FieldError{},
	}
}

func (v *Validator) Fail(field string, errorCode string, message string) {
	v.fields[field] = FieldError{
		Error:   errorCode,
		Message: message,
	}
}

func (v *Validator) Required(field string) {
	v.fields[field] = FieldError{
		Error:   ValidateErrRequired,
		Message: fmt.Sprintf("'%s' is required", field),
	}
}

func (v *Validator) Invalid(field string, message string) {
	v.fields[field] = FieldError{
		Error:   ValidateErrInvalidVal,
		Message: message,
	}
}

func (v *Validator) Err() *APIError {
	if len(v.fields) == 0 {
		return nil
	}
	return BadRequestFields(v.fields)
}
