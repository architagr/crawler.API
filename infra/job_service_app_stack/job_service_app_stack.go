package jobserviceappservice

import (
	"infra/common"
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
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
	env["DbConnectionString"] = props.JobAPIDB.GetConnectionString()
	env["DatabaseName"] = props.JobAPIDB.GetDbName()
	env["CollectionName"] = jsii.String(props.JobAPIDB.GetCollectionName())

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
	integration := common.BuildIntegration(&restApiProps)

	common.AddResource("getJobs", jobApi.Root(), []string{common.POST_METHOD}, integration, nil)
	common.AddResource("{jobId}",
		common.AddResource("getJobDetail", jobApi.Root(), []string{}, integration, nil),
		[]string{common.GET_METHOD}, integration, nil)

	common.AddResource("{keywords}",
		common.AddResource("courses", jobApi.Root(), []string{}, integration, nil),
		[]string{common.GET_METHOD}, integration, nil)

	common.AddResource("healthCheck", jobApi.Root(), []string{common.GET_METHOD}, integration, nil)
	return jobApi
}
