package chaincode

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this *ShipmentChaincode) getShipmentById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		errorMessage := "Incorrect number of arguments. Expecting 1"
		this.logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	id := args[0]

	this.logger.Debugf("Start getShipmentById for shipment ID %s", id)

	shipmentAsBytes, err := stub.GetState(id)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for shipment " + id + "\"}"

		return shim.Error(jsonResp)
	} else if shipmentAsBytes == nil {
		jsonResp := "{\"Error\":\"Shipment does not exist: " + id + "\"}"

		return shim.Error(jsonResp)
	}

	this.logger.Debugf("End getShipmentById for shipment ID %s: %s", id, string(shipmentAsBytes))

	return shim.Success(shipmentAsBytes)
}
