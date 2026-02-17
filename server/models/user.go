package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Role         string `json:"role"`
	Password     string `json:"-"` // Omit from JSON output
	PasswordHash string `json:"-"` // Omit from JSON output
}

// HashPassword hashes the user's plain text password and stores it in PasswordHash.
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	u.Password = "" // Clear plain text password after hashing
	return nil
}

// CheckPassword compares a plain text password with the user's stored hashed password.
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}
