package types

import (
	"fmt"
	gov "github.com/anathatech/project-anatha/x/governance"
)

const (
	ProposalTypeAddFeeExcludedMessage = "AddFeeExcludedMessage"
	ProposalTypeRemoveFeeExcludedMessage = "RemoveFeeExcludedMessage"
)

func init() {
	gov.RegisterProposalType(ProposalTypeAddFeeExcludedMessage)
	gov.RegisterProposalTypeCodec(AddFeeExcludedMessageProposal{}, "fee/AddFeeExcludedMessageProposal")
	gov.RegisterProposalType(ProposalTypeRemoveFeeExcludedMessage)
	gov.RegisterProposalTypeCodec(RemoveFeeExcludedMessageProposal{}, "fee/RemoveFeeExcludedMessageProposal")
}

// AddFeeExcludedMessageProposal
type AddFeeExcludedMessageProposal struct {
	Title       	string `json:"title" yaml:"title"`
	Description 	string `json:"description" yaml:"description"`
	MessageType 	string `json:"message_type" yaml:"message_type"`
}

func NewAddFeeExcludedMessageProposal(title string, description string, messageType string) gov.Content {
	return AddFeeExcludedMessageProposal{title, description, messageType}
}

// Implements Proposal Interface
var _ gov.Content = AddFeeExcludedMessageProposal{}

// nolint
func (p AddFeeExcludedMessageProposal) GetTitle() string       { return p.Title }
func (p AddFeeExcludedMessageProposal) GetDescription() string { return p.Description }
func (p AddFeeExcludedMessageProposal) ProposalRoute() string  { return RouterKey }
func (p AddFeeExcludedMessageProposal) ProposalType() string   { return ProposalTypeAddFeeExcludedMessage }
func (p AddFeeExcludedMessageProposal) ValidateBasic() error {
	return gov.ValidateAbstract(p)
}

func (p AddFeeExcludedMessageProposal) String() string {
	return fmt.Sprintf(`Add Fee-Excluded Message:
  Title:       %s
  Description: %s
  Message Type: %s
`, p.Title, p.Description, p.MessageType)
}

// RemoveFeeExcludedMessageProposal
type RemoveFeeExcludedMessageProposal struct {
	Title       	string `json:"title" yaml:"title"`
	Description 	string `json:"description" yaml:"description"`
	MessageType 	string `json:"message_type" yaml:"message_type"`
}

func NewRemoveFeeExcludedMessageProposal(title string, description string, messageType string) gov.Content {
	return RemoveFeeExcludedMessageProposal{title, description, messageType}
}

// Implements Proposal Interface
var _ gov.Content = RemoveFeeExcludedMessageProposal{}

// nolint
func (p RemoveFeeExcludedMessageProposal) GetTitle() string       { return p.Title }
func (p RemoveFeeExcludedMessageProposal) GetDescription() string { return p.Description }
func (p RemoveFeeExcludedMessageProposal) ProposalRoute() string  { return RouterKey }
func (p RemoveFeeExcludedMessageProposal) ProposalType() string   { return ProposalTypeRemoveFeeExcludedMessage }
func (p RemoveFeeExcludedMessageProposal) ValidateBasic() error {
	return gov.ValidateAbstract(p)
}

func (p RemoveFeeExcludedMessageProposal) String() string {
	return fmt.Sprintf(`Remove Fee-Excluded Message:
  Title:       %s
  Description: %s
  Message Type: %s
`, p.Title, p.Description, p.MessageType)
}