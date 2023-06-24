package main

import (
	jobapi "infra/JobAPI"
	"infra/config"
	distributionstack "infra/distribution_stack"
	loginservice "infra/login_service"

	"github.com/aws/aws-cdk-go/awscdk/v2"

	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	awscdk.StackProps
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	props := config.GetCommonProps(app)

	_, jobsRestApi := jobapi.NewJobAPILambdaStack(app, "JobAPIStack", &jobapi.JobAPILambdaStackProps{
		CommonProps: *props,
	})
	_, loginRestApi := loginservice.NewLoginAPILambdaStack(app, "LoginAPIStack", &loginservice.LoginAPILambdaStackProps{
		CommonProps: *props,
	})
	distributionstack.NewDistributionStackLambdaStack(app, "distributionStack", &distributionstack.DistributionStackambdaStackProps{
		CommonProps:  *props,
		LoginRestApi: loginRestApi,
		JobsRestApi:  jobsRestApi,
	})
	app.Synth(nil)
}
