package main

import (
	"io"
	"log"
	"net/http"
	"os"

	checkemail "./checkemail"
	login "./login"
	register "./register"
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
	http.ListenAndServe("0.0.0.0:7071", nil)

	// register email is sent to email sytax check, if its true,
	// go to register module
}
