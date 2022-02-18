package accountapiclient

import (
	"accountapiclient"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var CurrnetVersion string
var apiAddress string

type Client interface {
	CreateAccount(account accountapiclient.AccountData) (accountapiclient.AccountData, error)
	FetchAccount(accountId string) (accountapiclient.AccountData, error)
	DeleteAccount(accountId string) error
}

type ApiClientV1 struct{}

func CreateClient(baseUrl string) Client {
	CurrnetVersion = "v1"
	if baseUrl[len(baseUrl)-1] != '/' {
		baseUrl = baseUrl + "/"
	}

	apiAddress = baseUrl + CurrnetVersion
	return ApiClientV1{}
}

func sendRequest(method string, endpoint string, requestBody string) ([]byte, int) {
	request, e := http.NewRequest(method, apiAddress+endpoint, bytes.NewBuffer([]byte(requestBody)))
	if e != nil {
		log.Print(e.Error())
		return nil, 400
	}

	httpclient := &http.Client{}
	response, e := httpclient.Do(request)
	if e != nil {
		log.Printf("error calling remote address: '%s'", e.Error())
		return nil, 503
	}

	defer response.Body.Close()
	responsecontent, e := ioutil.ReadAll(response.Body)
	if e != nil {
		log.Printf("error reading http response: '%s'", e.Error())
	}
	return responsecontent, response.StatusCode
}

func (ApiClientV1) CreateAccount(account accountapiclient.AccountData) (accountapiclient.AccountData, error) {
	if (account == accountapiclient.AccountData{}) {
		return accountapiclient.AccountData{}, errors.New("invalid input")
	}

	requestbody, e := json.Marshal(accountapiclient.Account{
		AccountData: account,
	})
	if e != nil {
		log.Print(e.Error())
		return accountapiclient.AccountData{}, errors.New("invalid input structure")
	}

	var newaccount accountapiclient.Account
	response, statuccode := sendRequest("POST", "/organisation/accounts", string(requestbody))
	if statuccode >= 200 && statuccode <= 299 {
		json.NewDecoder(bytes.NewBuffer(response)).Decode(&newaccount)
		return newaccount.AccountData, nil
	} else if statuccode >= 400 && statuccode <= 499 {
		var apierrors accountapiclient.ApiErrors
		json.NewDecoder(bytes.NewBuffer(response)).Decode(&apierrors)
		return accountapiclient.AccountData{}, errors.New(apierrors.ErrorMessage)
	}
	return accountapiclient.AccountData{}, errors.New("somethign went wrong, try again")
}

func (ApiClientV1) FetchAccount(accountId string) (accountapiclient.AccountData, error) {
	var account accountapiclient.Account
	if !isValidUUID(accountId) {
		return accountapiclient.AccountData{}, errors.New("invalid account id")
	}

	responsebody, statuscode := sendRequest("GET", "/organisation/accounts/"+accountId, "")
	if statuscode >= 200 && statuscode <= 299 {
		json.NewDecoder(bytes.NewBuffer(responsebody)).Decode(&account)
	} else {
		var apiError accountapiclient.ApiErrors
		json.NewDecoder(bytes.NewBuffer(responsebody)).Decode(&apiError)
		return accountapiclient.AccountData{}, errors.New(apiError.ErrorMessage)
	}

	return account.AccountData, nil
}

func (ApiClientV1) DeleteAccount(accountId string) error {
	if !isValidUUID(accountId) {
		return errors.New("invalid account id")
	}

	_, statuscode := sendRequest("DELETE", "/organisation/accounts/"+accountId+"?version=0", "")
	if statuscode == 204 {
		return nil
	}

	if statuscode == 404 {
		return errors.New("account cannot be found")
	}
	return errors.New("something went wrong, try again")
}

func isValidUUID(accountId string) bool {
	id, e := uuid.Parse(accountId)
	if id == uuid.Nil || e != nil {
		return false
	}
	return true
}
