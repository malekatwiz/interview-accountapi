package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type ApiClient interface {
	SetupAccount() (Account, []string)
}

const baseUrl = "http://localhost:8080/" + ApiVersion
const ApiVersion = "v1"

func GetApiStatus() bool {
	_, statusCode, e := sendRequest("GET", "/health", "")
	if e == nil && statusCode == 200 {
		return true
	}
	return false
}

func sendRequest(method string, endpoint string, reqBody string) ([]byte, int, error) {
	req, e := http.NewRequest(method, baseUrl+endpoint, bytes.NewBuffer([]byte(reqBody)))
	if e != nil {
		log.Fatal(e.Error())
	}

	httpClient := &http.Client{}
	res, e := httpClient.Do(req)
	if e != nil {
		log.Fatal(e.Error())
	}

	defer res.Body.Close()
	resContent, e := ioutil.ReadAll(res.Body)
	if e != nil {
		log.Fatal(e)
	}

	return resContent, res.StatusCode, nil
}

func mapToCreateAccount(account OrganisationAccount) AccountRequest {
	if (account == OrganisationAccount{} || account.Country == "") {
		return AccountRequest{}
	}
	return AccountRequest{
		AccountData: AccountData{
			Id:             uuid.NewString(),
			OrganisationId: uuid.NewString(),
			Type:           "accounts",
			Version:        1,
			Attributes: AccountAttributes{
				Country:  account.Country,
				BankId:   account.BankId,
				BankCode: account.BankCode,
				Names:    []string{account.Name},
			},
		},
	}
}

func (orgAccount OrganisationAccount) SetupAccount() (Account, []string) {
	if (orgAccount == OrganisationAccount{}) {
		return Account{}, []string{"Invalid empty input."}
	}
	request := mapToCreateAccount(orgAccount)

	var errors []string
	var account AccountRequest
	reqBody, e := json.Marshal(request)
	if e != nil {
		log.Fatal(e)
		errors = append(errors, e.Error())
		return Account{account.AccountData.Id}, errors
	}

	resBody, statusCode, e := sendRequest("POST", "/organisation/accounts", string(reqBody))
	if e != nil {
		log.Fatal(e.Error()) // TODO: might be unnecessary
	}

	// casting from/to
	if statusCode >= 200 && statusCode <= 299 { // successful creation
		json.NewDecoder(bytes.NewBuffer(resBody)).Decode(&account)
	} else if statusCode >= 400 && statusCode <= 499 { // bad request
		var apiErr ApiError
		json.NewDecoder(bytes.NewBuffer(resBody)).Decode(&apiErr)
		errors = append(errors, strings.Split(apiErr.ErrorMessage, "\n")...)
	} else {
		errors = append(errors, "Something went worng, try again.")
	}
	return Account{Id: account.AccountData.Id}, errors
}

func CreateAccount(request CreateAccountRequest) AccountsApiResponse {
	var creationResult AccountsApiResponse
	reqBody, jsonErr := json.Marshal(request) // TODO: handle json errors.
	if jsonErr != nil {
		creationResult.Errors = append(creationResult.Errors, "Invalid JSON input")
		return creationResult
	}

	res, err := http.Post(baseUrl+"/organisation/accounts", "application/vnd.api+json", bytes.NewBuffer(reqBody))
	if err != nil || res.StatusCode < 200 || res.StatusCode > 299 {
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

func FetchAccount(resourceLocation string) AccountsApiResponse {
	var response AccountsApiResponse
	res, err := http.Get(baseUrl + resourceLocation)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
	}

	defer res.Body.Close()

	var account AccountRequest
	json.NewDecoder(res.Body).Decode(&account)
	response.Data = account
	return response
}

func DeleteAccount(resourceId string) error {
	req, e := http.NewRequest("DELETE", baseUrl+"/organisation/accounts/"+resourceId+"?version=0", bytes.NewBuffer(nil))
	if e != nil {
		return e
	}
	httpClient := &http.Client{}
	httpClient.Do(req)
	return nil
}
