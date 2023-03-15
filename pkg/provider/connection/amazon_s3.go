package connection

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"micro/pkg/configurator"
)

// NewS3Connection is a constructor will initialize connection to S3 server.
func NewS3Connection(config *configurator.Config) (*s3.S3, error) {
	s3Session, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AWSRegion),
		Credentials: credentials.NewStaticCredentials(config.AWSAccessKeyID, config.AWSSecretAccessKey, ""),
	})

	if err != nil {
		return nil, err
	}

	_, err = s3Session.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}

	return s3.New(s3Session), nil
}
