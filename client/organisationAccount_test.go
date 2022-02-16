package client

import (
	"fmt"
	"testing"
)

func TestIsReady_ReturnsFalse_WhenModelIsNotBuilt(t *testing.T) {
	var a OrganisationAccount
	if a.IsReady() == true {
		t.Errorf("expected False, got True")
	}
	t.Cleanup(a.Cleanup)
}

func TestIsReady_ReturnsTrue_WhenModelIsBuilt(t *testing.T) {
	var a OrganisationAccount
	a = a.WithName("M").Build("CA", "Malek")
	if !a.IsReady() {
		t.Errorf("expected True, got False")
	}
	t.Cleanup(a.Cleanup)
}

func TestWithName_AddUpTo3NamesOnly(t *testing.T) {
	var a OrganisationAccount
	a = a.WithName("John").WithName("Jane").WithName("Joe").WithName("Jenny").Build("CA", "M")
	if len(a.HolderName) > 4 {
		t.Errorf("Expected max 4 items, got %s", fmt.Sprint(len(a.HolderName)))
	}
	t.Cleanup(a.Cleanup)
}

func TestWithName_AddAdditionalAccountHolderNames(t *testing.T) {
	var a OrganisationAccount
	name1 := "John"
	name2 := "M"
	a = a.WithName(name1).Build("CA", name2)
	if len(a.HolderName) != 2 {
		t.Errorf("epected 2 items, got '%s'", fmt.Sprint(len(a.HolderName)))
	}
	if a.HolderName[0] != name1 {
		t.Errorf("expected value '%s', got '%s'", name1, a.HolderName[0])
	}

	if a.HolderName[1] != name2 {
		t.Errorf("expected value '%s', got '%s'", name2, a.HolderName[1])
	}
	t.Cleanup(a.Cleanup)
}

func TestWithBank_AssignBankInfo(t *testing.T) {
	var a OrganisationAccount
	code := "Code"
	id := "Id"
	a = a.WithBank(id, code).Build("CA", "M")
	if a.BankCode != code {
		t.Errorf("expected value '%s', got '%s'", code, a.BankCode)
	}
	if a.BankId != id {
		t.Errorf("expected value '%s', got '%s'", id, a.BankId)
	}
	t.Cleanup(a.Cleanup)
}

func TestWithIban_AssignsIbanValue(t *testing.T) {
	var a OrganisationAccount
	iban := "iban"
	a = a.WithBank("id", "code").WithIban(iban).Build("CA", "M")
	if a.Iban != iban {
		t.Errorf("expected value '%s', got '%s'", iban, a.Iban)
	}
	t.Cleanup(a.Cleanup)
}

func TestWithBic_AssignsBicValue(t *testing.T) {
	var a OrganisationAccount
	bic := "bic"
	a = a.WithBank("id", "code").WithBic(bic).Build("CA", "M")
	if a.Bic != bic {
		t.Errorf("expected value '%s', got '%s'", bic, a.Bic)
	}
	t.Cleanup(a.Cleanup)
}

func TestWithAlternativeName_AssignsUpTo3AlternativeNames(t *testing.T) {
	var a OrganisationAccount
	a = a.WithBank("id", "code").WithBic("bic").WithAlternativeName("John").WithAlternativeName("Jane").WithAlternativeName("Joe").WithAlternativeName("Jenny").Build("CA", "M")
	if len(a.AlternativeNames) > 3 {
		t.Errorf("expected 3 items, got '%s'", fmt.Sprint(len(a.AlternativeNames)))
	}
	t.Cleanup(a.Cleanup)
}

func TestWithCurrency_AssignCurrencyCode(t *testing.T) {
	var a OrganisationAccount
	currency := "CAD"
	a = a.WithBank("id", "code").WithBic("bic").WithAlternativeName("Jenny").WithCurrency(currency).Build("CA", "M")
	if a.CurrencyCode != currency {
		t.Errorf("expected value '%s', got '%s'", currency, a.CurrencyCode)
	}
	t.Cleanup(a.Cleanup)
}

func TestWithAccountNumber_AssignsAccountNumber(t *testing.T) {
	var a OrganisationAccount
	accountNumber := "A83L8342397HFAB"
	a = a.WithBank("id", "code").WithBic("bic").WithAlternativeName("Jenny").WithCurrency("CAD").WithAccountNumber(accountNumber).Build("CA", "M")
	if a.AccountNumber != accountNumber {
		t.Errorf("expected value '%s', got '%s'", accountNumber, a.AccountNumber)
	}
	t.Cleanup(a.Cleanup)
}
