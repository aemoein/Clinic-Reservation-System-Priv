package main

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	UserID   int
	Name     string
	Email    string
	Password string
	UserType string
}

// 1- Sign In
func SignIn(email, password string) (*User, error) {
	var user User
	err := DB.QueryRow(`
		SELECT userid, name, email, password, usertype 
		FROM users 
		WHERE email = ? AND password = ?
	`, email, password).Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.UserType)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		log.Fatal(err)
	}

	return &user, nil
}

// 2- Sign up
func SignUp(user User) error {
	_, err := DB.Exec(`
		INSERT INTO users (name, email, password, usertype) 
		VALUES (?, ?, ?, ?)
	`, user.Name, user.Email, user.Password, user.UserType)

	return err
}
