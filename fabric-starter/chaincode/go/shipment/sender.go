package main

type Sender struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Address   Address `json:"address"`
	Phone     string  `json:"phone"`
}
