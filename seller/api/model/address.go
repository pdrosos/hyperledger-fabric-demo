package model

type Address struct {
	Country  string `json:"country"`
	City     string `json:"city"`
	PostCode string `json:"postCode"`
	Address  string `json:"address"`
}
