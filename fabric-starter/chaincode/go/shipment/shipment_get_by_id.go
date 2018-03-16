package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this *ShipmentChaincode) getShipmentById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// validate arguments
	err = this.validateArgumentsNotEmpty(1, args)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]

	logger.Debugf("Start getShipmentById for shipment ID %s", id)

	shipmentAsBytes, err := stub.GetState(id)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Failed to get state for shipment ID %s", id))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	} else if shipmentAsBytes == nil {
		errorJson := this.errorJson(fmt.Sprintf("Shipment ID %s does not exist", id))
		logger.Error(errorJson)

		return pb.Response{Status: 404, Message: errorJson}
	}

	logger.Debugf("End getShipmentById for shipment ID %s: %s", id, string(shipmentAsBytes))

	return shim.Success(shipmentAsBytes)
}
