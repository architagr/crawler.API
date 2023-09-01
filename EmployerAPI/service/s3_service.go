package service

import (
	customerrors "EmployerAPI/custom_errors"
	"EmployerAPI/logger"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	s3Interface "github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type IS3Service interface {
	Put(_file multipart.File, fileName, mimetype string, size int64) error
	GetPreSignerUrl(filename string) (string, error)
}

type S3Service struct {
	s3Svc      s3Interface.S3API
	bucketName string
	logObj     logger.ILogger
}

func (s3Svc *S3Service) Put(_file multipart.File, fileName, mimetype string, size int64) error {
	// Create an S3 object with the specified bucket and key (filename)
	_, err := s3Svc.s3Svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s3Svc.bucketName),
		Key:           aws.String(fileName),
		Body:          _file,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(mimetype), // Specify the correct content type
	})
	if err != nil {
		s3Svc.logObj.Printf("Failed to upload image, error: %s\n", err.Error())
		return &customerrors.UploadFileException{}
	}
	return nil
}

func (s3Svc *S3Service) GetPreSignerUrl(filename string) (string, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(s3Svc.bucketName),
		Key:    aws.String(filename),
	}

	req, _ := s3Svc.s3Svc.GetObjectRequest(params)

	url, err := req.Presign(time.Duration(2 * time.Hour)) // Set link expiration time
	if err != nil {
		s3Svc.logObj.Printf("error in generating presigned url for %s, error: %s", filename, err.Error())
		return "", &customerrors.PreSignedUrlException{}
	}
	return url, nil
}

func InitS3Service(s3Session s3Interface.S3API, bucketName string, logObj logger.ILogger) IS3Service {
	return &S3Service{
		bucketName: bucketName,
		s3Svc:      s3Session,
		logObj:     logObj,
	}
}
