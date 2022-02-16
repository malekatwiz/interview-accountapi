package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type ApiClient interface {
	CreateNewAccount(orgAccount OrganisationAccount) (Account, []string)
	FetchAccount(accountId string) (AccountRequest, []string)
	DeleteAccount(accountId string) error
}

type ApiClientConnection struct{}

func InitializeClient(baseUrl string, apiVersion string) ApiClientConnection {
	apiUrl = baseUrl + "/" + apiVersion
	return ApiClientConnection{}
}

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
				Names:    account.HolderName[:],
			},
		},
	}
}

func (ApiClientConnection) CreateNewAccount(orgAccount OrganisationAccount) (Account, []string) {
	var errors []string
	if !orgAccount.IsReady() {
		return Account{}, append(errors, "Invalid empty input.")
	}
	request := mapToCreateAccount(orgAccount)

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

func (ApiClientConnection) DeleteAccount(accountId string) error {
	id, e := uuid.Parse(accountId)
	if e != nil || id == uuid.Nil {
		return errors.New("invalid account id")
	}
	_, statusCode, e := sendRequest("DELETE", "/organisation/accounts/"+accountId+"?version=1", "")
	if e != nil {
		log.Fatal(e)
		return e
	}

	if statusCode == 204 {
		return nil
	}

	if statusCode == 404 {
		return errors.New("account is not found")
	}

	return errors.New("something went wrong, try again")
}
