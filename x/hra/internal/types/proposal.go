package types

import (
	"fmt"

	gov "github.com/anathatech/project-anatha/x/governance"
)

const (
	ProposalTypeRegisterBlockchainId 	= "RegisterBlockchainId"
	ProposalTypeRemoveBlockchainId 		= "RemoveBlockchainId"
)

func init() {
	gov.RegisterProposalType(ProposalTypeRegisterBlockchainId)
	gov.RegisterProposalTypeCodec(RegisterBlockchainIdProposal{}, "hra/RegisterBlockchainIdProposal")
	gov.RegisterProposalType(ProposalTypeRemoveBlockchainId)
	gov.RegisterProposalTypeCodec(RemoveBlockchainIdProposal{}, "hra/RemoveBlockchainIdProposal")
}

// RegisterBlockchainIdProposal
type RegisterBlockchainIdProposal struct {
	Title       	string `json:"title" yaml:"title"`
	Description 	string `json:"description" yaml:"description"`
	BlockchainId 	string `json:"blockchain_id" yaml:"blockchain_id"`
}

func NewRegisterBlockchainIdProposal(title string, description string, blockchainId string) gov.Content {
	return RegisterBlockchainIdProposal{title, description, blockchainId}
}

// Implements Proposal Interface
var _ gov.Content = RegisterBlockchainIdProposal{}

// nolint
func (p RegisterBlockchainIdProposal) GetTitle() string       { return p.Title }
func (p RegisterBlockchainIdProposal) GetDescription() string { return p.Description }
func (p RegisterBlockchainIdProposal) ProposalRoute() string  { return RouterKey }
func (p RegisterBlockchainIdProposal) ProposalType() string   { return ProposalTypeRegisterBlockchainId }
func (p RegisterBlockchainIdProposal) ValidateBasic() error {
	if err := validateBlockchainId(p.BlockchainId); err != nil {
		return err
	}
	return gov.ValidateAbstract(p)
}

func (p RegisterBlockchainIdProposal) String() string {
	return fmt.Sprintf(`Register Blockchain Id:
  Title:       %s
  Description: %s
  Blockchain Id: %s
`, p.Title, p.Description, p.BlockchainId)
}

// RemoveBlockchainIdProposal
type RemoveBlockchainIdProposal struct {
	Title       	string `json:"title" yaml:"title"`
	Description 	string `json:"description" yaml:"description"`
	BlockchainId 	string `json:"blockchain_id" yaml:"blockchain_id"`
}

func NewRemoveBlockchainIdProposal(title string, description string, blockchainId string) gov.Content {
	return RemoveBlockchainIdProposal{title, description, blockchainId}
}

// Implements Proposal Interface
var _ gov.Content = RemoveBlockchainIdProposal{}

// nolint
func (p RemoveBlockchainIdProposal) GetTitle() string       { return p.Title }
func (p RemoveBlockchainIdProposal) GetDescription() string { return p.Description }
func (p RemoveBlockchainIdProposal) ProposalRoute() string  { return RouterKey }
func (p RemoveBlockchainIdProposal) ProposalType() string   { return ProposalTypeRemoveBlockchainId }
func (p RemoveBlockchainIdProposal) ValidateBasic() error 	{
	if err := validateBlockchainId(p.BlockchainId); err != nil {
		return err
	}
	return gov.ValidateAbstract(p)
}

func (p RemoveBlockchainIdProposal) String() string {
	return fmt.Sprintf(`Remove Blockchain Id:
  Title:       %s
  Description: %s
  Blockchain Id: %s
`, p.Title, p.Description, p.BlockchainId)
}