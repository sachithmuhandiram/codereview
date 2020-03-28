package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
	"os"
	//"strings"
	"io/ioutil"
	"encoding/json"
	"bytes"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	logs "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Globle variable
var checkEMAIL = os.Getenv("CHECKEMAIL")
var userLOGIN = os.Getenv("LOGIN")
var userREGISTER = os.Getenv("REGISTER")
var passwordRESET = os.Getenv("PASSWORDRESET")
var resposeURL = os.Getenv("RESPONSEURL")

// Response struct
type resposeObj struct{
	UID 	string `json :"uid"`
	Service string `json : "service"`
	Function string `json : "function"`
	Pack string 	`json : "pack"`
	Status string 	`json : "status"`
}

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

// This is the main module, this should update into an API gateway.
// Initial step, just routing functionality will be used.
// Running on localhost 7070

// Struct to hold UUID which is attached and passed to all API-Gateway calls
type UUID struct {
	apiUuid uuid.UUID
}

func main() {

	apiID := &UUID{apiUuid: generateUUID()}

	logs.WithFields(logs.Fields{
		"package":  "API-Gateway",
		"function": "main",
		"uuid":     apiID,
	}).Info("API - Gateway started at 7070")

	http.HandleFunc("/getemail", apiID.validatemail)
	http.HandleFunc("/response", reportResponse)
	http.HandleFunc("/login", apiID.userLogin)
	http.HandleFunc("/register", apiID.registerUser)
	http.HandleFunc("/passwordreset", apiID.sendPasswordResetEmail)
	http.HandleFunc("/updatepassword", apiID.updatePassword)

	http.ListenAndServe(":7070", nil)
}

