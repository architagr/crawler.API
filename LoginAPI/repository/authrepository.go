package repository

import (
	"LoginAPI/models"
)

type IAuthRepository interface {
}
type authRepository struct {
	collectionObj ICollection[models.LoginDetails]
}

var authRepoObj IAuthRepository

func InitAuthRepo(conn IConnection, databaseName, collection string) (IAuthRepository, error) {

	if authRepoObj != nil {
		return authRepoObj, nil
	}
	doc, err := InitCollection[models.LoginDetails](conn, databaseName, collection)
	if err != nil {
		return nil, err
	}
	return &authRepository{
		collectionObj: doc,
	}, nil
}
