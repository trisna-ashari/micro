package ping

import (
	"github.com/gin-gonic/gin"
	"micro/transport/rest/dependency"
	"micro/transport/rest/presenter"
	"net/http"
)

// Handler holds the dependency.
type Handler struct {
	Dependency *dependency.Dependency
}

// Ping will handle ping request.
// @Summary Uses to ping
// @Description ping.
// @Tags Ping API
// @Accept  json
// @Produce application/json
// @Param Accept-Language header string false "Fill with language code" Enums(en, id) default(id)
// @Param Set-Request-Id header string false "Fill with request id"
// @Success 200 {object} presenter.Success{data=ping.Response}
// @Failure 400 {object} presenter.Error
// @Failure 401 {object} presenter.Error
// @Failure 403 {object} presenter.Error
// @Failure 404 {object} presenter.Error
// @Failure 422 {object} presenter.Error
// @Failure 500 {object} presenter.Error
// @Router /ping [get]
func (h *Handler) Ping(c *gin.Context) {
	c.Status(http.StatusOK)
	presenter.NewSuccessPresenter(c, &Response{Status: "OK"}, "pong").JSON()
}
