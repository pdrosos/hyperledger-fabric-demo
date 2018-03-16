package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

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

	// get shipment from state
	trackingCode := args[0]

	logger.Debugf("Start changeShipmentStateAndLocation for shipment %s", trackingCode)

	shipmentAsBytes, err := stub.GetState(trackingCode)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for shipment " + trackingCode + "\"}"
		logger.Error(jsonResp)

		return shim.Error(jsonResp)
	} else if shipmentAsBytes == nil {
		jsonResp := "{\"Error\":\"Shipment does not exist: " + trackingCode + "\"}"
		logger.Error(jsonResp)

		return shim.Error(jsonResp)
	}

	// unmarshal and change shipment state, location and updated at date
	shipment := Shipment{}
	err = json.Unmarshal(shipmentAsBytes, shipment)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to unmarshal shipment %s bytes", trackingCode))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	lastState := args[1]
	shipment.LastState = lastState

	lastLocation := Address{}
	err = json.Unmarshal([]byte(args[2]), &lastLocation)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to unmarshal shipment %s LastLocation data", trackingCode))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	shipment.LastLocation = &lastLocation
	shipment.UpdatedAt = time.Now().UTC()

	// marshal shipment
	shipmentJSONAsBytes, err := json.Marshal(shipment)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to marshal shipment %s as JSON: %s", trackingCode, err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	// save shipment to state
	err = stub.PutState(trackingCode, shipmentJSONAsBytes)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to put shipment %s JSON bytes to state: %s", trackingCode, err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	logger.Debugf("End changeShipmentStateAndLocation for shipment %s", trackingCode)

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
