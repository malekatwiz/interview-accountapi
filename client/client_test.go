package client

import "testing"

func TestApiHealth(t *testing.T) {
	apiStatus := GetApiStatus()
	if !apiStatus {
		t.Fail()
	}
}
