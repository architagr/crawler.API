package aws

import (
	"UserAPI/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	s3Interface "github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var s3Svc s3Interface.S3API = nil

func GetS3Session(config config.IConfig) s3Interface.S3API {
	if s3Svc == nil {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config: aws.Config{
				Region: aws.String(config.GetAwsRegion()),
			},
		}))
		s3Svc = s3.New(sess)
	}
	return s3Svc
}
