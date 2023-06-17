package service

import (
	"UserAPI/models"
	"UserAPI/repository"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserProfileService interface {
	SaveUserProfile(user *models.UserDetail) (string, error)
	GetUserProfile(username string) (*models.UserDetail, error)
	SaveImagetoAWS(_file multipart.File, fileName string, size int64)
	GetUserImageURL(filename string) (string, error)
}

type userProfileService struct {
	repo repository.IUserRepository
}

var userServiceObj IUserProfileService

func InitUserService(repoObj repository.IUserRepository) IUserProfileService {
	if userServiceObj == nil {
		userServiceObj = &userProfileService{
			repo: repoObj,
		}
	}
	return userServiceObj
}

func (s *userProfileService) SaveUserProfile(user *models.UserDetail) (string, error) {
	existingUser := new(models.UserDetail)
	var err error
	existingUser, err = s.repo.GetById(user.Id)

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
	userList, err := s.repo.Get(_filter, 1, 1)
	if err != nil {
		return "", err
	} else if len(userList) > 0 {
		return "", errors.New("Email Id or User name already available")
	}

	if (existingUser == &models.UserDetail{} || user.Id == "") {
		userId, err := s.repo.AddSingle(*user)
		if err != nil {
			return "", err
		}
		return userId, nil
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

		err := s.repo.UpdateSingle(update, user.Id)
		if err != nil {
			return "", err
		}
		return user.Id, nil
	}
}

func (s *userProfileService) GetUserProfile(username string) (*models.UserDetail, error) {
	if username != "" {
		_filter := bson.M{"username": username}
		user, err := s.repo.Get(_filter, 1, 0)
		if err != nil {
			return nil, err
		}
		return &user[0], nil
	}
	return nil, errors.New("user id is not valid")
}

func (s *userProfileService) SaveImagetoAWS(_file multipart.File, fileName string, size int64) {
	// Specify your AWS region and S3 bucket name
	//region := "Asia Pacific (Mumbai) ap-south-1"
	bucketName := "jobcrawler.portalimages"

	// Create an AWS session
	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String(region)},
	// )
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an S3 client
	svc := s3.New(sess)

	// Create an S3 object with the specified bucket and key (filename)
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(fileName),
		Body:          _file,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String("image/jpeg"), // Specify the correct content type
	})
	if err != nil {
		fmt.Println("Failed to upload image:", err)
		return
	}

	fmt.Println("Image uploaded successfully!")
}

func (s *userProfileService) GetUserImageURL(filename string) (string, error) {
	bucketName := "jobcrawler.portalimages"

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an S3 client
	svc := s3.New(sess)

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	}

	req, _ := svc.GetObjectRequest(params)

	url, err := req.Presign(time.Duration(2 * time.Hour)) // Set link expiration time
	if err != nil {
		return "", err
	}
	return url, nil
}
