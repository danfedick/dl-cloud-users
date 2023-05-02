package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	UserID         string `json:"user_id"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	UserAgent      string `json:"user_agent"`
	Company        string `json:"company"`
	Email          string `json:"email"`
	Team           string `json:"team"`
	Location       string `json:"location"`
	CreditCard     string `json:"credit_card"`
	SocialSecurity string `json:"social_security"`
}

func main() {
	dbNamePtr := flag.String("database", "", "Database name (required)")
	hostPtr := flag.String("hostname", "", "Database host (required)")
	portPtr := flag.String("port", "", "Database port (required)")
	passwordPtr := flag.String("password", "", "Database password (required)")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user, err := getRandomUser(*dbNamePtr, *hostPtr, *portPtr, *passwordPtr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userJSON, err := json.MarshalIndent(user, "", " ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(userJSON)
	})

	http.ListenAndServe(":8888", nil)
}

func getRandomUser(dbName, host, port, password string) (*User, error) {
	connStr := fmt.Sprintf("dbname=%s host=%s port=%s password=%s sslmode=disable", dbName, host, port, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var user User

	err = db.QueryRow("SELECT * FROM users ORDER BY random() LIMIT 1").Scan(&user.ID, &user.Name, &user.UserID, &user.Address, &user.Phone, &user.UserAgent, &user.Company, &user.Email, &user.Team, &user.Location, &user.CreditCard, &user.SocialSecurity)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
