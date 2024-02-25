package secretsmanager

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecret(cfg aws.Config, name string) map[string]interface{} {
	client := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &name,
	}
	result, err := client.GetSecretValue(context.TODO(), input)

	if err != nil {
		panic(err)
	}

	secretValue := string(*result.SecretString)

	var secretMap map[string]interface{}
	json.Unmarshal([]byte(secretValue), &secretMap)

	return secretMap
}

func ListSecrets(cfg aws.Config) *secretsmanager.ListSecretsOutput {
	client := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.ListSecretsInput{}

	result, err := client.ListSecrets(context.TODO(), input)

	if err != nil {
		panic(err)
	}
	return result
}
