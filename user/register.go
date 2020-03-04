package main

import (
	"log"
	"net/http"
	"net/url"

	//checkemail "../"
	logs "github.com/sirupsen/logrus"
)

// database connection
// func dbConn() (db *sql.DB) {
// 	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/codereview_users")

// 	if err != nil {
// 		logs.WithFields(logs.Fields{
// 			"Service":  "User Service",
// 			"Package":  "register",
// 			"function": "dbConn",
// 			"error":    err,
// 		}).Error("Failed to connect to database")
// 	}
// 	return db
// }

// UserRegister function just insert new user to users table
func UserRegister(res http.ResponseWriter, req *http.Request) {

	// get user form from user register form
	// insert data to DB
	// First step would be Firstname, lastname and password..
	/*
	* encrypting password from frontend and decrypt at this end...
	* Password matching ( re entering)
	* Inserting to db ( firstname,lastname,email,password,registered_at)
	 */

	requestID := req.FormValue("uid")
	firstName := req.FormValue("first_name")
	lastName := req.FormValue("last_name")
	email := req.FormValue("email")
	password := req.FormValue("password")

	logs.WithFields(logs.Fields{
		"Service":  "User Service",
		"package":  "register",
		"function": "UserRegister",
		"uuid":     requestID,
		"email":    email,
	}).Info("Received data to insert to users table")

	// check user entered same email address
	hasAccount := Checkmail(email, requestID)

	if hasAccount != true {

		db := dbConn()

		// Inserting token to login_token table
		insertUser, err := db.Prepare("INSERT INTO users (email,first_name,last_name,password) VALUES(?,?,?,?)")
		if err != nil {
			logs.WithFields(logs.Fields{
				"Service":  "User Service",
				"package":  "register",
				"function": "UserRegister",
				"uuid":     requestID,
				"Error":    err,
			}).Error("Couldnt prepare insert statement for users table")
		}
		insertUser.Exec(email, firstName, lastName, password)

		// Inserting email to emails table

		insertEmail, err := db.Prepare("INSERT INTO emails (email,isActive) VALUES(?,?)")
		if err != nil {
			logs.WithFields(logs.Fields{
				"Service":  "User Service",
				"package":  "register",
				"function": "UserRegister",
				"uuid":     requestID,
				"Error":    err,
			}).Error("Couldnt prepare insert statement for emails table")
		}
		insertEmail.Exec(email, 1)

		_, err = http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID}, "service": {"User Service"},
			"function": {"UserRegister"}, "package": {"Register"}, "status": {"1"}})

		if err != nil {
			log.Println("Error response sending")
		}

		defer db.Close()
		return
	} // user has an account

	logs.WithFields(logs.Fields{
		"Service":  "User Service",
		"package":  "register",
		"function": "UserRegister",
		"uuid":     requestID,
		"email":    email,
	}).Error("User has an account for this email")

	_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID}, "service": {"User Service"},
		"function": {"sendLoginEmail"}, "package": {"Check Email"}, "status": {"0"}})

	if err != nil {
		log.Println("Error response sending")
	}
}
