package service

import (
	"jobcrawler.api/config"
	"jobcrawler.api/models"
	"jobcrawler.api/repository/collection"
	"jobcrawler.api/repository/connection"
)

type IJobService interface {
	GetJobs(filter *models.JobFilter, pageSize, pageNumber int16) (*models.GetJobResponse, error)
}

type JobService struct {
	collectionObj collection.ICollection[models.JobDetails]
}

func GetJobServiceObj() (IJobService, error) {
	env := config.GetConfig()
	connObj, err := connection.InitConnection(env.GetDatabaseConnectionString(), 10)
	if err != nil {
		return nil, err
	}
	doc, err := collection.InitCollection[models.JobDetails](connObj, env.GetDatabaseName(), env.GetCollectionName())
	if err != nil {
		return nil, err
	}
	return &JobService{
		collectionObj: doc,
	}, nil
}

func (svc *JobService) GetJobs(filter *models.JobFilter, pageSize, pageNumber int16) (*models.GetJobResponse, error) {
	data, err := svc.collectionObj.Get(filter, int64(pageSize), int64(pageNumber))
	if err != nil {
		return nil, err
	}
	return &models.GetJobResponse{
		Jobs:       data,
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}, nil
}
