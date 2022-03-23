// Copyright (c) 2022 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package apierror

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func GetAPIError(c *gin.Context) {
	v := NewValidator()
	v.Required("amount")
	v.Invalid("identity", "Not a valid identity")
	response := v.Err()
	c.AbortWithStatusJSON(400, response)
}

func TestAPIError_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/test-api-error", GetAPIError)

	req, _ := http.NewRequest(http.MethodGet, "/test-api-error", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resErr APIError
	err := json.Unmarshal(w.Body.Bytes(), &resErr)
	assert.NoError(t, err)
	assert.Equal(t, 400, resErr.Status)
	assert.Equal(t, "Request validation failed", resErr.Message())
	assert.Equal(t, "Request validation failed: amount: 'amount' is required [required]; identity: Not a valid identity [invalid_value]", resErr.Detail())
	assert.Equal(t, FieldError{
		Code:    ValidateErrRequired,
		Message: "'amount' is required",
	}, resErr.Err.Fields["amount"])
	assert.Equal(t, FieldError{
		Code:    ValidateErrInvalidVal,
		Message: "Not a valid identity",
	}, resErr.Err.Fields["identity"])
}
