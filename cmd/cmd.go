package cmd

import (
	"github.com/urfave/cli/v2"
	"micro/domain/seeds"
	"micro/persistence"
	"micro/pkg/configurator"
	"micro/pkg/domain/registry"
	"micro/pkg/domain/seed"
	"micro/pkg/logger"
	"micro/transport/grpc/server"
	"net/http"
)

// NewCli is a constructor will initialize cli.
func NewCli() *cli.App {
	c := cli.NewApp()

	return c
}

// NewCommand construct a CLI commands.
func NewCommand(
	config *configurator.Config,
	registry *registry.Registry,
	dbClient *persistence.DBClient,
	httpClient *http.Client,
	fileStorageClient *persistence.FileStorageClient,
	logger *logger.Logger,
) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "db:migrate",
			Usage: "run database migration",
			Action: func(c *cli.Context) error {
				err := registry.AutoMigrate(dbClient.DB)
				if err != nil {
					logger.Log.Errorf("Error executing command, err: %v", err)
				}

				return nil
			},
		},
		{
			Name:  "db:init",
			Usage: "run predefined database initial seeder",
			Action: func(c *cli.Context) error {
				factory := seeds.NewInitialFactory()
				seeder := seed.New(factory)
				err := seeder.Seed(dbClient.DB)
				if err != nil {
					logger.Log.Errorf("Error executing command, err: %v", err)
				}

				return nil
			},
		},
		{
			Name:  "grpc:start",
			Usage: "Start the grpc server",
			Action: func(c *cli.Context) error {
				err := registry.AutoMigrate(dbClient.DB)
				if err != nil {
					logger.Log.Errorf("Error executing command, err: %v", err)
				}

				grpcServer := server.New(
					server.WithConfig(config),
					server.WithLogger(logger),
					server.WithDBClient(dbClient),
					server.WithHTTPClient(httpClient),
					server.WithFileStorageClient(fileStorageClient),
				)

				errRun := grpcServer.Init()
				if errRun != nil {
					return errRun
				}
				return nil
			},
		},
	}
}
