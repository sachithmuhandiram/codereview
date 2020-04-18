package main

import (
  //  "os"
    "testing"   
    "log"

    "net/http"
    "net/http/httptest"
   // "bytes"
   // "encoding/json"
    "strconv"
)


func TestEmptyTable(t *testing.T) {
    
    req, _ := http.NewRequest("GET", "/products", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body != "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}


func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}