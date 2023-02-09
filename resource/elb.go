package resource

import (
	"fmt"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	elbv2 "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/jsii-runtime-go"
)

type ListenerOptionType struct {
	Listener    elbv2.ApplicationListener
	TargetGroup elbv2.ApplicationTargetGroup
}
type LoadBalancerResourceType struct {
	Alb             elbv2.ApplicationLoadBalancer
	ListenerOptions struct {
		Blue  ListenerOptionType
		Green ListenerOptionType
	}
}

func (s *ResourceService) NewLoadBalancer(
	nResource NetworkResourceReturnType,
	service ecs.FargateService,
) LoadBalancerResourceType {
	alb := elbv2.NewApplicationLoadBalancer(s.stack, jsii.String("CdkAlb"),
		&elbv2.ApplicationLoadBalancerProps{
			LoadBalancerName: jsii.String("CdkAlb"),
			Vpc:              nResource.Vpc,
			VpcSubnets:       &ec2.SubnetSelection{Subnets: nResource.Vpc.PublicSubnets()},
			SecurityGroup:    nResource.Sg,
			InternetFacing:   jsii.Bool(true),
		},
	)

	blueTargetGroup := s.NewTargetGroup("Blue", nResource.Vpc, service)
	greenTargetGroup := s.NewTargetGroup("Green", nResource.Vpc, service)

	blueListener := alb.AddListener(jsii.String("CdkBlueListener"), &elbv2.BaseApplicationListenerProps{
		Open:                jsii.Bool(true),
		Port:                jsii.Number(80),
		Protocol:            elbv2.ApplicationProtocol_HTTP,
		DefaultTargetGroups: &[]elbv2.IApplicationTargetGroup{blueTargetGroup},
	})

	greenListener := alb.AddListener(jsii.String("CdkGreenListener"), &elbv2.BaseApplicationListenerProps{
		Open:                jsii.Bool(true),
		Port:                jsii.Number(8080),
		Protocol:            elbv2.ApplicationProtocol_HTTP,
		DefaultTargetGroups: &[]elbv2.IApplicationTargetGroup{greenTargetGroup},
	})

	return LoadBalancerResourceType{
		Alb: alb,
		ListenerOptions: struct {
			Blue  ListenerOptionType
			Green ListenerOptionType
		}{
			Blue: ListenerOptionType{
				Listener:    blueListener,
				TargetGroup: blueTargetGroup,
			},
			Green: ListenerOptionType{
				Listener:    greenListener,
				TargetGroup: greenTargetGroup,
			},
		},
	}
}

func (s *ResourceService) NewTargetGroup(tgType string, vpc ec2.Vpc, service ecs.FargateService) elbv2.ApplicationTargetGroup {
	return elbv2.NewApplicationTargetGroup(s.stack, jsii.String(fmt.Sprintf("Cdk%sTargetGroup", tgType)),
		&elbv2.ApplicationTargetGroupProps{
			Port:            jsii.Number(8080),
			TargetGroupName: jsii.String(fmt.Sprintf("Cdk%sTargetGroup", tgType)),
			TargetType:      elbv2.TargetType_IP,
			Vpc:             vpc,
			HealthCheck: &elbv2.HealthCheck{
				Interval: cdk.Duration_Seconds(jsii.Number(300)),
				Path:     jsii.String("/hc"),
				Port:     jsii.String("8080"),
			},
			Targets: &[]elbv2.IApplicationLoadBalancerTarget{service},
		},
	)
}
