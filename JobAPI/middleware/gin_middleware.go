package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	region              = "ap-south-1"           // Update with your AWS region
	userPoolID          = "ap-south-1_PsMRSTJ4p" // Update with your Cognito user pool ID
	accessTokenEndpoint = "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
)

type PublicKeyResponse struct {
	Keys []PublicKey `json:"keys"`
}

type PublicKey struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
}

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

func (*ginMiddeleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		tokenString := accessToken[len("Bearer "):]

		// Verify the access token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Provide the secret key used for signing the tokens
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("Invalid access token")
			}

			publicKey, err := getPublicKey(kid)
			if err != nil {
				return nil, err
			}

			return publicKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			c.Abort()
			return
		}

		// Access token is valid, continue with the request
		c.Next()
	}
}

func getPublicKey(kid string) (*rsa.PublicKey, error) {
	endpoint := fmt.Sprintf(accessTokenEndpoint, region, userPoolID)

	// Fetch the public key from Cognito
	response, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve public key: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to retrieve public key. Status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	var publicKeyResponse PublicKeyResponse
	err = json.Unmarshal(body, &publicKeyResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse public key response: %v", err)
	}

	for _, key := range publicKeyResponse.Keys {
		if key.Kid == kid {
			modulusBytes, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, fmt.Errorf("Failed to decode modulus: %v", err)
			}

			exponentBytes, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, fmt.Errorf("Failed to decode exponent: %v", err)
			}

			modulus := &big.Int{}
			modulus.SetBytes(modulusBytes)

			exponent := &big.Int{}
			exponent.SetBytes(exponentBytes)

			pubKey := &rsa.PublicKey{
				N: modulus,
				E: int(exponent.Int64()),
			}

			return pubKey, nil
		}
	}

	return nil, fmt.Errorf("Public key not found for kid: %s", kid)
}

func InitGinMiddelware() IMiddleware[gin.HandlerFunc] {
	return &ginMiddeleware{}
}
