package registry

import (
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
)

// Entity is a struct represent of the entity name.
type Entity struct {
	Entity interface{}
}

// Table is a struct represent of the table name of the entity.
type Table struct {
	Name interface{}
}

// Registry is a struct holds the collections of Entity and Table.
type Registry struct {
	Entities []Entity
	Table    []Table
}

// Interface provides contract which is need to be implemented.
type Interface interface {
	AutoMigrate(db *gorm.DB) error
	ResetDatabase(db *gorm.DB) error
}

// AutoMigrate is a function uses to run auto migrate based on the schema of the Entity.
func (r *Registry) AutoMigrate(db *gorm.DB) error {
	var err error

	for _, model := range r.Entities {
		err = db.AutoMigrate(model.Entity)
		if err != nil {
			log.Fatal(err)
		}
	}

	return err
}

// ResetDatabase is a function uses to reset all table on the current database.
// It will be used for at test mode.
func (r *Registry) ResetDatabase(db *gorm.DB) error {
	var err error

	if os.Getenv("APP_ENV") == "production" {
		return nil
	}

	for _, table := range r.Table {
		errDrop := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table.Name)).Error
		if errDrop != nil {
			err = errDrop
			log.Fatal(err)
		}
	}

	for _, model := range r.Entities {
		errMigrate := db.AutoMigrate(model.Entity)
		if errMigrate != nil {
			err = errMigrate
			log.Fatal(err)
		}
	}

	return err
}
