package authappstack

import (
	"fmt"
	"infra/config"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	awscognito "github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	constructs "github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AuthStackProps struct {
	config.CommonProps
	UserPool awscognito.IUserPool
}

func NewAuthStack(scope constructs.Construct, id string, props *AuthStackProps) (awscdk.Stack, apigateway.IAuthorizer) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	return stack, buildCognitoAuthorizer(stack, props)
}

func buildCognitoAuthorizer(stack awscdk.Stack, props *AuthStackProps) apigateway.IAuthorizer {

	return apigateway.NewCognitoUserPoolsAuthorizer(stack, jsii.String(fmt.Sprintf("%s-cognito-authorizer", props.StackNamePrefix)), &apigateway.CognitoUserPoolsAuthorizerProps{
		CognitoUserPools: &[]awscognito.IUserPool{
			props.UserPool,
		},
		IdentitySource: jsii.String("method.request.header.Authorizer"),
		AuthorizerName: jsii.String(fmt.Sprintf("%s-cognito-authorizer", props.StackNamePrefix)),
	})
}
