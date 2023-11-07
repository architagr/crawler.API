package service

import (
	"EmployerAPI/filters"
	"EmployerAPI/logger"
	"EmployerAPI/models"
	"EmployerAPI/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if jobDetail.Id != "" {
		update := bson.M{"$set": bson.M{
			"title":       jobDetail.Title,
			"description": jobDetail.Description,
			"competency":  jobDetail.Competency,
			"category":    jobDetail.Category,
			"payby":       jobDetail.PayBy,
			"minamount":   jobDetail.MinAmount,
			"maxamount":   jobDetail.MaxAmount,
			"rate":        jobDetail.Rate,
			"modeofwork":  jobDetail.ModeOfWork,
			"jobType":     jobDetail.JobType,
			"gender":      jobDetail.Gender,
			"experience":  jobDetail.Experience,
			"deadline":    jobDetail.Deadline,
			"country":     jobDetail.Country,
			"city":        jobDetail.City,
			"companyid":   jobDetail.CompanyId,
		}}
		err := s.repo.UpdateSingle(update, jobDetail.Id)
		if err != nil {
			s.logObj.Printf("Error while updating Employer job, error: %s\n", err.Error())
			return nil, err
		}
	} else {
		jobId, err := s.repo.AddSingle(*jobDetail)
		if err != nil {
			return nil, err
		}
		jobDetail.Id = jobId
	}
	return jobDetail, nil
}

func (s *employerService) GetJobs(filterData *models.JobFilter) (*models.GetJobResponse, error) {
	var filter filters.IFilter = nil
	_filter := bson.M{}
	if filterData != nil {
		if filterData.EmployerId != "" {
			filter = filters.InitEmployerIdFilter(filter, filters.AND, filters.EQUAL, filterData.EmployerId)
		}
		if filterData.JobId != "" {
			objectId, err := primitive.ObjectIDFromHex(filterData.JobId)
			if err != nil {
				s.logObj.Printf("error while converting id to hex %s, error: %s\n", filterData.JobId, err.Error())
				return nil, err
			}
			filter = filters.InitIdFilter(filter, filters.AND, filters.EQUAL, objectId)
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
