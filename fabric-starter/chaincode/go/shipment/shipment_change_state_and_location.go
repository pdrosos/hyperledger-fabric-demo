package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this ShipmentChaincode) changeShipmentStateAndLocation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		errorMessage := "Incorrect number of arguments. Expecting 3"
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	for index, element := range args {
		if len(element) <= 0 {
			errorMessage := fmt.Sprintf("Argument %s must be a non-empty string", index+1)
			logger.Error(errorMessage)

			return shim.Error(errorMessage)
		}
	}

	id := args[0]

	logger.Debugf("Start changeShipmentStateAndLocation for shipment ID %s", id)

	shipmentAsBytes, err := stub.GetState(id)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for shipment " + id + "\"}"
		logger.Error(jsonResp)

		return shim.Error(jsonResp)
	} else if shipmentAsBytes == nil {
		jsonResp := "{\"Error\":\"Shipment does not exist: " + id + "\"}"
		logger.Error(jsonResp)

		return shim.Error(jsonResp)
	}

	shipment := Shipment{}
	err = json.Unmarshal(shipmentAsBytes, shipment)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to unmarshal shipment ID %s bytes", id)
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	lastState := args[1]
	shipment.LastState = lastState

	lastLocation := Address{}
	err = json.Unmarshal([]byte(args[2]), &lastLocation)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to unmarshal shipment ID %s LastLocation data", id)
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	shipment.LastLocation = &lastLocation

	shipmentJSONAsBytes, err := json.Marshal(shipment)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to marshal shipment ID %s as JSON: %s", id, err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	// save shipment to state
	err = stub.PutState(id, shipmentJSONAsBytes)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to put shipment ID %s JSON bytes to state: %s", id, err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	logger.Debugf("End changeShipmentStateAndLocation for shipment ID %s", id)

	return shim.Success(nil)
}
