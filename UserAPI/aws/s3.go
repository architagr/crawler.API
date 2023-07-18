package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	s3Interface "github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var s3Svc s3Interface.S3API = nil

func GetS3Session() s3Interface.S3API {
	if s3Svc == nil {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		s3Svc = s3.New(sess)
	}
	return s3Svc
}
