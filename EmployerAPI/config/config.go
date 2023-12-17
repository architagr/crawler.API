package config

import "os"

var ConfigContainerKey = "config"

type IConfig interface {
	GetDatabaseConnectionString() string
	GetDatabaseName() string
	GetJobCollectionName() string
	GetCompanyCollectionName() string
	GetAvatarImageBucketName() string
	GetAwsRegion() string
	IsLambda() bool
}
type Config struct {
	databaseConnectionString string
	databaseName             string
	jobCollectionName        string
	companyCollectionName    string
	isLambda                 bool
	avatarImageBucketName    string
	awsRegion                string
}

var env IConfig

const (
	databaseConnectionStringKey = "DbConnectionString"
	databaseNameKey             = "DatabaseName"
	jobCollectionKey            = "JobCollectionName"
	companyCollectionKey        = "CompanyCollectionName"
	avatarImageBucketNameKey    = "AvatarImageBucketName"
	awsRegionKey                = "AWS_REGION"
	isLambdaEnvKey              = "LAMBDA_TASK_ROOT"
)

func InitConfig() {
	_, ok := os.LookupEnv(isLambdaEnvKey)

	env = &Config{
		databaseConnectionString: os.Getenv(databaseConnectionStringKey),
		databaseName:             os.Getenv(databaseNameKey),
		jobCollectionName:        os.Getenv(jobCollectionKey),
		companyCollectionName:    os.Getenv(companyCollectionKey),
		avatarImageBucketName:    os.Getenv(avatarImageBucketNameKey),
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

func (e *Config) GetJobCollectionName() string {
	return e.jobCollectionName
}

func (e *Config) GetCompanyCollectionName() string {
	return e.companyCollectionName
}

func (e *Config) IsLambda() bool {
	return e.isLambda
}

func (e *Config) GetAvatarImageBucketName() string {
	return e.avatarImageBucketName
}
func (e *Config) GetAwsRegion() string {
	return e.awsRegion
}
