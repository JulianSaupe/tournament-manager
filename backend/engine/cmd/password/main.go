package main

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/big"
)

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars  = "0123456789"
	specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	passwordLen  = 12
)

func generatePassword() (string, error) {
	chars := lowerChars + upperChars + numberChars + specialChars
	password := make([]byte, passwordLen)

	for i := 0; i < passwordLen; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %v", err)
		}
		password[i] = chars[num.Int64()]
	}

	return string(password), nil
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedBytes), nil
}

func main() {
	//password, err := generatePassword()
	//if err != nil {
	//	fmt.Printf("Error generating password: %v\n", err)
	//	return
	//}

	password := "1234"

	hashedPassword, err := hashPassword(password)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		return
	}

	fmt.Printf("Generated password: %s\n", password)
	fmt.Printf("Hashed password: %s\n", hashedPassword)
}
