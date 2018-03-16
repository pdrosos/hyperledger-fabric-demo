package main

import (
	//"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

func (this *ShipmentChaincode) getAllShipments(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("Start getAllShipments")

	//if len(args) < 2 {
	//	return shim.Error("Incorrect number of arguments. Expecting 2")
	//}
	//
	//startDate := args[0]
	//endDate := args[1]

	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to get shipments: %s", err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	defer resultsIterator.Close()

	var shipments Shipments = make([]ShipmentRecord, 0)

	// buffer is a JSON array containing QueryResults
	//var buffer bytes.Buffer
	//buffer.WriteString("[")

	//bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			errorJson := this.errorJson(
				fmt.Sprintf("Unable to get shipments next result: %s", err.Error()),
			)
			logger.Error(errorJson)

			return shim.Error(errorJson)
		}

		//// Add a comma before array members, suppress it for the first array member
		//if bArrayMemberAlreadyWritten == true {
		//	buffer.WriteString(",")
		//}
		//buffer.WriteString("{\"Key\":")
		//buffer.WriteString("\"")
		//buffer.WriteString(response.Key)
		//buffer.WriteString("\"")
		//
		//buffer.WriteString(", \"Record\":")
		//// Record is a JSON object, so we write as-is
		//buffer.WriteString(string(response.Value))
		//buffer.WriteString("}")
		//bArrayMemberAlreadyWritten = true

		shipmentRecord := ShipmentRecord{
			Key:   response.Key,
			Value: string(response.Value),
		}

		shipments = append(shipments, shipmentRecord)
	}
	//buffer.WriteString("]")

	shipmentsBytes, err := json.Marshal(shipments)
	if err != nil {
		errorJson := this.errorJson(
			fmt.Sprintf("Unable to marshal shipments: %s", err.Error()),
		)
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	//logger.Debugf("End getShipmentHistory for shipment %s: %s", id, buffer.String())
	logger.Debugf("End getAllShipments: %s", string(shipmentsBytes))

	//return shim.Success(buffer.Bytes())
	return shim.Success(shipmentsBytes)
}
