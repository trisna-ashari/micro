package view

import (
	"errors"
	"github.com/gin-gonic/gin"
	"micro/domain/entity"
	"micro/transport/rest/dependency"
	"micro/transport/rest/presenter"
	"net/http"
	"time"
)

// Handler holds the dependency.
type Handler struct {
	Dependency *dependency.Dependency
}

// ViewCategory will handle find category request.
// @Summary Uses to find category request
// @Description Document category.
// @Tags Document Category API
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
// @Router /api/v1/document-categories/:id [get]
func (h *Handler) ViewCategory(c *gin.Context) {
	var payload Request
	err := c.ShouldBindUri(&payload)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("error.common.bad_request"))
		return
	}

	category, err := h.Dependency.DBClient.DocumentCategory.FindDocumentCategory(c.Request.Context(), &entity.DocumentCategory{ID: payload.ID})
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, errors.New("error.common.not_found"))
		return
	}

	response := &Response{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		Size:      category.Size,
		MimeTypes: category.MimeTypes,
		Desc:      category.Description,
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
	}

	c.Status(http.StatusOK)
	presenter.NewSuccessPresenter(c, response.WithoutCreatedAt(), "success.view_category").JSON()
}
