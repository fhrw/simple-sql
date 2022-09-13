package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type service struct {
	db *sql.DB
}

func main() {
	srv := service{}

	// set up db connection
	db, err := sql.Open("sqlite3", "mydb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	srv.db = db

	http.HandleFunc("/", srv.handleRoot)
	http.HandleFunc("/allocate", srv.handleAllocate)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func (s *service) handleRoot(w http.ResponseWriter, _ *http.Request) {
	fmt.Printf("got / request\n")

	// make db call
	query := "select * from students"
	students, err := s.db.Query(query)
	if err != nil {
		panic(err)
	}
	io.WriteString(w, "This is my website!\n")
}

func (s *service) handleAllocate(w http.ResponseWriter, _ *http.Request) {
	fmt.Printf("got /allocate request\n")

	// make db call

	io.WriteString(w, "some allocation here\n")
}
