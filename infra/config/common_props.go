package config

import (
	"fmt"
	"os"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	acm "github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	route53 "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DatabaseModel struct {
	connectionString, dbname, collectionName string
}

func (db *DatabaseModel) GetConnectionString() string {
	return db.connectionString
}

func (db *DatabaseModel) GetDbName() string {
	return db.dbname
}
func (db *DatabaseModel) GetCollectionName() string {
	return db.collectionName
}

type CommonStackProps struct {
	IsLocal string
	Stage   *apigateway.StageOptions
}
type DomainDetails struct {
	RecordName, Url string
}

type Domain struct {
	BaseApi                                     string
	JobApiDomain, UserApiDomain, LoginApiDomain DomainDetails
}
type InfraEnv struct {
	StackNamePrefix string
	ApiBasePath     string
	Domains         Domain
	CommonStackProps
	HostedZoneId   string
	CertificateArn string
}

type CommonProps struct {
	JobAPIDB, LoginAPIDB, UserAPIDB DatabaseModel
	awscdk.StackProps
	InfraEnv
}

const baseDomain = "hiringfunda.com"

func GetCommonProps(app awscdk.App) *CommonProps {
	return &CommonProps{
		JobAPIDB: DatabaseModel{
			connectionString: "mongodb+srv://webscrapper:WebScrapper123@cluster0.xzvihm7.mongodb.net/?retryWrites=true&w=majority",
			dbname:           "webscrapper",
			collectionName:   "test",
		},
		LoginAPIDB: DatabaseModel{
			connectionString: "mongodb+srv://webscrapper:WebScrapper123@cluster0.xzvihm7.mongodb.net/?retryWrites=true&w=majority",
			dbname:           "webscrapper",
			collectionName:   "",
		},
		UserAPIDB: DatabaseModel{
			connectionString: "mongodb+srv://webscrapper:WebScrapper123@cluster0.xzvihm7.mongodb.net/?retryWrites=true&w=majority",
			dbname:           "webscrapper",
			collectionName:   "",
		},
		StackProps: env(app),
		InfraEnv: InfraEnv{
			HostedZoneId:    "Z069835117JUXI2FCKK2F",
			CertificateArn:  "arn:aws:acm:ap-southeast-1:638580160310:certificate/39af19ea-f694-4524-83df-b043ba457278",
			StackNamePrefix: "crawler-api",
			ApiBasePath:     fmt.Sprintf("api.%s", baseDomain),
			Domains: Domain{
				BaseApi: baseDomain,
				JobApiDomain: DomainDetails{
					RecordName: "job-api",
					Url:        fmt.Sprintf("api.%s", baseDomain),
				},
				LoginApiDomain: DomainDetails{
					RecordName: "login-api",
					Url:        fmt.Sprintf("login-api.%s", baseDomain),
				},
				UserApiDomain: DomainDetails{
					RecordName: "user-api",
					Url:        fmt.Sprintf("user-api.%s", baseDomain),
				},
			},
			CommonStackProps: CommonStackProps{
				IsLocal: os.Getenv("isLocal"),
				Stage: &apigateway.StageOptions{
					StageName: jsii.String("Dev"),
				},
			},
		},
	}
}

func AllowCors() *apigateway.ResourceOptions {
	return &apigateway.ResourceOptions{
		DefaultCorsPreflightOptions: GetCorsPreflightOptions(),
		DefaultMethodOptions:        &apigateway.MethodOptions{},
	}
}

func GetCorsPreflightOptions() *apigateway.CorsOptions {
	return &apigateway.CorsOptions{
		AllowOrigins:     apigateway.Cors_ALL_ORIGINS(),
		AllowMethods:     apigateway.Cors_ALL_METHODS(),
		AllowHeaders:     jsii.Strings("Content-Type", "Authorization", "X-Amz-Date", "X-Api-Key"),
		AllowCredentials: jsii.Bool(true),
	}
}

func CreateAcmCertificate(stack awscdk.Stack, scope constructs.Construct, props *InfraEnv) acm.ICertificate {
	return acm.Certificate_FromCertificateArn(stack, jsii.String("clientApiCertificate"), &props.CertificateArn)
}

func GetHostedZone(scope constructs.Construct, id *string, props InfraEnv) route53.IHostedZone {
	return route53.PublicHostedZone_FromHostedZoneAttributes(scope, id, &route53.HostedZoneAttributes{
		HostedZoneId: jsii.String(props.HostedZoneId),
		ZoneName:     jsii.String(props.Domains.BaseApi),
	})
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env(app awscdk.App) awscdk.StackProps {
	accountId := fmt.Sprint(app.Node().TryGetContext(jsii.String("ACCOUNT_ID")))
	region := fmt.Sprint(app.Node().TryGetContext(jsii.String("REGION")))
	project := fmt.Sprint(app.Node().TryGetContext(jsii.String("PROJECT")))
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	return awscdk.StackProps{
		Env: &awscdk.Environment{
			Account: jsii.String(accountId),
			Region:  jsii.String(region),
		},
		Tags: &map[string]*string{
			"project": jsii.String(project),
		},
	}

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
