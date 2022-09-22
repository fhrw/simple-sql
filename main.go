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

type Constraint struct {
	Id   int
	Slot string
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
	http.HandleFunc("/addConstraint", srv.handleAddConstraint)

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
	var stu Student

	err := json.NewDecoder(r.Body).Decode(&stu)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query, err := s.db.Prepare("insert into students(name) values(?)")
	if err != nil {
		panic(err)
	}
	res, err := query.Exec(stu.Name)
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	io.WriteString(w, "wrote to db "+fmt.Sprint(id)+"\n")
}

func (s *service) handleAddConstraint(w http.ResponseWriter, r *http.Request) {
	var cstr Constraint

	err := json.NewDecoder(r.Body).Decode(&cstr)
	if err != nil {
		panic(err)
	}

	query, err := s.db.Prepare("insert into constraints(id, name) values(?,?)")
	if err != nil {
		panic(err)
	}

	res, err := query.Exec(cstr.Id, cstr.Slot)
	if err != nil {
		panic(err)
	}
	_ = res

	io.WriteString(w, "wrote constraint to db\n")
}
