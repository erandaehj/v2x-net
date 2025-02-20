package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"testing"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/require"
)

// Mock ChaincodeStubInterface
type MockChaincodeStub struct {
	shim.ChaincodeStubInterface
	data map[string][]byte
}

func (m *MockChaincodeStub) PutState(key string, value []byte) error {
	m.data[key] = value
	return nil
}

func (m *MockChaincodeStub) GetState(key string) ([]byte, error) {
	return m.data[key], nil
}

func (m *MockChaincodeStub) DelState(key string) error {
	delete(m.data, key)
	return nil
}

// Other methods from ChaincodeStubInterface would be here (e.g., Invoke, SetEvent, etc.)

// Mock TransactionContextInterface
type MockTransactionContext struct {
	contractapi.TransactionContextInterface
	stub *MockChaincodeStub
}

func (m *MockTransactionContext) GetStub() shim.ChaincodeStubInterface {
	return m.stub
}

// TestCreateBid tests the CreateBid function
func TestAuctionLifecycle(t *testing.T) {
	contract := new(AuctionContract)
	mockStub := &MockChaincodeStub{data: make(map[string][]byte)}
	mockContext := &MockTransactionContext{stub: mockStub}

	// Create a bid with ID "SLOT001"
	bidDuration := 5
	revealDuration := 5
	assetName := "SLOT001"
	nonce := "nonce"
	bidVal := 1000
	clientId := "client1"
	hash := sha256.New()
	hash.Write([]byte(strconv.Itoa(bidVal) + nonce))
	hashValue := hash.Sum(nil)
	hashBid := hex.EncodeToString(hashValue)

	err := contract.InitAuction(mockContext, assetName, bidDuration, revealDuration)
	require.NoError(t, err)

	err = contract.PlaceBid(mockContext, assetName, clientId, hashBid)
	require.NoError(t, err)

	err = contract.PlaceAsk(mockContext, assetName, clientId, 500)
	require.NoError(t, err)

	time.Sleep(6 * time.Second)

	err = contract.RevealBid(mockContext, assetName, clientId, bidVal, nonce)
	require.NoError(t, err)

	time.Sleep(6 * time.Second)

	err = contract.AwardSlot(mockContext, assetName)
	require.NoError(t, err)

}
