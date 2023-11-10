package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitializeDB() {
	db, err := sql.Open("mysql", "root:7R26@llg4grb$&@tcp(127.0.0.1:3306)/CRS")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")

	DB = db
}

func CreateTables() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			userid INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(50) NOT NULL,
			usertype ENUM('doctor', 'patient') NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS appointments (
			appointment_id INT AUTO_INCREMENT PRIMARY KEY,
			doctor_id INT NOT NULL,
			patient_id INT,
			appointment_date DATE NOT NULL,
			start_time TIME NOT NULL,
			end_time TIME NOT NULL,
			FOREIGN KEY (doctor_id) REFERENCES users(userid),
			FOREIGN KEY (patient_id) REFERENCES users(userid)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tables created successfully")
}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Disconnected from the database")
}
