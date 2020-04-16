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
			"package":  "User Service",
			"function": "dbConn",
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
			"package":  "User Service",
			"function": "dbConn",
			"error":    err,
		}).Error("Failed to connect to database")
	}
	return userDB
}

func main(){

	

	job := func() {
		log.Println("Time's up!")
	   }
	   scheduler.Every(30).Seconds().Run(job)
	   scheduler.Every().Day().Run(job)
	   scheduler.Every().Sunday().At("08:30").Run(job)
	   readLoginToken()
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

	// GMT time
	t := time.Now()
	utc := t.In(time.UTC)

	log.Println("UTC time: ",utc)
    for rows.Next() {
        err := rows.Scan(&loginToken,&createdAt)
        if err != nil {
            log.Println(err)
		}
		// t1.Sub(t2).Hours()
		if utc.Sub(createdAt).Minutes() > 10 {
			log.Println("this", loginToken,createdAt)
		}else{
			log.Println("There are new records")
		}

    }

	// tokens older than 10mins send to deleteLoginToken(token) 
	
	defer rows.Close()
}