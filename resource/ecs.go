package resource

import (
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
)

func (s *ResourceService) NewEcs(nResource NetworkResourceReturnType) ecs.FargateService {
	repository := ecr.Repository_FromRepositoryName(s.stack, jsii.String("ecr-association-id"), jsii.String("cdk-ecr"))

	logBucket := s3.Bucket_FromBucketName(s.stack, jsii.String("cdk-log-bucket-association"), jsii.String("cdk-log-bucket-for-ecs-q"))

	myLogGroup := logs.LogGroup_FromLogGroupName(s.stack, jsii.String("CdkClusterLogAssociation"), jsii.String("CdkLogGroup"))

	s.NewVpcEndPoint(nResource)

	cluster := ecs.NewCluster(s.stack, jsii.String("CdkCluster"),
		&ecs.ClusterProps{
			Vpc:         nResource.Vpc,
			ClusterName: jsii.String("CdkCluster"),
			ExecuteCommandConfiguration: &ecs.ExecuteCommandConfiguration{
				LogConfiguration: &ecs.ExecuteCommandLogConfiguration{
					CloudWatchLogGroup:          myLogGroup,
					CloudWatchEncryptionEnabled: jsii.Bool(true),
					S3Bucket:                    logBucket,
					S3EncryptionEnabled:         jsii.Bool(true),
					S3KeyPrefix:                 jsii.String("cdk-log"),
				},
				Logging: ecs.ExecuteCommandLogging_OVERRIDE,
			},
		},
	)

	taskDef := ecs.NewFargateTaskDefinition(s.stack, jsii.String("CdkTask"),
		&ecs.FargateTaskDefinitionProps{
			Family:          jsii.String("CdkTaskFamily"),
			Cpu:             jsii.Number(256),
			TaskRole:        iam.Role_FromRoleName(s.stack, jsii.String("taskRole"), jsii.String("ecsTaskExecutionRole"), nil),
			ExecutionRole:   iam.Role_FromRoleName(s.stack, jsii.String("execRole"), jsii.String("ecsTaskExecutionRole"), nil),
			RuntimePlatform: &ecs.RuntimePlatform{CpuArchitecture: ecs.CpuArchitecture_X86_64()},
		},
	)

	taskDef.AddContainer(jsii.String("CdkAddContainer"),
		&ecs.ContainerDefinitionOptions{
			Image:          ecs.ContainerImage_FromEcrRepository(repository, jsii.String("latest")),
			MemoryLimitMiB: jsii.Number(512),
			ContainerName:  jsii.String("CdkContainer"),
			PortMappings: &[]*ecs.PortMapping{
				{
					Protocol:      ecs.Protocol_TCP,
					HostPort:      jsii.Number(8080),
					ContainerPort: jsii.Number(8080),
				},
			},
			Logging:     ecs.LogDriver_AwsLogs(&ecs.AwsLogDriverProps{StreamPrefix: jsii.String("CdkLog"), LogGroup: myLogGroup}),
			Environment: &map[string]*string{"PORTS": jsii.String("8080"), "TEST_ENV": jsii.String("CDK_TASK")},
		},
	)

	service := ecs.NewFargateService(s.stack, jsii.String("CdkService"),
		&ecs.FargateServiceProps{
			Cluster:              cluster,
			DeploymentController: &ecs.DeploymentController{Type: "CODE_DEPLOY"},
			DesiredCount:         jsii.Number(1),
			ServiceName:          jsii.String("CdkService"),
			TaskDefinition:       taskDef,
			SecurityGroups:       &[]ec2.ISecurityGroup{nResource.Sg},
			VpcSubnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	return service
}
