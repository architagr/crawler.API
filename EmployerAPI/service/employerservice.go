package service

import (
	"EmployerAPI/logger"
	"EmployerAPI/models"
	"EmployerAPI/repository"
)

type IEmployerService interface {
	SaveJob(jobDetail *models.JobDetail) (*models.JobDetail, error)
}

type employerService struct {
	repo      repository.IJobRepository
	logObj    logger.ILogger
	s3Service IS3Service
}

var employerServiceObj IEmployerService

func InitJobService(repoObj repository.IJobRepository, s3Service IS3Service, logObj logger.ILogger) IEmployerService {
	if employerServiceObj == nil {
		employerServiceObj = &employerService{
			repo:      repoObj,
			s3Service: s3Service,
			logObj:    logObj,
		}
	}
	return employerServiceObj
}

func (s *employerService) SaveJob(jobDetail *models.JobDetail) (*models.JobDetail, error) {
	jobId, err := s.repo.AddSingle(*jobDetail)
	if err != nil {
		return nil, err
	}
	jobDetail.Id = jobId
	return jobDetail, nil
}
