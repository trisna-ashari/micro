package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"micro/cmd"
	"micro/docs"
	"micro/domain/registry"
	"micro/persistence"
	"micro/pkg/provider/connection"
	"micro/transport/rest/router"
	"net/http"

	"net"
	"os"
	"time"

	ginSwaggerFiles "github.com/swaggo/files"
	ginSwaggerHandler "github.com/swaggo/gin-swagger"

	"github.com/urfave/cli/v2"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"micro/pkg/configurator"
	"micro/pkg/logger"
	"micro/pkg/util"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file provided")
	}

	config := configurator.New(
		configurator.WithDBConfig(),
		configurator.WithGoogleCloudServiceConfig(),
		configurator.WithAmazonWebServiceConfig(),
		configurator.WithMinioConfig(),
		configurator.WithGCSConfig(),
		configurator.WithS3Config(),
		configurator.WithStorageConfig(),
		configurator.WithDatadogConfig(),
	)

	timeLoc, _ := time.LoadLocation(config.AppTimezone)
	time.Local = timeLoc

	var logStd *logger.Logger
	switch config.AppEnvironment {
	case "production", "staging", "development":
		logStd = logger.New(logger.NewProductionConfig())
	default:
		logStd = logger.New(logger.NewDevelopmentConfig())
	}

	if config.DataDogConfig.EnableTracer {
		tracer.Start(
			tracer.WithEnv(config.AppEnvironment+"-"+os.Getenv("DD_ENV_SUFFIX")),
			tracer.WithService(config.AppName),
			tracer.WithServiceVersion(config.AppVersion),
			tracer.WithAgentAddr(net.JoinHostPort(config.DataDogConfig.AgentHost, config.DataDogConfig.AgentPort)),
			tracer.WithProfilerCodeHotspots(true),
			tracer.WithRuntimeMetrics(),
		)
		defer tracer.Stop()
	}

	dbConn, errDBConn := connection.NewDBConnection(config)
	if errDBConn != nil {
		logStd.Log.Fatalf("Unable to connect database: %v", errDBConn)
	}

	dbClient := persistence.NewDBService(dbConn)

	fileStorageConn, errStorageConnection := connection.NewStorageConnection(config)
	if errStorageConnection != nil {
		logStd.Log.Error(errStorageConnection.Error())
	}
	fileStorageClient := persistence.NewFileStoreService(fileStorageConn.Driver)

	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.MaxIdleConns = 100
	httpTransport.MaxConnsPerHost = 100
	httpTransport.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Timeout:   60 * time.Second,
		Transport: httpTransport,
	}

	entityRegistry := registry.NewRegistry()

	// Swagger docs
	docs.SwaggerInfo.Host = os.Getenv("APP_SWAGGER_HOST")
	docs.SwaggerInfo.BasePath = ""

	app := cmd.NewCli()
	app.Commands = cmd.NewCommand(
		config,
		entityRegistry,
		dbClient,
		httpClient,
		fileStorageClient,
		logStd,
	)
	app.Action = func(c *cli.Context) error {
		errAutoMigrate := entityRegistry.AutoMigrate(dbConn)
		if errAutoMigrate != nil {
			logStd.Log.Fatalf("Unable to run database migration: %v", errAutoMigrate)
		}

		httpRouter := router.
			New(
				router.WithConfig(config),
				router.WithLogger(logStd),
				router.WithDBClient(dbClient),
				router.WithHTTPClient(httpClient),
				router.WithFileStorageClient(fileStorageClient),
			).
			Init()

		if config.AppEnvironment != "production" {
			httpRouter.GET("/swagger/*any", ginSwaggerHandler.WrapHandler(ginSwaggerFiles.Handler))
		}

		appPort := config.AppHTTPPort
		if appPort == "" {
			appPort = "6969"
		}

		shutdownTimeout := 10 * time.Second
		if err := util.RunHTTPServerWithGracefulShutdown(httpRouter, fmt.Sprintf(":%s", appPort), shutdownTimeout, logStd); err != nil {
			return err
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logStd.Log.Fatalf("Unable to run CLI command, err: %v", err)
	}
}
