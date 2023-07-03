package loginserviceappstack

import (
	"fmt"
	"infra/config"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	awscognito "github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/jsii-runtime-go"
)

type UserPoolLambdaStackProps struct {
	config.CommonProps
}

func BuildUserPool(stack awscdk.Stack, props *UserPoolLambdaStackProps) (userPool awscognito.IUserPool,
	userPoolClient awscognito.IUserPoolClient) {

	userPoolName := fmt.Sprintf("%s-userpool", props.StackNamePrefix)
	userClientName := fmt.Sprintf("%s-userpoolclient", props.StackNamePrefix)

	userPool = awscognito.NewUserPool(stack, &userPoolName, &awscognito.UserPoolProps{
		SelfSignUpEnabled: jsii.Bool(true),
		AutoVerify: &awscognito.AutoVerifiedAttrs{
			Email: jsii.Bool(true),
		},
		SignInAliases: &awscognito.SignInAliases{
			Email:    jsii.Bool(true),
			Username: jsii.Bool(true),
			Phone:    jsii.Bool(true),
		},
		PasswordPolicy: &awscognito.PasswordPolicy{
			MinLength:            jsii.Number(8),
			RequireDigits:        jsii.Bool(true),
			RequireLowercase:     jsii.Bool(true),
			RequireSymbols:       jsii.Bool(true),
			RequireUppercase:     jsii.Bool(true),
			TempPasswordValidity: awscdk.Duration_Days(jsii.Number(14)),
		},
		RemovalPolicy:   awscdk.RemovalPolicy_DESTROY,
		AccountRecovery: awscognito.AccountRecovery_PHONE_WITHOUT_MFA_AND_EMAIL,
	})

	userPoolClient = awscognito.NewUserPoolClient(stack, &userClientName, &awscognito.UserPoolClientProps{
		UserPool:       userPool,
		GenerateSecret: jsii.Bool(false),
		AuthFlows: &awscognito.AuthFlow{
			AdminUserPassword: jsii.Bool(true),
			UserPassword:      jsii.Bool(true),
		},
		SupportedIdentityProviders: &[]awscognito.UserPoolClientIdentityProvider{
			awscognito.UserPoolClientIdentityProvider_COGNITO(),
		},
	})
	return userPool, userPoolClient
}
