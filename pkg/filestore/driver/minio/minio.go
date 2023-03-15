package minio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	minio "github.com/minio/minio-go/v7"
	"micro/pkg/configurator"
	"micro/pkg/filestore"
	"micro/pkg/filestore/object"
	"micro/pkg/util"
)

// Driver is a struct represent dependencies needed to be initialized.
type Driver struct {
	client     *minio.Client
	config     *configurator.Config
	bucketName string
	pathPrefix string
}

// NewDriver is a constructor will initialize Driver.
func NewDriver(client *minio.Client, config *configurator.Config, bucketName string, pathPrefix string) *Driver {
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
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, filestore.TimeoutTime*time.Second)
	defer cancel()

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("inline; filename=%s", filepath.Base(objectPath)))

	path := util.MakePathWithPrefix(d.pathPrefix, objectPath)
	signedURL, err := d.client.PresignedGetObject(ctx, d.bucketName, path, filestore.ExpiredSignedURLTime*time.Minute, reqParams)
	if err != nil {
		return "", fmt.Errorf("request.GenerateGetObjectSignedURL: %v", err)
	}

	return signedURL.String(), nil
}

// GeneratePutObjectSignedURL is a method uses to generate PUT signed URL.
func (d *Driver) GeneratePutObjectSignedURL(m *object.Metadata) (string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, filestore.TimeoutTime*time.Second)
	defer cancel()

	path := util.MakePathWithPrefix(d.pathPrefix, m.Filepath())
	signedURL, err := d.client.PresignedPutObject(ctx, d.bucketName, path, filestore.ExpiredSignedURLTime*time.Minute)
	if err != nil {
		return "", fmt.Errorf("request.GeneratePutObjectSignedURL: %v", err)
	}

	return signedURL.String(), nil
}

// GetObject is a method uses to get an object.
func (d *Driver) GetObject(objectPath string) ([]byte, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, filestore.TimeoutTime*time.Second)
	defer cancel()

	opts := minio.GetObjectOptions{}

	path := util.MakePathWithPrefix(d.pathPrefix, objectPath)
	rc, err := d.client.GetObject(ctx, d.bucketName, path, opts)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
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
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, filestore.TimeoutTime*time.Second)
	defer cancel()

	srcOpts := minio.CopySrcOptions{
		Bucket: d.bucketName,
		Object: util.MakePathWithPrefix(d.pathPrefix, sourcePath),
	}

	dstOpts := minio.CopyDestOptions{
		Bucket: d.bucketName,
		Object: util.MakePathWithPrefix(d.pathPrefix, targetPath),
	}

	_, err := d.client.CopyObject(ctx, dstOpts, srcOpts)
	if err != nil {
		return fmt.Errorf("filestore.driver.s3.DuplicateObject: %v", err)
	}

	return nil
}

// DeleteObject is a method uses to delete an object.
func (d *Driver) DeleteObject(objectPath string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, filestore.TimeoutTime*time.Second)
	defer cancel()

	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	path := util.MakePathWithPrefix(d.pathPrefix, objectPath)
	err := d.client.RemoveObject(ctx, d.bucketName, path, opts)
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) putObjectDirectly(m *object.Metadata) (*object.Metadata, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, filestore.TimeoutTime*time.Second)
	defer cancel()

	opts := minio.PutObjectOptions{
		ContentType: m.ContentType,
	}

	path := util.MakePathWithPrefix(d.pathPrefix, m.Filepath())
	_, err := d.client.PutObject(ctx, d.bucketName, path, bytes.NewReader(m.Content), int64(len(m.Content)), opts)
	if err != nil {
		return m, fmt.Errorf("storage.PutObject: %s", err)
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
