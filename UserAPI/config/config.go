package config

import "os"

var ConfigContainerKey = "config"

type IConfig interface {
	GetDatabaseConnectionString() string
	GetDatabaseName() string
	GetCollectionName() string
}
type Config struct {
	databaseConnectionString string
	databaseName             string
	collectionName           string
}

var env IConfig

const (
	databaseConnectionStringKey = "DbConnectionString"
	databaseNameKey             = "DatabaseName"
	collectionNameKey           = "CollectionName"
)

func InitConfig() {
	env = &Config{
		databaseConnectionString: os.Getenv(databaseConnectionStringKey),
		databaseName:             os.Getenv(databaseNameKey),
		collectionName:           os.Getenv(collectionNameKey),
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
