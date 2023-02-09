package resource

import (
	gw "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	cg "github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/jsii-runtime-go"
)

func (s *ResourceService) NewGateway(
	handler LambdaResponseType,
	pool cg.UserPool,
) {
	mygw := gw.NewRestApi(s.stack, jsii.String("CdkApiGateway"),
		&gw.RestApiProps{
			RestApiName: jsii.String("CdkApiGateway"),
		},
	)

	authorizer := gw.NewCognitoUserPoolsAuthorizer(s.stack, jsii.String("CdkCognitoAuthorizer"),
		&gw.CognitoUserPoolsAuthorizerProps{
			AuthorizerName: jsii.String("CdkCognitoAuthorizer"),
			CognitoUserPools: &[]cg.IUserPool{
				pool,
			},
		},
	)

	mygw.Root().AddResource(jsii.String("hc"), nil).AddMethod(jsii.String("GET"),
		gw.NewLambdaIntegration(handler.HcHandler, nil), nil,
	)

	mygw.Root().AddResource(jsii.String("item"), nil).AddMethod(jsii.String("POST"),
		gw.NewLambdaIntegration(handler.ItemHanler,
			&gw.LambdaIntegrationOptions{
				IntegrationResponses: &[]*gw.IntegrationResponse{
					{StatusCode: jsii.String("200")},
				},
				RequestTemplates: &map[string]*string{
					"application/json": jsii.String("{ \"statusCode\": 200 }"),
				},
			},
		),
		&gw.MethodOptions{
			RequestModels: &map[string]gw.IModel{
				"application/json": NewItemModel(mygw),
			},
			RequestValidator: gw.NewRequestValidator(s.stack, jsii.String("CdkGwValidator"),
				&gw.RequestValidatorProps{
					RestApi:             mygw,
					ValidateRequestBody: jsii.Bool(true),
				},
			),
			AuthorizationType: gw.AuthorizationType_COGNITO,
			Authorizer:        authorizer,
		},
	)
}

func NewItemModel(mygw gw.RestApi) gw.Model {
	return mygw.AddModel(jsii.String("ItemModel"), &gw.ModelOptions{
		Schema: &gw.JsonSchema{
			Type: gw.JsonSchemaType_OBJECT,
			Properties: &map[string]*gw.JsonSchema{
				"itemId": {
					Type: gw.JsonSchemaType_STRING,
				},
				"name": {
					Type: gw.JsonSchemaType_STRING,
				},
			},
			Required: &[]*string{jsii.String("itemId")},
		},
	})
}
