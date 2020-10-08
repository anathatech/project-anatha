package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all treasury state that must be provided at genesis
type GenesisState struct {
	Treasury 	Treasury 			`json:"treasury" yaml:"treasury"`
	Params 		Params 				`json:"params" yaml:"params"`
	Operators 	[]sdk.AccAddress 	`json:"distribution_operators" yaml:"distribution_operators"`

	DisbursementQueue []Disbursement `json:"disbursement_queue" yaml:"disbursement_queue"`

	DisbursementReferences []ReferenceAmountInfo `json:"disbursement_references" yaml:"disbursement_references"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(treasury Treasury, params Params, operators []sdk.AccAddress, disbursements []Disbursement, references []ReferenceAmountInfo) GenesisState {
	return GenesisState{
		Treasury: 	treasury,
		Params: 	params,
		Operators: 	operators,
		DisbursementQueue: disbursements,
		DisbursementReferences: references,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Treasury: 	DefaultInitialTreasury(),
		Params: 	DefaultParams(),
		Operators: 	DefaultOperators(),
		DisbursementQueue: []Disbursement{},
		DisbursementReferences: []ReferenceAmountInfo{},
	}
}

// ValidateGenesis validates the treasury genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return nil
}
