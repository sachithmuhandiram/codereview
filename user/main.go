package main

import (
	"log"
	"net/http"
)

func main() {

	log.Println("User Service started")

	// web server
	http.HandleFunc("/checkemail", checkEmail)
	http.ListenAndServe("0.0.0.0:7071", nil)

	// register email is sent to email sytax check, if its true,
	// go to register module
}

// Check available user accounts from DB
// If user has an account, send login form sendLoginEmail()
// Else send register form sendRegisterEmail()
func checkEmail(res http.ResponseWriter, req *http.Request) {

}

// User doesnt have an account, send register form with token
func sendRegisterEmail() {

	email, err := http.Get("http://notification:7072/sendregisteremail")

	if err != nil {
		log.Println("Couldnt send register email, notification service sends an error : ", err)
	}

	log.Println("Sending register mail success ", email)
}

// User has an account, send login form
func sendLoginEmail() {

	email, err := http.Get("http://notification:7072/sendloginemail")

	if err != nil {
		log.Println("Couldnt send login email, notification service sends an error : ", err)
	}

	log.Println("Sending login mail success ", email)

}
