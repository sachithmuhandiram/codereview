package main

import (
	// "database/sql"
	 "log"
	 "net/http"
	// "net/url"
	// "regexp"
	// "strconv"
	 "time"
	// "os"
	// "fmt"
	// "io/ioutil"
	// "encoding/json"
	// "bytes"

	 _ "github.com/go-sql-driver/mysql"
	// "github.com/google/uuid"
	 logs "github.com/sirupsen/logrus"
	// "golang.org/x/crypto/bcrypt"
	// "github.com/dgrijalva/jwt-go"
)

// createSession function
func createSession(res http.ResponseWriter,req *http.Request){

	uid := req.FormValue("uid")
	authorized := req.FormValue("authorize")

	if authorized == "1"{
		user := req.FormValue("userid")
		expirationTime := time.Now().Add(5 * time.Minute)

		log.Println("Expire time : ",expirationTime)

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
		insertJWTResponse,err := InsertJWT(uid,user,jwt)
		
		if insertJWTResponse == true && err == nil{
			log.Println("JWT insert to insert to table")

			// removing logintoken
			expire := time.Now().Add(-7 * 24 * time.Hour)
			logincookie := http.Cookie{
				Name:    "logintoken",
				//Value:   "value",
				Expires: expire,
			}
			http.SetCookie(res, &logincookie)

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

func insertLoginToken(res http.ResponseWriter,req *http.Request){
	
	emailedLoginToken := req.FormValue("token")
	requestID := req.FormValue("uid")
	
	t := time.Now()
	t.Format("yyyy-MM-dd HH:mm:ss")


	db := dbConn()
	insertLoginToken, err := db.Prepare("INSERT INTO emailed_login_token(emailed_login_token,isActive,created_at) VALUES(?,?,?)")
        if err != nil {
            panic(err.Error())

            return //false,err
        }

    _,err = insertLoginToken.Exec(emailedLoginToken,1,t)
	if err != nil{
		logs.WithFields(logs.Fields{
			"Service":   "API Gateway",
			"Package":   "Internal Route Handler",
			"function":  "insertLoginToken",
		//	"userid":    userID,
			"requestID": requestID,
			"Error" : err,
		}).Error("Could not insert into valid emailed_login_token tokens")	
	}else{
	
		logs.WithFields(logs.Fields{
			"Service":   "API Gateway",
			"Package":   "Internal Route Handler",
			"function":  "insertLoginToken",
		//	"userid":    userID,
			"requestID": requestID,
			//"JWT" : jwtToken,
		}).Info("Insert into valid emailed_login_token tokens")	
	}
	
	defer db.Close()
}

func insertRegisterToken(res http.ResponseWriter,req *http.Request){
	
	email := req.FormValue("email")
	emailedRegisterToken := req.FormValue("token")
	requestID := req.FormValue("uid")
	
	t := time.Now()
	t.Format("yyyy-MM-dd HH:mm:ss")


	db := dbConn()
	insertRegisterToken, err := db.Prepare("INSERT INTO emailed_register_token(emailed_register_token,email,isActive,created_at) VALUES(?,?,?,?)")
        if err != nil {
            panic(err.Error())

            return //false,err
        }

    _,err = insertRegisterToken.Exec(emailedRegisterToken,email,1,t)
	if err != nil{
		logs.WithFields(logs.Fields{
			"Service":   "API Gateway",
			"Package":   "Internal Route Handler",
			"function":  "insertRegisterToken",
		//	"userid":    userID,
			"requestID": requestID,
			"Error" : err,
		}).Error("Could not insert into valid emailed_register_token tokens")	
	}else{
	
		logs.WithFields(logs.Fields{
			"Service":   "API Gateway",
			"Package":   "Internal Route Handler",
			"function":  "insertRegisterToken",
		//	"userid":    userID,
			"requestID": requestID,
			//"JWT" : jwtToken,
		}).Info("Insert into valid emailed_register_token tokens")	
	}
	
	defer db.Close()
}

func insertPasswordResetToken(res http.ResponseWriter, req *http.Request){
	emailedPassResetToken := req.FormValue("token")
	requestID := req.FormValue("uid")
	email := req.FormValue("email")
	
	t := time.Now()
	t.Format("yyyy-MM-dd HH:mm:ss")


	db := dbConn()
	insertPassResetToken, err := db.Prepare("INSERT INTO emailed_password_reset_token(emailed_passwort_reset_token,email,isActive,created_at) VALUES(?,?,?,?)")
        if err != nil {
            panic(err.Error())
			log.Println("Error : ",err)
            return //false,err
        }

    _,err = insertPassResetToken.Exec(emailedPassResetToken,email,1,t)
	if err != nil{
		logs.WithFields(logs.Fields{
			"Service":   "API Gateway",
			"Package":   "Internal Route Handler",
			"function":  "insertPasswordResetToken",
		//	"userid":    userID,
			"requestID": requestID,
			"Error" : err,
		}).Error("Could not insert into valid emailed_password_reset_token tokens")	
	}else{
	
		logs.WithFields(logs.Fields{
			"Service":   "API Gateway",
			"Package":   "Internal Route Handler",
			"function":  "insertPasswordResetToken",
		//	"userid":    userID,
			"requestID": requestID,
			//"JWT" : jwtToken,
		}).Info("Insert into valid emailed_password_reset_token tokens")	
	}
	
	defer db.Close()

}