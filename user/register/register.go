package register

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// database connection
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/codereview_users")

	if err != nil {
		log.Println("Can not open database connection - User Module")
	}
	return db
}

// Check available user accounts from DB
// If user has an account, send login form sendLoginEmail()
// Else send register form sendRegisterEmail()
func CheckEmail(res http.ResponseWriter, req *http.Request) {

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

	log.Println("Sent Login mail to ", email)

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

// func getemail() {
// 	// This will get the email address from first register window

// 	// call to checkemail
// 	email := "sachithnalaka@gmail.com"
// 	//hasAcct := checkemail(email)

// 	// send login page to email address
// 	if hasAcct {
// 		log.Printf("User has an account for email %s,send login view", email)
// 		// send an email with login page url and you have account
// 		sendLogin, err := http.Get("http://notification:7072/sendloginemail")

// 		if err != nil {
// 			log.Println("Couldnt send email, notification service sends an error : ", err)
// 		}

// 		log.Println("Sent login email ", sendLogin)
// 	} else {
// 		// send a link with token to register

// 		log.Printf("User doesnt have a profile for %s, send register view with token", email)

// 		token := generateToken(email)

// 		log.Println("Token", token) // for error preventing, remove after token being used

// 		// Send email with token

// 		send, err := http.Get("http://notification:7072/sendregisteremail")

// 		if err != nil {
// 			log.Println("Couldnt send email, notification service sends an error : ", err)
// 		}

// 		log.Println("Sent register email ", send)

// 	}

// }

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
		log.Println(err)
	}
	return string(hashedPass)
}
