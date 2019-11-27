package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"

	logs "github.com/sirupsen/logrus"
)

type emailDetails struct {
	From    string `json:"from"`
	Parse   string `json:"parse"`
	Toemail string `json:"toemail"`
}

func main() {

	logs.WithFields(logs.Fields{
		"package":  "Notification Service",
		"function": "main",
	}).Info("Notification Service started at 7072")

	http.HandleFunc("/sendregisteremail", sendRegisterEmail)
	http.HandleFunc("/sendloginemail", sendLoginEmail)
	http.ListenAndServe("0.0.0.0:7072", nil)
}

// This function reads the json file and pass values to SendNotification
func getCredintials() (string, string) {
	jsonFile, err := os.Open("emailData.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var email emailDetails
	json.Unmarshal(byteValue, &email)
	//log.Println("Received email : " + email.From)

	return email.From, email.Parse

}

func sendRegisterEmail(res http.ResponseWriter, req *http.Request) {

	email := req.FormValue("email")
	// This body value should have a token and it should be inserted to a db
	body := "This is register email"
	from, pass := getCredintials()

	msg := "From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: Register to the system\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{email}, []byte(msg))

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "Notification Service",
			"function": "sendRegisterEmail",
			"error":    err,
		}).Error("SMTP server failure")
		return
	}

	logs.WithFields(logs.Fields{
		"package":  "Notification Service",
		"function": "sendRegisterEmail",
		"email":    email,
	}).Info("Register email sent")

	// This should send true false, to calling function.
	// Eg : function may call for register page or may be for login
}

func sendLoginEmail(res http.ResponseWriter, req *http.Request) {

	email := req.FormValue("email")
	// This body value should have a token and it should be inserted to a db
	body := "This is login email"
	from, pass := getCredintials()

	msg := "From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: You have an account\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{email}, []byte(msg))

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "Notification Service",
			"function": "sendRegisterEmail",
			"error":    err,
		}).Error("SMTP server failure")
		return
	}

	logs.WithFields(logs.Fields{
		"package":  "Notification Service",
		"function": "sendLoginEmail",
		"email":    email,
	}).Info("Login email sent")

}
