// Copyright (c) 2022 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package apierror

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
)

func defaultErr() *APIError {
	return &APIError{
		Err: Err{
			Code:    ErrCodeInternal,
			Message: "Internal server error",
		},
		Status: 500,
	}
}

var DefaultErrStatic, _ = json.Marshal(defaultErr())

// ErrorHandler gets the first error from request context and formats it to an error response.
func ErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) < 1 {
		return
	}
	err := c.Errors[0].Err

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		apiErr = defaultErr()
		apiErr.Err.Message = err.Error()
	}
	apiErr.Path = c.Request.URL.String()

	// Gin only uses the first value from Accept header by default. Help him.
	//c.SetAccepted(c.Request.Header.Values("Accept")...)
	//
	//switch c.NegotiateFormat(ContentTypeV1) {
	//case ContentTypeV1:
	blob, err := json.Marshal(apiErr)
	if err != nil {
		c.Data(500, ContentTypeV1, DefaultErrStatic)
		return
	}
	c.Data(apiErr.Status, ContentTypeV1, blob)
	//default:
	// Fallback for older clients
	//	c.JSON(500, gin.H{"message": apiErr.Error()})
	//}
}
