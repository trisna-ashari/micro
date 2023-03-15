package repository

import (
	"context"
	"micro/domain/entity"
	"micro/pkg/parameter"
)

// DocumentRepositoryInterface need to be implements in persistence repository.
type DocumentRepositoryInterface interface {
	DeleteDocument(context.Context, *entity.Document) (*entity.Document, error)
	FindDocument(context.Context, *entity.Document) (*entity.Document, error)
	FindDocumentByIDAndCategoryID(context.Context, *entity.Document) (*entity.Document, error)
	FindDocumentByPath(context.Context, *entity.Document) (*entity.Document, error)
	FindDocumentByEntity(context.Context, *entity.Document) (*entity.Document, error)
	GetDocuments(context.Context, *parameter.SQLQueryParameters) (entity.Documents, *parameter.ResponseMetadata, error)
	SaveDocument(context.Context, *entity.Document) (*entity.Document, error)
	UpdateDocument(ctx context.Context, target *entity.Document, value *entity.Document) error
}
