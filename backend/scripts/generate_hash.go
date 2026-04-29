package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "Admin123!"

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		fmt.Println("Errore hashing password:", err)
		return
	}

	fmt.Println(string(hashedPassword))
}

// $2a$10$S.wIY7fVpIIQb0BI7RQGh.ZI8jCH9WXXxQjq.oHVcaB6ZbUGJk6TW
