package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/service"
)

func Register(channelClient *channel.Client) {
	router := mux.NewRouter()

	shipmentService := service.NewShipmentService(channelClient)
	shipmentHandler := NewShipmentHandler(shipmentService)

	router.Handle("/shipments", shipmentHandler.Create()).Methods("POST")

	router.Handle("/shipments", shipmentHandler.GetAll()).Methods("GET")

	router.Handle("/shipments/{trackingCode}", shipmentHandler.GetByTrackingCode()).Methods("GET")

	router.Handle("/shipments/history/{trackingCode}", shipmentHandler.GetHistory()).Methods("GET")

	// default route
	router.Handle("/", NewRootHandler()).Methods("GET", "HEAD")

	// not found
	//router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	http.Handle("/", router)
}
