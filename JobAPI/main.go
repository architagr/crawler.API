package main

import (
	"JobAPI/config"
	"JobAPI/controller"
	"JobAPI/repository"
	"JobAPI/service"
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var envVariables config.IConfig
var jobDetailsRepoObj repository.IJobDetailsRepository
var jobDetailsService service.IJobService
var jobControllerObj controller.IJobController
var ginLambda *ginadapter.GinLambda

func main() {
	// initConfig()
	// initRepository()
	// intitServices()
	// initControllers()
	initRoutes()
}

// Handler is the function that executes for every Request passed into the Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
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
	ginEngine := gin.Default()

	ginEngine.Use(CORSMiddleware())
	clientGroup := ginEngine.Group("/jobs")
	clientGroup.GET("/", func(c *gin.Context) {
		c.JSON(200, map[string]string{"data": "Welcome"})
	})
	if IsLambda() {
		ginLambda = ginadapter.New(ginEngine)
		lambda.Start(Handler)
	} else {
		ginEngine.Run(":8080")
	}

}
func IsLambda() bool {
	lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT")
	return lambdaTaskRoot != ""
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
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
