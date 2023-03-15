package filestore

import (
	"micro/pkg/filestore/object"
)

const (
	// ExpiredSignedURLTime represent signed URL expiration time in minute.
	ExpiredSignedURLTime = 15

	// TimeoutTime represent timout time in second.
	TimeoutTime = 10
)

// Interface is the interface that wraps FileStore interface.
type Interface interface {
	SignedURLInterface
	GetObjectInterface
	GetObjectURLInterface
	PutObjectInterface
	DuplicateObjectInterface
	DeleteObjectInterface
}

// SignedURLInterface is the interface that wraps generate signed URL method.
type SignedURLInterface interface {
	GenerateGetObjectSignedURL(objectPath string) (string, error)
	GeneratePutObjectSignedURL(object *object.Metadata) (string, error)
}

// GetObjectInterface is the interface that wraps the basic GetObject method.
type GetObjectInterface interface {
	GetObject(objectPath string) ([]byte, error)
}

// GetObjectURLInterface is the interface that wraps the basic GetObjectURL method.
type GetObjectURLInterface interface {
	GetObjectURL(objectPath string) (string, error)
}

// PutObjectInterface is the interface that wraps the basic PutObject method.
type PutObjectInterface interface {
	PutObject(object *object.Metadata) (*object.Metadata, error)
}

// DuplicateObjectInterface is the interface that wraps the basic DuplicateObject method.
type DuplicateObjectInterface interface {
	DuplicateObject(sourcePath string, targetPath string) error
}

// DeleteObjectInterface is the interface that wraps the basic DeleteObject method.
type DeleteObjectInterface interface {
	DeleteObject(objectPath string) error
}

// FileStore is a struct represent the storage driver.
type FileStore struct {
	Driver Interface
}

// New is a constructor will initialize FileStore.
func New(options ...Option) *FileStore {
	fileStore := &FileStore{Driver: nil}

	for _, opt := range options {
		opt(fileStore)
	}

	return fileStore
}
