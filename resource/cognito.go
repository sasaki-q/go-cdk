package resource

import (
	cg "github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/jsii-runtime-go"
)

func (s *ResourceService) NewUserPool() cg.UserPool {
	pool := cg.NewUserPool(s.stack, jsii.String("CdkUserpool"),
		&cg.UserPoolProps{
			UserPoolName: jsii.String("CdkUserpool"),
			PasswordPolicy: &cg.PasswordPolicy{
				RequireSymbols:   jsii.Bool(false),
				RequireUppercase: jsii.Bool(false),
			},
			Mfa: cg.Mfa_OFF,
		},
	)

	pool.AddClient(jsii.String("CdkPoolClient"), &cg.UserPoolClientOptions{
		AuthFlows: &cg.AuthFlow{
			AdminUserPassword: jsii.Bool(true),
			UserPassword:      jsii.Bool(true),
			UserSrp:           jsii.Bool(true),
		},
		GenerateSecret:     jsii.Bool(false),
		UserPoolClientName: jsii.String("CdkPoolClient"),
	})

	return pool
}
