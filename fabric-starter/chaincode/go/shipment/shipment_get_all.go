package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this ShipmentChaincode) getAllShipments(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("Start getAllShipments")



	logger.Debug("Start getAllShipments")

	return shim.Success(nil);
}
