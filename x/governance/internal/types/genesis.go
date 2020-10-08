package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	StartingProposalID 	uint64        `json:"starting_proposal_id" yaml:"starting_proposal_id"`
	Votes              	Votes         `json:"votes" yaml:"votes"`
	Proposals          	Proposals     `json:"proposals" yaml:"proposals"`
	VotingParams       	VotingParams  `json:"voting_params" yaml:"voting_params"`
	TallyParams        	TallyParams   `json:"tally_params" yaml:"tally_params"`
	Governors			[]sdk.AccAddress `json:"governors" yaml:"governors"`
}

func NewGenesisState(startingProposalID uint64, vp VotingParams, tp TallyParams, governors []sdk.AccAddress) GenesisState {
	return GenesisState{
		StartingProposalID: startingProposalID,
		VotingParams:       vp,
		TallyParams:        tp,
		Governors: 			governors,
	}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(
		DefaultStartingProposalID,
		DefaultVotingParams(),
		DefaultTallyParams(),
		DefaultGovernors(),
	)
}

func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

func (data GenesisState) IsEmpty() bool {
	return data.Equal(GenesisState{})
}

func ValidateGenesis(data GenesisState) error {
	threshold := data.TallyParams.Threshold
	if threshold.IsNegative() || threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("governance vote threshold should be positive and less or equal to one, is %s",
			threshold.String())
	}

	if len(data.Governors) == 0 {
		return fmt.Errorf("list of governors should not be empty")
	}
	
	return nil
}