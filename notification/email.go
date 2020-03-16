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

// Notifier interface just uses SendEmail function.
type Notifier interface {
	SendEmail(user *userEmailNotification, msg string)
}

type userEmailNotification struct {
	email             string
	requestID         string
	token             string
	emailNotification Notifier
}

type registerEmail struct {
}

type loginEmail struct {
}

type passwordResetEmail struct {
}

type emailDetails struct {
	From    string `json:"from"`
	Parse   string `json:"parse"`
	Toemail string `json:"toemail"`
}

func (user *userEmailNotification) notify(msg string) {
	user.emailNotification.SendEmail(user, msg)
}

func main() {

	logs.WithFields(logs.Fields{
		"package":  "Notification Service",
		"function": "main",
	}).Info("Notification Service started at 7072")

	http.HandleFunc("/sendemail", sendUserEmail)
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

func sendUserEmail(res http.ResponseWriter, req *http.Request) {

	email := req.FormValue("email")
	apiUUID := req.FormValue("uid")
	token := req.FormValue("token")
	notification := req.FormValue("nofitication")
	// This body value should have a token and it should be inserted to a db

	switch notification {
	case "register":
		notification1 := userEmailNotification{
			email:             email,
			requestID:         apiUUID,
			token:             token,
			emailNotification: registerEmail{},
		}

		notification1.notify("This is registering email")
	case "login":
		notification1 := userEmailNotification{
			email:             email,
			requestID:         apiUUID,
			token:             token,
			emailNotification: loginEmail{},
		}

		notification1.notify("This is login email")

	case "passwordreset":
		notification1 := userEmailNotification{
			email:             email,
			requestID:         apiUUID,
			token:             token,
			emailNotification: passwordResetEmail{},
		}

		notification1.notify("This is Password reset email")

	}
}

func (regiEmail registerEmail) SendEmail(user *userEmailNotification, msg string) {
	registerURL := os.Getenv("REGISTERURL") 

	body := msg + "\n" + registerURL + "=" +  user.token
	from, pass := getCredintials()

	emailMsg := "From: " + from + "\n" +
		"To: " + user.email + "\n" +
		"Subject: Register to the system\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{user.email}, []byte(emailMsg))

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "Notification Service",
			"function": "sendRegisterEmail",
			"error":    err,
			"uid":      user.requestID,
		}).Error("SMTP server failure")

		return
	}

	logs.WithFields(logs.Fields{
		"package":  "Notification Service",
		"function": "sendRegisterEmail",
		"email":    user.email,
		"uid":      user.requestID,
	}).Info("Register email sent")
}

// sending login email
func (loginEmail loginEmail) SendEmail(user *userEmailNotification, msg string) {
	
	loginURL := os.Getenv("LOGINURL") 
	body := msg + "\n" + loginURL + "\n This link valid only for 10 minutes"
	from, pass := getCredintials()

	emailMsg := "From: " + from + "\n" +
		"To: " + user.email + "\n" +
		"Subject: You have an account\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{user.email}, []byte(emailMsg))

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "Notification Service",
			"function": "sendRegisterEmail",
			"error":    err,
			"uid":      user.requestID,
		}).Error("SMTP server failure")
		return
	}

	logs.WithFields(logs.Fields{
		"package":  "Notification Service",
		"function": "sendLoginEmail",
		"email":    user.email,
		"uid":      user.requestID,
	}).Info("Login email sent")
}

func (passRestEmail passwordResetEmail) SendEmail(user *userEmailNotification, msg string) {

	body := msg + user.token
	from, pass := getCredintials()

	emailMsg := "From: " + from + "\n" +
		"To: " + user.email + "\n" +
		"Subject: Password reset\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{user.email}, []byte(emailMsg))

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "Notification Service",
			"function": "SendEmail - Password reset",
			"error":    err,
			"uid":      user.requestID,
		}).Error("SMTP server failure")
		return
	}

	logs.WithFields(logs.Fields{
		"package":  "Notification Service",
		"function": "SendEmail - Password reset",
		"email":    user.email,
		"uid":      user.requestID,
	}).Info("Password reset email sent")
}
