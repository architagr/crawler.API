package service

import (
	"encoding/json"
	"fmt"

	"jobcrawler.api/config"
	"jobcrawler.api/models"
	"jobcrawler.api/repository/collection"
	"jobcrawler.api/repository/connection"
)

type IUserService interface {
	Login(user *models.LoginDetails) (bool, error)
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

func (s *UserService) Login(user *models.LoginDetails) (bool, error) {
	existingUser, err := s.collectionObj.GetUserByUserName(user.UserName)
	if err != nil {
		return false, err
	}
	jsonString, _ := json.Marshal(existingUser)
	fmt.Print("user data: " + string(jsonString))
	if user.LoginType == "OTP" {
		//if user not available then insert in db
		if (existingUser == models.LoginDetails{}) {
			_, err := s.collectionObj.AddSingle(*user)
			if err != nil {
				return false, err
			}
		}
		//Validate OTP by fetching from Redis
		return true, nil
	} else if user.LoginType == "creadentials" {
		if (existingUser == models.LoginDetails{}) {
			return false, nil
		} else {
			if user.Password == existingUser.Password {
				return true, nil
			} else {
				return false, nil
			}
		}

	}
	return false, nil
}
