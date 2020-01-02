package main

import (
	"io"
	"log"
	"net/http"
	"os"

	checkemail "./checkemail"
	login "./login"
	passwordreset "./passwordreset"
	register "./register"
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
	http.HandleFunc("/checkemail", checkemail.CheckEmail)
	http.HandleFunc("/register", register.UserRegister)
	http.HandleFunc("/login", login.UserLogin)
	http.HandleFunc("/loginWithJWT", login.CheckUserLogin)
	http.HandleFunc("/updateaccount", checkUpdateRequest)
	http.ListenAndServe("0.0.0.0:7071", nil)

	// register email is sent to email sytax check, if its true,
	// go to register module
}

func checkUpdateRequest(res http.ResponseWriter, req *http.Request) {

	requestID := req.FormValue("uid")
	request := req.FormValue("request")
	token := req.FormValue("token")
	password := req.FormValue("password")

	logs.WithFields(logs.Fields{
		"package":  "User Service",
		"function": "checkUpdateRequest",
		"ApiUUID":  requestID,
	}).Info("Update password received to User service")

	switch request {
	case "updatepassword":
		validToken := passwordreset.CheckUpdatePasswordToken(requestID, token)

		if validToken != true {
			logs.WithFields(logs.Fields{
				"package":  "User Service",
				"function": "checkUpdateRequest",
				"ApiUUID":  requestID,
			}).Warning("Token is not valid for updatepassword")

			return
		}

		updatePassword := passwordreset.UpdatePassword(requestID, token, password)

		if updatePassword != true {
			logs.WithFields(logs.Fields{
				"package":  "User Service",
				"function": "checkUpdateRequest",
				"ApiUUID":  requestID,
			}).Error("Could not updatepassword in User module")

			return
		}
	} // updatepassword case

}
