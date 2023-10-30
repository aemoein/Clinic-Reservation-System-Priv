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

	insertUser, err := db.Exec("INSERT INTO users (username, email, password, user_type) VALUES (?, ?, ?, ?)", "Ahmed Elsayed", "aemoein@gmail.com", "12345@qwerty", "patient")
	if err != nil {
		panic(err.Error())
	}

	userID, _ := insertUser.LastInsertId()

	_, err = db.Exec("INSERT INTO patients (user_id, full_name, date_of_birth, gender) VALUES (?, ?, ?, ?)", userID, "Ahmed Elsayed", "2002-06-13", "male")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("User added to both users and patients tables successfully!")
}
