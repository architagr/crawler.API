package middleware

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

const secretKey = "mysecretkey"

func AuthMiddleware(c *fiber.Ctx) error {
	// Get the JWT token from the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.ErrUnauthorized
	}
	tokenString := authHeader[len("Bearer "):]
	// Verify the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return fiber.ErrUnauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fiber.ErrUnauthorized
	}
	c.Locals("username", claims["username"])
	return c.Next()
}
