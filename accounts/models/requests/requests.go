package Requests

type CreateAccountRequest struct {
	AccountData AccountData `json:"data"`
}

type AccountData struct {
	Type           string            `json:"type"`
	Id             string            `json:"id"`
	OrganisationId string            `json:"organisation_id"`
	Attributes     AccountAttributes `json:"attributes"`
}

type AccountAttributes struct {
	Country               string `json:"country"`
	Currency              string `json:"base_currency"`
	BankId                string `json:"bank_id"`
	BankCode              string `json:"bank_id_code"`
	Bic                   string `json:"bic"`
	AccountNumber         string `json:"account_number"`
	CustomerId            string `json:"customer_id"`
	Iban                  string `json:"iban"`
	AccountClassification string `json:"account_classification"`
	//ValidationType      string   `json:"validation_type"`
	//ReferenceMask       string   `json:"reference_mask"`
	//AcceptanceQualifier string   `json:"acceptance_qualifier"`
	Name []string `json:"name"`
	/*UserData            []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"user_defined_data"`*/
}
