package middleware

import (
	"EmployerAPI/enums"
	"EmployerAPI/logger"
	"EmployerAPI/models"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	}
	c.AbortWithStatusJSON(status, models.ErrorResponse{
		ErrorCode: errorCode,
		Message:   err.Error(),
	})
}
func (*ginMiddeleware) GetTokenInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed to parse token"})
			c.Abort()
			return
		}
		tokenString := strings.Split(authHeader, " ")[1]
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed to parse token"})
			c.Abort()
			return
		}
		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			c.Abort()
			return
		}
		userID, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in claims"})
			c.Abort()
			return
		}

		// Set the user ID in the context for use in downstream handlers
		c.Set("user_name", userID)

		c.Next()
	}
}
func InitGinMiddelware() IMiddleware[gin.HandlerFunc] {
	return &ginMiddeleware{}
}
