package persistence

import (
	"context"
	"micro/domain/entity"
	"micro/domain/repository"
	"micro/pkg/parameter"

	"gorm.io/gorm"
)

// DocumentRepo is a struct to store db connection.
type DocumentRepo struct {
	db *gorm.DB
}

// NewDocumentRepository will initialize DocumentRepo repository.
func NewDocumentRepository(db *gorm.DB) *DocumentRepo {
	return &DocumentRepo{db}
}

// DocumentRepo implements the repository.DocumentRepositoryInterface.
var _ repository.DocumentRepositoryInterface = &DocumentRepo{}

// DeleteDocument will delete Document from the database storage.
func (f *DocumentRepo) DeleteDocument(ctx context.Context, r *entity.Document) (*entity.Document, error) {
	var dataEntity entity.Document

	err := f.db.WithContext(ctx).Where("id = ?", r.ID).Take(&dataEntity).Delete(&dataEntity).Error
	if err != nil {
		return nil, err
	}

	return &dataEntity, nil
}

// FindDocument will find Document from the database storage.
func (f *DocumentRepo) FindDocument(ctx context.Context, r *entity.Document) (*entity.Document, error) {
	var dataEntity entity.Document

	err := f.db.WithContext(ctx).Where("id = ?", r.ID).Take(&dataEntity).Error
	if err != nil {
		return nil, err
	}
	return &dataEntity, nil
}

// FindDocumentByIDAndCategoryID will find Document from the database storage.
func (f *DocumentRepo) FindDocumentByIDAndCategoryID(ctx context.Context, r *entity.Document) (*entity.Document, error) {
	var dataEntity entity.Document

	err := f.db.WithContext(ctx).Where("id = ? AND category_id = ?", r.ID, r.CategoryID).Take(&dataEntity).Error
	if err != nil {
		return nil, err
	}
	return &dataEntity, nil
}

// FindDocumentByPath will find Document from the database storage.
func (f *DocumentRepo) FindDocumentByPath(ctx context.Context, r *entity.Document) (*entity.Document, error) {
	var dataEntity entity.Document

	err := f.db.WithContext(ctx).Where("path = ?", r.Path).Take(&dataEntity).Error
	if err != nil {
		return nil, err
	}
	return &dataEntity, nil
}

// FindDocumentByEntity is to find single row of data by entity as condition.
func (f *DocumentRepo) FindDocumentByEntity(ctx context.Context, document *entity.Document) (*entity.Document, error) {
	var dataEntity entity.Document

	err := f.db.WithContext(ctx).Take(&dataEntity, document).Error
	return &dataEntity, err
}

// GetDocuments will get Documents from the database storage.
func (f *DocumentRepo) GetDocuments(ctx context.Context, q *parameter.SQLQueryParameters) (entity.Documents, *parameter.ResponseMetadata, error) {
	var total int64
	var dataEntities entity.Documents

	errTotal := f.db.WithContext(ctx).Model(entity.Document{}).Where(q.QueryKey, q.QueryValue...).Where(q.DateRange).Count(&total).Limit(q.Limit).Offset(q.Offset).Find(&dataEntities).Error
	if errTotal != nil {
		return nil, nil, errTotal
	}

	errList := f.db.WithContext(ctx).Where(q.QueryKey, q.QueryValue...).Where(q.DateRange).Order(q.Order).Limit(q.Limit).Offset(q.Offset).Find(&dataEntities).Error
	if errList != nil {
		return nil, nil, errList
	}
	meta := parameter.NewMeta(q, total)

	return dataEntities, meta, nil
}

// SaveDocument will save Document from the database storage.
func (f *DocumentRepo) SaveDocument(ctx context.Context, r *entity.Document) (*entity.Document, error) {
	var dataEntity entity.Document
	dataEntity.ID = r.ID
	dataEntity.Name = r.Name
	dataEntity.OriginalName = r.OriginalName
	dataEntity.CategoryID = r.CategoryID
	dataEntity.Path = r.Path
	dataEntity.Type = r.Type
	dataEntity.Size = r.Size
	dataEntity.Token = r.Token

	err := f.db.WithContext(ctx).Create(&dataEntity).Error
	if err != nil {
		return nil, err
	}

	return &dataEntity, nil
}

// UpdateDocument is to update a single row of data.
func (f *DocumentRepo) UpdateDocument(ctx context.Context, target *entity.Document, value *entity.Document) error {
	value.ID = ""
	err := f.db.WithContext(ctx).Where(target).Updates(value).Find(&target, target).Error
	if err != nil {
		return err
	}
	return nil
}
