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
				handleAuthError(c, err.Err)
				return
			}
		}

		c.Next()
	}
}
func handleAuthError(c *gin.Context, err error) {
	errorCode := enums.ERROR_CODE_REQUEST_INTERNAL_ERROR
	status := http.StatusInternalServerError
	switch err.(type) {
	case *customerrors.AWSError:
		status = http.StatusInternalServerError
		errorCode = enums.ERROR_CODE_AWS
	case *customerrors.InvalidCredentialException:
		status = http.StatusUnauthorized
		errorCode = enums.ERROR_CODE_AUTH_INVALID_CREDENTIALS
	case *customerrors.InvalidPasswordException:
		status = http.StatusUnauthorized
		errorCode = enums.ERROR_CODE_AUTH_INVALID_PASSWORD
	case *customerrors.PasswordExpireException:
		status = http.StatusUnauthorized
		errorCode = enums.ERROR_CODE_AUTH_PASSWORD_EXPIRED
	case *customerrors.UpdatePasswordException:
		status = http.StatusBadRequest
		errorCode = enums.ERROR_CODE_AUTH_UPDATE_PASSWORD
	case *customerrors.CreateUserException:
		status = http.StatusBadRequest
		errorCode = enums.ERROR_CODE_AUTH_CREATE_USER
	case *customerrors.UsernameExistsException:
		status = http.StatusBadRequest
		errorCode = enums.ERROR_CODE_AUTH_USERNAME_EXISTS
	}
	c.AbortWithStatusJSON(status, models.ErrorResponse{
		ErrorCode: errorCode,
		Message:   err.Error(),
	})
}
func InitGinMiddelware() IMiddleware[gin.HandlerFunc] {
	return &ginMiddeleware{}
}
