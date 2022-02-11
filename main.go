package main

import (
	"fmt"
)

func main() {
	fmt.Println("App is starting..")
	/*apiIsHealthy := accountsclient.GetApiStatus()
	if !apiIsHealthy {
		fmt.Println("Accounts API is down, exiting..")
	}*/
	fmt.Println("Checking Accounts API health")
}
