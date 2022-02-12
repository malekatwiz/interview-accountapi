package client

import (
	"time"
)

type AccountAttributes struct {
	Country               string   `json:"country"`
	Currency              string   `json:"base_currency"`
	BankId                string   `json:"bank_id"`
	BankCode              string   `json:"bank_id_code"`
	Bic                   string   `json:"bic"`
	AccountNumber         string   `json:"account_number"`
	CustomerId            string   `json:"customer_id"`
	Iban                  string   `json:"iban"`
	AccountClassification string   `json:"account_classification,omitempty"`
	Names                 []string `json:"name"`
	AlternativeNames      []string `json:"alternative_names,omitempty"`
}

type AccountData struct {
	Type           string            `json:"type"`
	Id             string            `json:"id"`
	OrganisationId string            `json:"organisation_id"`
	Version        int               `json:"version"`
	CreatedOn      time.Time         `json:"createdOn"`
	ModifiedOn     time.Time         `json:"modified_on"`
	Attributes     AccountAttributes `json:"attributes"`
}

type AccountRequest struct {
	AccountData AccountData  `json:"data"`
	Links       AccountLinks `json:"links"`
}

type CreateAccountRequest struct {
	AccountData AccountData `json:"data"`
}

type AccountLinks struct {
	FirstAccount   string `json:"first"`
	LastAccount    string `json:"last"`
	CurrentAccount string `json:"self"`
}

type AccountsApiResponse struct {
	Data   AccountRequest
	Errors []string
}
