package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"micro/persistence"
	"micro/pkg/configurator"
	"micro/pkg/logger"
	"micro/transport/rest/dependency"
	"micro/transport/rest/handler/ping"
	"micro/transport/rest/middleware"
	"net/http"
	"strings"
)

// Router holds the dependency to initialize a new one.
type Router struct {
	config            *configurator.Config
	logger            *logger.Logger
	dbClient          *persistence.DBClient
	httpClient        *http.Client
	fileStorageClient *persistence.FileStorageClient
}

// New will initialize a new Router.
func New(options ...Option) *Router {
	router := &Router{}

	for _, opt := range options {
		opt(router)
	}

	return router
}

// Init will start the Router.
func (r *Router) Init() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	e := gin.Default()

	guardMiddleware := middleware.NewGuard(r.logger)
	guardMiddleware.SetAuthenticationMethod(middleware.AuthTypeBasicAppKey)
	errorMiddleware := middleware.NewError(r.config, r.logger)

	e.Use(errorMiddleware.ErrorHandler())

	dep := &dependency.Dependency{
		Config:            r.config,
		Logger:            r.logger,
		DBClient:          r.dbClient,
		HttpClient:        r.httpClient,
		FileStorageClient: r.fileStorageClient,
	}

	pingHandler := &ping.Handler{Dependency: dep}

	_ = e.Group("/api/v1", func(c *gin.Context) {
		if strings.Contains(c.Request.Referer(), "#") {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("common.error.uri_contains_illegal_character"))
			return
		}
	})

	e.GET("/ping", pingHandler.Ping)

	return e
}
