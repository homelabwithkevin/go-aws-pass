package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"go-aws-pass/internal/db"
	sm "go-aws-pass/internal/secretsmanager"

	"github.com/aws/aws-sdk-go-v2/config"
)

type Person struct {
	Email string
	Name  string
	Age   int
}

func ReadFromConsole(consoleType string) string {
	fmt.Println("Type /exit to exit")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "/exit" {
			os.Exit(0)
		}
		if consoleType == "table" {
			return text
		}
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return ""
}

func main() {
	table := "persons"

	d := db.CreateDatabase(table)
	email := "kevin@homelabwithkevin.com"

	p := Person{email, "kevin", 69}
	db.WriteToDatabase(d, table, db.Person(p))

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

	db.ReadFromDatabase(d, table, email)
}
