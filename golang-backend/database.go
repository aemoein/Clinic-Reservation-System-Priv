package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitializeDB() error {
	hostName := os.Getenv("DB_HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	pass := os.Getenv("PASSWORD")
	dbName := os.Getenv("DB_NAME")
	information := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, hostName, port, dbName)
	fmt.Print(information)
	db, err := sql.Open("mysql", information)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to check the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping the database: %w", err)
	}

	fmt.Println("Connected to the database")
	DB = db
	return nil
}

func CreateTables() error {
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
		return fmt.Errorf("failed to create 'users' table: %w", err)
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
		return fmt.Errorf("failed to create 'appointments' table: %w", err)
	}

	fmt.Println("Tables created successfully")
	return nil
}

func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing the database: %v", err)
		}
		fmt.Println("Disconnected from the database")
	}
}
