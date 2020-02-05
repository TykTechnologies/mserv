package aws

import (
	"fmt"
	"github.com/TykTechnologies/serverless/provider"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/pkg/errors"
	"reflect"
)

type AWSConf aws.Config

func init() {
	fmt.Println("registering AWS lambda")
	provider.RegisterProvider("aws-lambda", NewProvider)
}

func NewProvider() (provider.Provider, error) {
	return &Provider{}, nil
}

type Provider struct {
	aws.Config
}

func (p *Provider) Init(conf provider.Conf) error {

	c, ok := conf.(*AWSConf)
	if !ok {
		return fmt.Errorf("unable to resolve conf type: %v", reflect.TypeOf(conf))
	}

	p.Region = c.Region

	awsCfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return errors.Wrap(err, provider.ErrorLoadingDriverConfig)
	}

	awsCfg.Region = c.Region
	p.Config = awsCfg

	fmt.Println(awsCfg.Region)

	return nil
}

func (p Provider) List() ([]provider.Function, error) {

	service := lambda.New(p.Config)

	listFunctionsRequest := service.ListFunctionsRequest(nil)
	listFunctionsOutput, err := listFunctionsRequest.Send()

	if err != nil {
		return nil, errors.Wrap(err, provider.ErrorListingFunctions)
	}

	//lfoJs, _ := json.Marshal(listFunctionsOutput.Functions)
	//logrus.Debug(string(lfoJs))

	functions := make([]provider.Function, 0)
	for _, f := range listFunctionsOutput.Functions {

		myFunc := provider.Function{
			Name:    aws.StringValue(f.FunctionName),
			Version: aws.StringValue(f.Version),
		}

		functions = append(functions, myFunc)
	}

	return functions, nil
}

func (p Provider) Invoke(function provider.Function, requestBody []byte) (*provider.Response, error) {

	service := lambda.New(p.Config)

	if function.GetVersion() == "" {
		function.SetVersion("$LATEST")
	}

	input := lambda.InvokeInput{
		//ClientContext:  aws.String("index"), // need to investigate context
		FunctionName: aws.String(function.GetName()),
		//InvocationType: lambda.InvocationTypeEvent, // To make async
		LogType:   lambda.LogTypeTail,
		Payload:   requestBody,
		Qualifier: aws.String(function.GetVersion()),
	}

	request := service.InvokeRequest(&input)

	lambdaRes, err := request.Send()
	if err != nil {
		return nil, errors.Wrap(err, provider.ErrorInvokingFunction)
	}

	res := provider.Response{
		StatusCode: int(aws.Int64Value(lambdaRes.StatusCode)),
		Body:       lambdaRes.Payload,
	}

	if lambdaRes.FunctionError != nil {
		res.Error = errors.New(aws.StringValue(lambdaRes.FunctionError))
	}

	return &res, nil
}
