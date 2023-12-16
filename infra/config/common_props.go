package config

import (
	"fmt"
	"net/http"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"

	"github.com/aws/jsii-runtime-go"
)

type DatabaseModel[T any] struct {
	connectionString, dbname *string
	collectionName           T
}

func (db *DatabaseModel[T]) GetConnectionString() *string {
	return db.connectionString
}

func (db *DatabaseModel[T]) GetDbName() *string {
	return db.dbname
}
func (db *DatabaseModel[T]) GetCollectionName() T {
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
type EmployerCollection struct {
	jobCollectionName, companyCollectionName string
}

func (db *EmployerCollection) GetCompanyCollectionName() string {
	return db.companyCollectionName
}
func (db *EmployerCollection) GetJobCollectionName() string {
	return db.jobCollectionName
}

type StackNamePrefixModel struct {
	name string
}

func (obj *StackNamePrefixModel) PrependStackName(str string) string {
	return fmt.Sprintf("%s-%s", obj.name, str)
}

type CommonProps struct {
	awscdk.StackProps
	StackNamePrefix                 StackNamePrefixModel
	CurrentEnv                      string
	JobAPIDB, LoginAPIDB, UserAPIDB DatabaseModel[string]
	EmployerAPIDB                   DatabaseModel[EmployerCollection]
	InfraEnv
}

const baseDomain = "hiringfunda.com"

var currentEnv string = ""
var project string = ""
var connectionString string = ""
var dbname string = ""
var hostedZoneId string = ""
var certificateArn string = ""

func GetCommonProps(app awscdk.App) *CommonProps {
	stackProps := env(app)
	apiBasePath := getApiBasePath(currentEnv)
	stackNamePrefix := GetStringWithEnv(currentEnv, project)
	return &CommonProps{
		StackProps: stackProps,
		CurrentEnv: currentEnv,
		StackNamePrefix: StackNamePrefixModel{
			name: stackNamePrefix,
		},
		JobAPIDB: DatabaseModel[string]{
			connectionString: &connectionString,
			dbname:           &dbname,
			collectionName:   GetStringWithEnv(currentEnv, "jobDetails"),
		},
		LoginAPIDB: DatabaseModel[string]{
			connectionString: &connectionString,
			dbname:           &dbname,
			collectionName:   GetStringWithEnv(currentEnv, "authDetails"),
		},
		UserAPIDB: DatabaseModel[string]{
			connectionString: &connectionString,
			dbname:           &dbname,
			collectionName:   GetStringWithEnv(currentEnv, "userProfile"),
		},
		EmployerAPIDB: DatabaseModel[EmployerCollection]{
			connectionString: &connectionString,
			dbname:           &dbname,
			collectionName: EmployerCollection{
				jobCollectionName:     GetStringWithEnv(currentEnv, "employerJobDetails"),
				companyCollectionName: GetStringWithEnv(currentEnv, "companyDetails"),
			},
		},
		InfraEnv: InfraEnv{
			HostedZoneId:     hostedZoneId,
			CertificateArn:   certificateArn,
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
	connectionString = fmt.Sprint(app.Node().TryGetContext(jsii.String("DbConnectionString")))
	dbname = fmt.Sprint(app.Node().TryGetContext(jsii.String("DatabaseName")))
	hostedZoneId = fmt.Sprint(app.Node().TryGetContext(jsii.String("HostedZoneId")))
	certificateArn = fmt.Sprint(app.Node().TryGetContext(jsii.String("CertificateArn")))

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
