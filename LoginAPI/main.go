package main

import (
	"LoginAPI/config"
	"LoginAPI/controller"
	"LoginAPI/repository"
	"LoginAPI/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var envVariables config.IConfig
var authRepoObj repository.IAuthRepository
var authServiceObj service.IAuthService
var authControllerObj controller.IAuthController

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

	authRepoObj, err = repository.InitAuthRepo(mongodbConnection, envVariables.GetDatabaseName(), "authDetails")
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

func initRoutes() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Post("/createuser", authControllerObj.CreateUser)
	app.Post("/loginuser", authControllerObj.AuthenticateUser)

	log.Fatal(app.Listen(":8082"))
}
