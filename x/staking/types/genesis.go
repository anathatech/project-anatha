package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	Params               Params                `json:"params" yaml:"params"`
	LastTotalPower       sdk.Int               `json:"last_total_power" yaml:"last_total_power"`
	LastValidatorPowers  []LastValidatorPower  `json:"last_validator_powers" yaml:"last_validator_powers"`
	Validators           Validators            `json:"validators" yaml:"validators"`
	Delegations          Delegations           `json:"delegations" yaml:"delegations"`
	UnbondingDelegations []UnbondingDelegation `json:"unbonding_delegations" yaml:"unbonding_delegations"`
	Exported             bool                  `json:"exported" yaml:"exported"`
	StartingTicket       uint64                 `json:"starting_ticket" yaml:"starting_ticket"`
}

// LastValidatorPower required for validator set update logic
type LastValidatorPower struct {
	Address sdk.ValAddress
	Power   int64
}

// NewGenesisState creates a new GenesisState instanc e
func NewGenesisState(params Params, validators []Validator, delegations []Delegation, startingTicket uint64) GenesisState {
	return GenesisState{
		Params:      params,
		Validators:  validators,
		Delegations: delegations,
		StartingTicket: startingTicket,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
		StartingTicket: 1,
	}
}
