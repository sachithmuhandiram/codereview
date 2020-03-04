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

	requestID := req.FormValue("requestID")
	userID := req.FormValue("userid")
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

	_, err = http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID}, "service": {"User Service"},
		"function": {"UserLogin"}, "package": {"Login"}, "status": {"1"}})

	if err != nil {
		log.Println("Error response sending")
	}

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

// GenerateJWT takes eventID as a parameter and time for JWT
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
