package connection

import (
	"fmt"
	"strings"
)

// PostgresDSN is a struct for Postgres database connection DSN configuration.
type PostgresDSN struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  bool
	Timezone string
}

// ToString export to DSN string format.
func (dsn PostgresDSN) ToString() string {
	var s []string

	if dsn.Host != "" {
		s = append(s, fmt.Sprintf("host=%s", dsn.Host))
	} else {
		s = append(s, fmt.Sprintf("host=%s", "localhost"))
	}

	if dsn.Port != "" {
		s = append(s, fmt.Sprintf("port=%s", dsn.Port))
	} else {
		s = append(s, fmt.Sprintf("port=%s", "5432"))
	}

	if dsn.User != "" {
		s = append(s, fmt.Sprintf("user=%s", dsn.User))
	}

	if dsn.Password != "" {
		s = append(s, fmt.Sprintf("password=%s", dsn.Password))
	}

	if dsn.DBName != "" {
		s = append(s, fmt.Sprintf("dbname=%s", dsn.DBName))
	}

	if dsn.SSLMode {
		s = append(s, fmt.Sprintf("sslmode=%s", "require"))
	} else {
		s = append(s, fmt.Sprintf("sslmode=%s", "disable"))
	}

	if dsn.Timezone != "" {
		s = append(s, fmt.Sprintf("TimeZone=%s", dsn.Timezone))
	}

	return strings.Join(s, " ")
}
