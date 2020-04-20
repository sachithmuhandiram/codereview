package main

import(
	"log"
	"os"
	"database/sql"
	"time"
	"github.com/carlescere/scheduler"
	"net/http"
	logs "github.com/sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

// API gateway Database
var gatewayDB	= os.Getenv("MYSQLDBGATEWAY")
// User database
var mysqlDB	= os.Getenv("MYSQLDBUSERS")

func gatewayDBConn() (db *sql.DB) {
	gatewayDB, err := sql.Open("mysql", gatewayDB)

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "Database Helper Service",
			"function": "GatewayDBConn",
			"error":    err,
		}).Error("Failed to connect to database")
	}
	return gatewayDB
}

// User database connection
func userDBConn() (db *sql.DB) {
	userDB, err := sql.Open("mysql", mysqlDB)

	if err != nil {
		logs.WithFields(logs.Fields{
			"package":  "Database Helper Service",
			"function": "userDBConn",
			"error":    err,
		}).Error("Failed to connect to database")
	}
	return userDB
}

func main(){

	   scheduler.Every(30).Seconds().Run(readLoginToken)
	   scheduler.Every(30).Seconds().Run(readActiveJWTToken)
	//    scheduler.Every().Day().Run(job)
	//    scheduler.Every().Sunday().At("08:30").Run(job)
	//    readLoginToken()
	   http.ListenAndServe("0.0.0.0:7073", nil)
}

// Reading from API gateway login_token table
func readLoginToken(){
	// reads login tokens and checks whether they are older than 10mins
	gatewaydb := gatewayDBConn()

	var loginToken string 
	var createdAt time.Time 

	rows, err := gatewaydb.Query("select login_token,created_at from login_token")
    if err != nil {
        log.Println(err)
    }

	// UTC time, as we store time as UTC
	t := time.Now()
	utc := t.In(time.UTC)

    for rows.Next() {
        err := rows.Scan(&loginToken,&createdAt)
        if err != nil {
            log.Println(err)
		}
		// t1.Sub(t2).Hours()
		if utc.Sub(createdAt).Minutes() > 10 {
			// older tokens
			deleteOldLoginToken(loginToken)
		}
    }
	defer rows.Close()
}

// This will delete un-used old login tokens
func deleteOldLoginToken(token string){
	gatewaydb := gatewayDBConn()

	delTkn, err := gatewaydb.Prepare("DELETE FROM login_token WHERE login_token=?")
    if err != nil {
        log.Println(err.Error())
    }
    delTkn.Exec(token)

    logs.WithFields(logs.Fields{
			"Service":  "Database Helper Service",
			"function":  "deleteOldLoginToken",
		}).Info("Successfully deleted login token")

	defer gatewaydb.Close()
}

// Reading from API gateway activeJWTtokens table
func readActiveJWTToken(){
	// reads login tokens and checks whether they are older than 10mins
	gatewaydb := gatewayDBConn()

	var jwt string 
	var lastUpdated time.Time 

	rows, err := gatewaydb.Query("select jwt,last_update from activeJWTtokens")
    if err != nil {
        log.Println(err)
    }

	// UTC time, as we store time as UTC
	t := time.Now()
	utc := t.In(time.UTC)

    for rows.Next() {
        err := rows.Scan(&jwt,&lastUpdated)
        if err != nil {
            log.Println(err)
		}

		if utc.Sub(lastUpdated).Minutes() > 10 {
			deleteJWTToken(jwt)
		}
    }
	defer rows.Close()
}

// This will delete un-updated JWT tokens
func deleteJWTToken(jwttoken string){
	gatewaydb := gatewayDBConn()

	delTkn, err := gatewaydb.Prepare("DELETE FROM activeJWTtokens WHERE jwt=?")
    if err != nil {
        log.Println(err.Error())
    }
    delTkn.Exec(jwttoken)

    logs.WithFields(logs.Fields{
			"Service":  "Database Helper Service",
			"function":  "deleteJWTTokenn",
		}).Info("Successfully deleted jwt token")

	defer gatewaydb.Close()
}