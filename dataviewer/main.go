package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"

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

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbName"`
	WebPort  int    `yaml:"webPort"`
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func getRandomUser(db *sql.DB) (User, error) {
	var user User
	err := db.QueryRow("SELECT * FROM users ORDER BY random() LIMIT 1").Scan(&user.ID, &user.Name, &user.UserID, &user.Address, &user.Phone, &user.UserAgent, &user.Company, &user.Email, &user.Team, &user.Location, &user.CreditCard, &user.SocialSecurity)
	return user, err
}

func handleRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	user, err := getRandomUser(db)
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
}

func main() {
	configFilePtr := flag.String("config-file", "", "Path to the configuration file.")
	flag.Parse()

	configData, err := ioutil.ReadFile(*configFilePtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	rand.Seed(time.Now().UnixNano())

	connStr := fmt.Sprintf("dbname=%s host=%s port=%d user=%s password=%s sslmode=disable", config.DBName, config.Host, config.Port, config.Username, config.Password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, db)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, db)
	})
	http.HandleFunc("/health", healthCheck) // Add this line to register the health check handler
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.WebPort), nil))
}
