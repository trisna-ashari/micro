package persistence

import (
	"gorm.io/gorm"
	"micro/domain/repository"
)

// DBClient represent it self.
type DBClient struct {
	DB               *gorm.DB
	Document         repository.DocumentRepositoryInterface
	DocumentCategory repository.DocumentCategoryRepositoryInterface
}

// NewDBService will initialize db connection and return repositories.
func NewDBService(db *gorm.DB) *DBClient {
	return &DBClient{
		DB:               db,
		Document:         NewDocumentRepository(db),
		DocumentCategory: NewDocumentCategoryRepository(db),
	}
}
