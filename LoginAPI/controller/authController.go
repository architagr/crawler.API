package controller

import (
	customerrors "LoginAPI/custom_errors"
	"LoginAPI/models"
	"LoginAPI/service"
)

type IAuthController interface {
	CreateUser(newUser *models.LoginDetails) (*models.Token, *customerrors.AuthError)
	AuthenticateUser(loginRequest *models.LoginDetails) (*models.Token, *customerrors.AuthError)
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

func (ctlr *authController) CreateUser(newUser *models.LoginDetails) (*models.Token, *customerrors.AuthError) {
	return ctlr.service.CreateCognitoUser(newUser)
}

func (ctlr *authController) AuthenticateUser(loginRequest *models.LoginDetails) (*models.Token, *customerrors.AuthError) {
	return ctlr.service.LoginUser(loginRequest)
}
