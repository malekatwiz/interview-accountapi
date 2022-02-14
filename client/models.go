package client

type OrganisationAccount struct {
	Country  string
	BankId   string
	BankCode string
	Bic      string
	Name     string
}

type Account struct {
	Id string
}

type ApiError struct {
	ErrorMessage string `json:"error_message"`
}
