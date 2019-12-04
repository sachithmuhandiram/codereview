package main

import (
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/google/uuid"
	logs "github.com/sirupsen/logrus"
)

// This is the main module, this should update into an API gateway.
// Initial step, just routing functionality will be used.
// Running on localhost 7070

// Struct to hold UUID which is attached and passed to all API-Gateway calls
type UUID struct {
	apiUuid uuid.UUID
}

func main() {

	apiID := &UUID{apiUuid: generateUUID()}

	logs.WithFields(logs.Fields{
		"package":  "API-Gateway",
		"function": "main",
		"uuid":     apiID,
	}).Info("API - Gateway started at 7070")

	http.HandleFunc("/getemail", apiID.validatemail)
	http.HandleFunc("/response", reportResponse)

	http.ListenAndServe(":7070", nil)
}

// This will validate email address has valid syntax
func (apiID *UUID) validatemail(res http.ResponseWriter, req *http.Request) {

	validatemailID := apiID.apiUuid
	// Check method
	if req.Method != "POST" {
		logs.WithFields(logs.Fields{
			"package":  "API - Gateway",
			"function": "validatemail",
			"uuid":     validatemailID,
		}).Error("Request method is not POST")

		_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {validatemailID.String()}, "service": {"API Gateway"},
			"function": {"validatemail"}, "package": {"main"}, "status": {"0"}})

		if err != nil {
			log.Println("Error response sending")
		}
		//http.Redirect(res, req, "/", http.StatusSeeOther) // redirect back to register
	} else {
		// Method is POST

		email := req.FormValue("email") //"sachithnalaka@gmail.com" // parse form and get email

		validEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`) // regex to validate email address

		if validEmail.MatchString(email) {
			logs.WithFields(logs.Fields{
				"package":  "API-Gateway",
				"function": "validatemail",
				"uuid":     validatemailID,
			}).Info("Valid email format received")

			logs.WithFields(logs.Fields{
				"package":  "API-Gateway",
				"function": "validatemail",
				"email":    email,
				"uuid":     validatemailID, // Later this should change for function-wise uuid
			}).Info("Email will pass to User - Service")

			_, err := http.PostForm("http://user:7071/checkemail", url.Values{"email": {email}, "uid": {validatemailID.String()}})

			if err != nil {
				logs.WithFields(logs.Fields{
					"package":  "API-Gateway",
					"function": "validatemail",
					"email":    email,
					"error":    err,
					"uuid":     validatemailID,
				}).Error("Error posting data to User - Service")

				_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {validatemailID.String()}, "service": {"API Gateway"},
					"function": {"validatemail"}, "package": {"main"}, "status": {"0"}})

				if err != nil {
					log.Println("Error response sending")
				}
			}

		} else {
			logs.WithFields(logs.Fields{
				"package":  "API-Gateway",
				"function": "validatemail",
				"email":    email,
				"uuid":     validatemailID,
			}).Error("Wrong email format received")

			_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {validatemailID.String()}, "service": {"API Gateway"},
				"function": {"validatemail"}, "package": {"main"}, "status": {"0"}})

			if err != nil {
				log.Println("Error response sending")
			}
			// Return to register window
			//return false
		}
	} // Method checking if loop
}

func generateUUID() uuid.UUID {
	// Generating UUID
	uuid, err := uuid.NewUUID()

	if err != nil {
		if err != nil {
			logs.WithFields(logs.Fields{
				"package":  "API-Gateway",
				"function": "generateUUID",
				"error":    err,
			}).Error("Couldnt generate a UUID")
		}

		// Return to error page
		//http.Redirect(loginResponse, loginRequest, "/", http.StatusSeeOther)
	}
	return uuid //.String()
}

// Response for a request is recorded.
func reportResponse(res http.ResponseWriter, req *http.Request) {

	responseID := req.FormValue("uid")
	service := req.FormValue("service")
	function := req.FormValue("function")
	pack := req.FormValue("package")
	status := req.FormValue("status")

	logs.WithFields(logs.Fields{
		"ResponseService": service,
		"ResponsePackage": pack,
		"ResponseFunc":    function,
		"responseID":      responseID,
		"status":          status,
	}).Info("Response received for the request")
}
