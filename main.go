package main

import (
	"apiclient/client"
	"fmt"
)

func main() {
	fmt.Println("App is starting..")
	apiClient := client.InitializeClient("http://localhost:8080/", "v1")

	var myAccount client.OrganisationAccount
	myAccount = myAccount.WithBank("CACPA", "").WithCurrency("CAD").Build("CA", "Malek A")
	_, e := apiClient.CreateNewAccount(myAccount)
	if e != nil {
		fmt.Printf("Failed to create new account, %s ", e[:])
	}

	//c.FetchAccount("")
	//c.DeleteAccount("")
}
