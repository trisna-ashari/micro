package registry

import (
	"micro/domain/entity"
	"micro/pkg/domain/registry"
)

// CollectEntities will return collections of replication entity.
func CollectEntities() []registry.Entity {
	return []registry.Entity{
		{Entity: entity.Document{}},
		{Entity: entity.DocumentCategory{}},
	}
}

// CollectTableNames will return collections of replication table name.
func CollectTableNames() []registry.Table {
	var Document entity.Document
	var DocumentCategory entity.DocumentCategory
	return []registry.Table{
		{Name: Document.TableName()},
		{Name: DocumentCategory.TableName()},
	}
}

// NewRegistry will initialize registry.Registry.
// Return []registry.Entity and []registry.Table.
// The registry is uses to auto migrate and reset database when running test mode.
func NewRegistry() *registry.Registry {
	var entityRegistry []registry.Entity
	var tableRegistry []registry.Table
	entityRegistry = append(entityRegistry, CollectEntities()...)
	tableRegistry = append(tableRegistry, CollectTableNames()...)

	return &registry.Registry{
		Entities: entityRegistry,
		Table:    tableRegistry,
	}
}
