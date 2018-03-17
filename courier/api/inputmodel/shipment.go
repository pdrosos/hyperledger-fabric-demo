package inputmodel

type ShipmentRecord struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Shipments []ShipmentRecord
