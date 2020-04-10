package main

import(
	logs "github.com/sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"time"
	"log"
	"database/sql"
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
	log.Println("Assigned JWT : ",assignedjwt)
	log.Println("Available JWT",jwt)
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