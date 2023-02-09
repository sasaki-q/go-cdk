package resource

import (
	codebuild "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	codedeploy "github.com/aws/aws-cdk-go/awscdk/v2/awscodedeploy"
	pipeline "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipeline"
	actions "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipelineactions"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	elbv2 "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
)

type NewCicdProps struct {
	LResource LoadBalancerResourceType
	Service   ecs.FargateService
}

func (s *ResourceService) NewPipeline(e NewCicdProps) {
	artifactBucket := awss3.NewBucket(s.stack, jsii.String("ArtifactBucket"), nil)

	sourceOutput := pipeline.NewArtifact(jsii.String("SourceArtifact"))
	buildOutput := pipeline.NewArtifact(jsii.String("BuildArtifact"))

	sourceAction := actions.NewCodeStarConnectionsSourceAction(&actions.CodeStarConnectionsSourceActionProps{
		ActionName:    jsii.String("CdkCodeSourceAction"),
		RunOrder:      jsii.Number(1),
		Repo:          Repository,
		Owner:         Owner,
		Branch:        Branch,
		TriggerOnPush: jsii.Bool(true),
		Output:        sourceOutput,
		ConnectionArn: ConnectionArn,
	})

	buildRole := s.NewBuildRole()

	build := codebuild.NewProject(s.stack, jsii.String("CdkCodeBuild"), &codebuild.ProjectProps{
		BuildSpec: codebuild.BuildSpec_FromSourceFilename(jsii.String("conf/build.yml")),
		Environment: &codebuild.BuildEnvironment{
			BuildImage:  codebuild.LinuxBuildImage_AMAZON_LINUX_2(),
			ComputeType: codebuild.ComputeType_SMALL,
			Privileged:  jsii.Bool(true),
		},
		EnvironmentVariables: &map[string]*codebuild.BuildEnvironmentVariable{
			"AWS_DEFAULT_REGION": {Value: Region},
			"AWS_ACCOUNT_ID":     {Value: AccountId},
			"REPOSITORY_NAME":    {Value: RepoName},
			"REPOSITORY_URI":     {Value: RepoUri},
			"DOCKER_ACCOUNT":     {Value: DockerAccount},
			"DOCKER_PASSWORD":    {Value: DockerPassword},
			"TASK_NAME":          {Value: TaskName},
		},
		Source: codebuild.Source_GitHub(&codebuild.GitHubSourceProps{
			Identifier:  jsii.String("Cdk_Code_Build"),
			Repo:        Repository,
			Owner:       Owner,
			BranchOrRef: Branch,
		}),
		Role: buildRole,
	})

	buildAction := actions.NewCodeBuildAction(
		&actions.CodeBuildActionProps{
			ActionName: jsii.String("CdkCodeBuildAction"),
			Type:       actions.CodeBuildActionType_BUILD,
			Project:    build,
			RunOrder:   jsii.Number(1),
			Input:      sourceOutput,
			Outputs:    &[]pipeline.Artifact{buildOutput},
		},
	)

	deployGroup := codedeploy.NewEcsDeploymentGroup(s.stack, jsii.String("CdkDeployGroup"),
		&codedeploy.EcsDeploymentGroupProps{
			BlueGreenDeploymentConfig: &codedeploy.EcsBlueGreenDeploymentConfig{
				BlueTargetGroup: elbv2.ApplicationTargetGroup_FromTargetGroupAttributes(
					s.stack, jsii.String("CdkBlueApplicationTargetGroup"), &elbv2.TargetGroupAttributes{
						TargetGroupArn:   e.LResource.ListenerOptions.Blue.TargetGroup.TargetGroupArn(),
						LoadBalancerArns: e.LResource.Alb.LoadBalancerArn(),
					},
				),
				GreenTargetGroup: elbv2.ApplicationTargetGroup_FromTargetGroupAttributes(
					s.stack, jsii.String("CdkApplicationGreenTargetGroup"), &elbv2.TargetGroupAttributes{
						TargetGroupArn:   e.LResource.ListenerOptions.Green.TargetGroup.TargetGroupArn(),
						LoadBalancerArns: e.LResource.Alb.LoadBalancerArn(),
					},
				),
				Listener:     e.LResource.ListenerOptions.Blue.Listener,
				TestListener: e.LResource.ListenerOptions.Green.Listener,
			},
			DeploymentConfig: codedeploy.EcsDeploymentConfig_CANARY_10PERCENT_5MINUTES(),
			Service:          e.Service,
		},
	)

	deployAction := actions.NewCodeDeployEcsDeployAction(
		&actions.CodeDeployEcsDeployActionProps{
			ActionName:                  jsii.String("CdkCodeDeployAction"),
			RunOrder:                    jsii.Number(1),
			DeploymentGroup:             deployGroup,
			AppSpecTemplateFile:         pipeline.NewArtifactPath(sourceOutput, jsii.String("conf/deploy.yml")),
			TaskDefinitionTemplateInput: buildOutput,
		},
	)

	pipeline.NewPipeline(s.stack, jsii.String("CdkPipeline"),
		&pipeline.PipelineProps{
			ArtifactBucket: artifactBucket,
			PipelineName:   jsii.String("CdkPipeline"),
			Stages: &[]*pipeline.StageProps{
				{StageName: jsii.String("Source"), Actions: &[]pipeline.IAction{sourceAction}},
				{StageName: jsii.String("Build"), Actions: &[]pipeline.IAction{buildAction}},
				{StageName: jsii.String("Deploy"), Actions: &[]pipeline.IAction{deployAction}},
			},
		},
	)
}
