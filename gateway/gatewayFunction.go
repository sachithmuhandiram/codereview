package main

// Insert into valid token
func InsertJWT(res http.ResponseWriter,req * http.Request) {
	
	req.ParseForm()
	requestID := req.FormValue("uid")
	userID := req.FormValue("userid")//req.FormValue("userid")
	jwtToken := req.FormValue("jwt")
	
	t := time.Now()
	t.Format("yyyy-MM-dd HH:mm:ss")


	db := dbConn()
	insertToken, err := db.Prepare("INSERT INTO activeJWTtokens(user_id,jwt,created_at,last_update) VALUES(?,?,?,?)")
        if err != nil {
            panic(err.Error())

            return 
        }

    _,err = insertToken.Exec(userID, jwtToken,t,t)
	if err != nil{
		log.Println("Error occured : ",err)
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

}