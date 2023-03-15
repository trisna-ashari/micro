package persistence

import (
	"context"
	"micro/domain/entity"
	"micro/domain/repository"
	"micro/pkg/parameter"

	"gorm.io/gorm"
)

// DocumentCategoryRepo is a struct to store db connection.
type DocumentCategoryRepo struct {
	db *gorm.DB
}

// NewDocumentCategoryRepository will initialize DocumentCategoryRepo repository.
func NewDocumentCategoryRepository(db *gorm.DB) *DocumentCategoryRepo {
	return &DocumentCategoryRepo{db}
}

// DocumentCategoryRepo implements the repository.DocumentCategoryRepositoryInterface.
var _ repository.DocumentCategoryRepositoryInterface = &DocumentCategoryRepo{}

// DeleteDocumentCategory will delete Document category from the database storage.
func (f *DocumentCategoryRepo) DeleteDocumentCategory(ctx context.Context, r *entity.DocumentCategory) (*entity.DocumentCategory, error) {
	var dataEntity entity.DocumentCategory

	err := f.db.WithContext(ctx).Where("id = ?", r.ID).Take(&dataEntity).Delete(&dataEntity).Error
	if err != nil {
		return nil, err
	}

	return &dataEntity, nil
}

// FindDocumentCategory will find Document category from the database storage.
func (f *DocumentCategoryRepo) FindDocumentCategory(ctx context.Context, r *entity.DocumentCategory) (*entity.DocumentCategory, error) {
	var dataEntity entity.DocumentCategory

	err := f.db.WithContext(ctx).Where("id = ?", r.ID).Take(&dataEntity).Error
	if err != nil {
		return nil, err
	}
	return &dataEntity, nil
}

// FindDocumentCategoryBySlug will find Document category from the database storage.
func (f *DocumentCategoryRepo) FindDocumentCategoryBySlug(ctx context.Context, r *entity.DocumentCategory) (*entity.DocumentCategory, error) {
	var dataEntity entity.DocumentCategory

	err := f.db.WithContext(ctx).Where("slug = ?", r.Slug).Take(&dataEntity).Error
	if err != nil {
		return nil, err
	}
	return &dataEntity, nil
}

// GetDocumentCategories will get Document categories from the database storage.
func (f *DocumentCategoryRepo) GetDocumentCategories(ctx context.Context, q *parameter.SQLQueryParameters) (entity.DocumentCategories, *parameter.ResponseMetadata, error) {
	var total int64
	var dataEntities entity.DocumentCategories

	errTotal := f.db.WithContext(ctx).Model(entity.DocumentCategory{}).Where(q.QueryKey, q.QueryValue...).Where(q.DateRange).Count(&total).Limit(q.Limit).Offset(q.Offset).Find(&dataEntities).Error
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

// SaveDocumentCategory will save Document category into the database storage.
func (f *DocumentCategoryRepo) SaveDocumentCategory(ctx context.Context, r *entity.DocumentCategory) (*entity.DocumentCategory, error) {
	var dataEntity entity.DocumentCategory
	dataEntity.ID = r.ID
	dataEntity.Name = r.Name
	dataEntity.Slug = r.Slug
	dataEntity.Size = r.Size
	dataEntity.MimeTypes = r.MimeTypes
	dataEntity.Description = r.Description

	err := f.db.WithContext(ctx).Create(&dataEntity).Error
	if err != nil {
		return nil, err
	}

	return &dataEntity, nil
}

// UpdateDocumentCategory will update Document category in the database storage.
func (f *DocumentCategoryRepo) UpdateDocumentCategory(ctx context.Context, r *entity.DocumentCategory) (*entity.DocumentCategory, error) {
	var dataEntity entity.DocumentCategory
	dataEntity.Name = r.Name
	dataEntity.Slug = r.Slug
	dataEntity.Size = r.Size
	dataEntity.MimeTypes = r.MimeTypes
	dataEntity.Description = r.Description

	err := f.db.WithContext(ctx).Where("id = ?", r.ID).Updates(dataEntity).First(&dataEntity).Error
	if err != nil {
		return nil, err
	}

	return &dataEntity, nil
}
