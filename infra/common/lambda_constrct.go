package common

import (
	"fmt"
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awss3assets "github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/jsii-runtime-go"
)

type LambdaConstructProps struct {
	config.CommonProps
	Stack       awscdk.Stack
	Id          string
	Handler     string
	Service     string
	Name        string
	Description string
	Env         map[string]*string
}

func BuildLambda(props *LambdaConstructProps) lambda.IFunction {
	functionId := fmt.Sprintf("%s-%s", props.StackNamePrefix, props.Id)
	functionName := fmt.Sprintf("%s-%s", props.StackNamePrefix, props.Name)
	functionCodePath := jsii.String(fmt.Sprintf("./../%s/main.zip", props.Service))

	return lambda.NewFunction(props.Stack, &functionId, &lambda.FunctionProps{
		Environment:  &props.Env,
		Runtime:      lambda.Runtime_GO_1_X(),
		Handler:      &props.Handler,
		Code:         lambda.Code_FromAsset(functionCodePath, &awss3assets.AssetOptions{}),
		FunctionName: &functionName,
		Description:  &props.Description,
	})
}
