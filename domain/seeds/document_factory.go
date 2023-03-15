package seeds

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"micro/domain/entity"
	"micro/pkg/domain/seed"
)

// InitFactory is a struct uses to hold []seed.Seed.
type InitFactory struct {
	seeders []seed.Seed
}

type documentCategoryFactory struct {
	category *entity.DocumentCategory
}

var seeders = []*documentCategoryFactory{
	{
		category: &entity.DocumentCategory{
			ID:          uuid.New().String(),
			Slug:        "original",
			Name:        "Original",
			Description: "",
			MimeTypes:   "application/pdf",
			Size:        10 * 10 * 10 * 1024,
		},
	},
	{
		category: &entity.DocumentCategory{
			ID:          uuid.New().String(),
			Slug:        "sign",
			Name:        "Sign",
			Description: "",
			MimeTypes:   "application/pdf",
			Size:        10 * 10 * 10 * 1024,
		},
	},
}

func (is *InitFactory) generateDocumentCategoriesSeeder() *InitFactory {
	for _, p := range seeders {
		cp := p
		is.seeders = append(is.seeders, seed.Seed{
			Name: "Create initial Document category",
			Run: func(db *gorm.DB) error {
				_, errDB := createDocumentCategory(db, cp)
				return errDB
			},
		})
	}

	return is
}

// createDocumentCategory will create predefined category and insert into DB.
func createDocumentCategory(db *gorm.DB, document *documentCategoryFactory) (*documentCategoryFactory, error) {
	var DocumentCategoryExists entity.DocumentCategory

	err := db.Where("slug = ?", document.category.Slug).Take(&DocumentCategoryExists).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := db.Create(document.category).Error
			if err != nil {
				return document, err
			}

			return document, err
		}

		return document, err
	}
	return document, err
}

func newInitFactory() *InitFactory {
	return &InitFactory{seeders: make([]seed.Seed, 0)}
}

// NewInitialFactory is a constructor will initialize InitFactory.
func NewInitialFactory() []seed.Seed {
	initialSeeds := newInitFactory()

	initialSeeds.generateDocumentCategoriesSeeder()

	return initialSeeds.seeders
}
