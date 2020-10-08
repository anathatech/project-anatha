package utils

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
)

type BlockchainIdProposalJSON struct {
	Title 			string `json:"title" yaml:"title"`
	Description 	string `json:"description" yaml:"description"`
	BlockchainId 	string `json:"blockchain_id" yaml:"blockchain_id"`
}

func ParseBlockchainIdProposalJSON(cdc *codec.Codec, proposalFile string) (BlockchainIdProposalJSON, error) {
	proposal := BlockchainIdProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
