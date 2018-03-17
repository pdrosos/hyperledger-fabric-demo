package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this *ShipmentChaincode) createShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("Start createShipment")

	var err error

	// validate transaction creator
	err = this.requireSellerCreator(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	// validate arguments
	err = this.validateArgumentsNotEmpty(10, args)
	if err != nil {
		return shim.Error(err.Error())
	}

	trackingCode := args[0]
	courier := args[1]

	sender := Sender{}
	err = json.Unmarshal([]byte(args[2]), &sender)
	if err != nil {
		errorJson := this.errorJson("Unable to unmarshal Sender data")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	recipient := Recipient{}
	err = json.Unmarshal([]byte(args[3]), &recipient)
	if err != nil {
		errorJson := this.errorJson("Unable to unmarshal Recipient data")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	weightInGrams, err := strconv.Atoi(args[4])
	if err != nil {
		errorJson := this.errorJson("5th argument weightInGrams must be a numeric string")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	shippingType := args[5]

	size := Size{}
	err = json.Unmarshal([]byte(args[6]), &size)
	if err != nil {
		errorJson := this.errorJson("Unable to unmarshal Size data")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	content := args[7]

	isFragile, err := strconv.ParseBool(args[8])
	if err != nil {
		errorJson := this.errorJson("9th argument isFragile must be a boolean representing string")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	lastState := args[9]

	// check if shipment already exists
	shipmentAsBytes, err := stub.GetState(trackingCode)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Failed to get shipment %s: %s", trackingCode, err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	} else if shipmentAsBytes != nil {
		errorJson := this.errorJson(fmt.Sprintf("Shipment %s already exists", trackingCode))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	// create new shipment and marshal it to JSON
	shipment := NewShipment(
		trackingCode,
		courier,
		sender,
		recipient,
		weightInGrams,
		shippingType,
		size,
		content,
		isFragile,
		lastState,
	)

	shipmentJSONAsBytes, err := json.Marshal(shipment)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to marshal shipment as JSON: %s", err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	// save shipment to state
	err = stub.PutState(trackingCode, shipmentJSONAsBytes)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to put shipment JSON bytes to state: %s", err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	//  Index the shipment to enable tracking code queries
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~createdAtDate~trackingCode.
	//  This will enable very efficient state range queries based on composite keys matching indexName~createdAtDate~*
	indexName := "createdAtDate-trackingCode"
	createdAtDate := shipment.CreatedAt.Format("2006-01-02 15:04:05")
	createdAtDateTrackingCodeIndexKey, err := stub.CreateCompositeKey(indexName, []string{createdAtDate, shipment.TrackingCode})
	if err != nil {
		errorJson := this.errorJson(
			fmt.Sprintf("Unable to create shipment createdAtDate-trackingCode composite key: %s", err.Error()),
		)
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the shipment.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	err = stub.PutState(createdAtDateTrackingCodeIndexKey, value)
	if err != nil {
		errorJson := this.errorJson(
			fmt.Sprintf("Unable to put shipment %s index createdAtDate-trackingCode to state: %s", shipment.TrackingCode, err.Error()),
		)
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	// send shipment-created event to the SDK
	stub.SetEvent("shipment-created", shipmentJSONAsBytes)

	logger.Debug("End createShipment")

	return shim.Success(nil)
}

func (this *ShipmentChaincode) requireSellerCreator(stub shim.ChaincodeStubInterface) error {
	// only seller organization can create shipments
	creatorBytes, err := this.getTransactionCreator(stub)
	if err != nil {
		return err
	}

	name, org := this.getCreator(creatorBytes)
	if org != "seller" {
		errorJson := this.errorJson(
			fmt.Sprintf("Only seller organization can create shipments, called by %s@s", name, org),
		)
		logger.Error(errorJson)

		return errors.New(errorJson)
	}

	return nil
}
