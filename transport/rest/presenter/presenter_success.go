package presenter

import "github.com/gin-gonic/gin"

// Success is success output presenter.
type Success struct {
	c *gin.Context

	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Meta    interface{} `json:"meta,omitempty"`
}

// NewSuccessPresenter will initialize a new Success.
func NewSuccessPresenter(c *gin.Context, data interface{}, message string) *Success {
	return &Success{
		c:       c,
		Code:    c.Writer.Status(),
		Data:    data,
		Message: message,
	}
}

// WithMeta will append meta to Success.
func (s *Success) WithMeta(meta interface{}) *Success {
	s.Meta = meta

	return s
}

// JSON return as json.
func (s *Success) JSON() {
	s.c.JSON(s.Code, s)
}
