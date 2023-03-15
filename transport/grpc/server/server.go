package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"micro/persistence"
	"micro/pkg/configurator"
	"micro/pkg/logger"
	"micro/pkg/util"
	"micro/transport/grpc/dependency"
	"micro/transport/grpc/handler/healthcheck"
	"micro/transport/grpc/handler/v1/documentcategory"
	"micro/transport/grpc/interceptor/recovery"
	"net/http"
)

// Server holds the dependency to initialize a new one.
type Server struct {
	config            *configurator.Config
	logger            *logger.Logger
	dbClient          *persistence.DBClient
	httpClient        *http.Client
	fileStorageClient *persistence.FileStorageClient
}

// New will initialize a new Server.
func New(options ...Option) *Server {
	server := &Server{}

	for _, opt := range options {
		opt(server)
	}

	return server
}

// Init will start the Server.
func (s *Server) Init() error {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				resp, err = handler(ctx, req)
				if err != nil {
					return resp, err
				}

				return resp, err
			},
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				err := handler(srv, ss)

				return err
			},
		),
	)

	dep := &dependency.Dependency{
		Config:            s.config,
		Logger:            s.logger,
		DBClient:          s.dbClient,
		HttpClient:        s.httpClient,
		FileStorageClient: s.fileStorageClient,
	}

	healthCheckHandler := &healthcheck.Handler{Dependency: dep}
	documentCategoryHandler := &documentcategory.Handler{Dependency: dep}

	health.RegisterHealthServer(server, healthCheckHandler)

	// register gRPC handler
	documentcategory.RegisterDocumentCategoryServiceServer(server, documentCategoryHandler)

	// gRPC Server Reflection provides information about publicly-accessible gRPC services on a server,
	// and assists clients at runtime to construct RPC requests and responses without precompiled service information.
	// It is used by gRPCurl, which can be used to introspect server protos and send/receive test RPCs.
	// https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md
	reflection.Register(server)

	return util.RunGRPCServerWithGracefulShutdown(server, fmt.Sprintf(":%s", s.config.AppGRPCPort), s.logger)
}
