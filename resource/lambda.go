package resource

import (
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

type LambdaResponseType struct {
	HcHandler  lambda.Function `json:"hcHandler"`
	ItemHanler lambda.Function `json:"itemHandler"`
}

func (s *ResourceService) NewLambdaFunction() LambdaResponseType {
	hcHandler := lambda.NewFunction(s.stack, jsii.String("CdkHcHandler"),
		&lambda.FunctionProps{
			FunctionName: jsii.String("CdkHcHandlerFunction"),
			MemorySize:   jsii.Number(512),
			Handler:      jsii.String("hc"),
			Runtime:      lambda.Runtime_GO_1_X(),
			Code:         lambda.Code_FromAsset(jsii.String("./handler/bin/hc"), nil),
			Environment:  &map[string]*string{"Status": jsii.String("Healthy")},
		},
	)

	itemHandler := lambda.NewFunction(s.stack, jsii.String("CdkItemHandler"),
		&lambda.FunctionProps{
			FunctionName: jsii.String("CdkItemHandlerFunction"),
			MemorySize:   jsii.Number(512),
			Handler:      jsii.String("item"),
			Runtime:      lambda.Runtime_GO_1_X(),
			Code:         lambda.Code_FromAsset(jsii.String("./handler/bin/item"), nil),
		},
	)

	return LambdaResponseType{hcHandler, itemHandler}
}
