package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecret(client *secretsmanager.Client, name string) map[string]interface{} {
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

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	client := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.ListSecretsInput{}
	result, err := client.ListSecrets(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	for _, v := range result.SecretList {
		name := string(*v.Name)

		secrets := GetSecret(client, name)

		fmt.Printf("\n----------------------------\n")
		fmt.Printf("Secret Name: %s", name)
		fmt.Printf("\n----------------------------\n")

		for i, v := range secrets {
			fmt.Printf("Secret Key: %s %v", i, "\n")
			fmt.Printf("Secret Value: %s %v", v, "\n\n")
		}
	}
}
