package utils

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"io/ioutil"
)

type FeeExclusionProposalJSON struct {
	Title 			string `json:"title" yaml:"title"`
	Description 	string `json:"description" yaml:"description"`
	MessageType 	string `json:"message_type" yaml:"message_type"`
}

func ParseFeeExclusionProposalJSON(cdc *codec.Codec, proposalFile string) (FeeExclusionProposalJSON, error) {
	proposal := FeeExclusionProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
