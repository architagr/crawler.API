package main

import (
	"EmployerAPI/aws"
	"EmployerAPI/config"
	"EmployerAPI/controller"
	"EmployerAPI/logger"
	"EmployerAPI/repository"
	"EmployerAPI/routers"
	"EmployerAPI/service"
	"flag"

	s3Interface "github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var (
	port = flag.Int("port", 8083, "this value is used when we run the service on local")
)

var envVariables config.IConfig
var employerRepoObj repository.IJobRepository
var employerService service.IEmployerService
var employerController controller.IEmployerController
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
	routers.InitGinRouters(employerController, logObj).StartApp(*port)

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

	employerRepoObj, err = repository.InitEmployerRepository(mongodbConnection, envVariables.GetDatabaseName(), "employerProfile")
	if err != nil {
		panic(err)
	}
}

func intitServices() {
	s3Service = service.InitS3Service(s3Session, envVariables.GetAvatarImageBucketName(), logObj)

	employerService = service.InitJobService(employerRepoObj, s3Service, logObj)
}

func initControllers() {
	employerController = controller.InitEmployerController(employerService, logObj)
}
