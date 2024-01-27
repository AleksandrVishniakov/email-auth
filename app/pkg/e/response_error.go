package e

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type ResponseError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func NewResponseError(code int, message string) *ResponseError {
	return &ResponseError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

func (e *ResponseError) Abort(c *gin.Context) {
	log.Printf("%d: %s\n", e.Code, e.Message)
	c.AbortWithStatusJSON(e.Code, e)
}
