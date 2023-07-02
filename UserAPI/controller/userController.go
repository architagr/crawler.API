package controller

import (
	customerrors "UserAPI/custom_errors"
	"UserAPI/logger"
	"UserAPI/models"
	"UserAPI/service"
	"path/filepath"
)

type IUserController interface {
	SaveUserProfile(userDetails *models.UserDetail) (*models.UserDetail, error)
	SaveUserImage(updateAvatarRequest *models.UpdateAvatarRequest) error
	GetUserProfile(email string) (*models.UserDetail, error)
}

type userController struct {
	service service.IUserProfileService
	logObj  logger.ILogger
}

var userControllerObj IUserController

func InitUserController(serviceObj service.IUserProfileService, logObj logger.ILogger) IUserController {
	if userControllerObj == nil {
		userControllerObj = &userController{
			service: serviceObj,
			logObj:  logObj,
		}
	}
	return userControllerObj
}

func (ctlr *userController) SaveUserProfile(userDetails *models.UserDetail) (*models.UserDetail, error) {
	return ctlr.service.SaveUserProfile(userDetails)
}

func (ctrl *userController) SaveUserImage(updateAvatarRequest *models.UpdateAvatarRequest) error {

	src, err := updateAvatarRequest.Image.Open()
	if err != nil {
		ctrl.logObj.Printf("error in opening file for %+v, error: %s", updateAvatarRequest, err)
		return &customerrors.FileOpenException{}
	}

	defer src.Close()

	userId := updateAvatarRequest.Id

	// Create a unique filename
	ext := filepath.Ext(updateAvatarRequest.Image.Filename)
	filename := userId + ext

	mimetype := updateAvatarRequest.Image.Header.Get("Content-Type")
	ctrl.service.SaveImagetoAWS(src, userId, filename, mimetype, updateAvatarRequest.Image.Size)

	return nil
}

func (ctrl *userController) GetUserProfile(email string) (*models.UserDetail, error) {

	result, err := ctrl.service.GetUserProfile(email)
	if err != nil {
		return nil, err
	}

	return result, nil
}
