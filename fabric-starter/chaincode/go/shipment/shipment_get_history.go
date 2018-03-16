package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type ShipmentHistoryRecord struct {
	TxId      string
	Value     string
	Timestamp string
	IsDelete  bool
}

type ShipmentHistory []ShipmentHistoryRecord

func (this *ShipmentChaincode) getShipmentHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// validate arguments
	err = this.validateArgumentsNotEmpty(1, args)
	if err != nil {
		return shim.Error(err.Error())
	}

	trackingCode := args[0]

	logger.Debugf("Start getShipmentHistory for shipment %s", trackingCode)

	resultsIterator, err := stub.GetHistoryForKey(trackingCode)
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to get shipment %s history: %s", trackingCode, err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
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
			errorJson := this.errorJson(
				fmt.Sprintf("Unable to get shipment %s next result: %s", trackingCode, err.Error()),
			)
			logger.Error(errorJson)

			return shim.Error(errorJson)
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

		historyRecord := ShipmentHistoryRecord{
			TxId:      response.TxId,
			Value:     string(response.Value),
			Timestamp: time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String(),
			IsDelete:  response.IsDelete,
		}

		history = append(history, historyRecord)
	}

	//buffer.WriteString("]")

	if history == nil {
		errorJson := this.errorJson(fmt.Sprintf("Shipment %s does not exist", trackingCode))
		logger.Error(errorJson)

		return pb.Response{Status: 404, Message: errorJson}
	}

	historyBytes, err := json.Marshal(history)
	if err != nil {
		errorJson := this.errorJson(
			fmt.Sprintf("Unable to marshal shipment %s history: %s", trackingCode, err.Error()),
		)
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	//logger.Debugf("End getShipmentHistory for shipment %s: %s", trackingCode, buffer.String())
	logger.Debugf("End getShipmentHistory for shipment %s: %s", trackingCode, string(historyBytes))

	//return shim.Success(buffer.Bytes())
	return shim.Success(historyBytes)
}
