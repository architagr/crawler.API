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

func BuildUserPool(stack awscdk.Stack, id string, props *UserPoolLambdaStackProps) (userPool awscognito.IUserPool,
	userPoolClient awscognito.IUserPoolClient) {

	userPool = awscognito.NewUserPool(stack, jsii.String("UserPool"), &awscognito.UserPoolProps{
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
		AccountRecovery: awscognito.AccountRecovery_PHONE_WITHOUT_MFA_AND_EMAIL,
	})

	userPoolClient = awscognito.NewUserPoolClient(stack, jsii.String("UserPoolClient"), &awscognito.UserPoolClientProps{
		UserPool:       userPool,
		GenerateSecret: jsii.Bool(false),
		AuthFlows: &awscognito.AuthFlow{
			AdminUserPassword: jsii.Bool(true),
		},
	})

	// identityPool := awscognito.NewCfnIdentityPool(stack, jsii.String("IdentityPool"), &awscognito.CfnIdentityPoolProps{
	// 	AllowUnauthenticatedIdentities: jsii.Bool(false),
	// 	CognitoIdentityProviders: []map[string]string{
	// 		{
	// 			"clientId":     *userPoolClient.UserPoolClientId(),
	// 			"providerName": *userPool.UserPoolProviderName(),
	// 		},
	// 	},
	// })
	fmt.Printf("User pool %s\n", *userPool.UserPoolId())
	fmt.Printf("Client ID %s\n", *userPoolClient.UserPoolClientId())
	// fmt.Printf("Client Secret %+v\n", userPoolClient.UserPoolClientSecret())

	// fmt.Printf("Identity pool %s\n", identityPool.LogicalId())
	return userPool, userPoolClient
}
