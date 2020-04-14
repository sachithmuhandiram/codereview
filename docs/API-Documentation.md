# API document

This will describe publically exposed APIs in the system, their inputs and outputs. Also a short description about how they work internally.

1. `gateway:port/login`<br>
    Method : `GET`

    This will return login page to user.

2. `gateway:port/userlogin`<br>
    Method : `POST`
    
    Inputs :

    * Username : email address
    * Password : any string

    Both inputs are mandetory. Route will first check whether given email address has a valid email syntax  using regular expression. If so, then checks whether password field is empty. If so returns back to `gateway:port/login`.