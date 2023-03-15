package object

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"micro/pkg/util"
	"mime/multipart"
	"path/filepath"
	"text/template"
	"time"

	"github.com/google/uuid"

	"github.com/gabriel-vasile/mimetype"
)

// PutMethod represent how the object will be uploaded.
type PutMethod int

// PDFOverwriteMode represent how the system acts when facing PDF file that contains digital signatures
type PDFOverwriteMode int

const (
	// DirectPut represent the object will be uploaded directly.
	DirectPut PutMethod = iota

	// SignedURLPut represent the object will be uploaded via generated PUT signed URL.
	SignedURLPut
)

const (
	// Unspecified represent overwrite mode is not specified,
	// it will trigger error to ask user decisions about the PDF file.
	Unspecified PDFOverwriteMode = iota

	// KeepSignatures represent overwrite mode to keep all signatures
	// which is already placed in the PDF file. So, the PDF processor
	// will skip watermarking (QRCode stamp) process.
	KeepSignatures

	// RemoveSignatures represent overwrite mode to remove all signatures
	// which is already placed in the PDF file. So, the PDF processor
	// will perform watermarking (QRCode stamp) process.
	RemoveSignatures
)

// Metadata is a struct represent the object metadata.
type Metadata struct {
	ID           string
	Token        string
	Name         string
	NamePrefix   string
	NameSuffix   string
	OriginalName string
	ContentType  string
	Slug         string
	Date         string
	Extension    string
	SourcePath   string
	Password     string
	Size         int64
	Content      []byte
	Source       string

	IncludeSlug bool
	IncludeDate bool

	CustomPath string

	// QRCodeLogoPath holds
	QRCodeLogoPath string

	// PutMethod defines how object will be uploaded (DirectPut or SignedURLPut).
	PutMethod PutMethod

	// PutSignedURL holds generated signed URL. Its mandatory when PutMethod is SignedURLPut.
	PutSignedURL string

	// PDFOverwrite holds overwrite mode (true/false) to handle PDF that already have digital signatures.
	// Which is true to keep digital signatures and false to remove all digital signatures
	// before the PDF file is processed in the next step.
	// If PDFOverwrite false, it should not be watermarked with QRCode and the opposite.
	PDFOverwrite PDFOverwriteMode
}

// NewFromFilePath is a function to generate Metadata from the given file path.
func NewFromFilePath(filePath string, slug string, opts ...Option) *Metadata {
	objectReader, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	objectMime := mimetype.Detect(objectReader)
	objectMetadata := &Metadata{
		Name:           uuid.New().String(),
		OriginalName:   filepath.Base(filePath),
		Slug:           slug,
		Date:           time.Now().Format("2006/01/02"),
		Size:           int64(len(objectReader)),
		Content:        objectReader,
		ContentType:    objectMime.String(),
		Extension:      objectMime.Extension(),
		PutMethod:      SignedURLPut,
		SourcePath:     filePath,
		PDFOverwrite:   Unspecified,
		IncludeSlug:    false,
		IncludeDate:    false,
		QRCodeLogoPath: fmt.Sprintf("%s/core/assets/images/logo-qrcode.png", util.RootDir()),
	}

	for _, opt := range opts {
		opt(objectMetadata)
	}

	return objectMetadata
}

// NewFromMultipartFileHeader is a function to generate Metadata from the multipart form file header.
func NewFromMultipartFileHeader(file *multipart.FileHeader, path string, opts ...Option) *Metadata {
	objectOpen, err := file.Open()
	if err != nil {
		log.Print(err)
	}
	defer objectOpen.Close()

	objectSize := file.Size
	objectReader := make([]byte, objectSize)
	_, err = objectOpen.Read(objectReader)
	if err != nil {
		log.Print(err)
	}

	objectMime := mimetype.Detect(objectReader)
	objectMetadata := &Metadata{
		Name:           uuid.New().String(),
		OriginalName:   file.Filename,
		Slug:           path,
		Date:           time.Now().Format("2006/01/02"),
		Size:           objectSize,
		Content:        objectReader,
		ContentType:    objectMime.String(),
		Extension:      objectMime.Extension(),
		PutMethod:      SignedURLPut,
		PDFOverwrite:   Unspecified,
		IncludeSlug:    false,
		IncludeDate:    false,
		QRCodeLogoPath: fmt.Sprintf("%s/assets/images/logo-qrcode.png", util.RootDir()),
	}

	for _, opt := range opts {
		opt(objectMetadata)
	}

	return objectMetadata
}

