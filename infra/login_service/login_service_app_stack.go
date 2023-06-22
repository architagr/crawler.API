package loginservice

import (
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awss3assets "github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
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

	loginFunction := lambda.NewFunction(stack, jsii.String("login-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("LoginAPI"),
		Code:         lambda.Code_FromAsset(jsii.String("./../LoginAPI/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("login-lambda-fn"),
	})

	loginApi := apigateway.NewLambdaRestApi(stack, jsii.String("LoginApi"), &apigateway.LambdaRestApiProps{
		DeployOptions:               props.Stage,
		Handler:                     loginFunction,
		RestApiName:                 jsii.String("LoginRestApi"),
		Proxy:                       jsii.Bool(false),
		Deploy:                      jsii.Bool(true),
		DisableExecuteApiEndpoint:   jsii.Bool(false),
		EndpointTypes:               &[]apigateway.EndpointType{apigateway.EndpointType_EDGE},
		DefaultCorsPreflightOptions: config.GetCorsPreflightOptions(),
	})

	integration := apigateway.NewLambdaIntegration(loginFunction, &apigateway.LambdaIntegrationOptions{})

	baseApi := loginApi.Root()

	addResource("healthCheck", baseApi, []string{GET_METHOD}, integration)
	addResource("login", baseApi, []string{POST_METHOD}, integration)
	addResource("register", baseApi, []string{POST_METHOD}, integration)
	return loginApi
}

func addResource(path string, api apigateway.IResource, methods []string, integration apigateway.LambdaIntegration) apigateway.IResource {
	a := api.AddResource(jsii.String(path), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: config.GetCorsPreflightOptions(),
	})
	for _, method := range methods {
		addMethod(method, a, integration)
	}
	return a
}
func addMethod(method string, api apigateway.IResource, integration apigateway.LambdaIntegration) {
	api.AddMethod(jsii.String(method), integration, nil)
}
