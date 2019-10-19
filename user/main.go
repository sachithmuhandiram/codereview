package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
)

// database connection
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/codereview_users")

	if err != nil {
		log.Println("Can not open database connection - User Module")
	}
	return db
}

func main() {

	log.Println("User Service started")

	// web server
	http.HandleFunc("/checkemail", checkEmail)
	http.ListenAndServe("0.0.0.0:7071", nil)

	// register email is sent to email sytax check, if its true,
	// go to register module
}

// Check available user accounts from DB
// If user has an account, send login form sendLoginEmail()
// Else send register form sendRegisterEmail()
func checkEmail(res http.ResponseWriter, req *http.Request) {

	email := req.FormValue("email")

	log.Println("User module received email address : ", email)

	hasAcct := checkemail(email)

	// Account associates with email
	if hasAcct {
		log.Println("User has an account for ", email)
		sendLoginEmail(email)

	} else {
		log.Println("Email is not accociate with an account")
		sendRegisterEmail(email)
	}

}

// User doesnt have an account, send register form with token
func sendRegisterEmail(email string) {

	_, err := http.PostForm("http://notification:7072/sendregisteremail", url.Values{"email": {email}})

	if err != nil {
		log.Println("Couldnt send register email, notification service sends an error : ", err)
	}

	log.Println("Sent register mail to ", email)
}

// User has an account, send login form
func sendLoginEmail(email string) {

	//_, err := http.Get("http://notification:7072/sendloginemail")
	_, err := http.PostForm("http://notification:7072/sendloginemail", url.Values{"email": {email}})

	if err != nil {
		log.Println("Couldnt send login email, notification service sends an error : ", err)
	}

	log.Println("Sent mail to ", email)

}

func checkemail(email string) bool {
	// check DB whether we alreayd have a user for this email
	db := dbConn()

	var registeredEmail bool

	// This will return a true or false
	row := db.QueryRow("select exists(select id from emails where email=?)", email)

	err := row.Scan(&registeredEmail)
	if err != nil {
		log.Println("Error fetching data from emails table")
	}

	if registeredEmail {
		return true
	}

	defer db.Close()
	return false

}
