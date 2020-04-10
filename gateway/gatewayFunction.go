package main

import(
	logs "github.com/sirupsen/logrus"
	"time"
	"log"
)

// Insert into valid token
func InsertJWT(requestID,userID,jwtToken string) (bool,error) {
	
	t := time.Now()
	t.Format("yyyy-MM-dd HH:mm:ss")


	db := dbConn()
	insertToken, err := db.Prepare("INSERT INTO activeJWTtokens(user_id,jwt,created_at,last_update) VALUES(?,?,?,?)")
        if err != nil {
            panic(err.Error())

            return false,err
        }

    _,err = insertToken.Exec(userID, jwtToken,t,t)
	if err != nil{
		log.Println("Error occured : ",err)
		return false,err
	}else{
	
		logs.WithFields(logs.Fields{
			"Service":   "User Service",
			"Package":   "Login",
			"function":  "InsertJWT",
			"userid":    userID,
			"requestID": requestID,
			//"JWT" : jwtToken,
		}).Info("Insert into valid jwt tokens")	
	}
	
	defer db.Close()
	return true,nil // success
}