package model

import (
	"time"
)

type Shipment struct {
	TrackingCode  string    `json:"trackingCode"`
	Courier       string    `json:"courier"`
	Sender        Sender    `json:"sender"`
	Recipient     Recipient `json:"recipient"`
	WeightInGrams int       `json:"weightInGrams"`
	Size
	Content             string     `json:"content"`
	ShippingType        string     `json:"shippingType"`
	IsFragile           bool       `json:"isFragile"`
	LastState           string     `json:"lastState"`
	LastLocation        *Address   `json:"lastLocation"`
	IsInCourierFacility bool       `json:"isInCourierFacility"`
	CreatedAt           *time.Time `json:"createdAt"`
	UpdatedAt           *time.Time `json:"updatedAt"`
}

type ShipmentRecord struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Shipments []ShipmentRecord
