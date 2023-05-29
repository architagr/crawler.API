package main

import (
	"JobAPI/go/pkg/mod/github.com/gofiber/fiber/v2@v2.44.0/middleware/cors"
	"JobAPI/go/pkg/mod/github.com/gofiber/fiber@v1.14.6"
	"log"
)

var envVariables config.IConfig
var jobDetailsRepoObj repository.IJobDetailsRepository
var jobDetailsService service.IJobService
var jobControllerObj controller.IJobController

func main() {
	initConfig()
	initRepository()
	intitServices()
	initControllers()
	initRoutes()
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
	jobControllerObj = controller.InitJobController(jobDetailsService)
}

func initRoutes() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/test/:name", jobControllerObj.Test)

	app.Post("/getJobs", jobControllerObj.getJobs)

	log.Fatal(app.Listen(":8080"))
}
