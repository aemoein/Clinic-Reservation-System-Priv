package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Database connection setup
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7R26@llg4grb$&@tcp(127.0.0.1:3306)/Clinical_Reservation")
	if err != nil {
		panic(err.Error())
	}
	return db
}

// Handle patient signup
func patientSignup(w http.ResponseWriter, r *http.Request) {
	// Parse POST data
	r.ParseForm()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	dateOfBirth := r.FormValue("date_of_birth")
	gender := r.FormValue("gender")

	// Insert user into the users table
	db := dbConn()
	insertUser, err := db.Exec("INSERT INTO users (username, email, password, user_type) VALUES (?, ?, ?, ?)", username, email, password, "patient")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	userID, _ := insertUser.LastInsertId()

	// Insert patient into the patients table
	_, err = db.Exec("INSERT INTO patients (user_id, full_name, date_of_birth, gender) VALUES (?, ?, ?, ?)", userID, username, dateOfBirth, gender)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating patient", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Patient signup successful")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signup/patient", patientSignup).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

/*package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:7R26@llg4grb$&@tcp(127.0.0.1:3306)/Clinical_Reservation")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insertUser, err := db.Exec("INSERT INTO users (username, email, password, user_type) VALUES (?, ?, ?, ?)", "Ahmed Elsayed", "aemoein@gmail.com", "12345@qwerty", "patient")
	if err != nil {
		panic(err.Error())
	}

	userID, _ := insertUser.LastInsertId()

	_, err = db.Exec("INSERT INTO patients (user_id, full_name, date_of_birth, gender) VALUES (?, ?, ?, ?)", userID, "Ahmed Elsayed", "2002-06-13", "male")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("User added to both users and patients tables successfully!")*/
