package main

import(
	logs "github.com/sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"time"
	"log"
	"database/sql"
)

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

func addLoginJWT(uuid,loginJWT string) bool{
	db := dbConn()
	t := time.Now()
	t.Format("yyyy-MM-dd HH:mm:ss")

	insertLoginToken, err := db.Prepare("INSERT INTO login_token(login_token,created_at) VALUES(?,?)")
        if err != nil {
            panic(err.Error())

            return false
        }

    _,err = insertLoginToken.Exec(loginJWT,t)
	if err != nil{
		log.Println("Error occured : ",err)
		return false
	}else{
	
		logs.WithFields(logs.Fields{
			"Service":   "API gateway",
			"function":  "addLoginJWT",
			"requestID": uuid,
			//"JWT" : jwtToken,
		}).Info("Insert into Login jwt tokens")	

		return true
	}
	
	defer db.Close()

	return false
}

// Insert into valid token
func InsertJWT(requestID,userID,jwtToken string) (bool,error) {
	
	loc, _ := time.LoadLocation("Asia/Colombo")
	t := time.Now().In(loc)
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

// Check JWT and user with activeJWT
func checkJWT(user,jwt string)(bool,error){

	db := dbConn()
	var assignedjwt string // to use with select table

	row := db.QueryRow("select jwt from activeJWTtokens where user_id=?", user)
	err := row.Scan(&assignedjwt)

	if err != nil {
		if err == sql.ErrNoRows {
			logs.WithFields(logs.Fields{
				"Service":   "API Gateway",
				"Package":   "Helper",
				"function":  "checkJWT",
				"userid":    user,
			}).Error("No JWT available for the user")

			return false,err
		} else {
			logs.WithFields(logs.Fields{
				"Service":   "API Gateway",
				"Package":   "Helper",
				"function":  "checkJWT",
				"userid":    user,
			}).Error("Couldnt fetch activeJWTtokens table")

			return false,err
		}
	} // querying database table if

	if assignedjwt == jwt{

		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "checkJWT",
			"user" : user,
		}).Info("User has a valid JWT")

		return true,nil
	}else{

		logs.WithFields(logs.Fields{
			"package":  "API-Gateway",
			"function": "checkJWT",
			"user" : user,
		}).Info("User does not have a valid JWT")

		return false,nil	
	}
}
// Update JWT
func updateUserActivity(user string){
	loc, _ := time.LoadLocation("Asia/Colombo")
	t := time.Now().In(loc)
	t.Format("yyyy-MM-dd HH:mm:ss")


	db := dbConn()
	updateActivity, err := db.Prepare("UPDATE activeJWTtokens SET last_update =? WHERE user_id=?")
        if err != nil {
            log.Println(err.Error())

            return
        }

    _,err = updateActivity.Exec(t,user)
	if err != nil{
		log.Println("Error occured : ",err)
		return 
	}else{
	
		logs.WithFields(logs.Fields{
			"Service":   "API Gateway",
			"function":  "updateUserActivity",
			"userid":    user,
			//"JWT" : jwtToken,
		}).Info("Updated user activity")	
	}
	
	defer db.Close()
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