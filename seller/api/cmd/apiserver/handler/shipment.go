package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/cmd/apiserver/response"
	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/logger"
	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/model"
	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/service"
	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/viewmodel"
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
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	shipment := &model.Shipment{}
	decodeError := decoder.Decode(shipment)
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

	err := this.shipmentService.Create(shipment)
	if err != nil {
		logger.Log.WithError(err).Errorf("Unable to create shipment")

		response.ResponseError(
			rw,
			http.StatusInternalServerError,
			"InternalServerError",
			"Can not create shipment",
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
