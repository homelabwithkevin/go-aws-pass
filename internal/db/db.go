package db

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
)

type Person struct {
	Email string
	Name  string
	Age   int
}

func CreateDatabase(table string) *memdb.MemDB {
	fmt.Printf("Creating database with table %s", table)
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			table: {
				Name: table,
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

func WriteToDatabase(db *memdb.MemDB, table string, person Person) {
	txn := db.Txn(true)
	txn.Insert(table, person)
	txn.Commit()
}

func ReadFromDatabase(db *memdb.MemDB, table string, id string) {
	txn := db.Txn(false)
	txn.Abort()
	raw, err := txn.First(table, "id", id)
	if err != nil {
		panic(err)
	}
	fmt.Println(raw)
}
