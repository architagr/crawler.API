package loginserviceappstack

import (
	"infra/common"
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	awscognito "github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type LoginAPILambdaStackProps struct {
	config.CommonProps
}

func NewLoginAPILambdaStack(scope constructs.Construct, id string, props *LoginAPILambdaStackProps) (awscdk.Stack, apigateway.LambdaRestApi, awscognito.IUserPool) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)
	userPool, userPoolClient := BuildUserPool(stack, &UserPoolLambdaStackProps{
		CommonProps: props.CommonProps,
	})

	loginRestApi := buildLambda(stack, props, userPool, userPoolClient)
	return stack, loginRestApi, userPool
}

func buildLambda(stack awscdk.Stack, props *LoginAPILambdaStackProps, userPool awscognito.IUserPool,
	userPoolClient awscognito.IUserPoolClient) apigateway.LambdaRestApi {

	env := make(map[string]*string)
	env["LoginAPIDbConnectionString"] = props.LoginAPIDB.GetConnectionString()
	env["LoginAPIDatabaseName"] = props.LoginAPIDB.GetDbName()
	env["LoginCollectionName"] = jsii.String(props.LoginAPIDB.GetCollectionName())
	env["UserPoolId"] = userPool.UserPoolId()
	env["ClientId"] = userPoolClient.UserPoolClientId()

	loginFunction := common.BuildLambda(&common.LambdaConstructProps{
		CommonProps: props.CommonProps,
		Id:          "login-lambda",
		Handler:     "LoginAPI", // TODO: get this from makefile
		Service:     "LoginAPI",
		Name:        "login-lambda-fn",
		Description: "This function helps in all API related to login",
		Env:         env,
		Stack:       stack,
	})
	// userPool.Grant(loginFunction)
	loginFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: &[]*string{
			jsii.String("cognito-idp:*"),
		},
		Resources: &[]*string{
			userPool.UserPoolArn(),
		},
	}))
	restApiProps := common.RestApiProps{
		CommonProps: props.CommonProps,
		Stack:       stack,
		Id:          "LoginApi",
		Handler:     loginFunction,
		Name:        "LoginRestApi",
	}

	loginApi := common.BuildRestApi(&restApiProps)

	integration := common.BuildIntegration(&restApiProps)

	baseApi := loginApi.Root()

	common.AddResource("healthCheck", baseApi, []string{common.GET_METHOD}, integration, nil)
	common.AddResource("login", baseApi, []string{common.POST_METHOD}, integration, nil)
	common.AddResource("register", baseApi, []string{common.POST_METHOD}, integration, nil)

	return loginApi
}
