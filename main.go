package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"time"
	"encoding/json"
)

type ApartementRegister struct {
}

type Renter struct {
	name    string
	surname string
	movedIn time.Time
}

type Block struct {
	id       string
	street   string
	number   string
	renters  []Renter
	nOfRooms string
}

//cache of blocks id
var blocks map[string]bool

func createId(street string, number string) string {
	return fmt.Sprintf("%s%d", street, number)
}

//retrieve a block on the ledger
func getBlock(stub shim.ChaincodeStubInterface, key string) (*Block, error) {
	var block Block
	block_marshalled, err := stub.GetState(key)
	err = json.Unmarshal(block_marshalled, &block)
	return &block, err
}

//save a block on the ledger
func putBlock(stub shim.ChaincodeStubInterface, key string, block *Block) error {
	block_marshalled, _ := json.Marshal(*block)
	return stub.PutState(key, block_marshalled)
}



//Initialisation of the Chaincode
func (m *ApartementRegister) Init(stub shim.ChaincodeStubInterface) peer.Response {
	blocks = make(map[string]bool)
	return shim.Success([]byte("Successfully initialized Chaincode."))
}

//Entry Point of an invocation
func (m *ApartementRegister) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, para := stub.GetFunctionAndParameters()

	switch(function) {
	case "queryRenter":
		if len(para) < 3 {
			return shim.Error("not enough arguments for queryRenter. 3 required")
		} else {
			return queryRenter(stub, para[0], para[1], para[2])
		}
	case "registerRenter":
		if len(para) < 3 {
			return shim.Error("not enough arguments for registerRenter. 4 required")
		} else {
			return registerNewRenter(stub, para[0], para[1], para[2], para[3])
		}
	case "newBlock":
		return newBlock(stub, para[0], para[1], para[2])
	case "blocksCount":
		return blocksCount()
	case "rentersCount":
		if len(para) < 2 {
			return shim.Error("not enough arguments for rentersCount. 2 required")
		} else {
			return rentersCount(stub, para[0], para[1])
		}
	case "findEmptyBlock":
		return findEmptyBlock(stub)
	}
	return shim.Error(fmt.Sprintf("No function %s implemented", function))
}

func main() {
	if err := shim.Start(new(ApartementRegister)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
