package main

import (
	"log"
	"regexp"
)

func main() {

	log.Println("User Service started")

	// web server

	// register email is sent to email sytax check, if its true,
	// go to register module
}

// email sytax validator
func validateEmail() bool {
	// This will validate email address has valid syntax

	email := "sachithnalaka@gmail.com" // parse form and get email

	validEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`) // regex to validate email address

	if validEmail.MatchString(email) {
		log.Println("Valid email address format received")
		return true
	}

	log.Println("Wrong email address format")
	return false
}
