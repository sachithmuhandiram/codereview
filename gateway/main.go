package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// For generateUUID function return
type UUID [16]byte

type apiTracer struct {
	tracer opentracing.Tracer
}

// Tracing function
func startTracing(service string) opentracing.Tracer {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, _, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer
}

// This is the main module, this should update into an API gateway.
// Initial step, just routing functionality will be used.
// Running on localhost 7070

func main() {

	log.Println("API gateway started at port : 7070")
	apitracer := apiTracer{startTracing("Validate Email")}

	http.HandleFunc("/getemail", apitracer.validatemail)
	http.HandleFunc("/response", getRespose)

	http.ListenAndServe(":7070", nil)
}

// This will validate email address has valid syntax (apitracer apiTracer)
func (api *apiTracer) validatemail(res http.ResponseWriter, req *http.Request) {

	validateEmailSpan := api.tracer.StartSpan("Email address checking")

	// Get UUID
	uuid := getUUID()

	// Check method
	if req.Method != "POST" {
		log.Panic("Email form data is not Post")
		//http.Redirect(res, req, "/", http.StatusSeeOther) // redirect back to register
	}

	email := req.FormValue("email") //"sachithnalaka@gmail.com" // parse form and get email

	validEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`) // regex to validate email address

	if validEmail.MatchString(email) {
		log.Println("Valid email address format received")
		log.Println("Email is passed to user module to validate ")
		httpClient := http.Client{}

		form := url.Values{}
		form.Set("uuid", uuid)
		form.Set("email", email)

		//b := bytes.NewBufferString(form.Encode())
		userserviceUrl := "http://user:7071/checkemail"

		checkEmailReq, err := http.NewRequest("POST", userserviceUrl, strings.NewReader(form.Encode()))
		if err != nil {
			log.Println(err)
		}

		ext.SpanKindRPCClient.Set(validateEmailSpan)
		ext.HTTPUrl.Set(validateEmailSpan, userserviceUrl)
		ext.HTTPMethod.Set(validateEmailSpan, "POST")

		validateEmailSpan.Tracer().Inject(
			validateEmailSpan.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := httpClient.Do(checkEmailReq)
		//_, err = http.PostForm("http://user:7071/checkemail", url.Values{"uuid": {uuid}, "email": {email}})

		if err != nil {
			log.Println("Couldnt verify email address user service sends an error : ", err)
		}
		defer resp.Body.Close()

	} else {
		log.Println("Wrong email address format")
		// Return to register window
		//return false
	}

	//defer gatewayTracerClose.Close()
	defer validateEmailSpan.Finish()
}

func getUUID() string {
	// Generating UUID
	uuid, err := uuid.NewUUID()

	if err != nil {
		log.Println("Couldnt generate a UUID for the request")

		// Return to error page
		//http.Redirect(loginResponse, loginRequest, "/", http.StatusSeeOther)
	}
	return uuid.String()
}

func getRespose(res http.ResponseWriter, req *http.Request) {

	uuid := req.FormValue("uuid")
	serviceID := req.FormValue("serviceID")
	resultCode := req.FormValue("resultCode")

	log.Println("UUID,ServiceID,resultCode", uuid, serviceID, resultCode)
}
