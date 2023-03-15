package parameter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParameterNewHTTPParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	var err error

	c, r := gin.CreateTestContext(w)
	v1 := r.Group("/api/v1/external")
	v1.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
		parameters := NewHTTPParameters(c)
		fmt.Println(parameters.QueryParameters.EqualsQueryString)
		fmt.Println(parameters.QueryParameters.LikesQueryString)
		fmt.Println(parameters.QueryParameters.NotEqualsQueryString)
	})

	c.Request, err = http.NewRequest(http.MethodGet, "/api/v1/external/test", nil)
	q := c.Request.URL.Query()
	q.Add("equal[name]", "foo")
	q.Add("equal[name]", "bar")
	q.Add("equal[name]", "foo %26 bar")
	q.Add("not[name]", "zoo")
	q.Add("like[name]", "mia")
	c.Request.URL.RawQuery = q.Encode()

	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r.ServeHTTP(w, c.Request)
}
