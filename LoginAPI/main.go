package main

import (
	"LoginAPI/config"
	"LoginAPI/controller"
	"LoginAPI/logger"
	"LoginAPI/repository"
	"LoginAPI/routers"
	"LoginAPI/service"
)

var envVariables config.IConfig
var authRepoObj repository.IAuthRepository
var authServiceObj service.IAuthService
var authControllerObj controller.IAuthController
var logObj logger.ILogger

func main() {
	initLogger()
	initConfig()
	initRepository()
	intitServices()
	initControllers()
	routers.InitGinRouters(authControllerObj, logObj).StartApp()
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

func intitServices() {
	authServiceObj = service.InitAuthService(authRepoObj)
}

func initControllers() {
	authControllerObj = controller.InitAuthController(authServiceObj)
}
