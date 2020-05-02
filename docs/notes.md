# Notes

This document holds some important golang methods used during this project.

## How to pass data to REST API

When micro service wants to send data to another microservice, it uses 

`http.PostForm()` function. `func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error)`

Used sample :

`_, err := http.PostForm("http://user:7071/checkemail", url.Values{"email": {email}})`

Here it uses `key:value` combination.

[PostForm Example](https://github.com/solderjs/http-examples/blob/master/go/post-form.go)

[Golang Interfaces](https://medium.com/better-programming/a-real-world-example-of-go-interfaces-98e89b2ddb67)

[Passing structs to golang function](https://stackoverflow.com/questions/29805583/how-can-i-pass-struct-to-function-as-parameter-in-go-lang)

[Handling JSON objects in POST](https://stackoverflow.com/questions/15672556/handling-json-post-request-in-go)

[Implementing JWT based authentication in Golang](https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/)

[JWT authentication sample SO](https://stackoverflow.com/questions/36236109/go-and-jwt-simple-authentication)

[Extract JWT claims](https://stackoverflow.com/questions/39859244/how-to-extract-the-claims-from-jwt-token)

[Middleware checking and sending to appropriate route](https://stackoverflow.com/questions/48899277/what-is-the-best-way-to-send-back-the-response-from-the-middle-ware-in-golang)

[API design and testing](https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql)

[Deleting a cookie](https://stackoverflow.com/questions/29499843/how-do-you-delete-a-cookie-with-go-and-http-package)

[Getting time data from MySQL](https://stackoverflow.com/questions/29341590/how-to-parse-time-from-database)

[Comparing time difference using golang](https://stackoverflow.com/questions/40260599/difference-between-two-time-time-objects)

[Testing Docker-based Go Applications with Go Test](https://medium.com/@jmarhee/testing-docker-based-go-applications-with-go-test-and-drone-io-5ae5db7d3bcf)

[Testing golang Handlers](https://blog.questionable.services/article/testing-http-handlers-go/)

[Fix error with go-test and gcc](https://stackoverflow.com/questions/55442155/why-are-go-tests-suddenly-requiring-gcc)

## Environment variables

[Docker-compose environment variables](https://docs.docker.com/compose/environment-variables/)

[Setting-up timezones](https://stackoverflow.com/questions/47498297/docker-compose-container-timezone)

## Security concerns 

[JWT security issues](https://medium.com/swlh/hacking-json-web-tokens-jwts-9122efe91e4a)