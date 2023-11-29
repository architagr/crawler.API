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
var companyRepoObj repository.ICompanyRepository
var jobRepoObj repository.IJobRepository
var employerJobService service.IEmployerService
var companyService service.ICompanyService
var employerController controller.IEmployerController
var companyController controller.ICompanyController
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
	routers.InitGinRouters(employerController, companyController, logObj).StartApp(*port)

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

	companyRepoObj, err = repository.InitCompanyRepository(mongodbConnection, envVariables.GetDatabaseName(), envVariables.GetCompanyCollectionName())
	if err != nil {
		panic(err)
	}

	jobRepoObj, err = repository.InitJobRepository(mongodbConnection, envVariables.GetDatabaseName(), envVariables.GetJobCollectionName())
	if err != nil {
		panic(err)
	}
}

func intitServices() {
	s3Service = service.InitS3Service(s3Session, envVariables.GetAvatarImageBucketName(), logObj)

	employerJobService = service.InitJobService(jobRepoObj, s3Service, logObj)

	companyService = service.InitCompanyService(companyRepoObj, s3Service, logObj)
}

func initControllers() {
	employerController = controller.InitEmployerController(employerJobService, companyService, logObj)
	companyController = controller.InitCompanyController(companyService, logObj)
}
