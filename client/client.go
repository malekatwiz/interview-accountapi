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
	CreateNewAccount(orgAccount OrganisationAccount) (Account, []string)
}

type ApiClientConnection struct{}

func InitializeClient(baseUrl string, apiVersion string) ApiClientConnection {
	apiUrl = baseUrl + "/" + apiVersion
	return ApiClientConnection{}
}

const baseUrl = "http://localhost:8080/" + ApiVersion
const ApiVersion = "v1"

var apiUrl string

func sendRequest(method string, endpoint string, reqBody string) ([]byte, int, error) {
	req, e := http.NewRequest(method, apiUrl+endpoint, bytes.NewBuffer([]byte(reqBody)))
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

func (ApiClientConnection) CreateNewAccount(orgAccount OrganisationAccount) (Account, []string) {
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
		return Account{}, errors
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

	if len(errors) > 0 {
		return Account{}, errors
	}
	return Account{Id: account.AccountData.Id}, errors
}

func (ApiClientConnection) FetchAccount(accountId string) (AccountRequest, []string) {
	var account AccountRequest
	id, e := uuid.Parse(accountId)
	if e != nil || id == uuid.Nil {
		return account, []string{"invalid account id"}
	}

	resBody, statusCode, e := sendRequest("GET", "/organisation/accounts/"+accountId, "")
	if e != nil {
		log.Fatal(e.Error())
		return account, strings.Split(e.Error(), "\n")
	}
	if statusCode >= 200 && statusCode <= 299 { // successful creation
		json.NewDecoder(bytes.NewBuffer(resBody)).Decode(&account)
	} else {
		var apiErr ApiError
		json.NewDecoder(bytes.NewBuffer(resBody)).Decode(&apiErr)
		return account, strings.Split(apiErr.ErrorMessage, "\n")
	}
	return account, nil
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
