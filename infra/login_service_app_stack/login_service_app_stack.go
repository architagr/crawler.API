package loginserviceappstack

import (
	"infra/common"
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type LoginAPILambdaStackProps struct {
	config.CommonProps
}

func NewLoginAPILambdaStack(scope constructs.Construct, id string, props *LoginAPILambdaStackProps) (awscdk.Stack, apigateway.LambdaRestApi) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)
	loginRestApi := buildLambda(stack, scope, props)
	return stack, loginRestApi
}
func buildLambda(stack awscdk.Stack, scope constructs.Construct, props *LoginAPILambdaStackProps) apigateway.LambdaRestApi {
	userPool, userPoolClient := BuildUserPool(stack, *jsii.String("testUserPool"), &UserPoolLambdaStackProps{
		CommonProps: props.CommonProps,
	})

	env := make(map[string]*string)
	env["DbConnectionString"] = jsii.String(props.LoginAPIDB.GetConnectionString())
	env["DatabaseName"] = jsii.String(props.LoginAPIDB.GetDbName())
	env["LoginCollectionName"] = jsii.String(props.LoginAPIDB.GetCollectionName())
	env["GIN_MODE"] = jsii.String("release")
	env["UserPoolId"] = userPool.UserPoolId()
	env["ClientId"] = userPoolClient.UserPoolClientId()
	// env["ClientSecret"] = userPoolClient.UserPoolClientSecret().ToString()

	loginFunction := common.BuildLambda(&common.LambdaConstructProps{
		CommonProps: props.CommonProps,
		Id:          "login-lambda",
		Handler:     "LoginAPI",
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

	common.AddResource("healthCheck", baseApi, []string{common.GET_METHOD}, integration)
	common.AddResource("login", baseApi, []string{common.POST_METHOD}, integration)
	common.AddResource("register", baseApi, []string{common.POST_METHOD}, integration)

	return loginApi
}
