package utils

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"io/ioutil"
)

type (

	PlanJSON struct {
		Name 	string 	`json:"name" yaml:"name"`
		Height 	int64 	`json:"height" yaml:"height"`
		Info 	string 	`json:"info" yaml:"info"`
	}

	SoftwareUpgradeProposalJSON struct {
		Title       string   `json:"title" yaml:"title"`
		Description string   `json:"description" yaml:"description"`
		Plan        PlanJSON `json:"plan" yaml:"plan"`
	}

	SoftwareUpgradeProposalReq struct {
		BaseReq 	rest.BaseReq 	`json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Plan        PlanJSON       `json:"plan" yaml:"plan"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
	}
)

func NewPlanJSON(name string, height int64, info string) PlanJSON {
	return PlanJSON{name, height, info}
}

func ParseSoftwareUpgradeProposalJSON(cdc *codec.Codec, proposalFile string) (SoftwareUpgradeProposalJSON, error) {
	proposal := SoftwareUpgradeProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
