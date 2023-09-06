package config

import "os"

var ConfigContainerKey = "config"

type IConfig interface {
	GetDatabaseConnectionString() string
	GetDatabaseName() string
	GetEmployerCollectionName() string
	GetJobCollectionName() string
	GetAvatarImageBucketName() string
	IsLambda() bool
}
type Config struct {
	databaseConnectionString string
	databaseName             string
	employerCollectionName   string
	jobCollectionName        string
	isLambda                 bool
	avatarImageBucketName    string
}

var env IConfig

const (
	databaseConnectionStringKey = "DbConnectionString"
	databaseNameKey             = "DatabaseName"
	employerCollectionKey       = "EmployerCollection"
	jobCollectionKey            = "JobCollection"
	avatarImageBucketNameKey    = "AvatarImageBucketName"
	isLambdaEnvKey              = "LAMBDA_TASK_ROOT"
)

func InitConfig() {
	_, ok := os.LookupEnv(isLambdaEnvKey)

	env = &Config{
		databaseConnectionString: os.Getenv(databaseConnectionStringKey),
		databaseName:             os.Getenv(databaseNameKey),
		employerCollectionName:   os.Getenv(employerCollectionKey),
		jobCollectionName:        os.Getenv(jobCollectionKey),
		avatarImageBucketName:    os.Getenv(avatarImageBucketNameKey),
		isLambda:                 ok,
	}
}

func GetConfig() IConfig {
	if env == nil {
		InitConfig()
	}
	return env
}

func (e *Config) GetDatabaseConnectionString() string {
	return e.databaseConnectionString
}

func (e *Config) GetDatabaseName() string {
	return e.databaseName
}

func (e *Config) GetEmployerCollectionName() string {
	return e.employerCollectionName
}

func (e *Config) GetJobCollectionName() string {
	return e.jobCollectionName
}

func (e *Config) IsLambda() bool {
	return e.isLambda
}

func (e *Config) GetAvatarImageBucketName() string {
	return e.avatarImageBucketName
}
