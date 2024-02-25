package systemsmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func CreateParameter(cfg aws.Config, name string, password string) {
	fmt.Printf("\n Creating Parameter: %s \n", name)

	client := ssm.NewFromConfig(cfg)

	input := ssm.PutParameterInput{
		Name:      &name,
		Value:     &password,
		Type:      "SecureString",
		Overwrite: aws.Bool(true),
	}

	_, err := client.PutParameter(context.TODO(), &input)

	if err != nil {
		panic(err)
	}
}
