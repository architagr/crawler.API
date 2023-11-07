package repository

import (
	"EmployerAPI/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

type ICompanyRepository interface {
	AddSingle(company models.Company) (string, error)
	Get(filter bson.M, pageSize int64, startPage int64) ([]models.Company, error)
	UpdateSingle(data primitive.M, Id string) error
}

type companyRepository struct {
	collectionObj ICollection[models.Company]
}

var companyRepoObject ICompanyRepository

func InitCompanyRepository(conn IConnection, databaseName, collection string) (ICompanyRepository, error) {
	if jobRepoObject != nil {
		return companyRepoObject, nil
	}
	doc, err := InitCollection[models.Company](conn, databaseName, collection)
	if err != nil {
		return nil, err
	}
	companyRepoObject = &companyRepository{
		collectionObj: doc,
	}
	return companyRepoObject, nil
}

func (repo *companyRepository) GetById(Id string) (*models.Company, error) {
	if Id != "" {
		data, err := repo.collectionObj.GetById(Id)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
	return nil, nil
}

func (repo *companyRepository) Get(filter bson.M, pageSize int64, startPage int64) ([]models.Company, error) {
	data, err := repo.collectionObj.Get(filter, pageSize, startPage)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *companyRepository) AddSingle(job models.Company) (string, error) {
	data, err := repo.collectionObj.AddSingle(job)
	if err != nil {
		return "", err
	}
	return data.(primitive.ObjectID).Hex(), nil
}

func (repo *companyRepository) UpdateSingle(data primitive.M, Id string) error {
	err := repo.collectionObj.UpdateSingle(data, Id)
	return err
}
