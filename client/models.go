package client

type Account struct {
	Id string
}

type ApiError struct {
	ErrorMessage string `json:"error_message"`
}
