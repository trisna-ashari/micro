package connection

import (
	"errors"
	"fmt"
	"micro/pkg/configurator"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/lib/pq"

	consoleLog "log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	sqlTracer "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	dbLogger "gorm.io/gorm/logger"
)

const (
	driverMysql    = "mysql"
	driverPostgres = "postgres"
)

// DSNCollection holds DSN string.
type DSNCollection struct {
	driver string
	dsn    string
}

// DSNCollections holds multiple DSN.
type DSNCollections []DSNCollection

// ToGormDialects transform into multiple gorm dialects.
func (d DSNCollections) ToGormDialects() []gorm.Dialector {
	var gormDialects []gorm.Dialector
	for _, dsn := range d {
		if dsn.driver == driverPostgres {
			gormDialects = append(gormDialects, postgres.Open(dsn.dsn))
		}

		if dsn.driver == driverMysql {
			gormDialects = append(gormDialects, mysql.Open(dsn.dsn))
		}
	}

	return gormDialects
}

// NewDSN is a function uses to construct DSN for the sql database connection.
func NewDSN(config *configurator.Config) string {
	var driver string
	if config.TestMode == false {
		driver = config.DBConfig.DBDriver
	} else {
		driver = config.DBTestConfig.DBDriver
	}
	switch driver {
	case driverPostgres:
		if config.TestMode == false {
			dsn := &PostgresDSN{
				Host:     config.DBConfig.DBHost,
				Port:     config.DBConfig.DBPort,
				User:     config.DBConfig.DBUser,
				Password: config.DBConfig.DBPassword,
				DBName:   config.DBConfig.DBName,
				SSLMode:  false,
				Timezone: config.DBConfig.DBTimeZone,
			}

			return dsn.ToString()
		}

		dsn := &PostgresDSN{
			Host:     config.DBTestConfig.DBHost,
			Port:     config.DBTestConfig.DBPort,
			User:     config.DBTestConfig.DBUser,
			Password: config.DBTestConfig.DBPassword,
			DBName:   config.DBTestConfig.DBName,
			SSLMode:  false,
			Timezone: config.DBTestConfig.DBTimeZone,
		}

		return dsn.ToString()
	case driverMysql:
		if config.TestMode == false {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
				config.DBConfig.DBUser,
				config.DBConfig.DBPassword,
				config.DBConfig.DBHost,
				config.DBConfig.DBPort,
				config.DBConfig.DBName,
			)
			return dsn
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.DBTestConfig.DBUser,
			config.DBTestConfig.DBPassword,
			config.DBTestConfig.DBHost,
			config.DBTestConfig.DBPort,
			config.DBTestConfig.DBName,
		)
		return dsn
	default:
		return ""
	}
}

// NewReplicaDSNCollections is a function uses to construct multiple DSN for the sql database connection.
func NewReplicaDSNCollections(configs configurator.DBReplicasConfig) DSNCollections {
	var dsnCollection DSNCollections
	for _, config := range configs {
		driver := config.DBDriver
		switch driver {
		case driverPostgres:
			dsn := &PostgresDSN{
				Host:     config.DBHost,
				Port:     config.DBPort,
				User:     config.DBUser,
				Password: config.DBPassword,
				DBName:   config.DBName,
				SSLMode:  false,
				Timezone: config.DBTimeZone,
			}

			dsnCollection = append(dsnCollection, DSNCollection{
				driver: driverPostgres,
				dsn:    dsn.ToString(),
			})
		case driverMysql:
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
				config.DBUser,
				config.DBPassword,
				config.DBHost,
				config.DBPort,
				config.DBName,
			)

			dsnCollection = append(dsnCollection, DSNCollection{
				driver: driverPostgres,
				dsn:    dsn,
			})
		}
	}

	return dsnCollection
}

// NewDBConnection is a constructor will initialize sql database connection.
//gocyclo:ignore
func NewDBConnection(config *configurator.Config) (*gorm.DB, error) {
	newLogger := dbLogger.New(
		consoleLog.New(os.Stdout, "\r\n", consoleLog.LstdFlags),
		dbLogger.Config{
			SlowThreshold: time.Second,
			LogLevel:      dbLogger.Info,
			Colorful:      true,
		},
	)

	var disableFKConstraint bool
	if config.TestMode == false {
		disableFKConstraint = config.DBConfig.DisableForeignKeyConstraint
	} else {
		disableFKConstraint = config.DBTestConfig.DisableForeignKeyConstraint
	}

	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: disableFKConstraint,
		PrepareStmt:                              config.EnableCachePrepareStmt,
	}

	if config.DBConfig.DBLog {
		gormConfig.Logger = newLogger
	}

	var driver string
	if config.TestMode == false {
		driver = config.DBConfig.DBDriver
	} else {
		driver = config.DBTestConfig.DBDriver
	}

	switch driver {
	case driverPostgres:
		if config.DataDogConfig.EnableTracer {
			sqlTracer.Register("lib/pq", pq.Driver{},
				sqlTracer.WithServiceName(strings.ReplaceAll(config.AppName, "-http", "")+"-db__"+config.DBConfig.DBName),
				sqlTracer.WithDSN(config.DBConfig.DBName),
				sqlTracer.WithAnalytics(true),
			)
			sqlDb, err := sqlTracer.Open("lib/pq", NewDSN(config))

			db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDb}), gormConfig)
			if err != nil {
				return nil, err
			}

			db, err = setConnectionPool(db, config)
			if err != nil {
				return nil, err
			}

			return db, nil
		}

		db, err := gorm.Open(postgres.Open(NewDSN(config)), gormConfig)
		if err != nil {
			return nil, err
		}

		db, err = setConnectionPool(db, config)
		if err != nil {
			return nil, err
		}

		return db, nil

	case driverMysql:
		if config.DataDogConfig.EnableTracer {
			sqlTracer.Register("lib/pq", &stdlib.Driver{},
				sqlTracer.WithServiceName(strings.ReplaceAll(config.AppName, "-http", "")+"-db__"+config.DBConfig.DBName),
				sqlTracer.WithDSN(config.DBConfig.DBName),
				sqlTracer.WithAnalytics(true),
			)
			sqlDb, err := sqlTracer.Open("lib/pq", NewDSN(config))

			db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDb}), gormConfig)
			if err != nil {
				return nil, err
			}

			db, err = setConnectionPool(db, config)
			if err != nil {
				return nil, err
			}

			return db, nil
		}

		db, err := gorm.Open(mysql.Open(NewDSN(config)), gormConfig)
		if err != nil {
			return nil, err
		}

		db, err = setConnectionPool(db, config)
		if err != nil {
			return nil, err
		}

		return db, nil

	default:
		return nil, errors.New("common.error.database.driver.not_defined")
	}
}

func setConnectionPool(db *gorm.DB, config *configurator.Config) (*gorm.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleCons)
	sqlDB.SetMaxOpenConns(config.MaxOpenCons)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