// This will validate email address has valid syntax
func (apiID *UUID) validatemail(res http.ResponseWriter, req *http.Request) {

	validatemailID := apiID.apiUuid
	// Check method
	if req.Method != "POST" {
		logs.WithFields(logs.Fields{
			"package":  "API - Gateway",
			"function": "validatemail",
			"uuid":     validatemailID,
		}).Error("Request method is not POST")

		response := resposeObj{UID:validatemailID.String(),Service:"API Gateway",Function:"validatemail",Pack:"main",Status:"0"}
		
		validatemailResponse,err := json.Marshal(response)

		if err != nil{
			log.Println("Error in marshaling data",err)
		}

		err = sendResponse(validatemailResponse)

		if err != nil {
			log.Println("Error response sending")
		}
		//http.Redirect(res, req, "/", http.StatusSeeOther) // redirect back to register
	} else {
		// Method is POST

		email := req.FormValue("email") //"sachithnalaka@gmail.com" // parse form and get email
		request := "hasaccount"
		validEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`) // regex to validate email address

		if validEmail.MatchString(email) {
			logs.WithFields(logs.Fields{
				"package":  "API-Gateway",
				"function": "validatemail",
				"uuid":     validatemailID,
			}).Info("Valid email format received")

			logs.WithFields(logs.Fields{
				"package":  "API-Gateway",
				"function": "validatemail",
				"email":    email,
				"uuid":     validatemailID, // Later this should change for function-wise uuid
			}).Info("Email will pass to User - Service")

			_, err := http.PostForm(checkEMAIL, url.Values{"email": {email}, "uid": {validatemailID.String()},"request":{request}})

			if err != nil {
				logs.WithFields(logs.Fields{
					"package":  "API-Gateway",
					"function": "validatemail",
					"email":    email,
					"error":    err,
					"uuid":     validatemailID,
				}).Error("Error posting data to User - Service")

				_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {validatemailID.String()}, "service": {"API Gateway"},
					"function": {"validatemail"}, "package": {"main"}, "status": {"0"}})

				if err != nil {
					log.Println("Error response sending")
				}
			}

		} else {
			logs.WithFields(logs.Fields{
				"package":  "API-Gateway",
				"function": "validatemail",
				"email":    email,
				"uuid":     validatemailID,
			}).Error("Wrong email format received")

			_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {validatemailID.String()}, "service": {"API Gateway"},
				"function": {"validatemail"}, "package": {"main"}, "status": {"0"}})

			if err != nil {
				log.Println("Error response sending")
			}
			// Return to register window
			//return false
		}
	} // Method checking if loop
}

// User login
func (apiID *UUID) userLogin(res http.ResponseWriter, req *http.Request) {

	//jwt := req.FormValue("jwt") //req.URL.Query()["jwt"]
	requestID := apiID.apiUuid
	userid := req.FormValue("email")
	password := req.FormValue("password")

	if password == "" {
		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "userLogin",
			"uuid":     apiID,
		}).Error("User Login request received,without password")
		return

	}

	//hashedPassword := hashPassword(password)

	logs.WithFields(logs.Fields{
		"package":  "API-Gateway",
		"function": "validuserLoginatemail",
		"uuid":     apiID,
	}).Info("User Login request received")

	_, err := http.PostForm(userLOGIN, url.Values{"userid": {userid}, "uid": {requestID.String()},
		"password": {password}})

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "userLogin",
			"userid":   userid,
			"error":    err,
			"uuid":     requestID,
		}).Error("Error posting data to User - Service")
	}

}

func generateUUID() uuid.UUID {
	// Generating UUID
	uuid, err := uuid.NewUUID()

	if err != nil {
		if err != nil {
			logs.WithFields(logs.Fields{
				"package":  "API-Gateway",
				"function": "generateUUID",
				"error":    err,
			}).Error("Couldnt generate a UUID")
		}

		// Return to error page
		//http.Redirect(loginResponse, loginRequest, "/", http.StatusSeeOther)
	}
	return uuid //.String()
}

// Register a User
func (apiID *UUID) registerUser(res http.ResponseWriter, req *http.Request) {

	responseID := apiID.apiUuid //req.FormValue("uid")
	firstName := req.FormValue("first_name")
	lastName := req.FormValue("last_name")
	email := req.FormValue("email")
	password := req.FormValue("password")
	conformPassword := req.FormValue("conformpassword")

	logs.WithFields(logs.Fields{
		"package":  "API-Gateway",
		"function": "registerUser",
		"email":    email,
		"uuid":     responseID,
	}).Info("API gateway received data to register a user account")

	if password != conformPassword {
		logs.WithFields(logs.Fields{
			"package":         "API-Gateway",
			"function":        "registerUser",
			"email":           email,
			"uuid":            responseID,
			"password":        password,
			"conformPassword": conformPassword,
		}).Error("Password and conform password mismatch")
		return
	}

	password = hashPassword(password)

	_, err := http.PostForm(userREGISTER, url.Values{"email": {email}, "uid": {responseID.String()},
		"first_name": {firstName}, "last_name": {lastName}, "password": {password}})

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "registerUser",
			"email":    email,
			"error":    err,
			"uuid":     responseID,
		}).Error("Error posting data to User - Service")
	}

}

func (apiID *UUID) sendPasswordResetEmail(res http.ResponseWriter, req *http.Request) {

	requestID := apiID.apiUuid
	email := req.FormValue("email")

	logs.WithFields(logs.Fields{
		"package":  "API Gateway",
		"function": "passwordReset",
		"ApiUUID":  requestID,
		"email":    email,
	}).Info("Password Reset received email address")

	_, err := http.PostForm(checkEMAIL, url.Values{"email": {email}, "uid": {requestID.String()}, "request": {"passwordreset"}})

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "passwordReset",
			"email":    email,
			"error":    err,
			"uuid":     requestID,
		}).Error("Error posting data to User - Service")

		_, err := http.PostForm("http://localhost:7070/response", url.Values{"uid": {requestID.String()}, "service": {"API Gateway"},
			"function": {"passwordReset"}, "package": {"main"}, "status": {"0"}})

		if err != nil {
			log.Println("Error response sending")
		}
	}

}

func (apiID *UUID) updatePassword(res http.ResponseWriter, req *http.Request) {

	requestID := apiID.apiUuid
	token := req.FormValue("token")
	password := req.FormValue("password")
	conformPassword := req.FormValue("conformpassword")

	logs.WithFields(logs.Fields{
		"package":  "API Gateway",
		"function": "updatePassword",
		"ApiUUID":  requestID,
	}).Info("Update password received updated passwords")

	if password != conformPassword {
		logs.WithFields(logs.Fields{
			"package":         "API-Gateway",
			"function":        "updatePassword",
			"uuid":            requestID,
			"password":        password,
			"conformPassword": conformPassword,
		}).Error("Password and conform password mismatch")
		return
	}

	password = hashPassword(password)

	_, err := http.PostForm(passwordRESET, url.Values{"uid": {requestID.String()},
		"token": {token}, "password": {password}, "request": {"updatepassword"}})

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "registerUser",
			"error":    err,
			"uuid":     requestID,
		}).Error("Error posting data to User - Service from updatePassword")
	}

}

// Insert Req/Res to table
func storeDetails(uuid string, reqType, status bool) error {
	db := dbConn()

	t := time.Now()

	// insert token to gateway_req_res table
	insData, err := db.Prepare("INSERT INTO gateway_req_res (uuid,type,status,time) VALUES(?,?,?,?)")
	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "API Gateway",
			"function": "storeDetails",
			"Error":    err,
		}).Error("Couldnt prepare insert statement for gateway_req_res table")

		return err
	}
	insData.Exec(uuid, reqType, status, t)
	defer db.Close()
	return nil
}

// hashing password
func hashPassword(password string) string {
	// This will generate a token to

	bs := []byte(password) // convert UUID into a bytestream

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
// Sending a response
func sendResponse(req []byte)error{


	log.Println("Received response : ",string(req))
	request,err := http.NewRequest("POST",resposeURL, bytes.NewBuffer(req))
	request.Header.Set("Content-Type","application/json")

	client := &http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

	if err != nil{
		return err
	}

	return nil
}

// Response for a request is recorded.
func reportResponse(res http.ResponseWriter, req *http.Request) {

	data,_ := ioutil.ReadAll(req.Body)

	log.Printf("Request data from reportResponse function: %s\n",data)


    var request resposeObj
    err := json.Unmarshal(data, &request)
    
    if err != nil {
        log.Println("Error occured decoding JSON object",err)
    }

    log.Println("From response reporting",request.UID)

	logs.WithFields(logs.Fields{
		"ResponseService": request.Service,
		"ResponsePackage": request.Pack,
		"ResponseFunc":    request.Function,
		"responseID":      request.UID,
		"status":          request.Status,
	}).Info("Response received for the request")

	stat, _ := strconv.Atoi(request.Status) // convert string to int

	// stat 1 = success, 0 = failed
	if stat == 1 {
		storeData := storeDetails(request.UID, true, true)

		if storeData != nil {
			logs.WithFields(logs.Fields{
				"package":  "API Gateway",
				"function": "reportResponse",
				"Error":    storeData,
			}).Error("Response data insert to DB failed")
		}
	} else {
		storeData := storeDetails(request.UID, true, false)

		if storeData != nil {
			logs.WithFields(logs.Fields{
				"package":  "API Gateway",
				"function": "reportResponse",
				"Error":    storeData,
			}).Error("Response data insert to DB failed")
		}
	}

}
