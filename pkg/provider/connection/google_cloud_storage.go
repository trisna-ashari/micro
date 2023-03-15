package connection

import (
	"context"
	"micro/pkg/configurator"

	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
)

// NewGCSConnection is a constructor will initialize connection to GCS server.
func NewGCSConnection(config *configurator.Config) (*storage.Client, error) {
	ctx := context.Background()
	gcsClient, err := storage.NewClient(ctx, option.WithCredentialsFile(config.GoogleApplicationCredential))

	if err != nil {
		return nil, err
	}

	return gcsClient, nil
}
