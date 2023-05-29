package main

import (
	"UserAPI/config"
	"UserAPI/controller"
	"UserAPI/repository"
	"UserAPI/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var envVariables config.IConfig
var userRepoObj repository.IUserRepository
var userService service.IUserProfileService
var userController controller.IUserController

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

	userRepoObj, err = repository.InitUserRepository(mongodbConnection, envVariables.GetDatabaseName(), "userProfile")
	if err != nil {
		panic(err)
	}
}

func intitServices() {
	userService = service.InitUserService(userRepoObj)
}

func initControllers() {
	userController = controller.InitUserController(userService)
}

func initRoutes() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Post("/saveUserProfile", userController.SaveUserProfile)
	app.Post("/saveuserImage/:id", userController.SaveUserImage)
	app.Get("/getUserProfile/:userId", userController.GetUserProfile)

	log.Fatal(app.Listen(":8081"))
}
