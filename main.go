package main

import (
	"log"
	"net/http"
	"regexp"
)

// This is the main module, this should update into an API gateway.
// Initial step, just routing functionality will be used.
// Running on localhost 7070

func main() {

	log.Println("API gateway started at port : 7070")
	http.HandleFunc("/getemail", validatemail)

	http.ListenAndServe(":7070", nil)
}

// This will validate email address has valid syntax
func validatemail(res http.ResponseWriter, req *http.Request) {

	// Check method
	if req.Method != "POST" {
		log.Panic("Email form data is not Post")
		//http.Redirect(res, req, "/", http.StatusSeeOther) // redirect back to register
	}

	email := req.FormValue("email") //"sachithnalaka@gmail.com" // parse form and get email

	validEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`) // regex to validate email address

	if validEmail.MatchString(email) {
		log.Println("Valid email address format received")
		// Send email with token

		send, err := http.Get("http://notification:7072/sendemail")

		if err != nil {
			log.Println("Couldnt send email, notification service sends an error : ", err)
		}

		log.Println("Send an email ", send)

	} else {
		log.Println("Wrong email address format")
		// Return to register window
		//return false
	}
}
func checkEmail(res http.ResponseWriter, req *http.Request) {
	validate, err := http.Get("http://notification:7072")

	if err != nil {
		log.Println("Couldnt send request to add module", err)
	}

	log.Println(validate) // Just to verify we gets a response
}
