package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiOkResponse[T any] struct {
	Status    string    `json:"status"`
	Data      T         `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

type ApiErrorResponse[T any] struct {
	Status       string    `json:"status"`
	ErrorCode    int       `json:"code"`
	ErrorMessage string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
}

func StatusOkResponse[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, ApiOkResponse[T]{
		Status:    "success",
		Data:      data,
		Timestamp: time.Now(),
	})
}

func StatusErrorResponse[T any](c *gin.Context, errorCode int, errorMessage string) {
	c.JSON(http.StatusOK, ApiErrorResponse[T]{
		Status:       "error",
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		Timestamp:    time.Now(),
	})
}
