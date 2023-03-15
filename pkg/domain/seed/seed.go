package seed

import (
	"log"

	"gorm.io/gorm"
)

// Seed is a struct uses to holds the seeder name and the Run function to execute the seeder.
type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

// Seeder is a struct uses to holds collection of Seed.
type Seeder struct {
	Seeds []Seed
}

// New is a constructor will initialize Seeder.
func New(seeds []Seed) *Seeder {
	return &Seeder{Seeds: seeds}
}

// Seed is a functions uses to execute all available Seeder.
func (s *Seeder) Seed(db *gorm.DB) error {
	var err error

	for _, seed := range s.Seeds {
		errSeed := seed.Run(db)
		if errSeed != nil {
			err = errSeed
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}

	return err
}
