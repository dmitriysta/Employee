package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	//_ "github.com/lib/pq"
	"log"
)

type Employee struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	SecondName string `json:"secondname"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	Salary     string `json:"salary"`
	Department string `json:"department"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "qwerty123456"
	dbname   = "Employee"
)

var employees []Employee
var lastID int

func main() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL database!")

	port := ":80"
	log.Println("Server listen on port:", port)
	http.HandleFunc("/add", addEmployee)
	http.HandleFunc("/employee", getEmployeeByID)
	http.HandleFunc("/dismiss", dismissEmployee)
	http.HandleFunc("/change-salary", changeSalary)
	http.HandleFunc("/employee", changeDepartment)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Server: Could not listen and serve", err)
	}
}

func dismissEmployee(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	id, _ := strconv.Atoi(queryValues.Get("id"))

	sqlStatement := `DELETE FROM employees WHERE id=$1`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "Employee dismissed"})
}

func addEmployee(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	json.NewDecoder(r.Body).Decode(&employee)

	sqlStatement := `INSERT INTO employees (name, surname, second_name, address, phone, salary, department) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := db.QueryRow(sqlStatement, employee.Name, employee.Surname, employee.SecondName, employee.Address, employee.Phone, employee.Salary, employee.Department).Scan(&employee.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employee)
}

func getEmployeeByID(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	id, _ := strconv.Atoi(queryValues.Get("id"))

	var employee Employee
	err := db.QueryRow("SELECT * FROM employees WHERE id=$1", id).Scan(&employee.Id, &employee.Name, &employee.Surname, &employee.SecondName, &employee.Address, &employee.Phone, &employee.Salary, &employee.Department)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employee)
}

func changeSalary(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE employees SET salary=$1 WHERE id=$2", employee.Salary, employee.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Salary updated successfully"))
}

func changeDepartment(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE employees SET department=$1 WHERE id=$2", employee.Department, employee.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Department updated successfully"))
}

