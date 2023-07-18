package config

import "os"

var ConfigContainerKey = "config"

type IConfig interface {
	GetDatabaseConnectionString() string
	GetDatabaseName() string
	GetCollectionName() string
	GetUserPoolId() string
	GetClientId() string
	IsLambda() bool
}
type Config struct {
	databaseConnectionString string
	databaseName             string
	collectionName           string
	userPoolId               string
	clientId                 string
	isLambda                 bool
}

var env IConfig

const (
	databaseConnectionStringKey = "DbConnectionString"
	databaseNameKey             = "DatabaseName"
	collectionNameKey           = "LoginCollectionName"
	userPoolIdKey               = "UserPoolId"
	clientid                    = "ClientId"
	isLambdaEnvKey              = "LAMBDA_TASK_ROOT"
)

func InitConfig() {
	_, ok := os.LookupEnv(isLambdaEnvKey)

	env = &Config{
		databaseConnectionString: os.Getenv(databaseConnectionStringKey),
		databaseName:             os.Getenv(databaseNameKey),
		collectionName:           os.Getenv(collectionNameKey),
		userPoolId:               os.Getenv(userPoolIdKey),
		clientId:                 os.Getenv(clientid),
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

func (e *Config) GetUserPoolId() string {
	return e.userPoolId
}
func (e *Config) GetClientId() string {
	return e.clientId
}
