package connection

import (
	"errors"
	"micro/pkg/configurator"
	"micro/pkg/filestore"
	"micro/pkg/filestore/driver/gcs"
	"micro/pkg/filestore/driver/minio"
	"micro/pkg/filestore/driver/s3"
)

const (
	driverGCS   = "gcs"
	driverS3    = "s3"
	driverMinio = "minio"
	driverLocal = "local"
)

// NewStorageConnection is a constructor will initialize connection to the storage server.
func NewStorageConnection(config *configurator.Config) (*filestore.FileStore, error) {
	var driver string
	if config.TestMode == false {
		driver = config.StorageConfig.Driver
	} else {
		driver = config.StorageTestConfig.Driver
	}

	switch driver {
	case driverGCS:
		gcsClient, errGCSConn := NewGCSConnection(config)
		if errGCSConn != nil {
			return nil, errGCSConn
		}

		fileStoreDriver := gcs.NewDriver(gcsClient, config, config.GCSBucketName, config.GCSPathPrefix)

		return &filestore.FileStore{Driver: fileStoreDriver}, nil
	case driverS3:
		s3Client, errS3Conn := NewS3Connection(config)
		if errS3Conn != nil {
			return nil, errS3Conn
		}

		fileStoreDriver := s3.NewDriver(s3Client, config, config.S3BucketName, config.S3PathPrefix)

		return &filestore.FileStore{Driver: fileStoreDriver}, nil
	case driverMinio:
		minioClient, errMinioConn := NewMinioConnection(config)
		if errMinioConn != nil {
			return nil, errMinioConn
		}

		fileStoreDriver := minio.NewDriver(minioClient, config, config.MinioBucketName, config.MinioPathPrefix)

		return &filestore.FileStore{Driver: fileStoreDriver}, nil
	default:
		return nil, errors.New("error.pkg.core.provider.connection.need_to_specify_storage_driver")
	}
}
