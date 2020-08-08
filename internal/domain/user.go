package domain

import "time"

// User represents a user.
type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetTableName returns database table name
func (u User) GetTableName() string {
	return "users"
}

// GetID returns the user ID.
func (u User) GetID() string {
	return u.ID
}

// GetUsername returns the user name.
func (u User) GetUsername() string {
	return u.Email
}

// GetRole returns the user role
func (u User) GetRole() string {
	return "admin"
}
