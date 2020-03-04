package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	logs "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// database connection
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/codereview_users")

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "dbConn",
			"error":    err,
		}).Error("Failed to connect to database")
	}
	return db
}

// Function to modify as generic function
func ValidateEmail(res http.ResponseWriter, req *http.Request) {
	email := req.FormValue("email")
	apiUuid := req.FormValue("uid")

	logs.WithFields(logs.Fields{
		"package":  " User Service ",
		"function": " ValidateEmail ",
		"email":    email,
		"uuid":     apiUuid,
	}).Info("Validate email function received email address")

}

// Check available user accounts from DB
// If user has an account, send login form sendLoginEmail()
// Else send register form sendRegisterEmail()
func CheckEmail(res http.ResponseWriter, req *http.Request) {

	email := req.FormValue("email")
	apiUUID := req.FormValue("uid")
	request := req.FormValue("request")

	logs.WithFields(logs.Fields{
		"package":  " User Service ",
		"function": " CheckEmail ",
		"email":    email,
		"uuid":     apiUUID,
	}).Info("User Service received email address")

	switch request {

	case "hasaccount":
		hasAccount(email, apiUUID)
	case "passwordreset":
		passwordReset(email, apiUUID)

	} // switch statement

}

// Check whether user has an account for the given email
// Sends login or registering emails
func hasAccount(email, apiUUID string) {
	logs.WithFields(logs.Fields{
		"package":  " User Service ",
		"function": " hasAccount ",
		"email":    email,
		"uuid":     apiUUID,
	}).Info("User email address is being processed")

	hasAcct := Checkmail(email, apiUUID)

	// Account associates with email
	if hasAcct {
		logs.WithFields(logs.Fields{
			"package":  "Notification Service",
			"function": "CheckEmail",
			"email":    email,
			"uuid":     apiUUID,
		}).Info("This user has an account. send login email")

		sendLoginEmail(email, apiUUID)

	} else {
		logs.WithFields(logs.Fields{
			"package":  "Notification Service",
			"function": "CheckEmail",
			"email":    email,
			"uuid":     apiUUID,
		}).Info("This user doent have an account. send register email")

		sendRegisterEmail(email, apiUUID)
	}

}

// User doesnt have an account, send register form with token
func sendRegisterEmail(email string, apiUUID string) {
	db := dbConn()
	token := generateToken(apiUUID)

	// insert token to registering_token table
	insToken, err := db.Prepare("INSERT INTO registering_token (reg_token) VALUES(?)")
	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "Notification Service",
			"function": "sendRegisterEmail",
			"uuid":     apiUUID,
			"Error":    err,
		}).Error("Couldnt prepare insert statement for registering_token table")
	}
	insToken.Exec(token) //time.Now()

	// posting form to notification service
	_, err = http.PostForm("http://notification:7072/sendemail", url.Values{"email": {email}, "uuid": {apiUUID}, "token": {token}, "nofitication": {"register"}})

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "sendRegisterEmail",
			"error":    err,
			"uuid":     apiUUID,
		}).Error("Failed to connect to Notification Service")
	}

	// This is mistake, I should take response from notification service and then send the response to API gateway
	logs.WithFields(logs.Fields{
		"package":  "Notification Service",
		"function": "sendRegisterEmail",
		"email":    email,
		"uuid":     apiUUID,
	}).Info("Sent registering email to user")

	_, err = http.PostForm("http://localhost:7070/response", url.Values{"uid": {apiUUID}, "service": {"User Service"},
		"function": {"sendRegisterEmail"}, "package": {"Register"}, "status": {"1"}})

	if err != nil {
		log.Println("Error response sending")
	}
}

