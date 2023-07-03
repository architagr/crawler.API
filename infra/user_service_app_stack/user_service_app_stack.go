package userserviceappstack

import (
	"fmt"
	"infra/common"
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"

	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	constructs "github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type UserAPILambdaStackProps struct {
	config.CommonProps
	UserPoolArn string
}

func NewUserAPILambdaStack(scope constructs.Construct, id string, props *UserAPILambdaStackProps) (awscdk.Stack, apigateway.LambdaRestApi) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)
	loginRestApi := buildLambda(stack, scope, props)
	return stack, loginRestApi
}

func buildCognitoAuthorizer(stack awscdk.Stack, props *UserAPILambdaStackProps) apigateway.IAuthorizer {
	userPool := awscognito.UserPool_FromUserPoolArn(stack, jsii.String(fmt.Sprintf("%s-import-userpool-userapi", props.StackNamePrefix)), &props.UserPoolArn)

	return apigateway.NewCognitoUserPoolsAuthorizer(stack, jsii.String(fmt.Sprintf("%s-cognito-authorizer", props.StackNamePrefix)), &apigateway.CognitoUserPoolsAuthorizerProps{
		CognitoUserPools: &[]awscognito.IUserPool{
			userPool,
		},
		IdentitySource: jsii.String("method.request.header.Authorization"),
		AuthorizerName: jsii.String(fmt.Sprintf("%s-cognito-authorizer-userapi", props.StackNamePrefix)),
	})

}
func buildLambda(stack awscdk.Stack, scope constructs.Construct, props *UserAPILambdaStackProps) apigateway.LambdaRestApi {
	avatarBucket := BuildAvatarBucket(stack, &AvatarBucketStackProps{
		CommonProps: props.CommonProps,
	})

	authorizer := buildCognitoAuthorizer(stack, props)
	env := make(map[string]*string)
	env["DbConnectionString"] = jsii.String(props.UserAPIDB.GetConnectionString())
	env["DatabaseName"] = jsii.String(props.UserAPIDB.GetDbName())
	env["UserCollectionName"] = jsii.String(props.UserAPIDB.GetCollectionName())
	env["GIN_MODE"] = jsii.String("release")
	env["AvatarImageBucketName"] = avatarBucket.BucketName()

	userFunction := common.BuildLambda(&common.LambdaConstructProps{
		CommonProps: props.CommonProps,
		Id:          "user-lambda",
		Handler:     "UserAPI",
		Service:     "UserAPI",
		Name:        "user-lambda-fn",
		Description: "This function helps in all API related to user",
		Env:         env,
		Stack:       stack,
	})
	avatarBucket.GrantReadWrite(userFunction, nil)

	restApiProps := common.RestApiProps{
		CommonProps: props.CommonProps,
		Stack:       stack,
		Id:          "UserApi",
		Handler:     userFunction,
		Name:        "UserRestApi",
	}

	userApi := common.BuildRestApi(&restApiProps)

	integration := common.BuildIntegration(&restApiProps)

	baseApi := userApi.Root()

	common.AddResource("healthCheck", baseApi, []string{common.GET_METHOD}, integration, nil)
	profile := common.AddResource("profile", baseApi, []string{common.POST_METHOD}, integration, authorizer)

	common.AddResource("{id}",
		common.AddResource("image", profile, []string{}, integration, authorizer),
		[]string{common.PUT_METHOD}, integration, authorizer)
	common.AddResource("{username}", profile, []string{common.GET_METHOD}, integration, authorizer)

	return userApi
}
