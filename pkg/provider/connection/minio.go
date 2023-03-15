package connection

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"micro/pkg/configurator"
)

// NewMinioConnection is a constructor will initialize connection to minio server.
func NewMinioConnection(config *configurator.Config) (*minio.Client, error) {
	minioClient, err := minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
