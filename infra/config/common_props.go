package config

import (
	"fmt"
	"net/http"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"

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
	Stage *apigateway.StageOptions
}

type InfraEnv struct {
	CommonStackProps
	ApiBasePath    string
	BaseDomain     string
	HostedZoneId   string
	CertificateArn string
}

type CommonProps struct {
	awscdk.StackProps
	StackNamePrefix                 string
	CurrentEnv                      string
	JobAPIDB, LoginAPIDB, UserAPIDB DatabaseModel
	InfraEnv
}

const baseDomain = "hiringfunda.com"

var currentEnv string = ""
var project string = ""

func GetCommonProps(app awscdk.App) *CommonProps {
	stackProps := env(app)
	apiBasePath := getApiBasePath(currentEnv)
	stackNamePrefix := fmt.Sprintf("%s-%s", currentEnv, project)
	return &CommonProps{
		StackProps:      stackProps,
		CurrentEnv:      currentEnv,
		StackNamePrefix: stackNamePrefix,
		JobAPIDB: DatabaseModel{
			connectionString: "mongodb+srv://webscrapper:WebScrapper123@cluster0.xzvihm7.mongodb.net/?retryWrites=true&w=majority",
			dbname:           "webscrapper",
			collectionName:   "jobDetails",
		},
		LoginAPIDB: DatabaseModel{
			connectionString: "mongodb+srv://webscrapper:WebScrapper123@cluster0.xzvihm7.mongodb.net/?retryWrites=true&w=majority",
			dbname:           "webscrapper",
			collectionName:   "authDetails",
		},
		UserAPIDB: DatabaseModel{
			connectionString: "mongodb+srv://webscrapper:WebScrapper123@cluster0.xzvihm7.mongodb.net/?retryWrites=true&w=majority",
			dbname:           "webscrapper",
			collectionName:   "",
		},
		InfraEnv: InfraEnv{
			HostedZoneId: "Z069835117JUXI2FCKK2F",
			//CertificateArn:  "arn:aws:acm:ap-southeast-1:638580160310:certificate/39af19ea-f694-4524-83df-b043ba457278",
			CertificateArn:   "arn:aws:acm:us-east-1:638580160310:certificate/d70647bd-a714-4add-8d3b-787f55ab5213",
			ApiBasePath:      apiBasePath,
			BaseDomain:       baseDomain,
			CommonStackProps: CommonStackProps{},
		},
	}
}

func GetCorsPreflightOptions() *apigateway.CorsOptions {
	return &apigateway.CorsOptions{
		AllowOrigins: apigateway.Cors_ALL_ORIGINS(),
		// AllowMethods:     apigateway.Cors_ALL_METHODS(),
		// AllowHeaders:     jsii.Strings("Content-Type", "Authorization", "X-Amz-Date", "X-Api-Key"),
		// AllowCredentials: jsii.Bool(true),
		StatusCode: jsii.Number(http.StatusOK),
	}
}
func getApiBasePath(env string) string {
	return fmt.Sprintf("api-%s.%s", env, baseDomain)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env(app awscdk.App) awscdk.StackProps {
	accountId := fmt.Sprint(app.Node().TryGetContext(jsii.String("ACCOUNT_ID")))
	region := fmt.Sprint(app.Node().TryGetContext(jsii.String("REGION")))
	project = fmt.Sprint(app.Node().TryGetContext(jsii.String("PROJECT")))
	currentEnv = fmt.Sprint(app.Node().TryGetContext(jsii.String("ENV")))

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
			"env":     jsii.String(currentEnv),
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
