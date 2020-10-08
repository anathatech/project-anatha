package utils

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io/ioutil"
)

type GenericTreasuryProposalJSON struct {
	Title       string   		`json:"title" yaml:"title"`
	Description string   		`json:"description" yaml:"description"`
	Amount      sdk.Coins 		`json:"amount" yaml:"amount"`
}

func ParseGenericTreasuryProposalJSON(cdc *codec.Codec, proposalFile string) (GenericTreasuryProposalJSON, error) {
	proposal := GenericTreasuryProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
