package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupClient_ReturnsNoErrors_WhenCreationIsSuccessful(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(`{
			"data": {
				"attributes": {
					"alternative_names": null,
					"country": "CA",
					"name": [
						"M"
					]
				},
				"created_on": "2022-02-14T14:11:46.906Z",
				"id": "49dac5ee-6ffb-4bb3-a24d-9c36d4f4ca36",
				"modified_on": "2022-02-14T14:11:46.906Z",
				"organisation_id": "78398917-e6bd-4671-bc99-666c5015af99",
				"type": "accounts",
				"version": 0
			},
			"links": {
				"self": "/v1/organisation/accounts/49dac5ee-6ffb-4bb3-a24d-9c36d4f4ca36"
			}
		}`))
	}))
	defer server.Close()

	apiClient := InitializeClient(server.URL, "v1")
	r, l := apiClient.CreateNewAccount(OrganisationAccount{
		Country: "CA",
		Name:    "John Doe Inc.",
	})

	if (len(l) > 0 || r == Account{}) {
		t.Errorf("expected no errors but received %s", string(rune(len(l))))
	}
}

func TestSetupAccount_ReturnsEmptyAccountWithError_WhenInputIsEmpty(t *testing.T) {
	apiClient := InitializeClient("http://localhost:8080/", "v1")
	r, l := apiClient.CreateNewAccount(OrganisationAccount{})

	if (l == nil || r != Account{}) {
		t.Errorf("expected errors")
	}
}

func TestSetupAccount_ReturnsValidationErrors_WhenInputInvalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(400)
		rw.Write([]byte(""))
	}))
	defer server.Close()
	apiClient := InitializeClient(server.URL, "v1")
	_, l := apiClient.CreateNewAccount(OrganisationAccount{
		Country: "AAAA",
		Name:    "John Doe Inc.",
	})

	if len(l) <= 0 {
		t.Errorf("expected errors")
	}
}

func TestSetupAccount_ReturnsError_WhenCreateRequestFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(500)
		rw.Write([]byte(""))
	}))
	defer server.Close()
	apiClient := InitializeClient(server.URL, "v1")
	r, l := apiClient.CreateNewAccount(OrganisationAccount{
		Country: "CA",
		Name:    "John Doe Inc.",
	})

	if len(l) <= 0 {
		t.Errorf("expected errors, received %s", string(rune(len(l))))
	}

	if (r != Account{}) {
		t.Errorf("expected empty Account")
	}
}

/*
func TestSetupAccount_ReturnsCreatedAccountId_WhenMinimumInputIsValid(t *testing.T) {
	apiClient := InitializeClient("", "")
	apiClient.CreateNewAccount()
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

func TestSetupAccount_ReturnsEmptyAccountWithError_WhenInputIsEmpty(t *testing.T) {
	r, l := ApiClient.SetupAccount(OrganisationAccount{})

	if (l == nil || r != Account{}) {
		t.Fail()
	}
}

func TestSetupAccount_ReturnsError_WhenAccountInfoMissing(t *testing.T) {
	_, l := ApiClient.SetupAccount(OrganisationAccount{})
	if len(l) != 1 {
		t.Fail()
	}
}

func TestSetupAccount_ReturnsValidationErrors_WhenInputFieldInvalid(t *testing.T) {
	_, l := ApiClient.SetupAccount(OrganisationAccount{
		Country: "AAA",
	})
	if len(l) <= 0 {
		t.Fail()
	}
}*/
