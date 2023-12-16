package common

import (
	"fmt"
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/jsii-runtime-go"

	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
)

var (
	GET_METHOD  = "GET"
	POST_METHOD = "POST"
	PUT_METHOD  = "PUT"
)

type RestApiProps struct {
	config.CommonProps
	Stack   awscdk.Stack
	Id      string
	Handler lambda.IFunction
	Name    string
}

func BuildRestApi(props *RestApiProps) apigateway.LambdaRestApi {
	functionId := props.StackNamePrefix.PrependStackName(props.Id)
	functionName := props.StackNamePrefix.PrependStackName(props.Name)

	apiLogs := awslogs.NewLogGroup(props.Stack, &functionId, &awslogs.LogGroupProps{
		LogGroupName:  jsii.String(props.StackNamePrefix.PrependStackName(fmt.Sprintf("%s-Log", functionName))),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	deployOptions := &apigateway.StageOptions{
		StageName:            jsii.String(props.CurrentEnv),
		AccessLogDestination: apigateway.NewLogGroupLogDestination(apiLogs),
		LoggingLevel:         apigateway.MethodLoggingLevel_ERROR,
		DataTraceEnabled:     jsii.Bool(true),
		AccessLogFormat: apigateway.AccessLogFormat_JsonWithStandardFields(&apigateway.JsonWithStandardFieldProps{
			Caller:         jsii.Bool(true),
			HttpMethod:     jsii.Bool(true),
			Ip:             jsii.Bool(true),
			Protocol:       jsii.Bool(true),
			RequestTime:    jsii.Bool(true),
			ResourcePath:   jsii.Bool(true),
			ResponseLength: jsii.Bool(true),
			Status:         jsii.Bool(true),
			User:           jsii.Bool(true),
		}),
	}
	return apigateway.NewLambdaRestApi(props.Stack, &props.Id, &apigateway.LambdaRestApiProps{
		DeployOptions:               deployOptions,
		Handler:                     props.Handler,
		RestApiName:                 &functionName,
		Proxy:                       jsii.Bool(false),
		Deploy:                      jsii.Bool(true),
		DisableExecuteApiEndpoint:   jsii.Bool(false),
		EndpointTypes:               &[]apigateway.EndpointType{apigateway.EndpointType_REGIONAL},
		DefaultCorsPreflightOptions: config.GetCorsPreflightOptions(),
		CloudWatchRole:              jsii.Bool(true),
	})

}

func BuildIntegration(props *RestApiProps) apigateway.LambdaIntegration {
	return apigateway.NewLambdaIntegration(props.Handler, &apigateway.LambdaIntegrationOptions{
		Proxy: jsii.Bool(true),
	})
}
func AddResource(path string, api apigateway.IResource, methods []string, integration apigateway.LambdaIntegration, authorizer apigateway.IAuthorizer) apigateway.IResource {
	a := api.AddResource(jsii.String(path), &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: config.GetCorsPreflightOptions(),
	})
	for _, method := range methods {
		AddMethod(method, a, integration, authorizer)
	}
	return a
}
func AddMethod(method string, api apigateway.IResource, integration apigateway.LambdaIntegration, authorizer apigateway.IAuthorizer) {
	var options *apigateway.MethodOptions = nil
	if authorizer != nil {
		options = &apigateway.MethodOptions{
			Authorizer:          authorizer,
			AuthorizationType:   apigateway.AuthorizationType_COGNITO,
			AuthorizationScopes: jsii.Strings("aws.cognito.signin.user.admin"),
		}
	}
	api.AddMethod(jsii.String(method), integration, options)
}
