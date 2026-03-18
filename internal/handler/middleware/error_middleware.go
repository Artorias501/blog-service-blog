// Package middleware contains HTTP middleware for the blog service.
package middleware

import (
	"errors"
	"strings"

	"github.com/artorias501/blog-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorHandler is a middleware that handles errors and validation errors.
// It intercepts errors set by previous handlers and returns appropriate responses.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Handle validation errors
			var validationErrs validator.ValidationErrors
			if errors.As(err, &validationErrs) {
				handleValidationError(c, validationErrs)
				return
			}

			// Handle other errors
			response.InternalError(c, err.Error())
			return
		}
	}
}

// handleValidationError processes validation errors and returns detailed field errors.
func handleValidationError(c *gin.Context, errs validator.ValidationErrors) {
	fieldErrors := make([]response.ErrorDetail, 0, len(errs))

	for _, err := range errs {
		field := strings.ToLower(err.Field())
		message := getValidationErrorMessage(err)

		fieldErrors = append(fieldErrors, response.ErrorDetail{
			Field:   field,
			Message: message,
		})
	}

	response.ValidationError(c, fieldErrors)
}

// getValidationErrorMessage returns a human-readable error message for validation errors.
func getValidationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "this field is required"
	case "min":
		return "value must be at least " + err.Param() + " characters"
	case "max":
		return "value must be at most " + err.Param() + " characters"
	case "email":
		return "invalid email format"
	case "uuid":
		return "invalid UUID format"
	case "oneof":
		return "value must be one of: " + err.Param()
	default:
		return "invalid value"
	}
}
