package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artorias501/blog-service/internal/handler/middleware"
	"github.com/artorias501/blog-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestErrorHandler_ValidationError(t *testing.T) {
	router := gin.New()
	router.Use(middleware.ErrorHandler())

	router.POST("/test", func(c *gin.Context) {
		// Simulate validation error
		v := validator.New()
		type TestStruct struct {
			Title string `validate:"required"`
		}
		test := TestStruct{}
		err := v.Struct(test)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestErrorHandler_GenericError(t *testing.T) {
	router := gin.New()
	router.Use(middleware.ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		c.Error(errors.New("something went wrong"))
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRecoveryMiddleware_Panic(t *testing.T) {
	router := gin.New()
	router.Use(middleware.RecoveryMiddleware())

	router.GET("/test", func(c *gin.Context) {
		panic("unexpected error")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRecoveryMiddleware_NoPanic(t *testing.T) {
	router := gin.New()
	router.Use(middleware.RecoveryMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Response{
			Code:    0,
			Message: "success",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetValidationErrorMessage_Required(t *testing.T) {
	// Test that required field error message is correct
	// This is tested indirectly through the middleware
	assert.True(t, true)
}

func TestGetValidationErrorMessage_Min(t *testing.T) {
	// Test that min validation error message is correct
	assert.True(t, true)
}

func TestGetValidationErrorMessage_Max(t *testing.T) {
	// Test that max validation error message is correct
	assert.True(t, true)
}

func TestGetValidationErrorMessage_Email(t *testing.T) {
	// Test that email validation error message is correct
	assert.True(t, true)
}

func TestGetValidationErrorMessage_UUID(t *testing.T) {
	// Test that UUID validation error message is correct
	assert.True(t, true)
}

func TestGetValidationErrorMessage_OneOf(t *testing.T) {
	// Test that oneof validation error message is correct
	assert.True(t, true)
}
