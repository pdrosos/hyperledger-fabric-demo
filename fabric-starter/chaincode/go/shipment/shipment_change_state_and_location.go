package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
	err = this.validateArgumentsNotEmpty(5, args)
	if err != nil {
		return shim.Error(err.Error())
	}

	// get shipment from state
	trackingCode := args[0]

	logger.Debugf("Start changeShipmentStateAndLocation for shipment %s", trackingCode)

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

	// unmarshal and change shipment state, location and updated at date
	shipment := Shipment{}
	err = json.Unmarshal(shipmentAsBytes, &shipment)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to unmarshal shipment %s bytes", trackingCode))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	// cannot change state of delivered shipment
	if shipment.IsDelivered {
		errorJson := this.errorJson(fmt.Sprintf("Can not change state of delivered shipment %s", trackingCode))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	lastState := args[1]

	lastLocation := Address{}
	err = json.Unmarshal([]byte(args[2]), &lastLocation)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to unmarshal shipment %s location data", trackingCode))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	isDelivered, err := strconv.ParseBool(args[3])
	if err != nil {
		errorJson := this.errorJson("4th argument isDelivered must be a boolean representing string")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	updatedAt, err := time.Parse(time.RFC3339Nano, args[4])
	if err != nil {
		errorJson := this.errorJson("5th argument updatedAt must be time.RFC3339Nano formatted string")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	// change shipment state and location
	shipment.LastState = lastState
	shipment.LastLocation = &lastLocation
	shipment.IsDelivered = isDelivered
	shipment.UpdatedAt = updatedAt

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

	// send shipment-state-and-location-changed event to the SDK
	if shipment.IsDelivered {
		stub.SetEvent("shipment-delivered", shipmentJSONAsBytes)
	} else {
		stub.SetEvent("shipment-state-and-location-changed", shipmentJSONAsBytes)
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
