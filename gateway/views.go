package main

import(
	"net/http"
)
func registerView(res http.ResponseWriter,req *http.Request){
	http.ServeFile(res, req,"views/register.html")
}

func resetPasswordView(res http.ResponseWriter,req *http.Request){
	
}

func getEmailView(res http.ResponseWriter,req *http.Request){
	http.ServeFile(res, req,"views/getEmail.html")
}