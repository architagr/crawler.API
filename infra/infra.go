package main

import (
	"infra/config"
	distributionstack "infra/distribution_stack"
	employerserviceappstack "infra/employer_service_app_stack"
	jobserviceappstack "infra/job_service_app_stack"
	loginserviceappstack "infra/login_service_app_stack"
	userserviceappstack "infra/user_service_app_stack"

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

	_, loginRestApi, userPool := loginserviceappstack.NewLoginAPILambdaStack(app, props.StackNamePrefix.PrependStackName("LoginAPIStack"), &loginserviceappstack.LoginAPILambdaStackProps{
		CommonProps: *props,
	})
	_, jobsRestApi := jobserviceappstack.NewJobAPILambdaStack(app, props.StackNamePrefix.PrependStackName("JobAPIStack"), &jobserviceappstack.JobAPILambdaStackProps{
		CommonProps: *props,
		UserPoolArn: *userPool.UserPoolArn(),
	})
	_, userRestApi := userserviceappstack.NewUserAPILambdaStack(app, props.StackNamePrefix.PrependStackName("UserAPIStack"), &userserviceappstack.UserAPILambdaStackProps{
		CommonProps: *props,
		UserPoolArn: *userPool.UserPoolArn(),
	})
	_, employerRestApi := employerserviceappstack.NewEmployerAPILambdaStack(app, props.StackNamePrefix.PrependStackName("EmployerAPIStack"), &employerserviceappstack.EmployerAPILambdaStackProps{
		CommonProps: *props,
		UserPoolArn: *userPool.UserPoolArn(),
	})

	distributionstack.NewDistributionStackLambdaStack(app, props.StackNamePrefix.PrependStackName("DistributionStack"), &distributionstack.DistributionStackambdaStackProps{
		CommonProps:     *props,
		LoginRestApi:    loginRestApi,
		JobsRestApi:     jobsRestApi,
		UserRestApi:     userRestApi,
		EmployerRestApi: employerRestApi,
	})
	app.Synth(nil)
}
