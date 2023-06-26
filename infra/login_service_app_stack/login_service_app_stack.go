package loginserviceappstack

import (
	"infra/common"
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

var (
	GET_METHOD  = "GET"
	POST_METHOD = "POST"
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

	env := make(map[string]*string)
	env["DbConnectionString"] = jsii.String(props.LoginAPIDB.GetConnectionString())
	env["DatabaseName"] = jsii.String(props.LoginAPIDB.GetDbName())
	env["LoginCollectionName"] = jsii.String(props.LoginAPIDB.GetCollectionName())
	env["GIN_MODE"] = jsii.String("release")

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

	common.AddResource("healthCheck", baseApi, []string{GET_METHOD}, integration)
	common.AddResource("login", baseApi, []string{POST_METHOD}, integration)
	common.AddResource("register", baseApi, []string{POST_METHOD}, integration)
	return loginApi
}
