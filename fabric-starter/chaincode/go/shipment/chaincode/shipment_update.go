package chaincode

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this *ShipmentChaincode) updateShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	this.logger.Debug("Start updateShipment")

	// todo
	this.logger.Debug("TODO")

	this.logger.Debug("End updateShipment")

	return shim.Success(nil)
}
