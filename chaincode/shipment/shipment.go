package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type ShipmentChaincode struct {
}

type Shipment struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

/*
 * The Init method is called when the Chaincode is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (this *ShipmentChaincode) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the chaincode
 * The calling application program has also specified the particular chaincode function to be called, with arguments
 */
func (this *ShipmentChaincode) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryShipment" {
		return this.queryShipment(APIstub, args)
	} else if function == "createShipment" {
		return this.createShipment(APIstub, args)
	} else if function == "queryAllShipments" {
		return this.queryAllShipments(APIstub)
	} else if function == "changeShipmentLocation" {
		return this.changeShipmentLocation(APIstub, args)
	} else if function == "initLedger" {
		return this.initLedger(APIstub)
	}

	return shim.Error("Invalid Chaincode function name.")
}

func (this *ShipmentChaincode) queryShipment(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	shipmentAsBytes, _ := APIstub.GetState(args[0])

	return shim.Success(shipmentAsBytes)
}

func (this *ShipmentChaincode) createShipment(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	return shim.Success(nil)
}

func (this *ShipmentChaincode) queryAllShipments(APIstub shim.ChaincodeStubInterface) peer.Response {
	startKey := "SHIPMENT0"
	endKey := "SHIPMENT999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	fmt.Printf("- queryAllShipments\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (this *ShipmentChaincode) changeShipmentLocation(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	return shim.Success(nil)
}

func (this *ShipmentChaincode) initLedger(APIstub shim.ChaincodeStubInterface) peer.Response {
	shipments := []Shipment{
		Shipment{Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
		Shipment{Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad"},
		Shipment{Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
		Shipment{Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max"},
		Shipment{Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana"},
		Shipment{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
		Shipment{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
		Shipment{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
		Shipment{Make: "Tata", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
		Shipment{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	}

	i := 0
	for i < len(shipments) {
		fmt.Println("i is ", i)
		shipmentAsBytes, _ := json.Marshal(shipments[i])
		APIstub.PutState("Shipment" + strconv.Itoa(i), shipmentAsBytes)
		fmt.Println("Added", shipments[i])
		i = i + 1
	}

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Chaincode
	err := shim.Start(new(ShipmentChaincode))
	if err != nil {
		fmt.Printf("Error creating new Chaincode: %s", err)
	}
}
