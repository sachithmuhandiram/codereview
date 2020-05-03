package main

import (
	// "database/sql"
	// "log"
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

func insertLoginToken(res http.ResponseWriter,req *http.Request){
	
	emailedLoginToken := req.FormValue("token")
	requestID := req.FormValue("uid")
	
	t := time.Now()
	t.Format("yyyy-MM-dd HH:mm:ss")


	db := dbConn()
	insertLoginToken, err := db.Prepare("INSERT INTO emailed_login_token(emailed_login_token,created_at) VALUES(?,?)")
        if err != nil {
            panic(err.Error())

            return //false,err
        }

    _,err = insertLoginToken.Exec(emailedLoginToken,t)
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

}