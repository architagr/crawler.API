package controller

import (
	customerrors "EmployerAPI/custom_errors"
	"EmployerAPI/logger"
	"EmployerAPI/models"
	"EmployerAPI/service"
	"path/filepath"
)

type ICompanyController interface {
	SaveCompany(companyDetail *models.Company) (*models.Company, error)
	GetCompanies(filter *models.SearchFilter) (*models.GetCompanyResponse, error)
	SaveCompanyImage(updateAvatarRequest *models.CompanyImage) error
}

type companyController struct {
	service service.ICompanyService
	logObj  logger.ILogger
}

var companyControllerObj ICompanyController

func InitCompanyController(serviceObj service.ICompanyService, logObj logger.ILogger) ICompanyController {
	if companyControllerObj == nil {
		companyControllerObj = &companyController{
			service: serviceObj,
			logObj:  logObj,
		}
	}
	return companyControllerObj
}

func (ctlr *companyController) SaveCompany(companyDetail *models.Company) (*models.Company, error) {
	return ctlr.service.SaveCompany(companyDetail)
}

func (ctlr *companyController) GetCompanies(filter *models.SearchFilter) (*models.GetCompanyResponse, error) {
	return ctlr.service.GetCompanies(filter)
}

func (ctrl *companyController) SaveCompanyImage(updateAvatarRequest *models.CompanyImage) error {

	logoImage, err := updateAvatarRequest.Image.Open()
	if err != nil {
		ctrl.logObj.Printf("error in opening logo file for %+v, error: %s", updateAvatarRequest, err)
		return &customerrors.FileOpenException{}
	}

	defer logoImage.Close()

	companyId := updateAvatarRequest.Id

	// Create a unique filename
	ext := filepath.Ext(updateAvatarRequest.Image.Filename)
	filename := companyId + "_" + updateAvatarRequest.Type + ext

	mimetype := updateAvatarRequest.Image.Header.Get("Content-Type")
	err = ctrl.service.SaveImagetoAWS(logoImage, updateAvatarRequest.Type+"FileName", companyId, filename, mimetype, updateAvatarRequest.Image.Size)
	if err != nil {
		ctrl.logObj.Printf("error while saving the image, error: %s", err)
		return &customerrors.SaveImageException{}
	}

	return nil
}
