package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this *ShipmentChaincode) getShipmentByTrackingCode(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// validate arguments
	err = this.validateArgumentsNotEmpty(1, args)
	if err != nil {
		return shim.Error(err.Error())
	}

	trackingCode := args[0]

	logger.Debugf("Start getShipmentByTrackingCode for shipment %s", trackingCode)

	shipmentAsBytes, err := stub.GetState(trackingCode)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Failed to get state for shipment %s", trackingCode))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	} else if shipmentAsBytes == nil {
		errorJson := this.errorJson(fmt.Sprintf("Shipment %s does not exist", trackingCode))
		logger.Error(errorJson)

		return pb.Response{Status: 404, Message: errorJson}
	}

	logger.Debugf("End getShipmentByTrackingCode for shipment %s: %s", trackingCode, string(shipmentAsBytes))

	return shim.Success(shipmentAsBytes)
}
