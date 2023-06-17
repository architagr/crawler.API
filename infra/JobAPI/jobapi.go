package jobapi

import (
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	route53 "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	route53targets "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	awss3assets "github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

var (
	GET_METHOD  = "GET"
	POST_METHOD = "POST"
)

type JobAPILambdaStackProps struct {
	config.CommonProps
}

func NewJobAPILambdaStack(scope constructs.Construct, id string, props *JobAPILambdaStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	buildLambda(stack, scope, props)
	return stack
}
func buildLambda(stack awscdk.Stack, scope constructs.Construct, props *JobAPILambdaStackProps) {

	env := make(map[string]*string)
	env["DbConnectionString"] = jsii.String(props.JobAPIDB.GetConnectionString())
	env["DatabaseName"] = jsii.String(props.JobAPIDB.GetDbName())
	env["CollectionName"] = jsii.String(props.JobAPIDB.GetCollectionName())
	env["GIN_MODE"] = jsii.String("release")

	jobFunction := lambda.NewFunction(stack, jsii.String("job-lambda"), &lambda.FunctionProps{
		Environment:  &env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      jsii.String("JobAPI"),
		Code:         lambda.Code_FromAsset(jsii.String("./../JobAPI/main.zip"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("job-lambda-fn"),
	})

	jobApi := apigateway.NewLambdaRestApi(stack, jsii.String("JobApi"), &apigateway.LambdaRestApiProps{
		DeployOptions:               props.Stage,
		Handler:                     jobFunction,
		RestApiName:                 jsii.String("JobRestApi"),
		Proxy:                       jsii.Bool(false),
		Deploy:                      jsii.Bool(true),
		DisableExecuteApiEndpoint:   jsii.Bool(false),
		EndpointTypes:               &[]apigateway.EndpointType{apigateway.EndpointType_EDGE},
		DefaultCorsPreflightOptions: config.GetCorsPreflightOptions(),
		DomainName: &apigateway.DomainNameOptions{
			Certificate: config.CreateAcmCertificate(stack, scope, &props.InfraEnv),
			DomainName:  jsii.String(props.ApiBasePath),
		},
	})

	integration := apigateway.NewLambdaIntegration(jobFunction, &apigateway.LambdaIntegrationOptions{})

	baseApi := addResource("jobs", jobApi.Root(), []string{GET_METHOD}, integration)
	addResource("healthCheck", baseApi, []string{GET_METHOD}, integration)

	hostedZone := config.GetHostedZone(stack, jsii.String("JobHostedZone"), props.InfraEnv)

	route53.NewARecord(stack, jsii.String("JobArecord"), &route53.ARecordProps{
		RecordName: jsii.String(props.Domains.JobApiDomain.RecordName),
		Zone:       hostedZone,
		Target:     route53.RecordTarget_FromAlias(route53targets.NewApiGateway(jobApi)),
	})
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
