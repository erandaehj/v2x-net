package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AuctionContract struct {
	contractapi.Contract
}

type Bid struct {
	ClientID string `json:"clientID"`
	BidHash  string `json:"bidHash"`
	BidValue int    `json:"bidValue,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
	Revealed bool   `json:"revealed"`
}

type Auction struct {
	Asset     string          `json:"asset"`
	Bids      map[string]*Bid `json:"bids"`
	Asks      map[string]int  `json:"asks"`
	StartTime int64           `json:"startTime"`
	BidEnd    int64           `json:"bidEnd"`
	RevealEnd int64           `json:"revealEnd"`
	Awarded   bool            `json:"awarded"`
	Winner    string          `json:"winner,omitempty"`
}

func (s *AuctionContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *AuctionContract) InitAuction(ctx contractapi.TransactionContextInterface, asset string, bidDuration, revealDuration int) error {
	startTime := time.Now().Unix()
	auction := Auction{
		Asset:     asset,
		Bids:      make(map[string]*Bid),
		Asks:      make(map[string]int),
		StartTime: startTime,
		BidEnd:    startTime + int64(bidDuration),
		RevealEnd: startTime + int64(bidDuration) + int64(revealDuration),
		Awarded:   false,
	}
	auctionBytes, _ := json.Marshal(auction)
	return ctx.GetStub().PutState(asset, auctionBytes)
}

func (s *AuctionContract) PlaceBid(ctx contractapi.TransactionContextInterface, asset, clientID, bidHash string) error {
	auctionBytes, _ := ctx.GetStub().GetState(asset)
	if auctionBytes == nil {
		return fmt.Errorf("Auction not found")
	}

	var auction Auction
	json.Unmarshal(auctionBytes, &auction)
	if time.Now().Unix() > auction.BidEnd {
		return fmt.Errorf("Bidding phase over")
	}

	auction.Bids[clientID] = &Bid{ClientID: clientID, BidHash: bidHash}
	auctionBytes, _ = json.Marshal(auction)
	return ctx.GetStub().PutState(asset, auctionBytes)
}

func (s *AuctionContract) PlaceAsk(ctx contractapi.TransactionContextInterface, asset, clientID string, ask int) error {
	auctionBytes, _ := ctx.GetStub().GetState(asset)
	if auctionBytes == nil {
		return fmt.Errorf("Auction not found")
	}

	var auction Auction
	json.Unmarshal(auctionBytes, &auction)

	auction.Asks[clientID] = ask
	auctionBytes, _ = json.Marshal(auction)
	return ctx.GetStub().PutState(asset, auctionBytes)
}

func (s *AuctionContract) RevealBid(ctx contractapi.TransactionContextInterface, asset, clientID string, bidValue int, nonce string) error {
	auctionBytes, _ := ctx.GetStub().GetState(asset)
	if auctionBytes == nil {
		return fmt.Errorf("Auction not found")
	}

	var auction Auction
	json.Unmarshal(auctionBytes, &auction)
	if time.Now().Unix() <= auction.BidEnd || time.Now().Unix() > auction.RevealEnd {
		return fmt.Errorf("Not in the reveal phase")
	}

	bid, exists := auction.Bids[clientID]
	if !exists {
		return fmt.Errorf("No bid found")
	}

	hash := sha256.New()
	hash.Write([]byte(strconv.Itoa(bidValue) + nonce))
	hashValue := hash.Sum(nil)
	expectedHash := hex.EncodeToString(hashValue)
	if bid.BidHash != expectedHash {
		return fmt.Errorf("Hash mismatch")
	}

	bid.BidValue = bidValue
	bid.Nonce = nonce
	bid.Revealed = true
	auctionBytes, _ = json.Marshal(auction)
	return ctx.GetStub().PutState(asset, auctionBytes)
}

func (s *AuctionContract) AwardSlot(ctx contractapi.TransactionContextInterface, asset string) error {
	auctionBytes, _ := ctx.GetStub().GetState(asset)
	if auctionBytes == nil {
		return fmt.Errorf("Auction not found")
	}

	var auction Auction
	json.Unmarshal(auctionBytes, &auction)
	if time.Now().Unix() <= auction.RevealEnd {
		return fmt.Errorf("Reveal phase not over")
	}
	if auction.Awarded {
		return fmt.Errorf("Auction already awarded")
	}

	var highestBidder string
	var highestBid int
	for clientID, bid := range auction.Bids {
		if bid.Revealed && bid.BidValue > highestBid {
			highestBid = bid.BidValue
			highestBidder = clientID
		}
	}

	auction.Winner = highestBidder
	auction.Awarded = true
	auctionBytes, _ = json.Marshal(auction)
	return ctx.GetStub().PutState(asset, auctionBytes)
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(AuctionContract))
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
