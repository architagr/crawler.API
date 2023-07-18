package repository

import (
	"UserAPI/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

type IUserRepository interface {
	GetById(Id string) (*models.UserDetail, error)
	Get(filter bson.M, pageSize int64, startPage int64) ([]models.UserDetail, error)
	AddSingle(user models.UserDetail) (string, error)
	UpdateSingle(data primitive.M, Id string) error
}

type userRepository struct {
	collectionObj ICollection[models.UserDetail]
}

var userRepoObject IUserRepository

func InitUserRepository(conn IConnection, databaseName, collection string) (IUserRepository, error) {
	if userRepoObject != nil {
		return userRepoObject, nil
	}
	doc, err := InitCollection[models.UserDetail](conn, databaseName, collection)
	if err != nil {
		return nil, err
	}
	userRepoObject = &userRepository{
		collectionObj: doc,
	}
	return userRepoObject, nil
}

func (repo *userRepository) GetById(Id string) (*models.UserDetail, error) {
	if Id != "" {
		data, err := repo.collectionObj.GetById(Id)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
	return nil, nil
}

func (repo *userRepository) Get(filter bson.M, pageSize int64, startPage int64) ([]models.UserDetail, error) {
	data, err := repo.collectionObj.Get(filter, pageSize, startPage)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *userRepository) AddSingle(user models.UserDetail) (string, error) {
	data, err := repo.collectionObj.AddSingle(user)
	if err != nil {
		return "", err
	}
	return data.(primitive.ObjectID).Hex(), nil
}

func (repo *userRepository) UpdateSingle(data primitive.M, Id string) error {
	err := repo.collectionObj.UpdateSingle(data, Id)
	return err
}
