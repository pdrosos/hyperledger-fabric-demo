package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this *ShipmentChaincode) getShipmentById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		errorMessage := "Incorrect number of arguments. Expecting 1"
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	id := args[0]

	logger.Debugf("Start getShipmentById for shipment ID %s", id)

	shipmentAsBytes, err := stub.GetState(id)
	if err != nil {
		stateError := &Error{
			Error: fmt.Sprintf("Failed to get state for shipment ID %s", id),
		}
		jsonResponse, _ := json.Marshal(stateError)

		return shim.Error(string(jsonResponse))
	} else if shipmentAsBytes == nil {
		notFoundError := &Error{
			Error: fmt.Sprintf("Shipment ID %s does not exist", id),
		}
		jsonResponse, _ := json.Marshal(notFoundError)

		return pb.Response{Status: 404, Message: string(jsonResponse)}
	}

	logger.Debugf("End getShipmentById for shipment ID %s: %s", id, string(shipmentAsBytes))

	return shim.Success(shipmentAsBytes)
}
