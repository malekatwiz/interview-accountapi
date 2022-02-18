package accountapiclient

import (
	"accountapiclient"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

const host = "http://accountapi:8080/"

func TestCreateAccount_ReturnsCreatedAccount(t *testing.T) {
	client := CreateClient(host)

	country := "CA"
	version := int64(1)
	account, _ := client.CreateAccount(accountapiclient.AccountData{
		ID:             uuid.NewString(),
		OrganisationID: uuid.NewString(),
		Type:           "accounts",
		Version:        &version,
		Attributes: &accountapiclient.AccountAttributes{
			Country: &country,
			Name:    []string{"Malek"},
		},
	})

	if (account == accountapiclient.AccountData{}) {
		t.Fatalf("expected valid account type, received '%s'", fmt.Sprint(account))
	}

	t.Run("TestFetchAccount", func(t *testing.T) {
		foundAccount, e := client.FetchAccount(account.ID)
		if (e != nil || foundAccount == accountapiclient.AccountData{}) {
			t.Fatalf("expected account with no errors")
		}

		t.Run("TestDeleteAccount", func(t *testing.T) {
			e := client.DeleteAccount(account.ID)
			if e != nil {
				t.Fatalf("expected no errors, received '%s'", e.Error())
			}

			a, e := client.FetchAccount(account.ID)
			if (e == nil || a != accountapiclient.AccountData{}) {
				t.Fatalf("did not expect to find account '%s', supposed to be deleted", fmt.Sprint(account.ID))
			}
		})
	})
}

func TestCreateAccount_ReturnsError_WhenApiIsDown(t *testing.T) {
	// WARNING, it is expected for the request to timeout
	client := CreateClient("http://localhost:9292")
	country := "CA"
	version := int64(1)
	_, e := client.CreateAccount(accountapiclient.AccountData{
		ID:             uuid.NewString(),
		OrganisationID: uuid.NewString(),
		Type:           "accounts",
		Version:        &version,
		Attributes: &accountapiclient.AccountAttributes{
			Country: &country,
			Name:    []string{"Malek"},
		},
	})

	if e == nil {
		t.Errorf("expected remote host error")
	}
}
