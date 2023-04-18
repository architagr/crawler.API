package service

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"jobcrawler.api/models"
)

type IAuthService interface {
	GetJWT(user *models.LoginDetails) (string, error)
}

type AuthService struct {
	secretkey string
}

func AuthServiceObj() (IAuthService, error) {
	key := os.Getenv("JWTSecreatKey")
	return &AuthService{
		secretkey: key,
	}, nil
}

func (s *AuthService) GetJWT(user *models.LoginDetails) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.UserName
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString([]byte(s.secretkey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
