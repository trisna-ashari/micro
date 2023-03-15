package middleware

import (
	"github.com/gin-gonic/gin"
	"micro/pkg/configurator"
	"micro/pkg/logger"
	"micro/transport/rest/presenter"
)

// Error handle error.
type Error struct {
	Config *configurator.Config
	Logger *logger.Logger
}

// NewError will initialize a new Error middleware.
func NewError(config *configurator.Config, logger *logger.Logger) *Error {
	return &Error{
		Config: config,
		Logger: logger,
	}
}

// ErrorHandler will handle any error response.
func (e *Error) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset: utf-8")
		c.Next()

		if c.Errors.Last() == nil {
			return
		}

		err := c.Errors.Last().Err
		c.Errors = c.Errors[:0]

		if e.Config.AppEnvironment == "production" && c.Writer.Status() == 500 {
			c.JSON(c.Writer.Status(), &presenter.Error{
				Code:    c.Writer.Status(),
				Data:    nil,
				Message: err.Error(),
			})

			return
		}

		c.JSON(c.Writer.Status(), &presenter.Error{
			Code:    c.Writer.Status(),
			Data:    nil,
			Message: err.Error(),
		})
	}
}
