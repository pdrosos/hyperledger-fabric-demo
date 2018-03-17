package main

import (
	"time"
)

type Shipment struct {
	TrackingCode        string    `json:"trackingCode"`
	Courier             string    `json:"courier"`
	Sender              Sender    `json:"sender"`
	Recipient           Recipient `json:"recipient"`
	WeightInGrams       int       `json:"weightInGrams"`
	Size
	Content             string    `json:"content"`
	ShippingType        string    `json:"shippingType"`
	IsFragile           bool      `json:"isFragile"`
	LastState           string    `json:"lastState"`
	LastLocation        *Address  `json:"lastLocation"`
	IsInCourierFacility bool      `json:"isInCourierFacility"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

type ShipmentRecord struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Shipments []ShipmentRecord

func NewShipment(
	trackingCode string,
	courier string,
	sender Sender,
	recipient Recipient,
	weightInGrams int,
	shippingType string,
	size Size,
	content string,
	isFragile bool,
	lastState string,
	createdAt time.Time,
	updatedAt time.Time,
) *Shipment {
	shipment := &Shipment{
		TrackingCode:        trackingCode,
		Courier:             courier,
		Sender:              sender,
		Recipient:           recipient,
		WeightInGrams:       weightInGrams,
		ShippingType:        shippingType,
		Size:                size,
		Content:             content,
		IsFragile:           isFragile,
		LastState:           lastState,
		IsInCourierFacility: false,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}

	return shipment
}