// NewFromByteSlice is a function to generate Metadata from the given byte slice.
func NewFromByteSlice(fileBytes []byte, path string, opts ...Option) *Metadata {
	ObjectMime := mimetype.Detect(fileBytes)

	objectMetadata := &Metadata{
		ID:             uuid.New().String(),
		Name:           uuid.New().String(),
		OriginalName:   uuid.New().String() + ObjectMime.Extension(),
		Slug:           path,
		Date:           time.Now().Format("2006/01/02"),
		Size:           int64(len(fileBytes)),
		Content:        fileBytes,
		ContentType:    ObjectMime.String(),
		Extension:      ObjectMime.Extension(),
		PutMethod:      SignedURLPut,
		PDFOverwrite:   Unspecified,
		IncludeSlug:    false,
		IncludeDate:    false,
		QRCodeLogoPath: fmt.Sprintf("%s/assets/images/logo-qrcode.png", util.RootDir()),
	}

	for _, opt := range opts {
		opt(objectMetadata)
	}

	return objectMetadata
}

// Filename is a method uses to generate file name.
func (m *Metadata) Filename(opts ...Option) string {
	for _, opt := range opts {
		opt(m)
	}

	return m.generateFileName()
}

// Filepath is a method uses to generate file path with slug.
func (m *Metadata) Filepath(opts ...Option) string {
	for _, opt := range opts {
		opt(m)
	}

	if m.CustomPath != "" {
		return m.CustomPath
	}

	return m.generateFilePath()
}

// CopyFilepath is a method uses to copy file path.
func (m *Metadata) CopyFilepath(opts ...Option) string {
	copyMetadata := &Metadata{
		Name:        m.Name,
		NamePrefix:  m.NamePrefix,
		NameSuffix:  m.NameSuffix,
		Slug:        m.Slug,
		Date:        m.Date,
		Extension:   m.Extension,
		IncludeSlug: m.IncludeSlug,
		IncludeDate: m.IncludeDate,
	}

	for _, opt := range opts {
		opt(copyMetadata)
	}

	return copyMetadata.generateFilePath()
}

func (m *Metadata) generateFileName() string {
	filePathTemplate, _ := template.
		New("").
		Parse("{{.NamePrefix}}{{.Name}}{{.NameSuffix}}{{.Extension}}")

	bufferStr := bytes.Buffer{}
	_ = filePathTemplate.Execute(&bufferStr, m)

	return bufferStr.String()
}

func (m *Metadata) generateFilePath() string {
	filePathTemplate, _ := template.
		New("").
		Parse("{{.NamePrefix}}{{.Name}}{{.NameSuffix}}{{.Extension}}")

	if m.IncludeSlug {
		filePathTemplate, _ = template.
			New("").
			Parse("{{.Slug}}/{{.NamePrefix}}{{.Name}}{{.NameSuffix}}{{.Extension}}")
	}

	if m.IncludeDate {
		if m.Date == "" {
			m.Date = time.Now().Format("2006/01/02")
		}

		filePathTemplate, _ = template.
			New("").
			Parse("{{.Date}}/{{.NamePrefix}}{{.Name}}{{.NameSuffix}}{{.Extension}}")
	}

	if m.IncludeSlug && m.IncludeDate {
		filePathTemplate, _ = template.
			New("").
			Parse("{{.Slug}}/{{.Date}}/{{.NamePrefix}}{{.Name}}{{.NameSuffix}}{{.Extension}}")
	}

	bufferStr := bytes.Buffer{}
	_ = filePathTemplate.Execute(&bufferStr, m)

	return bufferStr.String()
}
