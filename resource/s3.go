package resource

import (
	"fmt"

	cf "github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	origins "github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
)

func (s *ResourceService) NewS3WithOAI() {
	bucket := s3.NewBucket(s.stack, jsii.String("CdkBucketWithOAI"),
		&s3.BucketProps{
			BucketName:        jsii.String("cdk-test-bucket-with-oai-q"),
			BlockPublicAccess: s3.BlockPublicAccess_BLOCK_ALL(),
			Encryption:        s3.BucketEncryption_S3_MANAGED,
		},
	)

	oai := cf.NewOriginAccessIdentity(s.stack, jsii.String("CdkOAI"),
		&cf.OriginAccessIdentityProps{
			Comment: jsii.String(fmt.Sprintf("%s's origin access identity", *bucket.BucketName())),
		},
	)

	policy := iam.NewPolicyStatement(
		&iam.PolicyStatementProps{
			Effect: iam.Effect_ALLOW,
			Actions: &[]*string{
				jsii.String("s3:Get*"),
				jsii.String("s3:List*"),
			},
			Resources: &[]*string{bucket.BucketArn()},
			Principals: &[]iam.IPrincipal{
				iam.NewArnPrincipal(jsii.String(*oai.Arn())),
			},
		},
	)

	bucket.AddToResourcePolicy(policy)

	cf.NewDistribution(s.stack, jsii.String("CdkCFDistribution"),
		&cf.DistributionProps{
			DefaultBehavior: &cf.BehaviorOptions{
				AllowedMethods: cf.AllowedMethods_ALLOW_ALL(),
				Origin: origins.NewS3Origin(bucket,
					&origins.S3OriginProps{
						OriginAccessIdentity: oai,
						OriginPath:           jsii.String("/test"),
						OriginShieldEnabled:  jsii.Bool(true),
					},
				),
			},
		},
	)
}
