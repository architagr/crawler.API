package main

import (
	"fmt"
	"infra/config"
	distributionstack "infra/distribution_stack"
	jobserviceappstack "infra/job_service_app_stack"
	loginserviceappstack "infra/login_service_app_stack"

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

	_, jobsRestApi := jobserviceappstack.NewJobAPILambdaStack(app, fmt.Sprintf("%s-JobAPIStack", props.StackNamePrefix), &jobserviceappstack.JobAPILambdaStackProps{
		CommonProps: *props,
	})
	_, loginRestApi := loginserviceappstack.NewLoginAPILambdaStack(app, fmt.Sprintf("%s-LoginAPIStack", props.StackNamePrefix), &loginserviceappstack.LoginAPILambdaStackProps{
		CommonProps: *props,
	})
	distributionstack.NewDistributionStackLambdaStack(app, fmt.Sprintf("%s-DistributionStack", props.StackNamePrefix), &distributionstack.DistributionStackambdaStackProps{
		CommonProps:  *props,
		LoginRestApi: loginRestApi,
		JobsRestApi:  jobsRestApi,
	})
	app.Synth(nil)
}
