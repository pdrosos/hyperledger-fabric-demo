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
	err = this.validateArgumentsNotEmpty(11, args)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]
	trackingCode := args[1]

	courier := Courier{}
	err = json.Unmarshal([]byte(args[2]), &courier)
	if err != nil {
		errorJson := this.errorJson("Unable to unmarshal Courier data")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	sender := Sender{}
	err = json.Unmarshal([]byte(args[3]), &sender)
	if err != nil {
		errorJson := this.errorJson("Unable to unmarshal Sender data")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	recipient := Recipient{}
	err = json.Unmarshal([]byte(args[4]), &recipient)
	if err != nil {
		errorJson := this.errorJson("Unable to unmarshal Recipient data")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	weightInGrams, err := strconv.Atoi(args[5])
	if err != nil {
		errorJson := this.errorJson("6th argument weightInGrams must be a numeric string")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	shippingType := args[6]

	size := Size{}
	err = json.Unmarshal([]byte(args[7]), &size)
	if err != nil {
		errorJson := this.errorJson("Unable to unmarshal Size data")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	content := args[8]

	isFragile, err := strconv.ParseBool(args[9])
	if err != nil {
		errorJson := this.errorJson("10th argument isFragile must be a boolean representing string")
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	lastState := args[10]

	// check if shipment already exists
	shipmentAsBytes, err := stub.GetState(id)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Failed to get shipment ID %s: %s", id, err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	} else if shipmentAsBytes != nil {
		errorJson := this.errorJson(fmt.Sprintf("Shipment ID %s already exists" + id))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	// create new shipment and marshal it to JSON
	shipment := NewShipment(
		id,
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
	err = stub.PutState(id, shipmentJSONAsBytes)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to put shipment JSON bytes to state: %s", err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	//  Index the shipment to enable tracking code queries
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~trackingCode~id.
	//  This will enable very efficient state range queries based on composite keys matching indexName~trackingCode~*
	indexName := "trackingCode-id"
	trackingCodeIdIndexKey, err := stub.CreateCompositeKey(indexName, []string{shipment.TrackingCode, shipment.Id})
	if err != nil {
		errorJson := this.errorJson(
			fmt.Sprintf("Unable to create shipment trackingCode-id composite key: %s", err.Error()),
		)
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the shipment.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	err = stub.PutState(trackingCodeIdIndexKey, value)
	if err != nil {
		errorJson := this.errorJson(
			fmt.Sprintf("Unable to put shipment ID %s index trackingCode-id to state: %s", shipment.Id, err.Error()),
		)
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

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
