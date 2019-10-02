package main

import (
	"log"
	"net/http"
)

func main() {

	log.Println("User Service started")

	// web server
	http.HandleFunc("/sendmail", sendEmail)
	http.ListenAndServe("0.0.0.0:7071", nil)

	// register email is sent to email sytax check, if its true,
	// go to register module
}

func sendEmail(res http.ResponseWriter, req *http.Request) {

	email, err := http.Get("http://notification:7072/sendemail")

	if err != nil {
		log.Println("Couldnt send email, notification service sends an error : ", err)
	}

	log.Println("Send an email ", email)

}
