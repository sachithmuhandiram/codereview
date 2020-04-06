# Notes

This documents holds some important golang methods used during this project.

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

[]()

## Environment variables

[Docker-compose environment variables](https://docs.docker.com/compose/environment-variables/)