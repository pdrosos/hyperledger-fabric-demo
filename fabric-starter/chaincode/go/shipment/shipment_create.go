package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this ShipmentChaincode) createShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("Start createShipment")

	var err error

	// only seller organization can create shipments
	creatorBytes, err := stub.GetCreator()
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to get transaction creator: %s", err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	name, org := getCreator(creatorBytes)
	if org != "seller" {
		errorMessage := fmt.Sprintf("Only seller organization can create shipments, called by %s@s", name, org)
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	// validate arguments
	if len(args) != 11 {
		errorMessage := "Incorrect number of arguments. Expecting 11"
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	for index, element := range args {
		if len(element) <= 0 {
			errorMessage := fmt.Sprintf("Argument %s must be a non-empty string", index + 1)
			logger.Error(errorMessage)

			return shim.Error(errorMessage)
		}
	}

	id := args[0]
	trackingCode := args[1]

	courier := Courier{}
	err = json.Unmarshal([]byte(args[2]), &courier)
	if err != nil {
		errorMessage := "Unable to unmarshal Courier data"
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	sender := Sender{}
	err = json.Unmarshal([]byte(args[3]), &sender)
	if err != nil {
		errorMessage := "Unable to unmarshal Sender data"
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	recipient := Recipient{}
	err = json.Unmarshal([]byte(args[4]), &recipient)
	if err != nil {
		errorMessage := "Unable to unmarshal Recipient data"
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	weightInGrams, err := strconv.Atoi(args[5])
	if err != nil {
		errorMessage := "6th argument weightInGrams must be a numeric string"
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	shippingType := args[6]

	size := Size{}
	err = json.Unmarshal([]byte(args[7]), &size)
	if err != nil {
		errorMessage := "Unable to unmarshal Size data"
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	content := args[8]

	isFragile, err := strconv.ParseBool(args[9])
	if err != nil {
		errorMessage := "10th argument isFragile must be a boolean representing string"
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	lastState := args[10]

	// check if shipment already exists
	shipmentAsBytes, err := stub.GetState(id)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get shipment ID %s: %s", id, err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	} else if shipmentAsBytes != nil {
		errorMessage := fmt.Sprintf("Shipment ID %s already exists" + id)

		return shim.Error(errorMessage)
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
		errorMessage := fmt.Sprintf("Unable to marshal shipment as JSON: %s", err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	// save shipment to state
	err = stub.PutState(id, shipmentJSONAsBytes)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to put shipment JSON bytes to state: %s", err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	//  Index the shipment to enable tracking code queries
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~trackingCode~id.
	//  This will enable very efficient state range queries based on composite keys matching indexName~trackingCode~*
	indexName := "trackingCode-id"
	trackingCodeIdIndexKey, err := stub.CreateCompositeKey(indexName, []string{shipment.TrackingCode, shipment.Id})
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to create shipment trackingCode-id composite key: %s", err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the shipment.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	err = stub.PutState(trackingCodeIdIndexKey, value)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to put shipment ID %s index trackingCode-id to state: %s", shipment.Id, err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	logger.Debug("End createShipment")

	return shim.Success(nil);
}
