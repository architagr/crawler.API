package aws

import (
	"LoginAPI/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	cognitoInterface "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

var svc cognitoInterface.CognitoIdentityProviderAPI = nil

func GetCognitoService(config config.IConfig) cognitoInterface.CognitoIdentityProviderAPI {
	if svc == nil {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config: aws.Config{
				Region: aws.String(config.GetAwsRegion()),
			},
		}))
		svc = cognitoidentityprovider.New(sess)
	}
	return svc
}

const (
	COGNITO_CHALLANGE_NAME_NEW_PASSWORD_REQUIRED = "NEW_PASSWORD_REQUIRED"
)
const (
	MESSAGE_ACTION_SUPRESS = "SUPPRESS"
)

const (
	USER_STATUS_FORCE_CHANGE_PASSWORD = "FORCE_CHANGE_PASSWORD"
)

const (
	AUTH_FLOW_ADMIN_USER_PASSWORD_AUTH = "ADMIN_USER_PASSWORD_AUTH"
	AUTH_FLOW_USER_PASSWORD_AUTH       = "USER_PASSWORD_AUTH"
)
