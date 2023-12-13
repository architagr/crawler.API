package controller

import (
	"EmployerAPI/logger"
	"EmployerAPI/models"
	"EmployerAPI/service"
)

type IEmployerController interface {
	SaveJob(jobDetail *models.JobDetail) (*models.JobDetail, error)
	GetJobs(filter *models.JobFilter) (*models.GetJobResponse, error)
}

type employerController struct {
	service        service.IEmployerService
	companyService service.ICompanyService
	logObj         logger.ILogger
}

var employerControllerObj IEmployerController

func InitEmployerController(serviceObj service.IEmployerService, companyServiceObj service.ICompanyService,
	logObj logger.ILogger) IEmployerController {
	if employerControllerObj == nil {
		employerControllerObj = &employerController{
			service:        serviceObj,
			companyService: companyServiceObj,
			logObj:         logObj,
		}
	}
	return employerControllerObj
}

func (ctlr *employerController) SaveJob(jobDetail *models.JobDetail) (*models.JobDetail, error) {
	return ctlr.service.SaveJob(jobDetail)
}

func (ctlr *employerController) GetJobs(filter *models.JobFilter) (*models.GetJobResponse, error) {
	return ctlr.service.GetJobs(filter)
}
