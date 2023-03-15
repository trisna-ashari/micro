package persistence

import "micro/pkg/filestore"

// FileStorageClient is a struct itself.
type FileStorageClient struct {
	Driver filestore.Interface
}

// NewFileStoreService will initialize connection to FileStorageClient.
func NewFileStoreService(driver filestore.Interface) *FileStorageClient {
	fileStorage := filestore.FileStore{Driver: driver}

	return &FileStorageClient{
		Driver: fileStorage.Driver,
	}
}
