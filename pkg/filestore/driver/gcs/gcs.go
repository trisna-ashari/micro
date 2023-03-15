package gcs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"

	"micro/pkg/configurator"
	"micro/pkg/filestore"
	"micro/pkg/filestore/object"
	"micro/pkg/util"
)

// Driver is a struct represent dependencies needed to be initialized.
type Driver struct {
	client     *storage.Client
	config     *configurator.Config
	bucketName string
	pathPrefix string
}

// NewDriver is a constructor will initialize Driver.
func NewDriver(client *storage.Client, config *configurator.Config, bucketName string, pathPrefix string) *Driver {
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
	jsonKey, err := ioutil.ReadFile(d.config.GoogleApplicationCredential)
	if err != nil {
		return "", fmt.Errorf("cannot read the JSON key file, err: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}

	path := util.MakePathWithPrefix(d.pathPrefix, objectPath)
	signedURL, err := storage.SignedURL(d.bucketName, path, opts)
	if err != nil {
		return "", fmt.Errorf("google.GenerateGetObjectSignedURL: %v", err)
	}

	return signedURL, nil
}

// GeneratePutObjectSignedURL is a method uses to generate PUT signed URL.
func (d *Driver) GeneratePutObjectSignedURL(m *object.Metadata) (string, error) {
	jsonKey, err := ioutil.ReadFile(d.config.GoogleApplicationCredential)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadFile: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}

	opts := &storage.SignedURLOptions{
		Scheme: storage.SigningSchemeV4,
		Method: "PUT",
		Headers: []string{
			fmt.Sprintf("Content-Type:%s", m.ContentType),
		},
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(filestore.ExpiredSignedURLTime * time.Minute),
	}

	path := util.MakePathWithPrefix(d.pathPrefix, m.Filepath())
	signedURL, err := storage.SignedURL(d.bucketName, path, opts)
	if err != nil {
		return "", fmt.Errorf("storage.GeneratePutObjectSignedURL: %v", err)
	}

	return signedURL, nil
}

// GetObject is a method uses to get an object.
func (d *Driver) GetObject(objectPath string) ([]byte, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, filestore.TimeoutTime*time.Second)
	defer cancel()

	path := util.MakePathWithPrefix(d.pathPrefix, objectPath)
	rc, err := d.client.Bucket(d.bucketName).Object(path).NewReader(ctx)
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
		return m, fmt.Errorf("filestore.driver.gcs.PutObject: %v", "invalid put method")
	}
}

// DuplicateObject is a method uses to duplicate an object to specific path.
func (d *Driver) DuplicateObject(sourcePath string, targetPath string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, filestore.TimeoutTime*time.Second)
	defer cancel()

	src := d.client.Bucket(d.bucketName).Object(util.MakePathWithPrefix(d.pathPrefix, sourcePath))
	dst := d.client.Bucket(d.bucketName).Object(util.MakePathWithPrefix(d.pathPrefix, targetPath))
	_, err := dst.CopierFrom(src).Run(ctx)
	if err != nil {
		return fmt.Errorf("filestore.driver.gcs.DuplicateObject: %v", err)
	}

	return nil
}

// DeleteObject is a method uses to delete an object.
func (d *Driver) DeleteObject(objectPath string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, filestore.TimeoutTime*time.Second)
	defer cancel()

	path := util.MakePathWithPrefix(d.pathPrefix, objectPath)
	o := d.client.Bucket(d.bucketName).Object(path)
	if err := o.Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (d *Driver) putObjectDirectly(m *object.Metadata) (*object.Metadata, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, filestore.TimeoutTime*time.Second)
	defer cancel()

	path := util.MakePathWithPrefix(d.pathPrefix, m.Filepath())
	storageWriter := d.client.Bucket(d.bucketName).Object(path).NewWriter(ctx)
	storageWriter.ContentType = m.ContentType
	if _, err := io.Copy(storageWriter, bytes.NewReader(m.Content)); err != nil {
		return m, fmt.Errorf("storage.Copy: %s", err)
	}

	if err := storageWriter.Close(); err != nil {
		return m, fmt.Errorf("storage.Close: %s", err)
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
