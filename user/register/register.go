package register

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func getemail() {
	// This will get the email address from first register window

	// call to checkemail
	email := "sachithnalaka@gmail.com"
	hasAcct := checkemail(email)

	// send login page to email address
	if hasAcct {
		log.Printf("User has an account for email %s,send login view", email)
		// send an email with login page url and you have account
	} else {
		// send a link with token to register

		log.Printf("User doesnt have a profile for %s, send register view with token", email)

		token := generateToken(email)

		log.Println("Token", token) // for error preventing, remove after token being used
		// send register form

		// Send email with token

		send, err := http.Get("http://notification:7072/sendemail")

		if err != nil {
			log.Println("Couldnt send email, notification service sends an error : ", err)
		}

		log.Println("Send an email ", send)

	}

}

func checkemail(email string) bool {
	// check DB whether we alreayd have a user for this email

	// true
	// false (no account)

	return false
}

func register() {
	// get user form from user register form
	// insert data to DB

}

// This is a custom token generation. This will use only for initial step
func generateToken(email string) string {
	// This will generate a token to

	bs := []byte(email) // convert email address into a bytestream

	hashedPass, err := bcrypt.GenerateFromPassword(bs, 8)

	if err != nil {
		log.Println(err)
	}
	return string(hashedPass)
}
