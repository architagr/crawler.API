package distributionstack

import (
	"infra/config"

	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	acm "github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	route53 "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	route53targets "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"

	constructs "github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DistributionStackambdaStackProps struct {
	config.CommonProps
	LoginRestApi apigateway.LambdaRestApi
	JobsRestApi  apigateway.LambdaRestApi
}

func NewDistributionStackLambdaStack(scope constructs.Construct, id string, props *DistributionStackambdaStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	certificate := acm.Certificate_FromCertificateArn(stack, jsii.String("ApiCertificate"), jsii.String(props.CertificateArn))
	hostedZone := GetHostedZone(stack, jsii.String("ApiHostedZone"), props)
	domain := apigateway.NewDomainName(stack, jsii.String("APiDomain"), &apigateway.DomainNameProps{
		DomainName:     jsii.String(props.ApiBasePath),
		SecurityPolicy: apigateway.SecurityPolicy_TLS_1_2,
		EndpointType:   apigateway.EndpointType_EDGE,
		Certificate:    certificate,
	})

	domain.AddBasePathMapping(props.LoginRestApi, &apigateway.BasePathMappingOptions{
		BasePath:      jsii.String("auth"),
		AttachToStage: jsii.Bool(true),
	})

	domain.AddBasePathMapping(props.JobsRestApi, &apigateway.BasePathMappingOptions{
		BasePath:      jsii.String("jobs"),
		AttachToStage: jsii.Bool(true),
	})
	route53.NewARecord(stack, jsii.String("APIArecord"), &route53.ARecordProps{
		RecordName: jsii.String("api"),
		Zone:       hostedZone,
		Target:     route53.RecordTarget_FromAlias(route53targets.NewApiGatewayDomain(domain)),
	})

	return stack
}

func GetHostedZone(scope constructs.Construct, id *string, props *DistributionStackambdaStackProps) route53.IHostedZone {
	return route53.PublicHostedZone_FromHostedZoneAttributes(scope, id, &route53.HostedZoneAttributes{
		HostedZoneId: jsii.String(props.HostedZoneId),
		ZoneName:     jsii.String(props.BaseDomain),
	})
}
