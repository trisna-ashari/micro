package server

import (
	"micro/persistence"
	"micro/pkg/configurator"
	"micro/pkg/logger"
	"net/http"
)

// Option return Server with Option.
type Option func(*Server)

// WithConfig is a function to set config to the Option.
func WithConfig(config *configurator.Config) Option {
	return func(r *Server) {
		r.config = config
	}
}

// WithLogger is a function to set logger to the Option.
func WithLogger(logger *logger.Logger) Option {
	return func(r *Server) {
		r.logger = logger
	}
}

// WithDBClient is a function to set DB client to the Option.
func WithDBClient(client *persistence.DBClient) Option {
	return func(r *Server) {
		r.dbClient = client
	}
}

// WithFileStorageClient is a function to set file storage client to the Option.
func WithFileStorageClient(client *persistence.FileStorageClient) Option {
	return func(r *Server) {
		r.fileStorageClient = client
	}
}

// WithHTTPClient is a function to set HTTP client to the Option.
func WithHTTPClient(client *http.Client) Option {
	return func(r *Server) {
		r.httpClient = client
	}
}
