package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this ShipmentChaincode) changeShipmentStateAndLocation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("Start changeShipmentStateAndLocation")



	logger.Debug("Start changeShipmentStateAndLocation")

	return shim.Success(nil);
}
