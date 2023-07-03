package routers

import (
	"LoginAPI/config"
	"LoginAPI/controller"
	"LoginAPI/logger"
	middlewarePkg "LoginAPI/middleware"
	"LoginAPI/models"
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

func InitGinRouters(authController controller.IAuthController, logObj logger.ILogger) IRouter {
	ginMiddleware = getMiddlewares()
	ginEngine := gin.Default()
	registerInitialCommonMiddleware(ginEngine, logObj)
	routerGroup := getInitialRouteGroup(ginEngine)

	routerGroup.GET("/healthCheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Login api is working",
		})
	})
	routerGroup.POST("/login", func(ginContext *gin.Context) {
		var loginRequest models.LoginDetails
		if err := ginContext.ShouldBind(&loginRequest); err != nil {
			logObj.Printf("wrong request body %+v", err)
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypeBind,
			})
			return
		}
		token, err := authController.AuthenticateUser(&loginRequest)
		if err != nil {

			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypePrivate,
			})
			return
		}
		ginContext.JSON(http.StatusOK, token)
	})
	routerGroup.POST("/register", func(ginContext *gin.Context) {
		var loginRequest models.LoginDetails
		if err := ginContext.ShouldBind(&loginRequest); err != nil {
			logObj.Printf("wrong request body %+v", err)
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypeBind,
			})
			return
		}
		token, err := authController.CreateUser(&loginRequest)
		if err != nil {
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypePrivate,
			})
			return
		}
		ginContext.JSON(http.StatusOK, token)
	})
	return &ginRouter{
		ginEngine: ginEngine,
		env:       config.GetConfig(),
		logObj:    logObj,
	}
}

func getInitialRouteGroup(ginEngine *gin.Engine) *gin.RouterGroup {
	return ginEngine.Group("/auth")
}
func registerInitialCommonMiddleware(ginEngine *gin.Engine, logObj logger.ILogger) {
	ginEngine.Use(ginMiddleware.GetErrorHandler(logObj))
	ginEngine.Use(ginMiddleware.GetCorsMiddelware())
}
func getMiddlewares() middlewarePkg.IMiddleware[gin.HandlerFunc] {
	return middlewarePkg.InitGinMiddelware()
}
