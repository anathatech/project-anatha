package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Vote struct {
	ProposalID uint64         `json:"proposal_id" yaml:"proposal_id"`
	Voter      sdk.AccAddress `json:"voter" yaml:"voter"`
	Option     VoteOption     `json:"option" yaml:"option"`
}

func NewVote(proposalID uint64, voter sdk.AccAddress, option VoteOption) Vote {
	return Vote{
		proposalID,
		voter,
		option,
	}
}

func (v Vote) String() string {
	return fmt.Sprintf("voter %s voted with option %s on proposal %d", v.Voter, v.Option, v.ProposalID)
}

type Votes []Vote

func (v Votes) String() string {
	if len(v) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Votes for Proposal %d:", v[0].ProposalID)
	for _, vot := range v {
		out += fmt.Sprintf("\n  %s: %s", vot.Voter, vot.Option)
	}
	return out
}

func (v Vote) Equals(comp Vote) bool {
	return v.Voter.Equals(comp.Voter) &&
		v.ProposalID == comp.ProposalID &&
		v.Option == comp.Option
}

func (v Vote) Empty() bool {
	return v.Equals(Vote{})
}

type VoteOption byte

const (
	OptionEmpty      VoteOption = 0x00
	OptionYes        VoteOption = 0x01
	OptionNo         VoteOption = 0x02
)

func VoteOptionFromString(str string) (VoteOption, error) {
	switch str {
	case "Yes":
		return OptionYes, nil

	case "No":
		return OptionNo, nil

	default:
		return VoteOption(0xff), fmt.Errorf("'%s' is not a valid vote option", str)
	}
}

func ValidVoteOption(option VoteOption) bool {
	if option == OptionYes ||
		option == OptionNo  {
		return true
	}
	return false
}

func (vo VoteOption) Marshal() ([]byte, error) {
	return []byte{byte(vo)}, nil
}

func (vo *VoteOption) Unmarshal(data []byte) error {
	*vo = VoteOption(data[0])
	return nil
}

func (vo VoteOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(vo.String())
}

func (vo *VoteOption) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := VoteOptionFromString(s)
	if err != nil {
		return err
	}

	*vo = bz2
	return nil
}

func (vo VoteOption) String() string {
	switch vo {
	case OptionYes:
		return "Yes"
	case OptionNo:
		return "No"

	default:
		return ""
	}
}

func (vo VoteOption) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(vo.String()))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(vo))))
	}
}