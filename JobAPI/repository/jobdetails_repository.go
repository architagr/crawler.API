package repository

import (
	"JobAPI/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type IJobDetailsRepository interface {
	GetJob(filter bson.M, pageSize, pageNumber int16) ([]models.JobDetails, error)
	GetJobDetail(id string) (*models.JobDetails, error)
}
type jobDetailsRepository struct {
	collectionObj ICollection[models.JobDetails]
}

var jobDetailsRepoObj IJobDetailsRepository

func InitJobDetailsRepo(conn IConnection, databaseName, collection string) (IJobDetailsRepository, error) {

	if jobDetailsRepoObj != nil {
		return jobDetailsRepoObj, nil
	}
	fmt.Println("collection name " + collection)
	fmt.Println("database name " + databaseName)
	doc, err := InitCollection[models.JobDetails](conn, databaseName, collection)
	if err != nil {
		return nil, err
	}
	jobDetailsRepoObj = &jobDetailsRepository{
		collectionObj: doc,
	}
	return jobDetailsRepoObj, nil
}
func (repo *jobDetailsRepository) GetJob(filter bson.M, pageSize, pageNumber int16) ([]models.JobDetails, error) {

	data, err := repo.collectionObj.Get(filter, int64(pageSize), int64(pageNumber))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *jobDetailsRepository) GetJobDetail(id string) (*models.JobDetails, error) {
	data, err := repo.collectionObj.GetById(id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
