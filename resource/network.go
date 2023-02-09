package resource

import (
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
)

type NetworkResourceReturnType struct {
	Vpc ec2.Vpc
	Sg  ec2.SecurityGroup
}

func (s *ResourceService) NewVpc() NetworkResourceReturnType {
	vpc := ec2.NewVpc(s.stack, jsii.String("CdkVpc"), &ec2.VpcProps{
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{
				Name:       jsii.String("Public"),
				CidrMask:   jsii.Number(24),
				SubnetType: ec2.SubnetType_PUBLIC,
			},
			{
				Name:       jsii.String("Private"),
				CidrMask:   jsii.Number(24),
				SubnetType: ec2.SubnetType_PRIVATE_ISOLATED,
			},
		},
	})

	sg := ec2.NewSecurityGroup(s.stack, jsii.String("CdkSg"),
		&ec2.SecurityGroupProps{
			Vpc:               vpc,
			AllowAllOutbound:  jsii.Bool(true),
			SecurityGroupName: jsii.String("CdkSg"),
		},
	)

	return NetworkResourceReturnType{vpc, sg}
}

func (s *ResourceService) NewVpcEndPoint(nResource NetworkResourceReturnType) {
	nResource.Vpc.AddGatewayEndpoint(jsii.String("CdkS3Endpoint"),
		&ec2.GatewayVpcEndpointOptions{
			Service: ec2.GatewayVpcEndpointAwsService_S3(),
			Subnets: &[]*ec2.SubnetSelection{{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED}},
		},
	)

	nResource.Vpc.AddInterfaceEndpoint(jsii.String("CdkEcrEndpoint"),
		&ec2.InterfaceVpcEndpointOptions{
			Service:           ec2.InterfaceVpcEndpointAwsService_ECR(),
			PrivateDnsEnabled: jsii.Bool(true),
			Subnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	nResource.Vpc.AddInterfaceEndpoint(jsii.String("CdkEcrDkrEndpoint"),
		&ec2.InterfaceVpcEndpointOptions{
			Service:           ec2.InterfaceVpcEndpointAwsService_ECR_DOCKER(),
			PrivateDnsEnabled: jsii.Bool(true),
			Subnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	nResource.Vpc.AddInterfaceEndpoint(jsii.String("CdkCloudWatchLogsEndpoint"),
		&ec2.InterfaceVpcEndpointOptions{
			Service:           ec2.InterfaceVpcEndpointAwsService_CLOUDWATCH_LOGS(),
			PrivateDnsEnabled: jsii.Bool(true),
			Subnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	nResource.Vpc.AddInterfaceEndpoint(jsii.String("CdkSecretsManagerEndpoint"),
		&ec2.InterfaceVpcEndpointOptions{
			Service:           ec2.InterfaceVpcEndpointAwsService_SECRETS_MANAGER(),
			PrivateDnsEnabled: jsii.Bool(true),
			Subnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)
}
