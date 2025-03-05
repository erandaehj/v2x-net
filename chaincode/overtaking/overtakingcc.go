package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing the overtaking process
type SmartContract struct {
	contractapi.Contract
}

// Vehicle represents basic information about a vehicle
type Vehicle struct {
	ID       string  `json:"id"`
	Role     string  `json:"role"` // e.g., "OV", "LV", "LDV", "FDV"
	Speed    float64 `json:"speed"`
	Position float64 `json:"position"`
	Lane     int     `json:"lane"`    // Lane position (1 for the left lane, 2 for the right lane)
	Status   string  `json:"status"`  // e.g., "Pending", "Approved", "Completed"
}

// OvertakeProposal contains parameters for evaluating an overtaking maneuver
type OvertakeProposal struct {
	Dv     float64 `json:"visibilityDistance"`
	Do     float64 `json:"overtakingDistance"`
	Vr     float64 `json:"relativeSpeed"`
	Vo     float64 `json:"oncomingTrafficSpeed"`
	Tr     float64 `json:"reactionTime"`
	Sm     float64 `json:"safetyMargin"`
	IsSafe bool    `json:"isSafe"`
}

// InitLedger initializes the ledger with sample data for each vehicle
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	vehicles := []Vehicle{
		{ID: "OV", Role: "OV", Speed: 30, Position: 0, Lane: 2, Status: "Idle"},
		{ID: "LV", Role: "LV", Speed: 25, Position: 10, Lane: 2, Status: "Idle"},
		{ID: "LDV", Role: "LDV", Speed: 30, Position: 20, Lane: 1, Status: "Idle"},
		{ID: "FDV", Role: "FDV", Speed: 30, Position: -10, Lane: 1, Status: "Idle"},
	}

	for _, vehicle := range vehicles {
		vehicleAsBytes, _ := json.Marshal(vehicle)
		err := ctx.GetStub().PutState(vehicle.ID, vehicleAsBytes)
		if err != nil {
			return fmt.Errorf("failed to initialize vehicle %s: %v", vehicle.ID, err)
		}
	}

	return nil
}

// InitiateOvertakeProposal starts an overtaking proposal by the OV
func (s *SmartContract) InitiateOvertakeProposal(ctx contractapi.TransactionContextInterface, dv, do, vr, vo, tr, sm float64) (bool, error) {
	// Create and store the proposal
	proposal := OvertakeProposal{
		Dv: dv,
		Do: do,
		Vr: vr,
		Vo: vo,
		Tr: tr,
		Sm: sm,
	}
	
	// Evaluate if it’s safe to overtake based on safety calculation and lane position
	isSafe, err := s.EvaluateSafety(ctx, proposal)
	if err != nil {
		return false, err
	}

	proposal.IsSafe = isSafe
	proposalBytes, _ := json.Marshal(proposal)
	err = ctx.GetStub().PutState("OVERTAKE_PROPOSAL", proposalBytes)
	if err != nil {
		return false, fmt.Errorf("failed to store overtake proposal: %v", err)
	}

	// Change the status of OV to "Pending" if the overtaking is deemed safe
	if isSafe {
		ov, err := s.GetVehicle(ctx, "OV")
		if err != nil {
			return false, err
		}
		ov.Status = "Pending"
		ovBytes, _ := json.Marshal(ov)
		ctx.GetStub().PutState(ov.ID, ovBytes)
	}

	return isSafe, nil
}

// EvaluateSafety checks if the overtaking maneuver is safe based on parameters and vehicle positions
func (s *SmartContract) EvaluateSafety(ctx contractapi.TransactionContextInterface, proposal OvertakeProposal) (bool, error) {
	// Retrieve the vehicles involved
	ov, err := s.GetVehicle(ctx, "OV")
	if err != nil {
		return false, err
	}
	lv, err := s.GetVehicle(ctx, "LV")
	if err != nil {
		return false, err
	}
	ldv, err := s.GetVehicle(ctx, "LDV")
	if err != nil {
		return false, err
	}
	fdv, err := s.GetVehicle(ctx, "FDV")
	if err != nil {
		return false, err
	}

	// Check if all vehicles are in the same lane for the overtaking maneuver without lane change
	if ov.Lane != lv.Lane || fdv.Lane != ldv.Lane {
		return false, fmt.Errorf("vehicles are not in the same lane, overtaking without lane change not possible")
	}

	// Calculate if there’s sufficient space for overtaking based on positions
	safetyDistance := (proposal.Dv - proposal.Do) / proposal.Vr > (proposal.Tr + proposal.Do/proposal.Vo + proposal.Sm)
	isPositionSafe := (lv.Position - ov.Position) > proposal.Do && (ldv.Position - lv.Position) > proposal.Do
	
	return safetyDistance && isPositionSafe, nil
}

// GetVehicle retrieves a vehicle from the world state by its ID
func (s *SmartContract) GetVehicle(ctx contractapi.TransactionContextInterface, id string) (*Vehicle, error) {
	vehicleAsBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read vehicle %s: %v", id, err)
	}
	if vehicleAsBytes == nil {
		return nil, fmt.Errorf("vehicle %s does not exist", id)
	}

	vehicle := new(Vehicle)
	_ = json.Unmarshal(vehicleAsBytes, vehicle)

	return vehicle, nil
}

// EndorseOvertakeRequest is called by LV to confirm that it is ready for the overtaking maneuver
func (s *SmartContract) EndorseOvertakeRequest(ctx contractapi.TransactionContextInterface) error {
	// Change the status of LV to "Approved" for overtaking
	lv, err := s.GetVehicle(ctx, "LV")
	if err != nil {
		return err
	}
	lv.Status = "Approved"

	lvBytes, _ := json.Marshal(lv)
	return ctx.GetStub().PutState(lv.ID, lvBytes)
}

// CommitOvertakingManeuver commits the overtaking maneuver if all endorsements are met
func (s *SmartContract) CommitOvertakingManeuver(ctx contractapi.TransactionContextInterface) error {
	// Ensure that both OV and LV have approved the overtaking process
	ov, err := s.GetVehicle(ctx, "OV")
	if err != nil {
		return err
	}
	lv, err := s.GetVehicle(ctx, "LV")
	if err != nil {
		return err
	}
	
	// Check if both OV and LV statuses are approved
	if ov.Status == "Approved" && lv.Status == "Approved" {
		// Commit the maneuver by updating their positions
		ov.Position += lv.Position
		ov.Status = "Completed"
		lv.Status = "Idle"

		ovBytes, _ := json.Marshal(ov)
		lvBytes, _ := json.Marshal(lv)
		err := ctx.GetStub().PutState(ov.ID, ovBytes)
		if err != nil {
			return fmt.Errorf("failed to update OV state: %v", err)
		}
		err = ctx.GetStub().PutState(lv.ID, lvBytes)
		if err != nil {
			return fmt.Errorf("failed to update LV state: %v", err)
		}

		return nil
	}

	return fmt.Errorf("conditions for overtaking maneuver not met")
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
