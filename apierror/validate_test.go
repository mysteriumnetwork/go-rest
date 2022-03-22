// Copyright (c) 2022 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package apierror

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	v := NewValidator()
	v.Required("amount")
	v.Invalid("id", "ID invalid")
	v.Fail("currency", "not_supported", "Unsupported currency")
	err := v.Err()

	assert.Equal(t, "Request validation failed", err.Message())
	assert.Equal(t, err.Err.Fields["amount"], FieldError{
		Error:   ValidateErrRequired,
		Message: "'amount' is required",
	})
	assert.Equal(t, err.Err.Fields["id"], FieldError{
		Error:   ValidateErrInvalidVal,
		Message: "ID invalid",
	})
	assert.Equal(t, err.Err.Fields["currency"], FieldError{
		Error:   "not_supported",
		Message: "Unsupported currency",
	})
}
