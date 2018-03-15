package chaincode

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (this *ShipmentChaincode) getAllShipments(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	this.logger.Debug("Start getAllShipments")

	this.logger.Debug("End getAllShipments")

	return shim.Success(nil)
}
