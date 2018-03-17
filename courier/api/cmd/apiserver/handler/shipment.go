package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/cmd/apiserver/response"
	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/inputmodel"
	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/logger"
	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/service"
	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/viewmodel"
)

type ShipmentHandler struct {
	shipmentService *service.ShipmentService
}

func NewShipmentHandler(shipmentService *service.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{
		shipmentService: shipmentService,
	}
}

func (this *ShipmentHandler) UpdateStateAndLocation() http.HandlerFunc {
	return this.updateStateAndLocation
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

func (this *ShipmentHandler) updateStateAndLocation(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackingCode := vars["trackingCode"]

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	shipmentState := &inputmodel.ShipmentState{}
	decodeError := decoder.Decode(shipmentState)
	if decodeError != nil {
		response.ResponseError(
			rw,
			http.StatusBadRequest,
			"BadRequest",
			"Can not decode json",
			make([]viewmodel.ErrorDetails, 0),
			decodeError,
		)

		return
	}

	// todo: validation

	err := this.shipmentService.UpdateStateAndLocation(trackingCode, shipmentState)
	if err != nil {
		logger.Log.WithError(err).Errorf("Unable to update shipment state and location")

		response.ResponseError(
			rw,
			http.StatusInternalServerError,
			"InternalServerError",
			"Can not update shipment state and location",
			make([]viewmodel.ErrorDetails, 0),
			err,
		)

		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write(nil)
}

func (this *ShipmentHandler) getAll(rw http.ResponseWriter, r *http.Request) {
	shipments, err := this.shipmentService.GetAll()
	if err != nil {
		response.ResponseError(
			rw,
			http.StatusInternalServerError,
			"InternalServerError",
			"Can not get shipments",
			make([]viewmodel.ErrorDetails, 0),
			err,
		)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(shipments)
	if err != nil {
		response.ResponseError(
			rw,
			http.StatusInternalServerError,
			"InternalServerError",
			"Can not encode json",
			make([]viewmodel.ErrorDetails, 0),
			err,
		)

		return
	}

	rw.Write(jsonResponse)
}

func (this *ShipmentHandler) getByTrackingCode(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackingCode := vars["trackingCode"]

	shipment, err := this.shipmentService.GetByTrackingCode(trackingCode)
	if err != nil {
		response.ResponseError(
			rw,
			http.StatusInternalServerError,
			"InternalServerError",
			"Can not get shipment",
			make([]viewmodel.ErrorDetails, 0),
			err,
		)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(shipment)
	if err != nil {
		response.ResponseError(
			rw,
			http.StatusInternalServerError,
			"InternalServerError",
			"Can not encode json",
			make([]viewmodel.ErrorDetails, 0),
			err,
		)

		return
	}

	rw.Write(jsonResponse)
}

func (this *ShipmentHandler) getHistory(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackingCode := vars["trackingCode"]

	shipmentHistory, err := this.shipmentService.GetHistory(trackingCode)
	if err != nil {
		response.ResponseError(
			rw,
			http.StatusInternalServerError,
			"InternalServerError",
			"Can not get shipment history",
			make([]viewmodel.ErrorDetails, 0),
			err,
		)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(shipmentHistory)
	if err != nil {
		response.ResponseError(
			rw,
			http.StatusInternalServerError,
			"InternalServerError",
			"Can not encode json",
			make([]viewmodel.ErrorDetails, 0),
			err,
		)

		return
	}

	rw.Write(jsonResponse)
}
