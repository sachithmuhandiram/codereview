package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	logs "github.com/sirupsen/logrus"
)

// database connection
// func dbConn() (db *sql.DB) {
// 	db, err := sql.Open("mysql", "root:7890@tcp(127.0.0.1:3306)/codereview_users")

// 	if err != nil {
// 		logs.WithFields(logs.Fields{
// 			"package":  "User Service",
// 			"function": "dbConn",
// 			"error":    err,
// 		}).Error("Failed to connect to database")
// 	}
// 	return db
// }

// CheckEmail verifies whether given email associates with an account
// func CheckEmail(res http.ResponseWriter, req *http.Request) {

// }

func CheckUpdatePasswordToken(requestID, token string) bool {
	db := dbConn()

	var resetToken bool

	log.Println("Token received : ", token)
	// This will return a true or false
	row := db.QueryRow("select exists(select email from passwordResetToken  where resettoken=?)", token)

	err := row.Scan(&resetToken)
	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "User Service",
			"function": "CheckEmail",
			"error":    err,
			"uuid":     requestID,
		}).Error("Failed to fetch data from user table")
	}

	log.Println("Password reset Toekn : ", resetToken)
	if resetToken {
		return true
	}

	defer db.Close()
	return false
}

func UpdatePassword(requestID, token, password string) bool {
	return false
}
