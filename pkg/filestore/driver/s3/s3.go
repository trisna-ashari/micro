package s3

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"micro/pkg/configurator"
	"micro/pkg/filestore"
	"micro/pkg/filestore/object"
	"micro/pkg/util"
)

// Driver is a struct represent dependencies needed to be initialized.
type Driver struct {
	client     *s3.S3
	config     *configurator.Config
	bucketName string
	pathPrefix string
}

// NewDriver is a constructor will initialize Driver.
func NewDriver(client *s3.S3, config *configurator.Config, bucketName string, pathPrefix string) *Driver {
	return &Driver{
		client:     client,
		config:     config,
		bucketName: bucketName,
		pathPrefix: pathPrefix,
	}
}

// Type assertions to make sure Driver already implement filestore.Interface.
var _ filestore.Interface = &Driver{}

// GenerateGetObjectSignedURL is a method uses to generate GET signed URL.
func (d *Driver) GenerateGetObjectSignedURL(objectPath string) (string, error) {
	opts := &s3.GetObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(objectPath),
	}

	req, _ := d.client.GetObjectRequest(opts)
	signedURL, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("request.GeneratePutObjectSignedURL: %v", err)
	}

	return signedURL, nil
}

// GeneratePutObjectSignedURL is a method uses to generate PUT signed URL.
func (d *Driver) GeneratePutObjectSignedURL(m *object.Metadata) (string, error) {
	path := util.MakePathWithPrefix(d.pathPrefix, m.Filepath())
	opts := &s3.PutObjectInput{
		Bucket:      aws.String(d.bucketName),
		Key:         aws.String(path),
		ContentType: aws.String(m.ContentType),
	}

	req, _ := d.client.PutObjectRequest(opts)
	signedURL, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("request.GeneratePutObjectSignedURL: %v", err)
	}

	return signedURL, nil
}

// GetObject is a method uses to get an object.
func (d *Driver) GetObject(objectPath string) ([]byte, error) {
	path := util.MakePathWithPrefix(d.pathPrefix, objectPath)
	opts := &s3.GetObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	}

	rc, err := d.client.GetObject(opts)
	if err != nil {
		return nil, fmt.Errorf("request.GetObject: %v", err)
	}
	defer rc.Body.Close()

	data, err := ioutil.ReadAll(rc.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetObjectURL is a method uses to get an object URL.
func (d *Driver) GetObjectURL(objectPath string) (string, error) {
	return d.GenerateGetObjectSignedURL(objectPath)
}

// PutObject is a method uses to upload an object.
func (d *Driver) PutObject(m *object.Metadata) (*object.Metadata, error) {
	switch m.PutMethod {
	case object.DirectPut:
		return d.putObjectDirectly(m)
	case object.SignedURLPut:
		return d.putObjectViaSignedURL(m)
	default:
		return m, errors.New("unknown put method")
	}
}

// DuplicateObject is a method uses to duplicate an object to specific path.
func (d *Driver) DuplicateObject(sourcePath string, targetPath string) error {
	opts := &s3.CopyObjectInput{
		Bucket:     aws.String(d.bucketName),
		CopySource: aws.String(fmt.Sprintf("%s/%s", d.bucketName, util.MakePathWithPrefix(d.pathPrefix, sourcePath))),
		Key:        aws.String(util.MakePathWithPrefix(d.pathPrefix, sourcePath)),
	}

	_, err := d.client.CopyObject(opts)
	if err != nil {
		return fmt.Errorf("filestore.driver.s3.DuplicateObject: %v", err)
	}

	return nil
}

// DeleteObject is a method uses to delete an object.
func (d *Driver) DeleteObject(objectPath string) error {
	path := util.MakePathWithPrefix(d.pathPrefix, objectPath)
	opts := &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(path),
	}

	_, err := d.client.DeleteObject(opts)
	if err != nil {
		return fmt.Errorf("request.DeleteObject: %v", err)
	}

	return nil
}

func (d *Driver) putObjectDirectly(m *object.Metadata) (*object.Metadata, error) {
	path := util.MakePathWithPrefix(d.pathPrefix, m.Filepath())
	opts := &s3.PutObjectInput{
		Bucket:             aws.String(d.bucketName),
		Key:                aws.String(path),
		Body:               bytes.NewReader(m.Content),
		ContentDisposition: aws.String("attachment"),
	}

	_, err := d.client.PutObject(opts)
	if err != nil {
		return m, fmt.Errorf("request.PutObject: %v", err)
	}

	return m, nil
}

func (d *Driver) putObjectViaSignedURL(m *object.Metadata) (*object.Metadata, error) {
	if m.PutSignedURL == "" {
		return m, errors.New("signed URL is empty")
	}

	httpClient := &http.Client{}
	request, err := http.NewRequest("PUT", m.PutSignedURL, bytes.NewReader(m.Content))
	if err != nil {
		return m, fmt.Errorf("httpClient.NewRequest: %s", err)
	}

	request.Header.Set("Content-Type", m.ContentType)
	_, err = httpClient.Do(request)
	if err != nil {
		return m, fmt.Errorf("httpClient.Do: %v", err)
	}

	return m, nil
}
