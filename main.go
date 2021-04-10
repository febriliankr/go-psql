package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Person struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

// postgres://vehymhxoblephe:a4f90fc5228989304adb753ad54971ce0d6cc41fcd77b3b785325f822ea2154f@ec2-107-22-245-82.compute-1.amazonaws.com:5432/d7vmobniqfeohe

const (
	host     = "ec2-107-22-245-82.compute-1.amazonaws.com"
	port     = 5432
	dbname   = "d7vmobniqfeohe"
	user     = "vehymhxoblephe"
	password = "a4f90fc5228989304adb753ad54971ce0d6cc41fcd77b3b785325f822ea2154f"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected! ✅")
	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM person")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var people []Person

	for rows.Next() {
		var person Person
		rows.Scan(&person.Name, &person.Nickname)
		people = append(people, person)
	}

	peopleBytes, _ := json.MarshalIndent(people, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err.Error())
		return
	}

	sqlStatement := `INSERT INTO person (name, nickname) VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, p.Name, p.Nickname)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func main() {

	http.HandleFunc("/", GETHandler)
	http.HandleFunc("/insert", POSTHandler)
	http.Handle("/client", http.FileServer(http.Dir("./static")))

	// Get the PORT from heroku env
	port := os.Getenv("PORT")

	// Verify if heroku provided the port or not
	if os.Getenv("PORT") == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(port, nil))
}
