package client

import (
	"testing"

	"github.com/google/uuid"
)

func TestApiHealth(t *testing.T) {
	apiStatus := GetApiStatus()
	if !apiStatus {
		t.Fail()
	}
}

func TestSetupAccount_ReturnsCreatedAccountId_WhenMinimumInputIsValid(t *testing.T) {
	r, l := ApiClient.SetupAccount(OrganisationAccount{
		Country: "CA",
		Name:    "John Doe Inc.",
	})

	if l != nil {
		t.Fail()
	}

	validId := uuid.MustParse(r.Id)
	if validId == uuid.Nil {
		t.Fail()
	}
}

func TestSetupAccount_Fails_WhenInputIsEmpty(t *testing.T) {
	r, l := ApiClient.SetupAccount(OrganisationAccount{})

	if l == nil || r.Id == uuid.Nil.String() {
		t.Fail()
	}
}

func TestSetupAccount_ReturnsValidationErrors_WhenCreationFails(t *testing.T) {
	_, l := ApiClient.SetupAccount(OrganisationAccount{})

	if len(l) == 0 {
		t.Fail()
	}
}
