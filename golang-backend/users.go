package main

import "database/sql"

// User represents a user in the system.
type User struct {
	UserID   int
	Name     string
	Email    string
	Password string
	UserType string
}

// CreateUser inserts a new user into the database.
func CreateUser(db *sql.DB, user *User) error {
	_, err := db.Exec("INSERT INTO users (username, email, user_type) VALUES (?, ?, ?)", user.Name, user.Email, user.UserType)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID retrieves a user by their ID.
func GetUserByID(db *sql.DB, userID int) (*User, error) {
	user := new(User)
	err := db.QueryRow("SELECT id, username, email, user_type FROM users WHERE id = ?", userID).Scan(&user.UserID, &user.Name, &user.Email, &user.UserType)
	if err != nil {
		return nil, err
	}
	return user, nil
}
