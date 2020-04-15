package main

import(
	"log"
	"github.com/carlescere/scheduler"
	"net/http"
)
func main(){

	

	job := func() {
		log.Println("Time's up!")
	   }
	   scheduler.Every(30).Seconds().Run(job)
	   scheduler.Every().Day().Run(job)
	   scheduler.Every().Sunday().At("08:30").Run(job)

	   http.ListenAndServe("0.0.0.0:7073", nil)
}