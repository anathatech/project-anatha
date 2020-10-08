package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"strings"
	"time"

	gov "github.com/cosmos/cosmos-sdk/x/gov"
)

const DefaultStartingProposalID uint64 = 1

func init() {
	RegisterProposalType(params.ProposalTypeChange)
}

type Proposal struct {
	gov.Content `json:"content" yaml:"content"` // Proposal content interface

	ProposalID  uint64         `json:"id" yaml:"id"`
	Status      ProposalStatus `json:"proposal_status" yaml:"proposal_status"`
	TallyResult TallyResult    `json:"tally_result" yaml:"tally_result"`

	VotingStartTime time.Time `json:"voting_start_time" yaml:"voting_start_time"`
	VotingEndTime   time.Time `json:"voting_end_time" yaml:"voting_end_time"`

	ExecutionTime time.Time `json:"proposal_execution_time" yaml:"proposal_execution_time"`
	Expedited     bool      `json:"expedited" yaml:"expedited"`
}

func NewProposal(content gov.Content, id uint64) Proposal {
	return Proposal{
		Content:       content,
		ProposalID:    id,
		TallyResult:   EmptyTallyResult(),
		ExecutionTime: time.Time{},
		Expedited:     false,
	}
}

func (p Proposal) String() string {
	return fmt.Sprintf(`Proposal %d:
  Title:              %s
  Type:               %s
  Status:             %s
  Voting Start Time:  %s
  Voting End Time:    %s
  Execution Time:     %s
  Expedited: 	      %s
  Description:        %s`,
		p.ProposalID, p.GetTitle(), p.ProposalType(),
		p.Status, p.VotingStartTime, p.VotingEndTime, p.ExecutionTime, p.Expedited , p.GetDescription(),
	)
}

type Proposals []Proposal

func (p Proposals) String() string {
	out := "ID - (Status) [Type] Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) [%s] %s\n",
			prop.ProposalID, prop.Status,
			prop.ProposalType(), prop.GetTitle())
	}
	return strings.TrimSpace(out)
}

type (
	ProposalStatus byte
)

const (
	StatusNil           ProposalStatus = 0x00
	StatusVotingPeriod  ProposalStatus = 0x01
	StatusExpediting	ProposalStatus = 0x05
	StatusPassed        ProposalStatus = 0x02
	StatusRejected      ProposalStatus = 0x03
	StatusFailed        ProposalStatus = 0x04
)

func ProposalStatusFromString(str string) (ProposalStatus, error) {
	switch str {
	case "VotingPeriod":
		return StatusVotingPeriod, nil

	case "Passed":
		return StatusPassed, nil

	case "Rejected":
		return StatusRejected, nil

	case "Failed":
		return StatusFailed, nil

	case "":
		return StatusNil, nil

	default:
		return ProposalStatus(0xff), fmt.Errorf("'%s' is not a valid proposal status", str)
	}
}

func ValidProposalStatus(status ProposalStatus) bool {
	if status == StatusVotingPeriod ||
		status == StatusPassed ||
		status == StatusRejected ||
		status == StatusFailed {
		return true
	}
	return false
}

func (status ProposalStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

func (status *ProposalStatus) Unmarshal(data []byte) error {
	*status = ProposalStatus(data[0])
	return nil
}

func (status ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

func (status *ProposalStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ProposalStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

func (status ProposalStatus) String() string {
	switch status {
	case StatusVotingPeriod:
		return "VotingPeriod"

	case StatusPassed:
		return "Passed"

	case StatusRejected:
		return "Rejected"

	case StatusFailed:
		return "Failed"

	default:
		return ""
	}
}

func (status ProposalStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}

var validProposalTypes = map[string]struct{}{
	ProposalTypeText: {},
	ProposalTypeAddGovernor: {},
	ProposalTypeRemoveGovernor: {},
}

func RegisterProposalType(ty string) {
	if _, ok := validProposalTypes[ty]; ok {
		panic(fmt.Sprintf("already registered proposal type: %s", ty))
	}

	validProposalTypes[ty] = struct{}{}
}

func ContentFromProposalType(title, desc, ty string) gov.Content {
	switch ty {
	case ProposalTypeText:
		return NewTextProposal(title, desc)

	default:
		return nil
	}
}

func IsValidProposalType(ty string) bool {
	_, ok := validProposalTypes[ty]
	return ok
}

const (
	ProposalTypeText string = "Text"
	ProposalTypeAddGovernor = "AddGovernor"
	ProposalTypeRemoveGovernor = "RemoveGovernor"
)

type TextProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
}

func NewTextProposal(title, description string) gov.Content {
	return TextProposal{
		title,
		description,
	}
}

var _ gov.Content = TextProposal{}

func (tp TextProposal) GetTitle() string { return tp.Title }
func (tp TextProposal) GetDescription() string { return tp.Description }
func (tp TextProposal) ProposalRoute() string { return RouterKey }
func (tp TextProposal) ProposalType() string { return ProposalTypeText }
func (tp TextProposal) ValidateBasic() error { return ValidateAbstract(tp) }

func (tp TextProposal) String() string {
	return fmt.Sprintf(`Text Proposal:
  Title:       %s
  Description: %s
`, tp.Title, tp.Description)
}

type AddGovernorProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Governor sdk.AccAddress `json:"governor" yaml:"governor"`
}

func NewAddGovernorProposal(title string, description string, governor sdk.AccAddress) gov.Content {
	return AddGovernorProposal{
		Title: title,
		Description: description,
		Governor: governor,
	}
}

var _ gov.Content = AddGovernorProposal{}

func (p AddGovernorProposal) GetTitle() string { return p.Title }
func (p AddGovernorProposal) GetDescription() string { return p.Description }
func (p AddGovernorProposal) ProposalRoute() string { return RouterKey }
func (p AddGovernorProposal) ProposalType() string { return ProposalTypeAddGovernor }
func (p AddGovernorProposal) ValidateBasic() error { return ValidateAbstract(p) }

func (p AddGovernorProposal) String() string {
	return fmt.Sprintf(`Add Governor Proposal:
  Title:       %s
  Description: %s
  Governor: %s
`, p.Title, p.Description, p.Governor)
}

type RemoveGovernorProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Governor sdk.AccAddress `json:"governor" yaml:"governor"`
}

func NewRemoveGovernorProposal(title string, description string, governor sdk.AccAddress) gov.Content {
	return RemoveGovernorProposal{
		Title: title,
		Description: description,
		Governor: governor,
	}
}

var _ gov.Content = RemoveGovernorProposal{}

func (p RemoveGovernorProposal) GetTitle() string { return p.Title }
func (p RemoveGovernorProposal) GetDescription() string { return p.Description }
func (p RemoveGovernorProposal) ProposalRoute() string { return RouterKey }
func (p RemoveGovernorProposal) ProposalType() string { return ProposalTypeRemoveGovernor }
func (p RemoveGovernorProposal) ValidateBasic() error { return ValidateAbstract(p) }

func (p RemoveGovernorProposal) String() string {
	return fmt.Sprintf(`Remove Governor Proposal:
  Title:       %s
  Description: %s
  Governor: %s
`, p.Title, p.Description, p.Governor)
}
