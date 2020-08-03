package main

import (
	// "bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
)

type SmartContract struct {
}

// User :  Define the User structure, with 5 properties.  Structure tags are used by encoding/json library
type User struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	DOB     string `json:"dob"`
	Gender  string `json:"gender"`
	Country string `json:"country"`
}

// func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
// 	return shim.Success(nil)
// }

var logger = flogging.MustGetLogger("fabUser_cc")

// Init : Method for INIT smart contract
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
	case "queryUser":
		return s.queryUser(APIstub, args)
	case "initLedger":
		return s.initLedger(APIstub)
	case "createUser":
		return s.createUser(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}

	// return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	UserAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(UserAsBytes)
}

func (s *SmartContract) createUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var User = User{Name: args[1], Email: args[2], Gender: args[3], Country: args[4]}

	UserAsBytes, _ := json.Marshal(User)
	APIstub.PutState(args[0], UserAsBytes)

	indexName := "owner~key"
	colorNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{User.Country, args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	APIstub.PutState(colorNameIndexKey, value)

	return shim.Success(UserAsBytes)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
