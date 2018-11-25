package main

import (
	"fmt"

	"gopkg.in/couchbase/gocb.v1"
)

type Beer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Brewery string `json:"brewery_id"`
}

func main() {
	cluster, _ := gocb.Connect("couchbase://localhost")
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "USERNAME",
		Password: "PASSWORD",
	})
	bucket, _ := cluster.OpenBucket("beer-sample", "")

	bucket.Manager("", "").CreatePrimaryIndex("", true, false)

	// Create the beer document
	beer := Beer{
		ID:      "Polar Ice",
		Name:    "Polar Ice",
		Brewery: "Polar",
	}

	// Insert the beer document
	bucket.Upsert("u:polarice", beer, 0)

	fmt.Printf("Inserted document '%s' \r\n", beer.ID)

	// Query the beer sample bucket and find the beer we just added.
	query := gocb.NewN1qlQuery("SELECT name FROM `beer-sample` WHERE brewery_id=$1")
	rows, err := bucket.ExecuteN1qlQuery(query, []interface{}{"Polar"})
	if err == nil {
		var row interface{}
		for rows.Next(&row) {
			fmt.Printf("Row: %v", row)
		}
	}
}
