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

func UserLogin(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	requestID := req.URL.Query().Get("uid")
	userID := req.URL.Query().Get("userid")//req.FormValue("userid")
	password := req.FormValue("password")

	logs.WithFields(logs.Fields{
		"Service":   "User Service",
		"Package":   "Login",
		"function":  "UserLogin",
		"userid":    userID,
		"requestID": requestID,
	}).Info("Login request received")

	db := dbConn()

	// compare password
	passwordMatch := comparePassword(requestID,userID,password)

	if passwordMatch == true{

		logs.WithFields(logs.Fields{
			"Service":   "User Service",
			"Package":   "Login",
			"function":  "UserLogin",
			"userid":    userID,
			"requestID": requestID,
		}).Info("Passwords match")

		http.Redirect(res, req, "http://localhost:7070/createsession?userid="+userID+"&authorize="+"1&uid="+requestID, http.StatusSeeOther)

		// send response to /gateway respose
		_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID}, "service": {"User Service"},
			"function": {"UserLogin"}, "package": {"Login"}, "status": {"1"}})

		if err != nil {
			log.Println("Error response sending")
		}

	}else{ // password do not match

		logs.WithFields(logs.Fields{
			"Service":   "User Service",
			"Package":   "Login",
			"function":  "UserLogin",
			"userid":    userID,
			"requestID": requestID,
		}).Error("Passwords do not match")


		_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID}, "service": {"User Service"},
			"function": {"UserLogin"}, "package": {"Login"}, "status": {"0"}})

		if err != nil {
			log.Println("Error response sending")
		}
		http.Redirect(res,req,"localhost:7070/login",http.StatusSeeOther)
		return

	} // password do not match

	defer db.Close()
}
// Password comparision
func comparePassword(requestID,userID,password string) bool{

	db := dbConn()
	var userpassword string // to use with select table

	row := db.QueryRow("select password from users where email=?", userID)
	err := row.Scan(&userpassword)

	if err != nil {
		if err == sql.ErrNoRows {
			logs.WithFields(logs.Fields{
				"Service":   "User Service",
				"Package":   "Login",
				"function":  "passwordComparision",
				"userid":    userID,
				"requestID": requestID,
			}).Error("No record available for the user")
		} else {
			logs.WithFields(logs.Fields{
				"Service":   "User Service",
				"Package":   "Login",
				"function":  "passwordComparision",
				"userid":    userID,
				"requestID": requestID,
			}).Error("Couldnt fetch users table")
		}
	} // querying database table if

	comparePassword := bcrypt.CompareHashAndPassword([]byte(userpassword), []byte(password))

	if comparePassword != nil {
		return false
	}
	defer db.Close()

	return true
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
