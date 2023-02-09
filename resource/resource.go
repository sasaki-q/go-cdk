package resource

import "github.com/aws/aws-cdk-go/awscdk/v2"

type ResourceService struct {
	stack awscdk.Stack
}

func NewResourceService(stack awscdk.Stack) ResourceService {
	return ResourceService{stack: stack}
}
