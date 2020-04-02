package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	logs "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// database connection
// func dbConn() (db *sql.DB) {
// 	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/codereview_users")

// 	if err != nil {
// 		logs.WithFields(logs.Fields{
// 			"Service":  "User Service",
// 			"Package":  "Login",
// 			"function": "dbConn",
// 			"error":    err,
// 		}).Error("Failed to connect to database")
// 	}
// 	return db
// }

func UserLogin(res http.ResponseWriter, req *http.Request) {

	requestID := req.FormValue("uid")
	userID := req.FormValue("userid")
	loginToken := req.FormValue("logintoken")
	password := req.FormValue("password")

	var userpassword string // to use with select table

	logs.WithFields(logs.Fields{
		"Service":   "User Service",
		"Package":   "Login",
		"function":  "UserLogin",
		"userid":    userID,
		"requestID": requestID,
	}).Info("Login request received")

	db := dbConn()

	// if user form doesnt have a logintoken then it rejects:
	validToken := checkLoginToken(requestID,loginToken)

	if validToken == false  {

		logs.WithFields(logs.Fields{
						"Service":   "User Service",
						"Package":   "Login",
						"function":  "UserLogin",
						"userid":    userID,
						"requestID": requestID,
	}).Warn("Login request does not have a login token")

	// response is sent to login page again , for now its status unauthorize
	res.WriteHeader(http.StatusUnauthorized)
	res.Write([]byte("500 - Something bad happened!"))

	_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID}, "service": {"User Service"},
			"function": {"UserLogin"}, "package": {"Login"}, "status": {"0"}})

		if err != nil {
			log.Println("Error response sending")
		}
		
	return

	}

	// else goes to password matching

	//getting user password from users table
	row := db.QueryRow("select password from users where email=?", userID)
	err := row.Scan(&userpassword)

	if err != nil {
		if err == sql.ErrNoRows {
			logs.WithFields(logs.Fields{
				"Service":   "User Service",
				"Package":   "Login",
				"function":  "UserLogin",
				"userid":    userID,
				"requestID": requestID,
			}).Error("No record available for the user")
		} else {
			logs.WithFields(logs.Fields{
				"Service":   "User Service",
				"Package":   "Login",
				"function":  "UserLogin",
				"userid":    userID,
				"requestID": requestID,
			}).Error("Couldnt fetch users table")
		}
	} // querying database table if

	comparePassword := bcrypt.CompareHashAndPassword([]byte(userpassword), []byte(password))

	if comparePassword != nil {
		logs.WithFields(logs.Fields{
			"Service":   "User Service",
			"Package":   "Login",
			"function":  "UserLogin",
			"userid":    userID,
			"requestID": requestID,
		}).Error("Passwords do not match")
		defer db.Close()

		_, err = http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID}, "service": {"User Service"},
			"function": {"UserLogin"}, "package": {"Login"}, "status": {"0"}})

		if err != nil {
			log.Println("Error response sending")
		}

		return

	}
	logs.WithFields(logs.Fields{
		"Service":   "User Service",
		"Package":   "Login",
		"function":  "UserLogin",
		"userid":    userID,
		"requestID": requestID,
	}).Info("Passwords match")

	// need to redirect tp /home
	jwtToken,jwtErr := GenerateJWT(loginToken,3600) // token valid for an hour

		if jwtErr != nil{
			logs.WithFields(logs.Fields{
			"Service":   "User Service",
			"Package":   "Login",
			"function":  "UserLogin",
			"userid":    userID,
			"requestID": requestID,
		}).Error("Generating jwt failed")

		}

	validTkn := insertToValidToken(userID,jwtToken,requestID)

	if validTkn != nil{
	
			logs.WithFields(logs.Fields{
			"Service":   "User Service",
			"Package":   "Login",
			"function":  "UserLogin",
			"userid":    userID,
			"requestID": requestID,
		}).Error("Couldnt insert JWT to table")
	}

	http.Redirect(res, req, "/", 301)
	/*
		Also to insert to db, logged user, password expire and token.
	*/

	_, err = http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID}, "service": {"User Service"},
		"function": {"UserLogin"}, "package": {"Login"}, "status": {"1"}})

	if err != nil {
		log.Println("Error response sending")
	}

	defer db.Close()

}

// Insert into valid token
func insertToValidToken(userID,jwtToken,requestID string) error{

	t := time.Now()
	// t.Format("20060102150405")
	db := dbConn()
	insertToken, err := db.Prepare("INSERT INTO activeJWTtokens(user_id,jwt,created_at,last_update) VALUES(?,?,?,?)")
        if err != nil {
            panic(err.Error())

            return err
        }

    insertToken.Exec(userID, jwtToken,t.Format("yyyy-MM-dd HH:mm:ss"),t.Format("yyyy-MM-dd HH:mm:ss"))

    logs.WithFields(logs.Fields{
		"Service":   "User Service",
		"Package":   "Login",
		"function":  "insertToValidToken",
		"userid":    userID,
		"requestID": requestID,
	}).Info("Insert into valid jwt tokens")

	defer db.Close()
	return nil  

}

// Checks login form's token.
func checkLoginToken(requestID,loginToken string) bool{
	var logintoken bool

	db := dbConn()
	row := db.QueryRow("SELECT EXISTS(SELECT login_token FROM login_token WHERE login_token=?)", loginToken)

	err := row.Scan(&logintoken)
	if err != nil {
		logs.WithFields(logs.Fields{
			"Service":   "User Service",
			"Package":   "Login",
			"function":  "checkLoginToken",
			"requestID": requestID,
			"error"		: err,
		}).Error("Failed to fetch data from login token table")
	}

	if logintoken {

		logs.WithFields(logs.Fields{
			"Service":   "User Service",
			"Package":   "Login",
			"function":  "checkLoginToken",
			"requestID": requestID,
		}).Info("Valid login token")

		
		/*
						Delete Valid Token
			For simplicity of the system, assums that this will delete that token from table
		*/
		deleteLoginToken(loginToken,requestID)

		return true
	}

	defer db.Close()
	return false

}

// Delete login token
func deleteLoginToken(loginToken,requestID string){
	db := dbConn()

	delTkn, err := db.Prepare("DELETE FROM login_token WHERE login_token=?")
    if err != nil {
        panic(err.Error())
    }
    delTkn.Exec(loginToken)

    logs.WithFields(logs.Fields{
			"Service":   "User Service",
			"Package":   "Login",
			"function":  "deleteLoginToken",
			"requestID": requestID,
		}).Info("Successfully deleted login token")

	defer db.Close()
}

func CheckUserLogin(res http.ResponseWriter, req *http.Request) {

	// get jwt and check with valid jwt table
	jwt := "get_jwt_from_cookies"

	validJWT := checkJWT(jwt)
	if validJWT != true {
		return
	}

}

func checkJWT(jwt string) bool {
	// read data from availableJWTtoken table
	return true

}

// GenerateJWT takes eventID as a parameter and time (minutes) for JWT
func GenerateJWT(initialToken string, validDuration int) (string, error) {

	loginKey := []byte(initialToken)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(validDuration))

	jwtToken, jwtErr := token.SignedString(loginKey)

	if jwtErr != nil {
		log.Println("Error creating jwt Token : ", jwtErr)
		return "", jwtErr
	}

	return jwtToken, nil
}
