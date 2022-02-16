package client

type OrganisationAccount struct {
	Country          string
	BankId           string
	BankCode         string
	Bic              string
	Iban             string
	HolderName       []string
	AlternativeNames []string
	CurrencyCode     string
	AccountNumber    string
}

var orgAccount OrganisationAccount

const maxNames = 4
const maxAltNames = 3

var isReady bool

func (OrganisationAccount) IsReady() bool {
	return isReady
}

// Maximum of 4 calls/Names extra items will be discarded.
func (OrganisationAccount) WithName(holderName string) OrganisationAccount {
	if len(orgAccount.HolderName) < maxNames-1 {
		orgAccount.HolderName = append(orgAccount.HolderName, holderName)
	}
	return OrganisationAccount{}
}

// Maximum of 3 calls/Names extra items will be discarded.
func (OrganisationAccount) WithAlternativeName(name string) OrganisationAccount {
	if len(orgAccount.AlternativeNames) < maxAltNames {
		orgAccount.AlternativeNames = append(orgAccount.AlternativeNames, name)
	}
	return OrganisationAccount{}
}

func (OrganisationAccount) WithBank(bankId string, bankCode string) OrganisationAccount {
	orgAccount.BankId = bankId
	orgAccount.BankCode = bankCode
	return OrganisationAccount{}
}

func (OrganisationAccount) WithCurrency(currencyCode string) OrganisationAccount {
	orgAccount.CurrencyCode = currencyCode
	return OrganisationAccount{}
}

func (OrganisationAccount) WithAccountNumber(accountNumber string) OrganisationAccount {
	orgAccount.AccountNumber = accountNumber
	return OrganisationAccount{}
}

func (OrganisationAccount) WithBic(bic string) OrganisationAccount {
	orgAccount.Bic = bic
	return OrganisationAccount{}
}

func (OrganisationAccount) WithIban(iban string) OrganisationAccount {
	orgAccount.Iban = iban
	return OrganisationAccount{}
}

func (OrganisationAccount) Build(countryCode string, holderName string) OrganisationAccount {
	orgAccount.Country = countryCode
	orgAccount.WithName(holderName)
	isReady = true
	return orgAccount
}

func (OrganisationAccount) Cleanup() {
	orgAccount = OrganisationAccount{}
	isReady = false
}
