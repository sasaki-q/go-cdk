package resource

import (
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/jsii-runtime-go"
)

func (s *ResourceService) NewBuildRole() iam.Role {
	statements := iam.NewPolicyStatement(&iam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("ecr:BatchCheckLayerAvailability"),
			jsii.String("ecr:CompleteLayerUpload"),
			jsii.String("ecr:GetAuthorizationToken"),
			jsii.String("ecr:InitiateLayerUpload"),
			jsii.String("ecr:PutImage"),
			jsii.String("ecr:UploadLayerPart"),
			jsii.String("ecs:DescribeTasks"),
			jsii.String("ecs:DescribeTaskDefinition"),
		},
		Effect:    iam.Effect_ALLOW,
		Resources: &[]*string{jsii.String("*")},
	})

	policy := iam.NewPolicy(s.stack, jsii.String("CdkCodeBuildPolicy"), &iam.PolicyProps{
		Statements: &[]iam.PolicyStatement{statements},
	})

	buildRole := iam.NewRole(s.stack, jsii.String("CdkCodeBuildRole"),
		&iam.RoleProps{
			Description: jsii.String("Service Role For Code Build"),
			RoleName:    jsii.String("CdkCodeBuildRole"),
			AssumedBy:   iam.NewServicePrincipal(jsii.String("codebuild.amazonaws.com"), nil),
		},
	)

	buildRole.AttachInlinePolicy(policy)

	return buildRole
}
