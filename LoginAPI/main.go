package main

import (
	"LoginAPI/aws"
	"LoginAPI/config"
	"LoginAPI/controller"
	"LoginAPI/logger"
	"LoginAPI/repository"
	"LoginAPI/routers"
	"LoginAPI/service"
	"flag"

	cognitoInterface "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

var (
	port = flag.Int("port", 8082, "this value is used when we run the service on local")
)
var envVariables config.IConfig
var authRepoObj repository.IAuthRepository
var authServiceObj service.IAuthService
var authControllerObj controller.IAuthController
var logObj logger.ILogger
var cognito cognitoInterface.CognitoIdentityProviderAPI

func main() {
	flag.Parse()
	cognito = aws.GetCognitoService()
	initLogger()
	initConfig()
	initRepository()
	initServices()
	initControllers()
	routers.InitGinRouters(authControllerObj, logObj).StartApp(*port)
}

func initLogger() {
	logObj = logger.InitConsoleLogger()
}
func initConfig() {
	envVariables = config.GetConfig()
}
func initRepository() {
	mongodbConnection, err := repository.InitConnection(envVariables.GetDatabaseConnectionString(), 10)
	if err != nil {
		panic(err)
	}

	authRepoObj, err = repository.InitAuthRepo(mongodbConnection, envVariables.GetDatabaseName(), envVariables.GetCollectionName())
	if err != nil {
		panic(err)
	}
}

func initServices() {
	authServiceObj = service.InitAuthService(authRepoObj, envVariables, cognito, logObj)
}

func initControllers() {
	authControllerObj = controller.InitAuthController(authServiceObj)
}
