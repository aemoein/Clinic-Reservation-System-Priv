package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7R26@llg4grb$&@tcp(127.0.0.1:3306)/Clinical_Reservation")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func patientSignup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	dateOfBirth := r.FormValue("date_of_birth")
	gender := r.FormValue("gender")

	db := dbConn()
	insertUser, err := db.Exec("INSERT INTO users (username, email, password, user_type) VALUES (?, ?, ?, ?)", username, email, password, "patient")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	userID, _ := insertUser.LastInsertId()

	_, err = db.Exec("INSERT INTO patients (user_id, full_name, date_of_birth, gender) VALUES (?, ?, ?, ?)", userID, username, dateOfBirth, gender)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating patient", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Patient signup successful")
}

func doctorSignup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	specialization := r.FormValue("specialization")

	db := dbConn()
	insertUser, err := db.Exec("INSERT INTO users (username, email, password, user_type) VALUES (?, ?, ?, ?)", username, email, password, "doctor")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	userID, _ := insertUser.LastInsertId()

	_, err = db.Exec("INSERT INTO doctors (user_id, full_name, specialization) VALUES (?, ?, ?, ?)", userID, username, specialization)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating patient", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "doctor signup successful")
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Check if the user exists in the database
	db := dbConn()

	var storedPassword, userType string
	err := db.QueryRow("SELECT password, user_type FROM users WHERE username = ?", email).Scan(&storedPassword, &userType)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored password
	if password != storedPassword {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Prepare the response
	response := map[string]string{
		"message":   "Sign in successful",
		"user_type": userType, // Include user type in the response
	}

	// Serialize the response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	// Set the content type and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signup/patient", patientSignup).Methods("POST")
	r.HandleFunc("/signup/doctor", doctorSignup).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

/*
package main

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

	insertUser, err := db.Exec("INSERT INTO users (username, email, password, user_type) VALUES (?, ?, ?, ?)", "Ahmed Elsayed", "ahmed33elsayed22@gmail.com", "12345", "patient")
	if err != nil {
		panic(err.Error())
	}

	userID, _ := insertUser.LastInsertId()

	_, err = db.Exec("INSERT INTO patients (user_id, full_name, date_of_birth, gender) VALUES (?, ?, ?, ?)", userID, "Ahmed Elsayed", "2002-06-13", "male")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("User added to both users and patients tables successfully!")
}*/
