package jobserviceappservice

import (
	"fmt"
	"infra/common"
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

type JobAPILambdaStackProps struct {
	config.CommonProps
	UserPoolArn string
}

func NewJobAPILambdaStack(scope constructs.Construct, id string, props *JobAPILambdaStackProps) (awscdk.Stack, apigateway.LambdaRestApi) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	jobRestApi := buildLambda(stack, scope, props)
	return stack, jobRestApi
}
func buildLambda(stack awscdk.Stack, scope constructs.Construct, props *JobAPILambdaStackProps) apigateway.LambdaRestApi {

	env := make(map[string]*string)
	env["DbConnectionString"] = jsii.String(props.JobAPIDB.GetConnectionString())
	env["DatabaseName"] = jsii.String(props.JobAPIDB.GetDbName())
	env["CollectionName"] = jsii.String(props.JobAPIDB.GetCollectionName())
	env["GIN_MODE"] = jsii.String("release")

	jobFunction := common.BuildLambda(&common.LambdaConstructProps{
		CommonProps: props.CommonProps,
		Id:          "job-lambda",
		Handler:     "JobAPI",
		Service:     "JobAPI",
		Name:        "job-lambda-fn",
		Description: "This function helps in all API related to jobs",
		Env:         env,
		Stack:       stack,
	})

	restApiProps := common.RestApiProps{
		CommonProps: props.CommonProps,
		Stack:       stack,
		Id:          "JobApi",
		Handler:     jobFunction,
		Name:        "JobRestApi",
	}

	jobApi := common.BuildRestApi(&restApiProps)
	authorizer := buildCognitoAuthorizer(stack, props)
	integration := common.BuildIntegration(&restApiProps)

	common.AddResource("getJobs", jobApi.Root(), []string{common.POST_METHOD}, integration, authorizer)
	common.AddResource("{jobId}",
		common.AddResource("getJobDetail", jobApi.Root(), []string{}, integration, authorizer),
		[]string{common.GET_METHOD}, integration, authorizer)

	common.AddResource("healthCheck", jobApi.Root(), []string{common.GET_METHOD}, integration, nil)
	return jobApi
}

func buildCognitoAuthorizer(stack awscdk.Stack, props *JobAPILambdaStackProps) apigateway.IAuthorizer {
	userPool := awscognito.UserPool_FromUserPoolArn(stack, jsii.String(fmt.Sprintf("%s-import-userpool-jobapi", props.StackNamePrefix)), &props.UserPoolArn)

	return apigateway.NewCognitoUserPoolsAuthorizer(stack, jsii.String(fmt.Sprintf("%s-cognito-authorizer-jobapi", props.StackNamePrefix)), &apigateway.CognitoUserPoolsAuthorizerProps{
		CognitoUserPools: &[]awscognito.IUserPool{
			userPool,
		},
		IdentitySource: jsii.String("method.request.header.Authorization"),
		AuthorizerName: jsii.String(fmt.Sprintf("%s-cognito-authorizer-jobapi", props.StackNamePrefix)),
	})

}
