package controller

import (
	"EmployerAPI/logger"
	"EmployerAPI/models"
	"EmployerAPI/service"
)

type IEmployerController interface {
	SaveJob(jobDetail *models.JobDetail) (*models.JobDetail, error)
}

type employerController struct {
	service service.IEmployerService
	logObj  logger.ILogger
}

var employerControllerObj IEmployerController

func InitEmployerController(serviceObj service.IEmployerService, logObj logger.ILogger) IEmployerController {
	if employerControllerObj == nil {
		employerControllerObj = &employerController{
			service: serviceObj,
			logObj:  logObj,
		}
	}
	return employerControllerObj
}

func (ctlr *employerController) SaveJob(jobDetail *models.JobDetail) (*models.JobDetail, error) {
	return ctlr.service.SaveJob(jobDetail)
}
