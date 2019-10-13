# Notes

This documents holds some important golang methods used during this project.

## How to pass data to REST API

When micro service wants to send data to another microservice, it uses 

`http.PostForm()` function. `func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error)`

Used sample :

`_, err := http.PostForm("http://user:7071/checkemail", url.Values{"email": {email}})`

Here it uses `key:value` combination.

[PostForm Example](https://github.com/solderjs/http-examples/blob/master/go/post-form.go)