package main

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

var logger = shim.NewLogger("ShipmentChaincode")

type ShipmentChaincode struct {
}

func (this *ShipmentChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Init")

	creatorBytes, err := stub.GetCreator()
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to get chaincode creator: %s", err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	name, org := this.getCreator(creatorBytes)
	logger.Debug(fmt.Sprintf("Chaincode creator: %s@%s", name, org))

	return shim.Success(nil)
}

func (this *ShipmentChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Invoke")
	creatorBytes, err := stub.GetCreator()
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to get transaction creator: %s", err.Error()))
		logger.Error(errorJson)

		return shim.Error(errorJson)
	}

	name, org := this.getCreator(creatorBytes)
	logger.Debug("Transaction creator " + name + "@" + org)

	function, args := stub.GetFunctionAndParameters()
	if function == "createShipment" {
		return this.createShipment(stub, args)
	} else if function == "updateShipment" { // todo
		return this.updateShipment(stub, args)
	} else if function == "changeShipmentStateAndLocation" {
		return this.changeShipmentStateAndLocation(stub, args)
	} else if function == "getShipmentById" {
		return this.getShipmentById(stub, args)
	} else if function == "getAllShipments" {
		return this.getAllShipments(stub, args)
	} else if function == "getShipmentHistory" {
		return this.getShipmentHistory(stub, args)
	}

	return pb.Response{Status: 403, Message: "Invalid invoke function name."}
}

func (this *ShipmentChaincode) getTransactionCreator(stub shim.ChaincodeStubInterface) ([]byte, error) {
	creatorBytes, err := stub.GetCreator()
	if err != nil {
		errorJson := this.errorJson(fmt.Sprintf("Unable to get transaction creator: %s", err.Error()))
		logger.Error(errorJson)

		return nil, errors.New(errorJson)
	}

	return creatorBytes, nil
}

func (this *ShipmentChaincode) getCreator(certificate []byte) (string, string) {
	data := certificate[strings.Index(string(certificate), "-----") : strings.LastIndex(string(certificate), "-----")+5]
	block, _ := pem.Decode([]byte(data))
	cert, _ := x509.ParseCertificate(block.Bytes)
	organization := cert.Issuer.Organization[0]
	commonName := cert.Subject.CommonName

	logger.Debug("commonName: " + commonName + ", organization: " + organization)

	organizationShort := strings.Split(organization, ".")[0]

	return commonName, organizationShort
}

func (this *ShipmentChaincode) validateArgumentsNotEmpty(argsCount int, args []string) error {
	if len(args) != argsCount {
		errorJson := this.errorJson(fmt.Sprintf("Incorrect number of arguments. Expecting %s", argsCount))
		logger.Error(errorJson)

		return errors.New(errorJson)
	}

	for index, element := range args {
		if len(element) <= 0 {
			errorJson := this.errorJson(fmt.Sprintf("Argument %s must be a non-empty string", index+1))
			logger.Error(errorJson)

			return errors.New(errorJson)
		}
	}

	return nil
}

func (this *ShipmentChaincode) errorJson(errorMessage string) string {
	stateError := Error{
		Error: errorMessage,
	}

	errorJson, _ := json.Marshal(stateError)

	return string(errorJson)
}

func main() {
	shipmentChaincode := new(ShipmentChaincode)
	err := shim.Start(shipmentChaincode)
	if err != nil {
		logger.Error(err.Error())
	}
}
