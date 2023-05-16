package service

import (
	"go.mongodb.org/mongo-driver/bson"
	"jobcrawler.api/config"
	"jobcrawler.api/models"
	"jobcrawler.api/repository/collection"
	"jobcrawler.api/repository/connection"
)

type IJobService interface {
	GetJobs(filter *models.JobFilter, pageSize, pageNumber int16) (*models.GetJobResponse, error)
	GetJobDetail(id string) (*models.JobDetails, error)
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

	_filter := bson.M{}
	if filter == nil {
		_filter = bson.M{}
	} else {
		if filter.Location != "" {
			_filter = bson.M{
				"$and": []bson.M{
					bson.M{"location": filter.Location},
					bson.M{
						"$or": []bson.M{
							bson.M{"title": filter.Keywords},
							bson.M{"companyname": filter.Keywords},
						},
					},
				},
			}
		} else {
			_filter = bson.M{
				"$or": []bson.M{
					bson.M{"title": filter.Keywords},
					bson.M{"companyname": filter.Keywords},
				},
			}
		}
	}
	data, err := svc.collectionObj.Get(_filter, int64(pageSize), int64(pageNumber))
	if err != nil {
		return nil, err
	}
	return &models.GetJobResponse{
		Jobs:       data,
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}, nil
}

func (svc *JobService) GetJobDetail(id string) (*models.JobDetails, error) {
	data, err := svc.collectionObj.GetById(id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
