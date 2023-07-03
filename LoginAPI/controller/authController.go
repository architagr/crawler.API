package controller

import (
	"LoginAPI/models"
	"LoginAPI/service"
)

type IAuthController interface {
	CreateUser(newUser *models.LoginDetails) (*models.Token, error)
	AuthenticateUser(loginRequest *models.LoginDetails) (*models.Token, error)
}

type authController struct {
	service service.IAuthService
}

var authControllerObj IAuthController

func InitAuthController(serviceObj service.IAuthService) IAuthController {
	if authControllerObj == nil {
		authControllerObj = &authController{
			service: serviceObj,
		}
	}
	return authControllerObj
}

func (ctlr *authController) CreateUser(newUser *models.LoginDetails) (*models.Token, error) {
	return ctlr.service.CreateCognitoUser(newUser)
}

func (ctlr *authController) AuthenticateUser(loginRequest *models.LoginDetails) (*models.Token, error) {
	return ctlr.service.LoginUser(loginRequest)
}
