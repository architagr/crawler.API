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

func InitGinRouters(employerController controller.IEmployerController, companyController controller.ICompanyController, logObj logger.ILogger) IRouter {
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
	createCompanyRoutes(routerGroup, companyController, logObj)
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

	jobGroup.POST("/get", func(ginContext *gin.Context) {
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

func createCompanyRoutes(group *gin.RouterGroup, companyController controller.ICompanyController, logObj logger.ILogger) {
	companyGroup := group.Group("company")

	companyGroup.POST("/save", func(ginContext *gin.Context) {
		var companyDetail models.Company
		if err := ginContext.ShouldBind(&companyDetail); err != nil {
			logObj.Printf("wrong request body %+v", err)
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypeBind,
			})
			return
		}
		userName, exists := ginContext.Get("user_name")
		if !exists {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
			ginContext.Abort()
			return
		}
		companyDetail.EmployerId = userName.(string)
		userData, err := companyController.SaveCompany(&companyDetail)
		if err != nil {
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypePrivate,
			})
			return
		}
		ginContext.JSON(http.StatusOK, userData)
	})

	companyGroup.PUT("/image/:id", func(ginContext *gin.Context) {

		var updateAvatarRequest models.CompanyImage
		if err := ginContext.ShouldBind(&updateAvatarRequest); err != nil {
			logObj.Printf("wrong request body %+v", err)
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypeBind,
			})
			return
		}
		if err := ginContext.ShouldBindUri(&updateAvatarRequest); err != nil {
			logObj.Printf("wrong request param %+v", err)
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypeBind,
			})
			return
		}
		err := companyController.SaveCompanyImage(&updateAvatarRequest)
		if err != nil {
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypePrivate,
			})
			return
		}
		ginContext.JSON(http.StatusOK, map[string]string{
			"message": "image uploaded",
		})
	})

	companyGroup.POST("/get", func(ginContext *gin.Context) {
		var filter models.SearchFilter
		if err := ginContext.ShouldBind(&filter); err != nil {
			logObj.Printf("wrong request body %+v", err)
			ginContext.Errors = append(ginContext.Errors, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypeBind,
			})
			return
		}
		userName, exists := ginContext.Get("user_name")
		if !exists {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
			ginContext.Abort()
			return
		}
		filter.EmployerId = userName.(string)
		userData, err := companyController.GetCompanies(&filter)
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
	ginEngine.Use(ginMiddleware.GetTokenInfo())
}
func getMiddlewares() middlewarePkg.IMiddleware[gin.HandlerFunc] {
	return middlewarePkg.InitGinMiddelware()
}
