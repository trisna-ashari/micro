package configurator

// Config represent config keys.
type Config struct {
	AppName        string
	AppVersion     string
	AppEnvironment string
	AppTimezone    string

	AppHTTPPort string
	AppGRPCPort string

	DBConfig
	DBReplicaTotal int
	DBReplicasConfig
	DBTestConfig

	GoogleCloudServiceConfig
	AmazonWebServiceConfig

	StorageTestConfig
	StorageConfig
	MinioConfig
	GCSConfig
	S3Config
	LocalFileConfig

	DataDogConfig

	DebugMode              bool
	TestMode               bool
	MaxIdleCons            int
	MaxOpenCons            int
	EnableCachePrepareStmt bool
}

// DBConfig represent db config keys.
type DBConfig struct {
	DBDriver                    string
	DBHost                      string
	DBPort                      string
	DBUser                      string
	DBName                      string
	DBPassword                  string
	DBTimeZone                  string
	DBLog                       bool
	DisableForeignKeyConstraint bool
}

// GoogleCloudServiceConfig represent google cloud config keys.
type GoogleCloudServiceConfig struct {
	GoogleApplicationCredential string
}

// AmazonWebServiceConfig represent google cloud config keys.
type AmazonWebServiceConfig struct {
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
}

// MinioConfig represent minio config keys.
type MinioConfig struct {
	MinioEndpoint   string
	MinioHost       string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
	MinioPathPrefix string
}

// GCSConfig represent google cloud storage config keys.
type GCSConfig struct {
	GCSBucketName string
	GCSPathPrefix string
}

// S3Config represent amazon S3 config keys.
type S3Config struct {
	S3BucketName string
	S3PathPrefix string
}

// LocalFileConfig represent local file storage config keys.
type LocalFileConfig struct {
	Path string
}

// StorageConfig represent storage driver config keys.
// There only three drivers: gcs, s3, and minio.
type StorageConfig struct {
	Driver string
}

// StorageTestConfig represent storage driver config keys.
// There only three drivers: gcs, s3, and minio.
type StorageTestConfig struct {
	Driver string
}

// DataDogConfig represent datadog tracer config keys.
type DataDogConfig struct {
	EnableTracer bool
	AgentHost    string
	AgentPort    string
	APIKey       string
}

// DBReplicaConfig represent db config keys.
type DBReplicaConfig struct {
	DBDriver                    string
	DBHost                      string
	DBPort                      string
	DBUser                      string
	DBName                      string
	DBPassword                  string
	DBTimeZone                  string
	DBLog                       bool
	DisableForeignKeyConstraint bool
}

// DBReplicasConfig holds multiple DBReplicaConfig
type DBReplicasConfig []*DBReplicaConfig

// DBTestConfig represent db test config keys.
type DBTestConfig struct {
	DBDriver                    string
	DBHost                      string
	DBPort                      string
	DBUser                      string
	DBName                      string
	DBPassword                  string
	DBTimeZone                  string
	DBLog                       bool
	DisableForeignKeyConstraint bool
}

// New is a function uses to initialize Config.
// Calling this function without Option will return default config.
func New(opts ...Option) *Config {
	config := &Config{
		AppHTTPPort: GetEnv("APP_HTTP_PORT", "6969"),
		AppGRPCPort: GetEnv("APP_GRPC_PORT", "4949"),
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}
