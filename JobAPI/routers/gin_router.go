package routers

import (
	"JobAPI/config"
	"JobAPI/controller"
	"JobAPI/logger"
	middlewarePkg "JobAPI/middleware"
	"JobAPI/models"
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginMiddleware middlewarePkg.IMiddleware[gin.HandlerFunc]

type ginRouter struct {
	ginEngine *gin.Engine
	env       config.IConfig
	logObj    logger.ILogger
}

var ginLambda *ginadapter.GinLambda

// Handler is the function that executes for every Request passed into the Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}
func (router *ginRouter) StartApp() {
	if router.env.IsLambda() {
		ginLambda = ginadapter.New(router.ginEngine)
		lambda.Start(Handler)
	} else {
		router.ginEngine.Run(":8080")
	}
}

func InitGinRouters(jobController controller.IJobController, logObj logger.ILogger) IRouter {
	ginMiddleware = getMiddlewares()
	ginEngine := gin.Default()
	registerInitialCommonMiddleware(ginEngine)
	routerGroup := getInitialRouteGroup(ginEngine)

	routerGroup.GET("/healthCheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Jobs api is working",
		})
	})

	routerGroup.GET("/", func(ginContext *gin.Context) {
		var request models.JobFilter
		if err := ginContext.ShouldBind(&request); err != nil {
			logObj.Printf("wrong query paramater %+v", err)
		}
		response, err := jobController.GetJobs(&request)
		if err != nil {
			logObj.Printf("error in getting jobs %+v", err)
			ginContext.JSON(http.StatusBadRequest, err)
			ginContext.Next()
		}
		ginContext.JSON(http.StatusOK, response)

	})

	return &ginRouter{
		ginEngine: ginEngine,
		env:       config.GetConfig(),
		logObj:    logObj,
	}
}

func getInitialRouteGroup(ginEngine *gin.Engine) *gin.RouterGroup {
	return ginEngine.Group("/jobs")
}
func registerInitialCommonMiddleware(ginEngine *gin.Engine) {
	ginEngine.Use(ginMiddleware.GetCorsMiddelware())
}
func getMiddlewares() middlewarePkg.IMiddleware[gin.HandlerFunc] {
	return middlewarePkg.InitGinMiddelware()
}
