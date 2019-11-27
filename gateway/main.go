package main

import (
	"net/http"
	"net/url"
	"regexp"

	logs "github.com/sirupsen/logrus"
)

// This is the main module, this should update into an API gateway.
// Initial step, just routing functionality will be used.
// Running on localhost 7070

func main() {

	logs.WithFields(logs.Fields{
		"package":  "API-Gateway",
		"function": "main",
	}).Info("API - Gateway started at 7070")

	http.HandleFunc("/getemail", validatemail)

	http.ListenAndServe(":7070", nil)
}

// This will validate email address has valid syntax
func validatemail(res http.ResponseWriter, req *http.Request) {

	// Check method
	if req.Method != "POST" {
		logs.WithFields(logs.Fields{
			"package":  "API - Gateway",
			"function": "validatemail",
		}).Error("Request method is not POST")
		//http.Redirect(res, req, "/", http.StatusSeeOther) // redirect back to register
	}

	email := req.FormValue("email") //"sachithnalaka@gmail.com" // parse form and get email

	validEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`) // regex to validate email address

	if validEmail.MatchString(email) {
		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "validatemail",
		}).Info("Valida email format received")

		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "validatemail",
			"email":    email,
		}).Info("Email will pass to User - Service")

		_, err := http.PostForm("http://user:7071/checkemail", url.Values{"email": {email}})

		if err != nil {
			logs.WithFields(logs.Fields{
				"package":  "API-Gateway",
				"function": "validatemail",
				"email":    email,
				"error":    err,
			}).Error("Error posting data to User - Service")
		}

	} else {
		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "validatemail",
			"email":    email,
		}).Error("Wrong email format received")
		// Return to register window
		//return false
	}
}

// func checkEmail(res http.ResponseWriter, req *http.Request) {
// 	validate, err := http.Get("http://notification:7072")

// 	if err != nil {
// 		log.Println("Couldnt send request to add module", err)
// 	}

// 	log.Println(validate) // Just to verify we gets a response
// }
