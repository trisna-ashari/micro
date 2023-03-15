package dependency

import (
	"micro/persistence"
	"micro/pkg/configurator"
	"micro/pkg/logger"
	"net/http"
)

// Dependency holds all dependency to run the rest transport.
type Dependency struct {
	Config            *configurator.Config
	Logger            *logger.Logger
	DBClient          *persistence.DBClient
	HttpClient        *http.Client
	FileStorageClient *persistence.FileStorageClient
}
