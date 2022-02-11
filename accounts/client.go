package accountsclient

import (
	"net/http"
)

const baseUrl = "http://localhost:8080/v1/"

func GetApiStatus() bool {
	response, err := http.Get(baseUrl + "health")
	if err == nil && response.StatusCode == 200 {
		return true
	}
	return false
}
