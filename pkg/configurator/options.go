package configurator

// Option return config with option.
type Option func(config *Config)

// WithDBConfig is a function uses to set DBConfig to the Config.
func WithDBConfig() Option {
	return func(config *Config) {
		config.DBConfig = DBConfig{
			DBDriver:                    GetEnv("DB_DRIVER", "mysql"),
			DBHost:                      GetEnv("DB_HOST", "localhost"),
			DBPort:                      GetEnv("DB_PORT", "3306"),
			DBUser:                      GetEnv("DB_USER", "root"),
			DBName:                      GetEnv("DB_NAME", "go_user"),
			DBPassword:                  GetEnv("DB_PASSWORD", ""),
			DBTimeZone:                  GetEnv("APP_TIMEZONE", "Asia/Jakarta"),
			DBLog:                       GetEnvAsBool("ENABLE_LOGGER", true),
			DisableForeignKeyConstraint: GetEnvAsBool("DISABLE_FOREIGN_KEY_CONSTRAINT", false),
		}
	}
}

// WithMinioConfig is a function uses to set MinioConfig to the Config.
func WithMinioConfig() Option {
	return func(config *Config) {
		config.MinioConfig = MinioConfig{
			MinioEndpoint:   GetEnv("MINIO_ENDPOINT", ""),
			MinioHost:       GetEnv("MINIO_HOST", ""),
			MinioAccessKey:  GetEnv("MINIO_ACCESS_KEY", ""),
			MinioSecretKey:  GetEnv("MINIO_SECRET_KEY", ""),
			MinioBucketName: GetEnv("MINIO_BUCKET_NAME", ""),
			MinioPathPrefix: GetEnv("MINIO_PATH_PREFIX", ""),
		}
	}
}

// WithGoogleCloudServiceConfig is a function uses to set GoogleCloudServiceConfig to the Config.
func WithGoogleCloudServiceConfig() Option {
	return func(config *Config) {
		config.GoogleCloudServiceConfig = GoogleCloudServiceConfig{
			GoogleApplicationCredential: GetEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
		}
	}
}

// WithAmazonWebServiceConfig is a function uses to set AmazonWebServiceConfig to the Config.
func WithAmazonWebServiceConfig() Option {
	return func(config *Config) {
		config.AmazonWebServiceConfig = AmazonWebServiceConfig{
			AWSAccessKeyID:     GetEnv("AWS_ACCESS_KEY_ID", ""),
			AWSSecretAccessKey: GetEnv("AWS_SECRET_ACCESS_KEY", ""),
			AWSRegion:          GetEnv("AWS_REGION", ""),
		}
	}
}

// WithGCSConfig is a function uses to set GCSConfig to the Config.
func WithGCSConfig() Option {
	return func(config *Config) {
		config.GCSConfig = GCSConfig{
			GCSBucketName: GetEnv("GCS_BUCKET_NAME", ""),
			GCSPathPrefix: GetEnv("GCS_PATH_PREFIX", ""),
		}
	}
}

// WithS3Config is a function uses to set S3Config to the Config.
func WithS3Config() Option {
	return func(config *Config) {
		config.S3Config = S3Config{
			S3BucketName: GetEnv("S3_BUCKET_NAME", ""),
			S3PathPrefix: GetEnv("S3_PATH_PREFIX", ""),
		}
	}
}

// WithDatadogConfig is a function uses to set datadog tracer provider configuration.
func WithDatadogConfig() Option {
	return func(config *Config) {
		config.DataDogConfig = DataDogConfig{
			EnableTracer: GetEnvAsBool("DD_ENABLE_TRACER", false),
			AgentHost:    GetEnv("DD_AGENT_HOST", "127.0.0.1"),
			AgentPort:    GetEnv("DD_AGENT_PORT", "8126"),
			APIKey:       GetEnv("DD_APP_KEY", ""),
		}
	}
}

// WithStorageConfig is a function uses to set StorageConfig to the Config.
func WithStorageConfig() Option {
	return func(config *Config) {
		config.StorageConfig = StorageConfig{
			Driver: GetEnv("STORAGE_DRIVER", "minio"),
		}
	}
}
