package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	memdb "github.com/hashicorp/go-memdb"
)

type Person struct {
	Email string
	Name  string
	Age   int
}

func createDatabase(table string) *memdb.MemDB {
	fmt.Printf("Creating database...")

	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			table: {
				Name: "person",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"age": {
						Name:    "age",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Age"},
					},
				},
			},
		},
	}
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	return db
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

func writeToDatabase(db *memdb.MemDB, table string, p Person) {
	txn := db.Txn(true)
	txn.Insert(table, p)
	txn.Commit()
}

func readFromDatabase(db *memdb.MemDB, table string, id string) {
	txn := db.Txn(false)
	txn.Abort()
	raw, err := txn.First(table, "id", id)
	if err != nil {
		panic(err)
	}
	fmt.Println(raw)
}

func main() {
	db := createDatabase("person")
	email := "kevin@homelabwithkevin.com"

	p := Person{email, "kevin", 69}
	writeToDatabase(db, "person", p)

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

	readFromDatabase(db, "person", email)

}
