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
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	logs "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
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

// Requests are authenticated
func authenticateToken(handlerFunc http.HandlerFunc) http.HandlerFunc {

	log.Println("Authentication function called")

	return func(res http.ResponseWriter, req *http.Request) {
		cookie, _ := req.Cookie("usertoken")

		if cookie == nil {
			log.Println("Cant find cookie :")
			http.Redirect(res,req,"/login",http.StatusSeeOther)
		}else{
			// Cookie is there, need to validate that cookie
			if cookie.Name == "usertoken"{
				jwtClaims,err := getUserFromJWT(cookie.Value)

				if err != nil{
					log.Println("There was an error getting user details ",err)
					http.Redirect(res,req,"/login",http.StatusSeeOther)		
				} // got jwt claims 
					
				user := jwtClaims["user"].(string)
				log.Println("User is : ",user)
				home(res,req)

			} // cookie name checking if loop

		} // there is a cookie
	} // function return
}

func main() {

	apiID := &UUID{apiUuid: generateUUID()}

	logs.WithFields(logs.Fields{
		"package":  "API-Gateway",
		"function": "main",
		"uuid":     apiID,
	}).Info("API - Gateway started at 7070")

	http.HandleFunc("/getemail", apiID.validatemail)
	http.HandleFunc("/home",authenticateToken(home))
	http.HandleFunc("/createsession",createSession)
	http.HandleFunc("/response", reportResponse)
	http.HandleFunc("/userlogin", apiID.userLogin)
	http.HandleFunc("/register", apiID.registerUser)
	http.HandleFunc("/passwordreset", apiID.sendPasswordResetEmail)
	http.HandleFunc("/updatepassword", apiID.updatePassword)
	http.HandleFunc("/login",login)

	http.ListenAndServe(":7070", nil)
}
// Home function
func home(res http.ResponseWriter,req *http.Request){
	fmt.Fprintf(res, "Welcome to Home Page!")
	
}
// Login function - This will return login form
func login(res http.ResponseWriter,req *http.Request){
	fmt.Fprintf(res, "This is login page")
}

// createSession function
func createSession(res http.ResponseWriter,req *http.Request){

	uid := req.FormValue("uid")
	authorized := req.FormValue("authorize")

	if authorized == "1"{
		user := req.FormValue("userid")
		expirationTime := time.Now().Add(5 * time.Minute)

		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "createSession",
			"uuid":     uid,
		}).Info("User details received to create a session")

		// // creating JWT for user
		jwt,jwtErr := GenerateJWT(user)

		if jwtErr != nil{
			log.Println("Error generating JWT, cant go further")
			return
		}
		// Insert JWT to table
		insertJWTResponse,err := http.PostForm("http://user:7071/insertJWT",url.Values{"uid":{uid},"userid":{user},"jwt":{jwt}})
		
		if insertJWTResponse.StatusCode == 200 && err == nil{
			log.Println("JWT sent to user service to insert to table")

			// Setting jwt cookie
			http.SetCookie(res, &http.Cookie{
				Name:    "usertoken",
				Value:   jwt,
				Expires: expirationTime,
			})

			http.Redirect(res,req,"/home",http.StatusSeeOther)

		}else{
			log.Println("JWT sending to user service failed : Abort",err)
			return
		}

	}else{ // authorized = 0

		log.Println("Something bad happened")
	}	
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

	loginToken := req.FormValue("logintoken")
	requestID := apiID.apiUuid
	userid := req.FormValue("email")
	password := req.FormValue("password")

	if password == "" {
		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "userLogin",
			"uuid":     requestID,
		}).Error("User Login request received,without password")
		return

	}

	// //hashedPassword := hashPassword(password)

	logs.WithFields(logs.Fields{
		"package":  "API-Gateway",
		"function": "validuserLoginatemail",
		"uuid":     requestID,
	}).Info("User Login request received")

	parameters := userLOGIN+"?userid="+ userid +"&uid="+requestID.String()+"&password="+password+"&logintoken="+loginToken
	
	http.Redirect(res, req, parameters, http.StatusSeeOther)
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
	//log.Println("Received response : ",string(req))
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
	var stat int
	var requestID string
	requestType := req.Header.Get("Content-Type")

	//log.Println("Content type is : ",requestType)
	switch(requestType){
	case "application/json" :
		data,_ := ioutil.ReadAll(req.Body)

		//log.Printf("Request data from reportResponse function: %s\n",data)
	    var request resposeObj
	    err := json.Unmarshal(data, &request)
	    
	    if err != nil {
	        log.Println("Error occured decoding JSON object",err)
	    }

	    requestID = request.UID

		logs.WithFields(logs.Fields{
			"ResponseService": request.Service,
			"ResponsePackage": request.Pack,
			"ResponseFunc":    request.Function,
			"responseID":      request.UID,
			"status":          request.Status,
		}).Info("Response received for the request")

		stat, _ = strconv.Atoi(request.Status) // convert string to int
	default: // form data
		responseID := req.FormValue("uid")
		service := req.FormValue("service")
		function := req.FormValue("function")
		pack := req.FormValue("package")
		status := req.FormValue("status")

		requestID = req.FormValue("uid")
		logs.WithFields(logs.Fields{
			"ResponseService": service,
			"ResponsePackage": pack,
			"ResponseFunc":    function,
			"responseID":      responseID,
			"status":          status,
		}).Info("Response received for the request")

		stat, _ = strconv.Atoi(status) // convert string to int

	} // response is json or form data
	

	// stat 1 = success, 0 = failed
	if stat == 1 {
		storeData := storeDetails(requestID, true, true)

		if storeData != nil {
			logs.WithFields(logs.Fields{
				"package":  "API Gateway",
				"function": "reportResponse",
				"Error":    storeData,
			}).Error("Response data insert to DB failed")
		}
	} else {
		storeData := storeDetails(requestID, true, false)

		if storeData != nil {
			logs.WithFields(logs.Fields{
				"package":  "API Gateway",
				"function": "reportResponse",
				"Error":    storeData,
			}).Error("Response data insert to DB failed")
		}
	}

}

// JWT
// GenerateJWT takes eventID as a parameter and time (minutes) for JWT
func GenerateJWT(user string) (string, error) {

	//loginKey := []byte(user)
	appSecretKey  := []byte("du-bi-du-bi-dub") // takes 531855448467 years to break using brute-force attack
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Duration(5))

	jwtToken, jwtErr := token.SignedString(appSecretKey)

	if jwtErr != nil {
		log.Println("Error creating jwt Token : ", jwtErr)
		return "", jwtErr
	}

	return jwtToken, nil
}

// Get claims from JWT
func getUserFromJWT(userJWT string)(jwt.MapClaims,error){

	    hmacSecretString := "du-bi-du-bi-dub"
        hmacSecret := []byte(hmacSecretString)
        token, err := jwt.Parse(userJWT, func(token *jwt.Token) (interface{}, error) {
             // check token signing method etc
             return hmacSecret, nil
        })

        if err != nil {
            return nil, err
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            return claims, nil
        } else {
            log.Printf("Invalid JWT Token")
            return nil, nil
        }

}
