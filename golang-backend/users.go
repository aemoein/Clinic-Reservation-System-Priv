package main

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	UserID   int    `json:"userid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"usertype"`
}

// 1- Sign In
func SignIn(email, password string) (*User, error) {
	var user User
	err := DB.QueryRow(`
		SELECT userid, name, email, password, usertype 
		FROM users 
		WHERE email = ? AND password = ?
	`, email, password).Scan(&user.UserID, &user.UserName, &user.Email, &user.Password, &user.UserType)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		log.Fatal(err)
	}

	log.Printf("User data before returning: %+v", user)
	return &user, nil
}

// 2- Sign up
func SignUp(user User) error {
	_, err := DB.Exec(`
		INSERT INTO users (name, email, password, usertype) 
		VALUES (?, ?, ?, ?)
	`, user.UserName, user.Email, user.Password, user.UserType)

	return err
}

func GetUsernameByID(userID int) (string, error) {
	var username string

	row := DB.QueryRow("SELECT name FROM users WHERE userid = ?", userID)
	err := row.Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found for ID: %d", userID)
		}
		return "", err
	}

	return username, nil
}
