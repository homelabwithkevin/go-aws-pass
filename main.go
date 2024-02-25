package main

import (
	"context"
	"fmt"

	"go-aws-pass/internal/db"
	sm "go-aws-pass/internal/secretsmanager"

	"github.com/aws/aws-sdk-go-v2/config"
)

type Person struct {
	Email string
	Name  string
	Age   int
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

	featureSecretsManager := false

	// Retrieve Secret if Feature Flag is True
	if featureSecretsManager {
		result := sm.ListSecrets(cfg)

		for _, v := range result.SecretList {
			name := string(*v.Name)

			secrets := sm.GetSecret(cfg, name)

			fmt.Printf("\n----------------------------\n")
			fmt.Printf("Secret Name: %s", name)
			fmt.Printf("\n----------------------------\n")

			for i, v := range secrets {
				fmt.Printf("Secret Key: %s %v", i, "\n")
				fmt.Printf("Secret Value: %s %v", v, "\n\n")
			}
		}
	}

	db.ReadFromDatabase(d, "person", email)
}
