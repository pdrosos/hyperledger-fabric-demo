package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/service"
)

func Register(channelClient *channel.Client) {
	router := mux.NewRouter()

	shipmentService := service.NewShipmentService(channelClient)
	shipmentHandler := NewShipmentHandler(shipmentService)

	router.Handle("/shipments/{trackingCode}/state", shipmentHandler.UpdateStateAndLocation()).Methods("PATCH")

	router.Handle("/shipments", shipmentHandler.GetAll()).Methods("GET")

	router.Handle("/shipments/{trackingCode}", shipmentHandler.GetByTrackingCode()).Methods("GET")

	router.Handle("/shipments/{trackingCode}/history", shipmentHandler.GetHistory()).Methods("GET")

	// default route
	router.Handle("/", NewRootHandler()).Methods("GET", "HEAD")

	http.Handle("/", router)
}
