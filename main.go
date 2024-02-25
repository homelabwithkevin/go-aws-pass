package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go-aws-pass/internal/db"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Person struct {
	Email string
	Name  string
	Age   int
}

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
	d := db.CreateDatabase("person")
	email := "kevin@homelabwithkevin.com"

	p := Person{email, "kevin", 69}
	db.WriteToDatabase(d, "person", db.Person(p))

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

	db.ReadFromDatabase(d, "person", email)

}
