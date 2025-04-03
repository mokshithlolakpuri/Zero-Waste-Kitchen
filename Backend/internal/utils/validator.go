package utils

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("password", validatePassword)
	_ = validate.RegisterValidation("notpast", validateNotPastDate)
}

// ValidateStruct validates any struct using validator
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// validatePassword custom validation for password
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

// validateNotPastDate validates that date is not in the past
func validateNotPastDate(fl validator.FieldLevel) bool {
	date := fl.Field().Interface().(time.Time)
	return !date.Before(time.Now())
}

// NormalizeEmail normalizes email addresses
func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

// IsValidUUID checks if string is a valid UUID
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
