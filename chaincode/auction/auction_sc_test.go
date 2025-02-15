package main

import (
	"testing"

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
func TestCreateBid(t *testing.T) {
	contract := new(AuctionContract)
	mockStub := &MockChaincodeStub{data: make(map[string][]byte)}
	mockContext := &MockTransactionContext{stub: mockStub}

	// Create a bid with ID "BID001" and amount 50
	err := contract.CreateBid(mockContext, "BID001", 50)
	require.NoError(t, err)

	// Query the bid back to verify it exists
	bid, err := contract.QueryBid(mockContext, "BID001")
	require.NoError(t, err)
	require.Equal(t, "BID001", bid.ID)
	require.Equal(t, 50, bid.Amount)
}

// TestCreateAsk tests the CreateAsk function
func TestCreateAsk(t *testing.T) {
	contract := new(AuctionContract)
	mockStub := &MockChaincodeStub{data: make(map[string][]byte)}
	mockContext := &MockTransactionContext{stub: mockStub}

	// Create an ask with ID "ASK001" and amount 60
	err := contract.CreateAsk(mockContext, "ASK001", 60)
	require.NoError(t, err)

	// Query the ask back to verify it exists
	ask, err := contract.QueryAsk(mockContext, "ASK001")
	require.NoError(t, err)
	require.Equal(t, "ASK001", ask.ID)
	require.Equal(t, 60, ask.Amount)
}

// TestExecuteTransaction tests matching a bid and an ask
func TestExecuteTransaction(t *testing.T) {
	contract := new(AuctionContract)
	mockStub := &MockChaincodeStub{data: make(map[string][]byte)}
	mockContext := &MockTransactionContext{stub: mockStub}

	// Create a bid and ask
	contract.CreateBid(mockContext, "BID001", 70)
	contract.CreateAsk(mockContext, "ASK001", 60)

	// Test executing the transaction
	result, err := contract.ExecuteTransaction(mockContext, "BID001", "ASK001")
	require.NoError(t, err)
	require.Contains(t, result, "Transaction executed")
}