// User has an account, send login form
func sendLoginEmail(email string, apiUUID string) {
	db := dbConn()
	token := generateToken(apiUUID)

	// Inserting token to login_token table
	insToken, err := db.Prepare("INSERT INTO login_token (login_token) VALUES(?)")
	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "Notification Service",
			"function": "sendLoginEmail",
			"uuid":     apiUUID,
			"Error":    err,
		}).Error("Couldnt prepare insert statement for login token table")
	}
	insToken.Exec(token) //, time.Now()
	// Sending login form to notification service
	_, err = http.PostForm("http://notification:7072/sendemail", url.Values{"email": {email}, "uuid": {apiUUID}, "token": {token}, "nofitication": {"login"}})

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "sendLoginEmail",
			"error":    err,
			"uuid":     apiUUID,
		}).Error("Failed to connect to Notification Service")
	}

	logs.WithFields(logs.Fields{
		"package":  "User Service",
		"function": "sendLoginEmail",
		"email":    email,
		"uuid":     apiUUID,
	}).Info("Sent login email to user")

	_, err = http.PostForm("http://localhost:7070/response", url.Values{"uid": {apiUUID}, "service": {"User Service"},
		"function": {"sendLoginEmail"}, "package": {"Check Email"}, "status": {"1"}})

	if err != nil {
		log.Println("Error response sending")
	}

}

func Checkmail(email string, uuid string) bool {
	// check DB whether we alreayd have a user for this email
	db := dbConn()

	var registeredEmail bool

	// This will return a true or false
	row := db.QueryRow("select exists(select id from emails where email=?)", email)

	err := row.Scan(&registeredEmail)
	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "CheckEmail",
			"error":    err,
			"uuid":     uuid,
		}).Error("Failed to fetch data from user table")
	}

	if registeredEmail {
		return true
	}

	defer db.Close()
	return false

}

// Password reset
func passwordReset(email, apiUUID string) {
	validEmail := Checkmail(email, apiUUID)

	if validEmail != true {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "passwordReset",
			"uuid":     apiUUID,
		}).Info("Email does not associate with an account")

		_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {apiUUID}, "service": {"User Service"},
			"function": {"passwordReset"}, "package": {"Check Email"}, "status": {"0"}})

		if err != nil {
			log.Println("Error response sending")
		}
		return
	}

	logs.WithFields(logs.Fields{
		"package":  "User Service",
		"function": "passwordReset",
		"uuid":     apiUUID,
	}).Info("This user has an account")

	token := generateToken(apiUUID)

	// Insert this token to "passwordResetToken" table
	insertedToTable := insertPasswordResetToken(email, apiUUID, token)
	// if its inserted to table, then send email

	if insertedToTable != true {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "passwordReset",
			"uuid":     apiUUID,
		}).Error("Failed to insert token to passwordResetToken table")

		return
	}
	_, err := http.PostForm("http://notification:7072/sendemail", url.Values{"email": {email}, "uuid": {apiUUID}, "token": {token}, "nofitication": {"passwordreset"}})

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "passwordReset",
			"error":    err,
			"uuid":     apiUUID,
		}).Error("Failed to connect to Notification Service")
	}
}

func insertPasswordResetToken(email, uuid, token string) bool {

	db := dbConn()

	logs.WithFields(logs.Fields{
		"package":  "User Service",
		"function": "insertPasswordResetToken",
		"uuid":     uuid,
	}).Info("Password reset token insert to passwordResetToken request received")

	insToken, err := db.Prepare("INSERT INTO passwordResetToken (email,resettoken) VALUES(?,?)")
	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "insertPasswordResetToken",
			"uuid":     uuid,
			"Error":    err,
		}).Error("Couldnt prepare insert statement for passwordResetToken table")
		defer db.Close()
		return false
	}
	insToken.Exec(email, token)

	logs.WithFields(logs.Fields{
		"package":  "User Service",
		"function": "insertPasswordResetToken",
		"uuid":     uuid,
	}).Info("Password reset token inserted to passwordResetToken table")
	return true
}

// This is a custom token generation. This will use only for initial step
func generateToken(uuid string) string {
	// This will generate a token to

	bs := []byte(uuid) // convert UUID into a bytestream

	hashedPass, err := bcrypt.GenerateFromPassword(bs, 8)

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "generateToken",
			"error":    err,
		}).Error("Failed to generate a token")
	}
	return string(hashedPass)
}

// This is a comment for testing