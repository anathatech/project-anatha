package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/config"
	gov "github.com/anathatech/project-anatha/x/governance"
)

const (
	ProposalTypeDevelopmentFundDistribution   = "DevelopmentFundDistribution"
	ProposalTypeSecurityTokenFundDistribution = "SecurityTokenFundDistribution"
)

type DevelopmentFundDistributionProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Amount 		sdk.Coins `json:"amount" yaml:"amount"`
	Recipient 	sdk.AccAddress `json:"recipient" yaml:"recipient"`
}

func NewDevelopmentFundDistributionProposal(title, description string, amount sdk.Coins, recipient sdk.AccAddress) gov.Content {
	return DevelopmentFundDistributionProposal{title, description, amount,recipient}
}

var _ gov.Content = DevelopmentFundDistributionProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypeDevelopmentFundDistribution)
	gov.RegisterProposalTypeCodec(DevelopmentFundDistributionProposal{}, "distribution/DevelopmentFundDistributionProposal")
	gov.RegisterProposalType(ProposalTypeSecurityTokenFundDistribution)
	gov.RegisterProposalTypeCodec(SecurityTokenFundDistributionProposal{}, "distribution/SecurityTokenFundDistributionProposal")
}

func (p DevelopmentFundDistributionProposal) GetTitle() string       { return p.Title }
func (p DevelopmentFundDistributionProposal) GetDescription() string { return p.Description }
func (p DevelopmentFundDistributionProposal) ProposalRoute() string  { return RouterKey }
func (p DevelopmentFundDistributionProposal) ProposalType() string   { return ProposalTypeDevelopmentFundDistribution }
func (p DevelopmentFundDistributionProposal) ValidateBasic() error {
	if ! p.Amount.IsValid() || p.Amount.AmountOf(config.DefaultDenom).IsZero() {
		return sdkerrors.ErrInvalidCoins
	}

	return gov.ValidateAbstract(p)
}

func (sup DevelopmentFundDistributionProposal) String() string {
	return fmt.Sprintf(`Development Fund Distribution Proposal:
  Title: 		%s
  Description: 	%s
  Amount: 		%s
  Recipient: 	%s
`, sup.Title, sup.Description, sup.Amount, sup.Recipient)
}

type SecurityTokenFundDistributionProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Recipients 	[]Recipients `json:"recipients" yaml:"recipients"`
}

type Recipients struct {
	Amount sdk.Coins `json:"amount" yaml:"amount"`
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
}

func NewSecurityTokenFundDistributionProposal(title, description string, recipients []Recipients) gov.Content {
	return SecurityTokenFundDistributionProposal{title, description ,recipients}
}

var _ gov.Content = SecurityTokenFundDistributionProposal{}

func (p SecurityTokenFundDistributionProposal) GetTitle() string       { return p.Title }
func (p SecurityTokenFundDistributionProposal) GetDescription() string { return p.Description }
func (p SecurityTokenFundDistributionProposal) ProposalRoute() string  { return RouterKey }
func (p SecurityTokenFundDistributionProposal) ProposalType() string   { return ProposalTypeSecurityTokenFundDistribution
}
func (p SecurityTokenFundDistributionProposal) ValidateBasic() error {
	for i := 0; i < len(p.Recipients); i++ {
		if ! p.Recipients[i].Amount.IsValid() || p.Recipients[i].Amount.AmountOf(config.DefaultDenom).IsZero() {
			return sdkerrors.ErrInvalidCoins
		}
	}

	return gov.ValidateAbstract(p)
}

func (sup SecurityTokenFundDistributionProposal) String() string {
	return fmt.Sprintf(`Security Token Fund Distribution Proposal:
  Title: 		%s
  Description: 	%s
  Recipients: 	%s
`, sup.Title, sup.Description, sup.Recipients)
}