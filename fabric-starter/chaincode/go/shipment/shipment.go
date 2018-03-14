package main

import (
	"time"
)

type Shipment struct {
	Id                  string    `json:"id"`
	TrackingCode        string    `json:"trackingCode"`
	Courier             Courier   `json:"courier"`
	Sender              Sender    `json:"sender"`
	Recipient           Recipient `json:"recipient"`
	WeightInGrams       int       `json:"weightInGrams"`
	ShippingType        string    `json:"shippingType"`
	Size                Size      `json:"size"`
	Content             string    `json:"content"`
	IsFragile           bool      `json:"isFragile"`
	LastState           string    `json:"lastState"`
	LastLocation        *Address  `json:"lastLocation"`
	IsInCourierFacility bool      `json:"isInCourierFacility"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

func NewShipment(
	id string,
	trackingCode string,
	courier Courier,
	sender Sender,
	recipient Recipient,
	weightInGrams int,
	shippingType string,
	size Size,
	content string,
	isFragile bool,
	lastState string,
) *Shipment {
	createdAt := time.Now().UTC()

	shipment := &Shipment{
		Id:                  id,
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
		UpdatedAt:           createdAt,
	}

	return shipment
}
