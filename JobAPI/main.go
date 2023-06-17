package main

import (
	"JobAPI/config"
	"JobAPI/controller"
	"JobAPI/logger"
	"JobAPI/repository"
	"JobAPI/routers"
	"JobAPI/service"
)

var envVariables config.IConfig
var jobDetailsRepoObj repository.IJobDetailsRepository
var jobDetailsService service.IJobService
var jobControllerObj controller.IJobController
var logObj logger.ILogger

func main() {
	initLogger()
	initConfig()
	initRepository()
	intitServices()
	initControllers()
	routers.InitGinRouters(jobControllerObj, logObj).StartApp()
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

	jobDetailsRepoObj, err = repository.InitJobDetailsRepo(mongodbConnection, envVariables.GetDatabaseName(), envVariables.GetCollectionName())
	if err != nil {
		panic(err)
	}
}

func intitServices() {
	jobDetailsService = service.InitJobService(jobDetailsRepoObj)
}

func initControllers() {
	jobControllerObj = controller.InitJobController(jobDetailsService, logObj)
}

// func GetJobs(c *fiber.Ctx) error {
// 	//var pageSize, pageNumber int64 = 10, 0 // todo: get this from querystring

// 	//get params from body
// 	filter := new(models.JobFilter)
// 	if err := c.BodyParser(filter); err != nil {
// 		return err
// 	}
// 	_filter, err := json.Marshal(filter)
// 	fmt.Println("filter: " + string(_filter))

// 	response, err := jobDetailsService.GetJobs(filter, filter.PageSize, filter.PageNumber)
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(response)
// }
