package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type ShipmentHistoryItem struct {
	TxId      string
	Value     string
	Timestamp string
	IsDelete  bool
}

type ShipmentHistory []ShipmentHistoryItem

func (this ShipmentChaincode) getShipmentHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		errorMessage := "Incorrect number of arguments. Expecting 1"
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	id := args[0]

	logger.Debugf("Start getShipmentHistory for shipment ID %s", id)

	resultsIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to get shipment ID %s history: %s", id, err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	defer resultsIterator.Close()

	var history ShipmentHistory

	// buffer is a JSON array containing historic values for the shipment
	//var buffer bytes.Buffer
	//buffer.WriteString("[")

	//bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			errorMessage := fmt.Sprintf("Unable to get shipment ID %s next result: %s", id, err.Error())
			logger.Error(errorMessage)

			return shim.Error(errorMessage)
		}

		//// Add a comma before array members, suppress it for the first array member
		//if bArrayMemberAlreadyWritten == true {
		//	buffer.WriteString(",")
		//}
		//
		//buffer.WriteString("{\"TxId\":")
		//buffer.WriteString("\"")
		//buffer.WriteString(response.TxId)
		//buffer.WriteString("\"")
		//
		//buffer.WriteString(", \"Value\":")
		//// if it was a delete operation on given key, then we need to set the
		////corresponding value null. Else, we will write the response.Value
		////as-is (as the Value itself a JSON shipment)
		//if response.IsDelete {
		//	buffer.WriteString("null")
		//} else {
		//	buffer.WriteString(string(response.Value))
		//}
		//
		//buffer.WriteString(", \"Timestamp\":")
		//buffer.WriteString("\"")
		//buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		//buffer.WriteString("\"")
		//
		//buffer.WriteString(", \"IsDelete\":")
		//buffer.WriteString("\"")
		//buffer.WriteString(strconv.FormatBool(response.IsDelete))
		//buffer.WriteString("\"")
		//
		//buffer.WriteString("}")
		//bArrayMemberAlreadyWritten = true

		historyItem := ShipmentHistoryItem{
			TxId:      response.TxId,
			Value:     string(response.Value),
			Timestamp: time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String(),
			IsDelete:  response.IsDelete,
		}

		history = append(history, historyItem)
	}

	//buffer.WriteString("]")

	historyBytes, err := json.Marshal(history)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to marshal shipment ID %s history: %s", id, err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	//logger.Debugf("End getShipmentHistory for shipment ID %s: %s", id, buffer.String())
	logger.Debugf("End getShipmentHistory for shipment ID %s: %s", id, string(historyBytes))

	//return shim.Success(buffer.Bytes())
	return shim.Success(historyBytes)
}
