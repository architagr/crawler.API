package service

import (
	"EmployerAPI/filters"
	"EmployerAPI/logger"
	"EmployerAPI/models"
	"EmployerAPI/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type IEmployerService interface {
	SaveJob(jobDetail *models.JobDetail) (*models.JobDetail, error)
	GetJobs(filterData *models.JobFilter) (*models.GetJobResponse, error)
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

func (s *employerService) GetJobs(filterData *models.JobFilter) (*models.GetJobResponse, error) {
	var filter filters.IFilter = nil
	_filter := bson.M{}
	if filterData != nil {
		if filterData.EmployerId != "" {
			filter = filters.InitEmployerIdFilter(filter, filters.AND, filters.EQUAL, filterData.EmployerId)
		}
	}
	if filter != nil {
		_filter = filter.Build()
	}
	data, err := s.repo.Get(_filter, int64(filterData.PageSize), int64(filterData.PageNumber))
	if err != nil {
		return nil, err
	}
	return &models.GetJobResponse{
		Jobs:       data,
		PageSize:   filterData.PageSize,
		PageNumber: filterData.PageNumber,
	}, nil
}
