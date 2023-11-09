package main

// User represents a user in the system.
type User struct {
	UserID   int
	Name     string
	Email    string
	Password string
	UserType string
}

func SignUp(user User) error {
	_, err := DB.Exec(`
		INSERT INTO users (name, email, password, usertype) 
		VALUES (?, ?, ?, ?)
	`, user.Name, user.Email, user.Password, user.UserType)

	return err
}
