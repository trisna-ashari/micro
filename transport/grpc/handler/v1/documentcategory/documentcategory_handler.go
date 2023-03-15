package documentcategory

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
	"micro/domain/entity"
	"micro/pkg/parameter"
	"micro/transport/grpc/dependency"
	"micro/transport/grpc/presenter"
	"time"
)

// Handler is a struct represent itself.
type Handler struct {
	Dependency *dependency.Dependency

	// It is for forward-compatibility, that if you changed your service files and added some new methods,
	// your binary doesn't fail if you don't implement the new methods in your server.
	// https://github.com/grpc/grpc-go/issues/3669
	UnimplementedDocumentCategoryServiceServer
}

func (h *Handler) DeleteDocumentCategory(ctx context.Context, request *DeleteDocumentCategoryRequest) (*DocumentCategoryDeleted, error) {
	category, err := h.Dependency.DBClient.DocumentCategory.FindDocumentCategory(ctx, &entity.DocumentCategory{
		ID: request.Id,
	})
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.NotFound, "error.document_category.not_found", nil).
			Error()
	}
	if err != nil {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.Internal, "error.common.internal_server_error", nil).
			Error()
	}

	category, err = h.Dependency.DBClient.DocumentCategory.DeleteDocumentCategory(ctx, &entity.DocumentCategory{
		ID: request.Id,
	})
	if err != nil {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.Internal, "error.common.internal_server_error", nil).
			Error()
	}

	return &DocumentCategoryDeleted{DeletedAt: category.DeletedAt.Time.Format(time.RFC3339)}, nil
}

func (h *Handler) FindDocumentCategory(ctx context.Context, request *FindDocumentCategoryRequest) (*DocumentCategory, error) {
	category, err := h.Dependency.DBClient.DocumentCategory.FindDocumentCategory(ctx, &entity.DocumentCategory{
		ID: request.Id,
	})
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.NotFound, "error.document_category.not_found", nil).
			Error()
	}
	if err != nil {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.Internal, "error.common.internal_server_error", nil).
			Error()
	}

	return &DocumentCategory{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		Size:      category.Size,
		MimeTypes: category.MimeTypes,
		Desc:      category.Description,
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (h *Handler) FindDocumentCategoryBySlug(ctx context.Context, request *FindDocumentCategoryBySlugRequest) (*DocumentCategory, error) {
	category, err := h.Dependency.DBClient.DocumentCategory.FindDocumentCategoryBySlug(ctx, &entity.DocumentCategory{
		Slug: request.Slug,
	})
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.NotFound, "error.document_category.not_found", nil).
			Error()
	}
	if err != nil {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.Internal, "error.common.internal_server_error", nil).
			Error()
	}

	return &DocumentCategory{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		Size:      category.Size,
		MimeTypes: category.MimeTypes,
		Desc:      category.Description,
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (h *Handler) GetDocumentCategories(ctx context.Context, request *GetDocumentCategoriesRequest) (*DocumentCategories, error) {
	var dataEntity entity.DocumentCategory

	reqParameters := request.Parameters
	rpcParameters := parameter.RPCParameters{
		SearchCondition: reqParameters.SearchCondition,
		Page:            int(reqParameters.Page),
		PerPage:         int(reqParameters.PerPage),
		OrderBy:         reqParameters.OrderBy,
		OrderMethod:     reqParameters.OrderMethod,
		Equal:           reqParameters.Equal,
		Not:             reqParameters.Not,
		Like:            reqParameters.Like,
		DateRangeBy:     reqParameters.DateRangeBy,
		DateStart:       reqParameters.DateStart,
		DateEnd:         reqParameters.DateEnd,
	}
	sqlParameters := rpcParameters.ToSQLQueryParameters()

	validationResult := sqlParameters.ValidateParameter(dataEntity.FilterableFields(), dataEntity.TimeFields())
	if len(validationResult) > 0 {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.InvalidArgument, "error.common.unprocessable_entity", nil).
			Error()
	}

	categories, meta, err := h.Dependency.DBClient.DocumentCategory.GetDocumentCategories(ctx, sqlParameters)
	if err != nil {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.Internal, "error.common.internal_server_error", nil).
			Error()
	}

	return &DocumentCategories{
		Data: func() []*DocumentCategory {
			var documentCategories []*DocumentCategory
			for _, category := range categories {
				documentCategories = append(documentCategories, &DocumentCategory{
					Id:        category.ID,
					Name:      category.Name,
					Slug:      category.Slug,
					Size:      category.Size,
					MimeTypes: category.MimeTypes,
					Desc:      category.Description,
					CreatedAt: category.CreatedAt.Format(time.RFC3339),
				})
			}

			return documentCategories
		}(),
		Meta: &DocumentCategoryMeta{
			Page:    int32(meta.Page),
			PerPage: int32(meta.PerPage),
			Total:   int32(meta.Total),
		},
	}, nil
}

func (h *Handler) SaveDocumentCategory(ctx context.Context, request *SaveDocumentCategoryRequest) (*DocumentCategory, error) {
	category, err := h.Dependency.DBClient.DocumentCategory.SaveDocumentCategory(ctx, &entity.DocumentCategory{
		ID:          uuid.New().String(),
		Slug:        request.Slug,
		Name:        request.Name,
		Description: request.Description,
		MimeTypes:   request.MimeTypes,
		Size:        request.Size,
	})
	if err != nil {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.Internal, "error.common.internal_server_error", nil).
			Error()
	}

	return &DocumentCategory{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		Size:      category.Size,
		MimeTypes: category.MimeTypes,
		Desc:      category.Description,
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (h *Handler) UpdateDocumentCategory(ctx context.Context, request *UpdateDocumentCategoryRequest) (*DocumentCategory, error) {
	category, err := h.Dependency.DBClient.DocumentCategory.FindDocumentCategory(ctx, &entity.DocumentCategory{
		ID: request.Id,
	})
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.NotFound, "error.document_category.not_found", nil).
			Error()
	}
	if err != nil {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.Internal, "error.common.internal_server_error", nil).
			Error()
	}

	category, err = h.Dependency.DBClient.DocumentCategory.UpdateDocumentCategory(ctx, &entity.DocumentCategory{
		ID:          category.ID,
		Slug:        request.Slug,
		Name:        request.Name,
		Description: request.Description,
		MimeTypes:   request.MimeTypes,
		Size:        request.Size,
	})
	if err != nil {
		return nil, presenter.
			NewErrorPresenter(ctx, codes.Internal, "error.common.internal_server_error", nil).
			Error()
	}

	return &DocumentCategory{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		Size:      category.Size,
		MimeTypes: category.MimeTypes,
		Desc:      category.Description,
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Type assertion ensure that Handler implements DocumentCategoryServiceServer.
var _ DocumentCategoryServiceServer = &Handler{}
