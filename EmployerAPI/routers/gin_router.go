package routers

import (
	"EmployerAPI/config"
	"EmployerAPI/controller"
	"EmployerAPI/logger"
	middlewarePkg "EmployerAPI/middleware"
	"EmployerAPI/models"
	"context"
	"fmt"
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
func (router *ginRouter) StartApp(port int) {
	if router.env.IsLambda() {
		ginLambda = ginadapter.New(router.ginEngine)
		lambda.Start(Handler)
	} else {
		router.ginEngine.Run(fmt.Sprintf(":%d", port))
	}
}

func InitGinRouters(employerController controller.IEmployerController, logObj logger.ILogger) IRouter {
	ginMiddleware = getMiddlewares()
	ginEngine := gin.Default()
	registerInitialCommonMiddleware(ginEngine, logObj)
	routerGroup := getInitialRouteGroup(ginEngine)

	routerGroup.GET("/healthCheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Employer api is working",
		})
	})
	createEmployerRoutes(routerGroup, employerController, logObj)
	return &ginRouter{
		ginEngine: ginEngine,
		env:       config.GetConfig(),
		logObj:    logObj,
	}
}

func createEmployerRoutes(group *gin.RouterGroup, employerController controller.IEmployerController, logObj logger.ILogger) {
	jobGroup := group.Group("job")

	jobGroup.POST("/save", func(ginContext *gin.Context) {
		var jobDetail models.JobDetail
		if err := ginContext.ShouldBind(&jobDetail); err != nil {
			logObj.Printf("wrong request body %+v", err)
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypeBind,
			})
			return
		}
		userData, err := employerController.SaveJob(&jobDetail)
		if err != nil {
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypePrivate,
			})
			return
		}
		ginContext.JSON(http.StatusOK, userData)
	})

	jobGroup.POST("/getJobs", func(ginContext *gin.Context) {
		var jobFilter models.JobFilter
		if err := ginContext.ShouldBind(&jobFilter); err != nil {
			logObj.Printf("wrong request body %+v", err)
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypeBind,
			})
			return
		}
		userData, err := employerController.GetJobs(&jobFilter)
		if err != nil {
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypePrivate,
			})
			return
		}
		ginContext.JSON(http.StatusOK, userData)
	})

}
func getInitialRouteGroup(ginEngine *gin.Engine) *gin.RouterGroup {
	return ginEngine.Group("/employer")
}
func registerInitialCommonMiddleware(ginEngine *gin.Engine, logObj logger.ILogger) {
	ginEngine.Use(ginMiddleware.GetErrorHandler(logObj))
	ginEngine.Use(ginMiddleware.GetCorsMiddelware())
}
func getMiddlewares() middlewarePkg.IMiddleware[gin.HandlerFunc] {
	return middlewarePkg.InitGinMiddelware()
}
