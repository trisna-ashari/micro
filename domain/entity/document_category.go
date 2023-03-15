package entity

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

// DocumentCategory represent schema of table categories.
type DocumentCategory struct {
	ID          string    `gorm:"size:36;not null;unique_index;primary_key"`
	Slug        string    `gorm:"size:100;not null;index;"`
	Name        string    `gorm:"size:100;not null;index;"`
	Description string    `gorm:"size:255;not null;"`
	MimeTypes   string    `gorm:"size:255;not null;"`
	Size        float64   `gorm:"bigint;not null;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt
}

var _ Interface = &DocumentCategory{}

// DocumentCategories represent multiple DocumentCategory.
type DocumentCategories []*DocumentCategory

// TableName return name of table.
func (fc *DocumentCategory) TableName() string {
	return "document_categories"
}

// FilterableFields return fields.
func (fc *DocumentCategory) FilterableFields() []interface{} {
	return []interface{}{"name", "slug"}
}

// TimeFields return fields.
func (fc *DocumentCategory) TimeFields() []interface{} {
	return []interface{}{"created_at", "updated_at", "deleted_at"}
}

// BeforeCreate handle uuid generation.
func (fc *DocumentCategory) BeforeCreate(tx *gorm.DB) error {
	generateUUID := uuid.New()
	if fc.ID == "" {
		fc.ID = generateUUID.String()
	}
	return nil
}
