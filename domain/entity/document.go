package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Document represent schema of table Documents.
type Document struct {
	ID           string    `gorm:"size:36;not null;uniqueIndex;primary_key;"`
	CategoryID   string    `gorm:"size:36;not null;index;"`
	OriginalName string    `gorm:"size:255;not null;"`
	Name         string    `gorm:"size:255;not null;"`
	Path         string    `gorm:"size:255;not null;"`
	Type         string    `gorm:"size:36;not null;"`
	Size         int64     `gorm:"not null;"`
	Token        string    `gorm:"size:300;"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt
}

var _ Interface = &DocumentCategory{}

// Documents represent multiple Document.
type Documents []*Document

// TableName return name of table.
func (f *Document) TableName() string {
	return "documents"
}

// FilterableFields return fields.
func (f *Document) FilterableFields() []interface{} {
	return []interface{}{"name", "slug"}
}

// TimeFields return fields.
func (f *Document) TimeFields() []interface{} {
	return []interface{}{"created_at", "updated_at", "deleted_at"}
}

// BeforeCreate handle uuid generation.
func (f *Document) BeforeCreate(tx *gorm.DB) error {
	generateUUID := uuid.New()
	if f.ID == "" {
		f.ID = generateUUID.String()
	}

	defaultTime := time.Time{}

	if f.CreatedAt == defaultTime {
		f.CreatedAt = time.Now()
	}

	if f.UpdatedAt == defaultTime {
		f.UpdatedAt = time.Now()
	}
	return nil
}
