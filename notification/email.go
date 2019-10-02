package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type emailDetails struct {
	From    string `json:"from"`
	Parse   string `json:"parse"`
	Toemail string `json:"toemail"`
}

func main() {

	log.Println("Email service started")

	http.HandleFunc("/sendemail", sendEmail)
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
	log.Println("Received email : " + email.From)

	return email.From, email.Parse

}

func sendEmail(res http.ResponseWriter, req *http.Request) {

	email := "sachithnalaka@gmail.com" // this is taken from request
	// This body value should have a token and it should be inserted to a db
	body := "hi hi"
	from, pass := getCredintials()

	msg := "From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{email}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("Email has been sent to : ", email)

	// This should send true false, to calling function.
	// Eg : function may call for register page or may be for login
}
