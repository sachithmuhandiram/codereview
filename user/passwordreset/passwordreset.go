package passwordreset

import "net/http"

// CheckEmail verifies whether given email associates with an account
func CheckEmail(res http.ResponseWriter, req *http.Request) {

}

func CheckUpdatePasswordToken(requestID, token string) bool {

	return false
}

func UpdatePassword(requestID, token, password string) bool {
	return false
}
