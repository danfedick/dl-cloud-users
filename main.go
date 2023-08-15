package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	Username          string `json:"username"`
	Groupname         string `json:"groupname"`
	AzureSubscription string `json:"azure_subscription_id"`
	AwsAccount        string `json:"aws_account_id"`
}

func main() {
	// Define flag arguments
	filePtr := flag.String("file", "", "Path to the users JSON file (required)")
	dbNamePtr := flag.String("database", "", "Database name (required)")
	hostPtr := flag.String("hostname", "", "Database host (required)")
	portPtr := flag.String("port", "", "Database port (required)")
	passwordPtr := flag.String("password", "", "Database password (required)")

	// Parse flag arguments
	flag.Parse()

	// Check if required flags are provided
	if *filePtr == "" || *dbNamePtr == "" || *hostPtr == "" || *portPtr == "" || *passwordPtr == "" {
		flag.PrintDefaults()
		log.Fatal("All flags are required.")
	}

	// Read users JSON file
	data, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		log.Fatal(err)
	}

	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to PostgreSQL
	connStr := fmt.Sprintf("user=postgres password=%s dbname=%s host=%s port=%s sslmode=disable", *passwordPtr, *dbNamePtr, *hostPtr, *portPtr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username TEXT,
        groupname TEXT,
        azure_subscription_id TEXT,
        aws_account_id TEXT
    )`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert users into the database
	for _, user := range users {
		_, err := db.Exec(`INSERT INTO users (
            username, groupname, azure_subscription_id, aws_account_id
        ) VALUES ($1, $2, $3, $4)`,
			user.Username, user.Groupname, user.AzureSubscription, user.AwsAccount)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Data loaded successfully.")
}
