package service

import (
	"JobAPI/models"
	"JobAPI/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type IJobService interface {
	GetJobs(filter *models.JobFilter, pageSize, pageNumber int16) (*models.GetJobResponse, error)
	GetJobDetail(id string) (*models.JobDetails, error)
}

type jobService struct {
	repo repository.IJobDetailsRepository
}

var jobServiceObj IJobService

func InitJobService(repoObj repository.IJobDetailsRepository) IJobService {
	if jobServiceObj == nil {
		jobServiceObj = &jobService{
			repo: repoObj,
		}
	}
	return jobServiceObj
}

func (svc *jobService) GetJobs(filter *models.JobFilter, pageSize, pageNumber int16) (*models.GetJobResponse, error) {

	_filter := bson.M{}
	if filter != nil {
		if filter.Location != "" {
			//_filter = bson.M{"location": filter.Location}
			_filter = bson.M{
				"$and": []bson.M{
					{"location": filter.Location},
					{
						"$or": []bson.M{
							{"title": filter.Keywords},
							{"companyName": filter.Keywords},
						},
					},
				},
			}
		} else {
			_filter = bson.M{
				"$or": []bson.M{
					{"title": filter.Keywords},
					{"companyName": filter.Keywords},
				},
			}
		}
	}
	data, err := svc.repo.GetJob(&_filter, pageSize, pageNumber)
	if err != nil {
		return nil, err
	}
	return &models.GetJobResponse{
		Jobs:       data,
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}, nil
}

func (svc *jobService) GetJobDetail(id string) (*models.JobDetails, error) {
	data, err := svc.repo.GetJobDetail(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}
