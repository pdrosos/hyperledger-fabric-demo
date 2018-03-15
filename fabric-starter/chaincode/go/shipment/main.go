package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/pdrosos/hyperledger-fabric-demo/fabric-starter/chaincode/go/shipment/chaincode"
)

func main() {
	logger := shim.NewLogger("ShipmentChaincode")

	shipmentChaincode := chaincode.NewShipmentChaincode(*logger)

	err := shim.Start(shipmentChaincode)
	if err != nil {
		logger.Error(err.Error())
	}
}

