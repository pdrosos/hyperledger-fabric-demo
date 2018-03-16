package handler

import (
	"net/http"

	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/service"
)

type ShipmentHandler struct {
	shipmentService *service.ShipmentService
}

func NewShipmentHandler(shipmentService *service.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{
		shipmentService: shipmentService,
	}
}

func (this *ShipmentHandler) Create() http.HandlerFunc {
	return this.create
}

func (this *ShipmentHandler) GetAll() http.HandlerFunc {
	return this.getAll
}

func (this *ShipmentHandler) GetByTrackingCode() http.HandlerFunc {
	return this.getByTrackingCode
}

func (this *ShipmentHandler) GetHistory() http.HandlerFunc {
	return this.getHistory
}

func (this *ShipmentHandler) create(rw http.ResponseWriter, r *http.Request) {

}

func (this *ShipmentHandler) getAll(rw http.ResponseWriter, r *http.Request) {

}

func (this *ShipmentHandler) getByTrackingCode(rw http.ResponseWriter, r *http.Request) {

}

func (this *ShipmentHandler) getHistory(rw http.ResponseWriter, r *http.Request) {

}
