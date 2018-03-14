package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("ShipmentChaincode")

type ShipmentChaincode struct {
}

func (this *ShipmentChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Init")

	creatorBytes, err := stub.GetCreator()
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to get chaincode creator: %s", err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	name, org := getCreator(creatorBytes)
	logger.Debug(fmt.Sprintf("Chaincode creator: %s@%s", name, org))

	return shim.Success(nil)
}

func (this *ShipmentChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Invoke")
	creatorBytes, err := stub.GetCreator()
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to get transaction creator: %s", err.Error())
		logger.Error(errorMessage)

		return shim.Error(errorMessage)
	}

	name, org := getCreator(creatorBytes)
	logger.Debug("Transaction creator " + name + "@" + org)

	function, args := stub.GetFunctionAndParameters()
	if function == "createShipment" {
		return this.createShipment(stub, args)
	} else if function == "updateShipment" {
		return this.updateShipment(stub, args)
	} else if function == "changeShipmentStateAndLocation" {
		return this.getShipmentHistory(stub, args)
	}  else if function == "getShipmentById" {
		return this.getShipmentById(stub, args)
	} else if function == "getAllShipments" {
		return this.getAllShipments(stub, args)
	} else if function == "getShipmentHistory" {
		return this.getShipmentHistory(stub, args)
	}

	return pb.Response{Status:403, Message:"Invalid invoke function name."}
}

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

func (this ShipmentChaincode) updateShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// todo

	return shim.Success(nil);
}

func (this ShipmentChaincode) changeShipmentStateAndLocation(stub shim.ChaincodeStubInterface, args []string) pb.Response {


	return shim.Success(nil);
}

func (this ShipmentChaincode) getShipmentById(stub shim.ChaincodeStubInterface, args []string) pb.Response {


	return shim.Success(nil);
}

func (this ShipmentChaincode) getAllShipments(stub shim.ChaincodeStubInterface, args []string) pb.Response {


	return shim.Success(nil);
}

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
			TxId: response.TxId,
			Value: string(response.Value),
			Timestamp: time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String(),
			IsDelete: response.IsDelete,
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

func getCreator(certificate []byte) (string, string) {
	data := certificate[strings.Index(string(certificate), "-----"): strings.LastIndex(string(certificate), "-----")+5]
	block, _ := pem.Decode([]byte(data))
	cert, _ := x509.ParseCertificate(block.Bytes)
	organization := cert.Issuer.Organization[0]
	commonName := cert.Subject.CommonName
	logger.Debug("commonName: " + commonName + ", organization: " + organization)

	organizationShort := strings.Split(organization, ".")[0]

	return commonName, organizationShort
}

func main() {
	err := shim.Start(new(ShipmentChaincode))
	if err != nil {
		logger.Error(err.Error())
	}
}
