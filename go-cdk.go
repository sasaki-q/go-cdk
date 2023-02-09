package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoCdkStackProps struct {
	awscdk.StackProps
}

func NewGoCdkStack(scope constructs.Construct, id string, props *GoCdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps

	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoCdkStack(app, "MyCdkStack", &GoCdkStackProps{
		awscdk.StackProps{
			Env: env(),
			Synthesizer: awscdk.NewDefaultStackSynthesizer(
				&awscdk.DefaultStackSynthesizerProps{
					// cdk bootstrap --bootstrap-bucket-name {} --profile {}
					FileAssetsBucketName: jsii.String("my-cdk-test-bucket-qq"),
				},
			),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Region: jsii.String("ap-northeast-1"),
	}
}
