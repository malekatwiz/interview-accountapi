package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const baseUrl = "http://localhost:8080/v1"

func GetApiStatus() bool {
	response, err := http.Get(baseUrl + "/health")
	if err == nil && response.StatusCode == 200 {
		return true
	}
	return false
}

func CreateAccount(request CreateAccountRequest) AccountsApiResponse {
	var creationResult AccountsApiResponse
	reqBody, jsonErr := json.Marshal(request) // TODO: handle json errors.
	if jsonErr != nil {
		creationResult.Errors = append(creationResult.Errors, "Invalid JSON input")
		return creationResult
	}

	res, err := http.Post(baseUrl+"/organisation/accounts", "application/vnd.api+json", bytes.NewBuffer(reqBody))
	if err != nil || res.StatusCode < 200 && res.StatusCode > 302 {
		creationResult.Errors = append(creationResult.Errors, err.Error())
		return creationResult
	}

	defer res.Body.Close()
	var data AccountRequest
	jErr := json.NewDecoder(res.Body).Decode(&data)
	if jErr != nil {
		creationResult.Errors = append(creationResult.Errors, jErr.Error())
	}

	creationResult.Data = data
	return creationResult
}
