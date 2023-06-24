package middleware

import (
	"LoginAPI/enums"
	"LoginAPI/errors"
	"LoginAPI/logger"
	"LoginAPI/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginMiddeleware struct {
}

func (*ginMiddeleware) GetCorsMiddelware() gin.HandlerFunc {
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

func (*ginMiddeleware) GetErrorHandler(logObj logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			logObj.Printf("error in api %s, error: %+v", c.Request.RequestURI, err)
			switch err.Type {
			case gin.ErrorTypeBind:
				c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
					ErrorCode: enums.ERROR_CODE_REQUEST_PARAM,
					Message:   "Request params are not valid",
				})
			default:
				switch err.Err.(type) {
				case *errors.AWSError:
					c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
						ErrorCode: enums.ERROR_CODE_AWS,
						Message:   "There is some issue in connecting to aws, please contact the Admin team.",
					})
				}

			}
		}

		c.Next()
	}
}
func InitGinMiddelware() IMiddleware[gin.HandlerFunc] {
	return &ginMiddeleware{}
}
