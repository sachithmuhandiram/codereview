package register

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

// Check available user accounts from DB
// If user has an account, send login form sendLoginEmail()
// Else send register form sendRegisterEmail()
func CheckEmail(res http.ResponseWriter, req *http.Request) {

	email := req.FormValue("email")
	apiUuid := req.FormValue("uid")

	logs.WithFields(logs.Fields{
		"package":  " Notification Service ",
		"function": " CheckEmail ",
		"email":    email,
		"uuid":     apiUuid,
	}).Info("User Service received email address")

	hasAcct := checkemail(email, apiUuid)

	// Account associates with email
	if hasAcct {
		logs.WithFields(logs.Fields{
			"package":  "Notification Service",
			"function": "CheckEmail",
			"email":    email,
			"uuid":     apiUuid,
		}).Info("This user has an account. send login email")

		sendLoginEmail(email, apiUuid)

	} else {
		logs.WithFields(logs.Fields{
			"package":  "Notification Service",
			"function": "CheckEmail",
			"email":    email,
			"uuid":     apiUuid,
		}).Info("This user doent have an account. send register email")

		sendRegisterEmail(email, apiUuid)
	}

}

// User doesnt have an account, send register form with token
func sendRegisterEmail(email string, apiUuid string) {

	_, err := http.PostForm("http://notification:7072/sendregisteremail", url.Values{"email": {email}, "uuid": {apiUuid}})

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "sendRegisterEmail",
			"error":    err,
			"uuid":     apiUuid,
		}).Error("Failed to connect to Notification Service")
	}

	logs.WithFields(logs.Fields{
		"package":  "Notification Service",
		"function": "sendRegisterEmail",
		"email":    email,
		"uuid":     apiUuid,
	}).Info("Sent registering email to user")

	_, err = http.PostForm("http://localhost:7070/response", url.Values{"uid": {apiUuid}})

	if err != nil {
		log.Println("Error response sending")
	}
}

// User has an account, send login form
func sendLoginEmail(email string, apiUuid string) {

	//_, err := http.Get("http://notification:7072/sendloginemail")
	_, err := http.PostForm("http://notification:7072/sendloginemail", url.Values{"email": {email}, "uuid": {apiUuid}})

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "sendLoginEmail",
			"error":    err,
			"uuid":     apiUuid,
		}).Error("Failed to connect to Notification Service")
	}

	logs.WithFields(logs.Fields{
		"package":  "User Service",
		"function": "sendLoginEmail",
		"email":    email,
		"uuid":     apiUuid,
	}).Info("Sent login email to user")

	_, err = http.PostForm("http://localhost:7070/response", url.Values{"uid": {apiUuid}, "service": {"User Service"},
		"function": {"sendLoginEmail"}, "package": {"Register"}})

	if err != nil {
		log.Println("Error response sending")
	}

}

func checkemail(email string, uuid string) bool {
	// check DB whether we alreayd have a user for this email
	db := dbConn()

	var registeredEmail bool

	// This will return a true or false
	row := db.QueryRow("select exists(select id from emails where email=?)", email)

	err := row.Scan(&registeredEmail)
	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "checkemail",
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

func register() {
	// get user form from user register form
	// insert data to DB

}

// This is a custom token generation. This will use only for initial step
func generateToken(email string) string {
	// This will generate a token to

	bs := []byte(email) // convert email address into a bytestream

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
