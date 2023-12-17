package config

import "os"

var ConfigContainerKey = "config"

type IConfig interface {
	GetDatabaseConnectionString() string
	GetDatabaseName() string
	GetCollectionName() string
	GetAwsRegion() string
	IsLambda() bool
}
type Config struct {
	databaseConnectionString string
	databaseName             string
	collectionName           string
	isLambda                 bool
	awsRegion                string
}

var env IConfig

const (
	databaseConnectionStringKey = "JobAPIDbConnectionString"
	databaseNameKey             = "JobAPIDatabaseName"
	collectionNameKey           = "JobCollectionName"
	awsRegionKey                = "AWS_REGION"
	isLambdaEnvKey              = "LAMBDA_TASK_ROOT"
)

func InitConfig() {
	_, ok := os.LookupEnv(isLambdaEnvKey)
	env = &Config{
		databaseConnectionString: os.Getenv(databaseConnectionStringKey),
		databaseName:             os.Getenv(databaseNameKey),
		collectionName:           os.Getenv(collectionNameKey),
		awsRegion:                os.Getenv(awsRegionKey),
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

func (e *Config) GetCollectionName() string {
	return e.collectionName
}

func (e *Config) IsLambda() bool {
	return e.isLambda
}
func (e *Config) GetAwsRegion() string {
	return e.awsRegion
}
