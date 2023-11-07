package service

import (
	customerrors "UserAPI/custom_errors"
	"UserAPI/filters"
	"UserAPI/logger"
	"UserAPI/models"
	"UserAPI/repository"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserProfileService interface {
	SaveUserProfile(user *models.UserDetail) (*models.UserDetail, error)
	GetUserProfile(username string) (*models.UserDetail, error)
	SaveImagetoAWS(_file multipart.File, userId, fileName, mimetype string, size int64) error
	GetUserImageURL(filename string) (string, error)
}

type userProfileService struct {
	repo      repository.IUserRepository
	logObj    logger.ILogger
	s3Service IS3Service
}

var userServiceObj IUserProfileService

func InitUserService(repoObj repository.IUserRepository, s3Service IS3Service, logObj logger.ILogger) IUserProfileService {
	if userServiceObj == nil {
		userServiceObj = &userProfileService{
			repo:      repoObj,
			s3Service: s3Service,
			logObj:    logObj,
		}
	}
	return userServiceObj
}

func (s *userProfileService) SaveUserProfile(user *models.UserDetail) (*models.UserDetail, error) {

	var filter filters.IFilter = nil
	_filter := bson.M{}
	if user != nil {
		if user.Email != "" {
			filter = filters.InitEmailFilter(filter, filters.AND, filters.EQUAL, user.Email)
		}
		if user.Phone != "" {
			filter = filters.InitPhoneFilter(filter, filters.OR, filters.EQUAL, user.Phone)
		}
		if user.Id != "" {
			//check if email id or mobile already exists in other profiles
			objectId, err := primitive.ObjectIDFromHex(user.Id)
			if err != nil {
				s.logObj.Printf("error while converting id to hex %s, error: %s\n", user.Id, err.Error())
				return nil, &customerrors.GetUserException{}
			}
			filter = filters.InitIdFilter(filter, filters.AND, filters.NOT_EQUAL, objectId)
		}
	}
	if filter != nil {
		_filter = filter.Build()
	}
	userList, err := s.repo.Get(_filter, 1, 1)
	if err != nil {
		s.logObj.Printf("error while getting user from db using filter %+v, error: %s\n", _filter, err.Error())
		return nil, &customerrors.GetUserException{}
	} else if len(userList) > 0 {
		s.logObj.Printf("user with same email, phone exist")
		return nil, &customerrors.UsernameExistException{}
	}

	if user.Id == "" {
		userId, err := s.repo.AddSingle(*user)
		if err != nil {
			s.logObj.Printf("Error while adding user: %+v, error: %s\n", user, err.Error())
			return nil, &customerrors.AddUserException{}
		}
		user.Id = userId
		return user, nil
	} else {
		existingUser := new(models.UserDetail)
		existingUser, _ = s.repo.GetById(user.Id)
		if existingUser == nil {
			return nil, &customerrors.UserNotFoundException{}
		}

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
		}}

		err = s.repo.UpdateSingle(update, user.Id)
		if err != nil {
			s.logObj.Printf("Error while updating user: %+v, error: %s\n", user, err.Error())
			return nil, &customerrors.UpdateUserException{}
		}
		return user, nil
	}
}

func (s *userProfileService) GetUserProfile(username string) (*models.UserDetail, error) {
	_filter := filters.InitUsernameFilter(nil, filters.AND, filters.EQUAL, username)
	users, err := s.repo.Get(_filter.Build(), 1, 0)
	if err != nil {
		return nil, &customerrors.GetUserException{}
	}
	if len(users) == 0 {
		return nil, &customerrors.UserNotFoundException{}
	}
	user := &users[0]
	if user.ImagePath != "" {
		user.ImagePath, err = s.GetUserImageURL(user.ImagePath)
		if err != nil {
			user.ImagePath = ""
		}
	}
	return user, nil
}

func (s *userProfileService) SaveImagetoAWS(_file multipart.File, userId, fileName, mimetype string, size int64) error {

	update := bson.M{"$set": bson.M{
		"imagepath": fileName,
	}}

	err := s.repo.UpdateSingle(update, userId)
	if err != nil {
		s.logObj.Printf("Error while updating userId: %+v, error: %s\n", userId, err.Error())
		return &customerrors.UpdateUserException{}
	}
	err = s.s3Service.Put(_file, fileName, mimetype, size)
	if err != nil {
		return err
	}
	return nil
}

func (s *userProfileService) GetUserImageURL(filename string) (string, error) {
	return s.s3Service.GetPreSignerUrl(filename)

}
