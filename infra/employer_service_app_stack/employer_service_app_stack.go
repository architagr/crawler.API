package employerserviceappstack

import (
	"infra/common"
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"

	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	constructs "github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type EmployerAPILambdaStackProps struct {
	config.CommonProps
	UserPoolArn string
}

func NewEmployerAPILambdaStack(scope constructs.Construct, id string, props *EmployerAPILambdaStackProps) (awscdk.Stack, apigateway.LambdaRestApi) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)
	loginRestApi := buildLambda(stack, scope, props)
	return stack, loginRestApi
}

func buildCognitoAuthorizer(stack awscdk.Stack, props *EmployerAPILambdaStackProps) apigateway.IAuthorizer {
	userPool := awscognito.UserPool_FromUserPoolArn(stack, jsii.String(props.StackNamePrefix.PrependStackName("import-userpool-userapi")), &props.UserPoolArn)

	return apigateway.NewCognitoUserPoolsAuthorizer(stack, jsii.String(props.StackNamePrefix.PrependStackName("cognito-authorizer")), &apigateway.CognitoUserPoolsAuthorizerProps{
		CognitoUserPools: &[]awscognito.IUserPool{
			userPool,
		},
		IdentitySource: jsii.String("method.request.header.Authorization"),
		AuthorizerName: jsii.String(props.StackNamePrefix.PrependStackName("cognito-authorizer-userapi")),
	})

}
func buildLambda(stack awscdk.Stack, scope constructs.Construct, props *EmployerAPILambdaStackProps) apigateway.LambdaRestApi {
	avatarBucket := BuildAvatarBucket(stack, &AvatarBucketStackProps{
		CommonProps: props.CommonProps,
	})
	collectionNames := props.EmployerAPIDB.GetCollectionName()

	authorizer := buildCognitoAuthorizer(stack, props)
	env := make(map[string]*string)
	env["EmployerAPIDbConnectionString"] = props.EmployerAPIDB.GetConnectionString()
	env["EmployerAPIDatabaseName"] = props.EmployerAPIDB.GetDbName()
	env["JobCollectionName"] = jsii.String(collectionNames.GetJobCollectionName())
	env["CompanyCollectionName"] = jsii.String(collectionNames.GetCompanyCollectionName())
	env["EmployerImageBucketName"] = avatarBucket.BucketName()

	userFunction := common.BuildLambda(&common.LambdaConstructProps{
		CommonProps: props.CommonProps,
		Id:          "employer-lambda",
		Handler:     "EmployerAPI", // TODO: get this from makefile
		Service:     "EmployerAPI",
		Name:        "employer-lambda-fn",
		Description: "This function helps in all API related to employer",
		Env:         env,
		Stack:       stack,
	})

	restApiProps := common.RestApiProps{
		CommonProps: props.CommonProps,
		Stack:       stack,
		Id:          "EmployerApi",
		Handler:     userFunction,
		Name:        "EmployerRestApi",
	}

	employerApi := common.BuildRestApi(&restApiProps)

	integration := common.BuildIntegration(&restApiProps)

	baseApi := employerApi.Root()

	common.AddResource("healthCheck", baseApi, []string{common.GET_METHOD}, integration, nil)

	jobsAPI := common.AddResource("job", baseApi, []string{}, integration, authorizer)

	common.AddResource("save", jobsAPI, []string{common.POST_METHOD}, integration, authorizer)
	common.AddResource("get", jobsAPI, []string{common.POST_METHOD}, integration, authorizer)

	companyAPI := common.AddResource("company", baseApi, []string{}, integration, authorizer)

	common.AddResource("save", companyAPI, []string{common.POST_METHOD}, integration, authorizer)
	common.AddResource("get", companyAPI, []string{common.POST_METHOD}, integration, authorizer)

	common.AddResource("{id}",
		common.AddResource("image", companyAPI, []string{}, integration, authorizer),
		[]string{common.PUT_METHOD}, integration, authorizer)

	return employerApi
}
