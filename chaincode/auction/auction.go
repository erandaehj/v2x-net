package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AuctionContract struct {
	contractapi.Contract
}

// Bid represents a bid in the auction
type Bid struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
}

// Ask represents an ask in the auction
type Ask struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
}

// InitLedger initializes the auction with some predefined asks and bids
func (s *AuctionContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// Initialize some dummy data
	bids := []Bid{
		{ID: "BID001", Amount: 50},
		{ID: "BID002", Amount: 80},
		{ID: "BID003", Amount: 100},
	}
	asks := []Ask{
		{ID: "ASK001", Amount: 60},
		{ID: "ASK002", Amount: 90},
		{ID: "ASK003", Amount: 110},
	}

	// Store bids
	for _, bid := range bids {
		bidJSON, err := json.Marshal(bid)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState("BID_"+bid.ID, bidJSON)
		if err != nil {
			return fmt.Errorf("failed to put bid %s to world state: %v", bid.ID, err)
		}
	}

	// Store asks
	for _, ask := range asks {
		askJSON, err := json.Marshal(ask)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState("ASK_"+ask.ID, askJSON)
		if err != nil {
			return fmt.Errorf("failed to put ask %s to world state: %v", ask.ID, err)
		}
	}

	return nil
}

// CreateBid creates a new bid
func (s *AuctionContract) CreateBid(ctx contractapi.TransactionContextInterface, bidID string, amount int) error {
	bid := Bid{
		ID:     bidID,
		Amount: amount,
	}

	bidJSON, err := json.Marshal(bid)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("BID_"+bidID, bidJSON)
}

// CreateAsk creates a new ask
func (s *AuctionContract) CreateAsk(ctx contractapi.TransactionContextInterface, askID string, amount int) error {
	ask := Ask{
		ID:     askID,
		Amount: amount,
	}

	askJSON, err := json.Marshal(ask)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("ASK_"+askID, askJSON)
}

// ExecuteTransaction tries to match a bid with an ask
func (s *AuctionContract) ExecuteTransaction(ctx contractapi.TransactionContextInterface, bidID string, askID string) (string, error) {
	// Retrieve the bid and ask from the state
	bidJSON, err := ctx.GetStub().GetState("BID_" + bidID)
	if err != nil || bidJSON == nil {
		return "", fmt.Errorf("bid %s not found", bidID)
	}

	askJSON, err := ctx.GetStub().GetState("ASK_" + askID)
	if err != nil || askJSON == nil {
		return "", fmt.Errorf("ask %s not found", askID)
	}

	var bid Bid
	err = json.Unmarshal(bidJSON, &bid)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal bid %s: %v", bidID, err)
	}

	var ask Ask
	err = json.Unmarshal(askJSON, &ask)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal ask %s: %v", askID, err)
	}

	// Perform the transaction if a match is found
	if bid.Amount >= ask.Amount {
		// A match happens, execute the transaction (e.g., transferring assets)
		err = ctx.GetStub().DelState("BID_" + bidID)
		if err != nil {
			return "", fmt.Errorf("failed to delete bid %s: %v", bidID, err)
		}
		err = ctx.GetStub().DelState("ASK_" + askID)
		if err != nil {
			return "", fmt.Errorf("failed to delete ask %s: %v", askID, err)
		}
		return fmt.Sprintf("Transaction executed between Bid %s and Ask %s for amount %d", bidID, askID, bid.Amount), nil
	}

	return "", fmt.Errorf("no match found for Bid %s and Ask %s", bidID, askID)
}

// QueryBid retrieves a bid by ID
func (s *AuctionContract) QueryBid(ctx contractapi.TransactionContextInterface, bidID string) (*Bid, error) {
	bidJSON, err := ctx.GetStub().GetState("BID_" + bidID)
	if err != nil || bidJSON == nil {
		return nil, fmt.Errorf("bid %s not found", bidID)
	}

	var bid Bid
	err = json.Unmarshal(bidJSON, &bid)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bid %s: %v", bidID, err)
	}

	return &bid, nil
}

// QueryAsk retrieves an ask by ID
func (s *AuctionContract) QueryAsk(ctx contractapi.TransactionContextInterface, askID string) (*Ask, error) {
	askJSON, err := ctx.GetStub().GetState("ASK_" + askID)
	if err != nil || askJSON == nil {
		return nil, fmt.Errorf("ask %s not found", askID)
	}

	var ask Ask
	err = json.Unmarshal(askJSON, &ask)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ask %s: %v", askID, err)
	}

	return &ask, nil
}
