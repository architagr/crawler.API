package main

import (
	"UserAPI/aws"
	"UserAPI/config"
	"UserAPI/controller"
	"UserAPI/logger"
	"UserAPI/repository"
	"UserAPI/routers"
	"UserAPI/service"
	"flag"

	s3Interface "github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var (
	port = flag.Int("port", 8081, "this value is used when we run the service on local")
)

var envVariables config.IConfig
var userRepoObj repository.IUserRepository
var userService service.IUserProfileService
var userController controller.IUserController
var s3Service service.IS3Service
var s3Session s3Interface.S3API
var logObj logger.ILogger

func main() {
	flag.Parse()
	initConfig()
	initS3Session()
	initLogger()
	initRepository()
	intitServices()
	initControllers()
	routers.InitGinRouters(userController, logObj).StartApp(*port)

}
func initS3Session() {
	s3Session = aws.GetS3Session()
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

	userRepoObj, err = repository.InitUserRepository(mongodbConnection, envVariables.GetDatabaseName(), "userProfile")
	if err != nil {
		panic(err)
	}
}

func intitServices() {
	s3Service = service.InitS3Service(s3Session, envVariables.GetAvatarImageBucketName(), logObj)

	userService = service.InitUserService(userRepoObj, s3Service, logObj)
}

func initControllers() {
	userController = controller.InitUserController(userService, logObj)
}
