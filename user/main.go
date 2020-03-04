package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	//checkemail "./checkemail"
	//login "./login"
	//passwordreset "./passwordreset"
	//register "./register"
	logs "github.com/sirupsen/logrus"
)

func main() {

	f, err := os.OpenFile("logs/usermodule.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	log.Println("User Service started")

	// web server
	http.HandleFunc("/checkemail", CheckEmail)
	http.HandleFunc("/register", UserRegister)
	http.HandleFunc("/login", UserLogin)
	http.HandleFunc("/loginWithJWT", CheckUserLogin)
	http.HandleFunc("/updateaccount", checkUpdateRequest)
	http.ListenAndServe("0.0.0.0:7071", nil)

	// register email is sent to email sytax check, if its true,
	// go to register module
}

func checkUpdateRequest(res http.ResponseWriter, req *http.Request) {

	requestID := req.FormValue("uid")
	request := req.FormValue("request")
	token := req.FormValue("token")

	logs.WithFields(logs.Fields{
		"package":  "User Service",
		"function": "checkUpdateRequest",
		"ApiUUID":  requestID,
	}).Info("Update password form received to User service")

	switch request {
	case "updatepassword":
		validToken := CheckUpdatePasswordToken(requestID, token)

		if validToken != true {
			logs.WithFields(logs.Fields{
				"package":  "User Service",
				"function": "checkUpdateRequest",
				"ApiUUID":  requestID,
			}).Warning("Token is not valid for updatepassword")

			_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID}, "service": {"User Service"},
				"function": {"checkUpdateRequest"}, "package": {"Password Reset"}, "status": {"0"}})

			if err != nil {
				log.Println("Error response sending")
			}

			return
		} // token is not valid

		password := req.FormValue("password")
		updatePassword := UpdatePassword(requestID, token, password)

		if updatePassword != true {
			logs.WithFields(logs.Fields{
				"package":  "User Service",
				"function": "checkUpdateRequest",
				"ApiUUID":  requestID,
			}).Error("Could not updatepassword in User module")

			_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID}, "service": {"User Service"},
				"function": {"checkUpdateRequest"}, "package": {"Password Reset"}, "status": {"0"}})

			if err != nil {
				log.Println("Error response sending")
			}
			return
		} // password update request faild to update table

		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "checkUpdateRequest",
			"ApiUUID":  requestID,
		}).Info("Password successfully updated")

		_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID}, "service": {"User Service"},
			"function": {"checkUpdateRequest"}, "package": {"Password Reset"}, "status": {"0"}})

		if err != nil {
			log.Println("Error response sending")
		}
	} // updatepassword case

}
