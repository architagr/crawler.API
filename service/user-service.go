package service

import (
	"jobcrawler.api/config"
	"jobcrawler.api/models"
	"jobcrawler.api/repository/collection"
	"jobcrawler.api/repository/connection"
)

type IUserService interface {
	Login(user *models.LoginDetails) (interface{}, error)
}

type UserService struct {
	collectionObj collection.ICollection[models.LoginDetails]
}

func UserServiceObj() (IUserService, error) {
	env := config.GetConfig()
	connObj, err := connection.InitConnection(env.GetDatabaseConnectionString(), 10)
	if err != nil {
		return nil, err
	}
	doc, err := collection.InitCollection[models.LoginDetails](connObj, env.GetDatabaseName(), "userDetails")
	if err != nil {
		return nil, err
	}
	return &UserService{
		collectionObj: doc,
	}, nil
}

func (s *UserService) Login(user *models.LoginDetails) (interface{}, error) {
	result, err := s.collectionObj.AddSingle(*user)
	if err != nil {
		return "", err
	}
	return result, nil
}
