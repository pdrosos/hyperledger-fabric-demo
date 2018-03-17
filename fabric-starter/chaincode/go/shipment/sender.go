package main

type Sender struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address
	Phone string `json:"phone"`
}
