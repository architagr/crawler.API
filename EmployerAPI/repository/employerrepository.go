package repository

import (
	"EmployerAPI/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

type IEmployerRepository interface {
	AddSingle(job models.JobDetail) (string, error)
	Get(filter bson.M, pageSize int64, startPage int64) ([]models.JobDetail, error)
}

type employerRepository struct {
	collectionObj ICollection[models.JobDetail]
}

var employerRepoObject IJobRepository

func InitEmployerRepository(conn IConnection, databaseName, collection string) (IJobRepository, error) {
	if jobRepoObject != nil {
		return jobRepoObject, nil
	}
	doc, err := InitCollection[models.JobDetail](conn, databaseName, collection)
	if err != nil {
		return nil, err
	}
	jobRepoObject = &jobRepository{
		collectionObj: doc,
	}
	return jobRepoObject, nil
}

func (repo *employerRepository) GetById(Id string) (*models.JobDetail, error) {
	if Id != "" {
		data, err := repo.collectionObj.GetById(Id)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
	return nil, nil
}

func (repo *employerRepository) Get(filter bson.M, pageSize int64, startPage int64) ([]models.JobDetail, error) {
	data, err := repo.collectionObj.Get(filter, pageSize, startPage)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *employerRepository) AddSingle(job models.JobDetail) (string, error) {
	data, err := repo.collectionObj.AddSingle(job)
	if err != nil {
		return "", err
	}
	return data.(primitive.ObjectID).Hex(), nil
}

func (repo *employerRepository) UpdateSingle(data primitive.M, Id string) error {
	err := repo.collectionObj.UpdateSingle(data, Id)
	return err
}
