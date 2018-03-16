package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this *ShipmentChaincode) changeShipmentStateAndLocation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// validate transaction creator
	err = this.requireCourier1Creator(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	// validate arguments
	err = this.validateArgumentsNotEmpty(3, args)
	if err != nil {
		return shim.Error(err.Error())
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
		errorJson := this.errorJson(fmt.Sprintf("Unable to unmarshal shipment ID %s bytes", id))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	lastState := args[1]
	shipment.LastState = lastState

	lastLocation := Address{}
	err = json.Unmarshal([]byte(args[2]), &lastLocation)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to unmarshal shipment ID %s LastLocation data", id))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	shipment.LastLocation = &lastLocation

	shipmentJSONAsBytes, err := json.Marshal(shipment)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to marshal shipment ID %s as JSON: %s", id, err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	// save shipment to state
	err = stub.PutState(id, shipmentJSONAsBytes)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to put shipment ID %s JSON bytes to state: %s", id, err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	logger.Debugf("End changeShipmentStateAndLocation for shipment ID %s", id)

	return shim.Success(nil)
}

func (this *ShipmentChaincode) requireCourier1Creator(stub shim.ChaincodeStubInterface) error {
	// only courier1 organization can change shipments state and location
	creatorBytes, err := this.getTransactionCreator(stub)
	if err != nil {
		return err
	}

	name, org := this.getCreator(creatorBytes)
	if org != "courier1" {
		errorJson := this.errorJson(fmt.Sprintf("Only courier1 organization can create shipments, called by %s@s", name, org))
		logger.Error(errorJson)

		return errors.New(errorJson)
	}

	return nil
}
