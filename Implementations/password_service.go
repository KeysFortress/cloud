package implementations

import (
	"fmt"
	"math/rand"
	"time"
)

type PasswordService struct {
}

const (
	lowerCaseLetters  = "abcdefghijklmnopqrstuvwxyz"
	upperCaseLetters  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharacters = "!@#$%^&*()-_=+,.?/:;{}[]~"
)

func (ps *PasswordService) GeneratePassword(length int, lowerCase, upperCase, unique, specialChars bool) (string, error) {
	if length < 16 || length > 512 {
		return "", fmt.Errorf("length must be between 16 and 512")
	}

	var validChars string
	if lowerCase {
		validChars += lowerCaseLetters
	}
	if upperCase {
		validChars += upperCaseLetters
	}
	if specialChars {
		validChars += specialCharacters
	}

	rand.NewSource(time.Now().UnixNano())

	password := make([]byte, length)
	for i := range password {
		password[i] = validChars[rand.Intn(len(validChars))]
	}

	if unique {
		password = ps.shuffle(password)
	}

	return string(password), nil
}

func (pr *PasswordService) shuffle(s []byte) []byte {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	return s
}
