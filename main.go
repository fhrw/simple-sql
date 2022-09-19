package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type service struct {
	db *sql.DB
}

type Student struct {
	Id   int
	Name string
}

func main() {
	fmt.Println("running...")
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
	http.HandleFunc("/list", srv.handleList)
	http.HandleFunc("/add", srv.handleAdd)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func (s *service) handleRoot(w http.ResponseWriter, _ *http.Request) {
	fmt.Printf("got / request\n")

	io.WriteString(w, "SQL server for timetable program!\n")
}

func (s *service) handleAllocate(w http.ResponseWriter, _ *http.Request) {
	fmt.Printf("got /allocate request\n")

	// make db call

	io.WriteString(w, "some allocation here\n")
}

func (s *service) handleList(w http.ResponseWriter, _ *http.Request) {
	fmt.Printf("got /list request\n")
	query := `SELECT * from STUDENTS`

	rows, err := s.db.Query(query)
	if err != nil {
		panic(err)
	}
	var allStudents = []*Student{}
	for rows.Next() {
		s := new(Student)
		rows.Scan(&s.Id, &s.Name)
		allStudents = append(allStudents, s)
	}

	if err := json.NewEncoder(w).Encode(allStudents); err != nil {
		panic(err)
	}
}

func (s *service) handleAdd(w http.ResponseWriter, r *http.Request) {
	// work out how to send JSON or something as POST request
	currStudent := Student{Id: 666, Name: "doris d"}
	insertQuery := `INSERT INTO students VALUES ($1, $2)`
	_, err := s.db.Exec(insertQuery, currStudent.Id, currStudent.Name)
	if err != nil {
		panic(err)
	}
}
