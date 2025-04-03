package models

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

const (
	saltLength = 16
	timeCost   = 3         // Increased from 1 for better security
	memoryCost = 64 * 1024 // 64 MB
	threads    = 4
	keyLength  = 32
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	FCMToken  string    `json:"fcm_token"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"` // New field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HashPassword hashes the user's password using Argon2
func (u *User) HashPassword() error {
	if len(u.Password) == 0 {
		return errors.New("password cannot be empty")
	}

	// Additional password complexity checks
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Generate a random salt
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("failed to generate salt: %v", err)
	}

	// Hash the password with Argon2
	hashedPassword := argon2.IDKey([]byte(u.Password), salt, timeCost, memoryCost, threads, keyLength)

	// Store the salt and hashed password as a single string
	u.Password = fmt.Sprintf("%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hashedPassword))
	return nil
}

// ComparePassword compares the provided password with the stored hash
func (u *User) ComparePassword(password string) bool {
	if len(u.Password) == 0 || len(password) == 0 {
		return false
	}

	// Split the stored password into salt and hash
	parts := strings.Split(u.Password, "$")
	if len(parts) != 2 {
		return false
	}

	// Decode the salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	// Hash the provided password with the same salt
	computedHash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, threads, keyLength)

	// Compare the computed hash with the stored hash
	return subtle.ConstantTimeCompare(storedHash, computedHash) == 1
}
