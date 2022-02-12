package main

import (
	"apiclient/client"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("App is starting..")
	apiIsHealthy := client.GetApiStatus()
	if !apiIsHealthy {
		fmt.Println("Accounts API is down, exiting..")
		os.Exit(1)
	}
	fmt.Println("Accounts API is healthy.")
	fmt.Println("Creating new account...")

	var req client.CreateAccountRequest
	req.AccountData.Id = uuid.New().String()
	req.AccountData.OrganisationId = uuid.New().String()
	req.AccountData.Type = "accounts"
	req.AccountData.Attributes.Country = "GB"
	req.AccountData.Attributes.Currency = "GBP"
	req.AccountData.Attributes.BankId = "400302"
	req.AccountData.Attributes.BankCode = "GBDSC"
	req.AccountData.Attributes.Names = append(req.AccountData.Attributes.Names, "Malek A")
	res := client.CreateAccount(req)
	fmt.Println("Account resource: " + res.Data.Links.CurrentAccount)

	fmt.Println("Fetching resource: " + res.Data.Links.CurrentAccount)
	fetchAccount := client.FetchAccount(res.Data.Links.CurrentAccount)
	if fetchAccount.Errors != nil {
		println("cannot find resource " + res.Data.Links.CurrentAccount)
	} else {
		println("Resource ID " + fetchAccount.Data.AccountData.Id + "found")
	}
}
