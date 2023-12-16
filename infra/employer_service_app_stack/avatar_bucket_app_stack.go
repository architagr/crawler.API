package employerserviceappstack

import (
	"infra/config"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
)

type AvatarBucketStackProps struct {
	config.CommonProps
}

func BuildAvatarBucket(stack awscdk.Stack, props *AvatarBucketStackProps) awss3.IBucket {
	avatarBucketName := props.StackNamePrefix.PrependStackName("employer-bucket")
	return awss3.NewBucket(stack, &avatarBucketName, &awss3.BucketProps{
		BucketName:    &avatarBucketName,
		AccessControl: awss3.BucketAccessControl_BUCKET_OWNER_FULL_CONTROL,
		Versioned:     jsii.Bool(true),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
}
