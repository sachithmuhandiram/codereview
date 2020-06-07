package main

import(
	"net/http"
	"html/template"
)
func registerView(res http.ResponseWriter,req *http.Request){
	http.ServeFile(res, req,"views/register.html")
}

func resetPasswordView(res http.ResponseWriter,req *http.Request){
	
}

func getEmailView(res http.ResponseWriter,req *http.Request){
	http.ServeFile(res, req,"views/getEmail.html")
}

func codeSubmitView(res http.ResponseWriter, req *http.Request){

	t, err := template.ParseFiles("views/codesubmit.html")
	if err != nil {
		log.Fatalln(err)
	}

	colours := map[string]string{
		"1": "blue",
		"2": "red",
		"3": "yellow",
	}

	err = t.Execute(res, colours)

	if err != nil {
		log.Fatalln(err)
	}
}