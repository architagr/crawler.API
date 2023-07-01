package middleware

import (
	customerrors "LoginAPI/custom_errors"
	"LoginAPI/enums"
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
				return
			default:
				switch err.Err.(type) {
				case *customerrors.AWSError:
					c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
						ErrorCode: enums.ERROR_CODE_AWS,
						Message:   "There is some issue in connecting to aws, please contact the Admin team.",
					})
					return
				case *customerrors.AuthError:
					handleAuthError(c, err.Err.(*customerrors.AuthError))
					return
				}
				return
			}
		}

		c.Next()
	}
}
func handleAuthError(c *gin.Context, err *customerrors.AuthError) {
	switch err.GetCode() {
	case enums.ERROR_CODE_AUTH_INVALID_CREDENTIALS:
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode: err.GetCode(),
			Message:   "Invalid credentials",
		})
	case enums.ERROR_CODE_AUTH_PASSWORD_EXPIRED:
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode: err.GetCode(),
			Message:   "Password expired, please update the password",
		})
	case enums.ERROR_CODE_AUTH_CREATE_USER:
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode: err.GetCode(),
			Message:   "error while creating user",
		})
	case enums.ERROR_CODE_AUTH_UPDATE_PASSWORD:
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode: err.GetCode(),
			Message:   "error while updating password",
		})
	case enums.ERROR_CODE_AUTH_USERNAME_EXISTS:
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode: err.GetCode(),
			Message:   "username already eixst.",
		})
	case enums.ERROR_CODE_AUTH_INVALID_PASSWORD:
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode: err.GetCode(),
			Message:   "invalid password",
		})
	}
}
func InitGinMiddelware() IMiddleware[gin.HandlerFunc] {
	return &ginMiddeleware{}
}
