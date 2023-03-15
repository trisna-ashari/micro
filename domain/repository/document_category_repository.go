package repository

import (
	"context"
	"micro/domain/entity"
	"micro/pkg/parameter"
)

// DocumentCategoryRepositoryInterface need to be implements in persistence repository.
type DocumentCategoryRepositoryInterface interface {
	DeleteDocumentCategory(context.Context, *entity.DocumentCategory) (*entity.DocumentCategory, error)
	FindDocumentCategory(context.Context, *entity.DocumentCategory) (*entity.DocumentCategory, error)
	FindDocumentCategoryBySlug(context.Context, *entity.DocumentCategory) (*entity.DocumentCategory, error)
	GetDocumentCategories(context.Context, *parameter.SQLQueryParameters) (entity.DocumentCategories, *parameter.ResponseMetadata, error)
	SaveDocumentCategory(context.Context, *entity.DocumentCategory) (*entity.DocumentCategory, error)
	UpdateDocumentCategory(context.Context, *entity.DocumentCategory) (*entity.DocumentCategory, error)
}
