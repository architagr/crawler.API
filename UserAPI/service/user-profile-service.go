package service

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"jobcrawler.api/config"
	"jobcrawler.api/models"
	"jobcrawler.api/repository/collection"
	"jobcrawler.api/repository/connection"
)

type IUserProfileService interface {
	SaveUserProfile(user *models.UserDetail) (string, error)
	GetUserProfile(username string) (*models.UserDetail, error)
}

type UserProfileService struct {
	collectionObj collection.ICollection[models.UserDetail]
}

func UserProfileServiceObj(collectionName string) (IUserProfileService, error) {
	env := config.GetConfig()
	connObj, err := connection.InitConnection(env.GetDatabaseConnectionString(), 10)
	if err != nil {
		return nil, err
	}
	doc, err := collection.InitCollection[models.UserDetail](connObj, env.GetDatabaseName(), collectionName)

	if err != nil {
		return nil, err
	}
	return &UserProfileService{
		collectionObj: doc,
	}, nil
}

func (s *UserProfileService) SaveUserProfile(user *models.UserDetail) (string, error) {
	existingUser := new(models.UserDetail)
	var err error
	if user.Id != "" {
		*existingUser, err = s.collectionObj.GetById(user.Id)
	}
	//check if email id or mobile already exists in other profiles
	objectId, err := primitive.ObjectIDFromHex(user.Id)
	_filter := bson.M{}
	if user == nil {
		_filter = bson.M{}
	} else {
		if user.Email != "" {
			_filter = bson.M{"email": user.Email}
		}
		if user.Phone != "" {
			_filter = bson.M{
				"$or": []bson.M{
					_filter,
					bson.M{"phone": user.Phone},
				},
			}
		}
		if user.Id != "" {
			_filter = bson.M{
				"$and": []bson.M{
					_filter,
					bson.M{"_id": bson.M{"$ne": objectId}},
				},
			}
		}
	}
	userList, err := s.collectionObj.Get(_filter, 1, 1)
	if err != nil {
		return "", err
	} else if len(userList) > 0 {
		return "", errors.New("Email Id or User name already available")
	}

	if (existingUser == &models.UserDetail{} || user.Id == "") {
		userId, err := s.collectionObj.AddSingle(*user)
		if err != nil {
			return "", err
		}
		return userId.(primitive.ObjectID).Hex(), nil
	} else {
		update := bson.M{"$set": bson.M{
			"name":          user.Name,
			"username":      user.UserName,
			"jobtitle":      user.JobTitle,
			"phone":         user.Phone,
			"email":         user.Email,
			"currentsalary": user.CurrentSalary,
			"experience":    user.Experience,
			"gender":        user.Gender,
			"age":           user.Age,
			"jobcategory":   user.JobCategory,
			"language":      user.Language,
			"description":   user.Description,
			"imagepath":     user.ImagePath,
		}}

		err := s.collectionObj.UpdateSingle(update, user.Id)
		if err != nil {
			return "", err
		}
		return user.Id, nil
	}
}

func (s *UserProfileService) GetUserProfile(username string) (*models.UserDetail, error) {
	if username != "" {
		_filter := bson.M{"username": username}
		user, err := s.collectionObj.Get(_filter, 1, 1)
		if err != nil {
			return nil, err
		}
		return &user[0], nil
	}
	return nil, errors.New("user id is not valid")
}
