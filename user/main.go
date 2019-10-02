package main

import (
	"log"
	"net/http"
)

func main() {

	log.Println("User Service started")

	// web server
	http.ListenAndServe("0.0.0.0:7071", nil)

	// register email is sent to email sytax check, if its true,
	// go to register module
}
