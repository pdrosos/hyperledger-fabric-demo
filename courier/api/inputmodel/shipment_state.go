package inputmodel

type ShipmentState struct {
	State    string `json:"state"`
	Location struct {
		Country  string `json:"country"`
		City     string `json:"city"`
		PostCode string `json:"postCode"`
		Address  string `json:"address"`
	} `json:"location"`
	IsDelivered bool `json:"isDelivered"`
}
