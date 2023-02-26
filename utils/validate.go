package utils

import (
	"net/mail"
	"strings"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidatePassword(password string) bool {
	return !CheckEmptyString(password)
}

func CheckEmptyString(str string) bool {
	return len(strings.Trim(str, " ")) == 0
}
