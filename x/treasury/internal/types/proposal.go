package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/config"
	gov "github.com/anathatech/project-anatha/x/governance"
)

const (
	ProposalTypeAddBuyBackLiquidity = "AddBuyBackLiquidity"
	ProposalTypeRemoveBuyBackLiquidity = "RemoveBuyBackLiquidity"
	ProposalTypeBurnDistributionProfits = "BurnDistributionProfits"
	ProposalTypeTransferFromDistributionProfitsToBuyBackLiquidity = "TransferFromDistributionProfitsToBuyBackLiquidity"
	ProposalTypeTransferFromTreasuryToSwapEscrow = "TransferFromTreasuryToSwapEscrow"
	ProposalTypeTransferSwapEscrowToBuyBack = "TransferSwapEscrowToBuyBack"
)

func init() {
	gov.RegisterProposalType(ProposalTypeAddBuyBackLiquidity)
	gov.RegisterProposalTypeCodec(AddBuyBackLiquidityProposal{}, "treasury/AddBuyBackLiquidityProposal")
	gov.RegisterProposalType(ProposalTypeRemoveBuyBackLiquidity)
	gov.RegisterProposalTypeCodec(RemoveBuyBackLiquidityProposal{}, "treasury/RemoveBuyBackLiquidityProposal")
	gov.RegisterProposalType(ProposalTypeBurnDistributionProfits)
	gov.RegisterProposalTypeCodec(BurnDistributionProfitsProposal{}, "treasury/BurnDistributionProfitsProposal")
	gov.RegisterProposalType(ProposalTypeTransferFromDistributionProfitsToBuyBackLiquidity)
	gov.RegisterProposalTypeCodec(TransferFromDistributionProfitsToBuyBackLiquidityProposal{}, "treasury/TransferFromDistributionProfitsToBuyBackLiquidityProposal")
	gov.RegisterProposalType(ProposalTypeTransferFromTreasuryToSwapEscrow)
	gov.RegisterProposalTypeCodec(TransferFromTreasuryToSwapEscrowProposal{}, "treasury/TransferFromTreasuryToSwapEscrowProposal")
	gov.RegisterProposalType(ProposalTypeTransferSwapEscrowToBuyBack)
	gov.RegisterProposalTypeCodec(TransferFromSwapEscrowToBuyBackProposal{}, "treasury/TransferFromSwapEscrowToBuyBackProposal")
}

var _ gov.Content = AddBuyBackLiquidityProposal{}

type AddBuyBackLiquidityProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Amount 		sdk.Coins `json:"amount" yaml:"amount"`
}

func NewAddBuyBackLiquidityProposal(title, description string, amount sdk.Coins) gov.Content {
	return AddBuyBackLiquidityProposal{title, description, amount}
}

func (p AddBuyBackLiquidityProposal) GetTitle() string       { return p.Title }
func (p AddBuyBackLiquidityProposal) GetDescription() string { return p.Description }
func (p AddBuyBackLiquidityProposal) ProposalRoute() string  { return RouterKey }
func (p AddBuyBackLiquidityProposal) ProposalType() string   { return ProposalTypeAddBuyBackLiquidity }
func (p AddBuyBackLiquidityProposal) ValidateBasic() error {
	if ! p.Amount.IsValid() || p.Amount.AmountOf(config.DefaultStableDenom).IsZero() {
		return sdkerrors.ErrInvalidCoins
	}

	return gov.ValidateAbstract(p)
}

func (p AddBuyBackLiquidityProposal) String() string {
	return fmt.Sprintf(`Add BuyBack Liquidity Proposal:
  Title: 		%s
  Description: 	%s
  Amount: 		%s
`, p.Title, p.Description, p.Amount)
}

var _ gov.Content = RemoveBuyBackLiquidityProposal{}

type RemoveBuyBackLiquidityProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Amount 		sdk.Coins `json:"amount" yaml:"amount"`
}

func NewRemoveBuyBackLiquidityProposal(title, description string, amount sdk.Coins) gov.Content {
	return RemoveBuyBackLiquidityProposal{title, description, amount}
}

func (p RemoveBuyBackLiquidityProposal) GetTitle() string       { return p.Title }
func (p RemoveBuyBackLiquidityProposal) GetDescription() string { return p.Description }
func (p RemoveBuyBackLiquidityProposal) ProposalRoute() string  { return RouterKey }
func (p RemoveBuyBackLiquidityProposal) ProposalType() string   { return ProposalTypeRemoveBuyBackLiquidity }
func (p RemoveBuyBackLiquidityProposal) ValidateBasic() error {
	if ! p.Amount.IsValid() || p.Amount.AmountOf(config.DefaultStableDenom).IsZero() {
		return sdkerrors.ErrInvalidCoins
	}

	return gov.ValidateAbstract(p)
}

func (p RemoveBuyBackLiquidityProposal) String() string {
	return fmt.Sprintf(`Remove BuyBack Liquidity Proposal:
  Title: 		%s
  Description: 	%s
  Amount: 		%s
`, p.Title, p.Description, p.Amount)
}

var _ gov.Content = BurnDistributionProfitsProposal{}

type BurnDistributionProfitsProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Amount 		sdk.Coins `json:"amount" yaml:"amount"`
}

func NewBurnDistributionProfitsProposal(title, description string, amount sdk.Coins) gov.Content {
	return BurnDistributionProfitsProposal{title, description, amount}
}

func (p BurnDistributionProfitsProposal) GetTitle() string       { return p.Title }
func (p BurnDistributionProfitsProposal) GetDescription() string { return p.Description }
func (p BurnDistributionProfitsProposal) ProposalRoute() string  { return RouterKey }
func (p BurnDistributionProfitsProposal) ProposalType() string   { return ProposalTypeBurnDistributionProfits }
func (p BurnDistributionProfitsProposal) ValidateBasic() error {
	if ! p.Amount.IsValid() || p.Amount.AmountOf(config.DefaultStableDenom).IsZero() {
		return sdkerrors.ErrInvalidCoins
	}

	return gov.ValidateAbstract(p)
}

func (p BurnDistributionProfitsProposal) String() string {
	return fmt.Sprintf(`Burn Distribution Profits Proposal:
  Title: 		%s
  Description: 	%s
  Amount: 		%s
`, p.Title, p.Description, p.Amount)
}

var _ gov.Content = TransferFromDistributionProfitsToBuyBackLiquidityProposal{}

type TransferFromDistributionProfitsToBuyBackLiquidityProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Amount 		sdk.Coins `json:"amount" yaml:"amount"`
}

func NewTransferFromDistributionProfitsToBuyBackLiquidityProposal(title, description string, amount sdk.Coins) gov.Content {
	return TransferFromDistributionProfitsToBuyBackLiquidityProposal{title, description, amount}
}

func (p TransferFromDistributionProfitsToBuyBackLiquidityProposal) GetTitle() string       { return p.Title }
func (p TransferFromDistributionProfitsToBuyBackLiquidityProposal) GetDescription() string { return p.Description }
func (p TransferFromDistributionProfitsToBuyBackLiquidityProposal) ProposalRoute() string  { return RouterKey }
func (p TransferFromDistributionProfitsToBuyBackLiquidityProposal) ProposalType() string   { return ProposalTypeTransferFromDistributionProfitsToBuyBackLiquidity }
func (p TransferFromDistributionProfitsToBuyBackLiquidityProposal) ValidateBasic() error {
	if ! p.Amount.IsValid() || p.Amount.AmountOf(config.DefaultStableDenom).IsZero() {
		return sdkerrors.ErrInvalidCoins
	}

	return gov.ValidateAbstract(p)
}

func (p TransferFromDistributionProfitsToBuyBackLiquidityProposal) String() string {
	return fmt.Sprintf(`Transfer From Distribution Profits To BuyBack Liquidity Proposal:
  Title: 		%s
  Description: 	%s
  Amount: 		%s
`, p.Title, p.Description, p.Amount)
}

var _ gov.Content = TransferFromTreasuryToSwapEscrowProposal{}

type TransferFromTreasuryToSwapEscrowProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Amount 		sdk.Coins `json:"amount" yaml:"amount"`
}

func NewTransferFromTreasuryToSwapEscrowProposal(title, description string, amount sdk.Coins) gov.Content {
	return TransferFromTreasuryToSwapEscrowProposal{title, description, amount}
}

func (p TransferFromTreasuryToSwapEscrowProposal) GetTitle() string       { return p.Title }
func (p TransferFromTreasuryToSwapEscrowProposal) GetDescription() string { return p.Description }
func (p TransferFromTreasuryToSwapEscrowProposal) ProposalRoute() string  { return RouterKey }
func (p TransferFromTreasuryToSwapEscrowProposal) ProposalType() string   { return ProposalTypeTransferFromTreasuryToSwapEscrow }
func (p TransferFromTreasuryToSwapEscrowProposal) ValidateBasic() error {
	if ! p.Amount.IsValid() || p.Amount.AmountOf(config.DefaultDenom).IsZero() {
		return sdkerrors.ErrInvalidCoins
	}

	return gov.ValidateAbstract(p)
}

func (p TransferFromTreasuryToSwapEscrowProposal) String() string {
	return fmt.Sprintf(`Transfer from Treasury to Swap Escrow Proposal:
  Title: 		%s
  Description: 	%s
  Amount: 		%s
`, p.Title, p.Description, p.Amount)
}

var _ gov.Content = TransferFromSwapEscrowToBuyBackProposal{}

type TransferFromSwapEscrowToBuyBackProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Amount 		sdk.Coins `json:"amount" yaml:"amount"`
}

func NewTransferFromSwapEscrowToBuyBackProposal(title, description string, amount sdk.Coins) gov.Content {
	return TransferFromSwapEscrowToBuyBackProposal{title, description, amount}
}

func (p TransferFromSwapEscrowToBuyBackProposal) GetTitle() string       { return p.Title }
func (p TransferFromSwapEscrowToBuyBackProposal) GetDescription() string { return p.Description }
func (p TransferFromSwapEscrowToBuyBackProposal) ProposalRoute() string  { return RouterKey }
func (p TransferFromSwapEscrowToBuyBackProposal) ProposalType() string   { return ProposalTypeTransferSwapEscrowToBuyBack }
func (p TransferFromSwapEscrowToBuyBackProposal) ValidateBasic() error {
	if ! p.Amount.IsValid() || p.Amount.AmountOf(config.DefaultDenom).IsZero() {
		return sdkerrors.ErrInvalidCoins
	}

	return gov.ValidateAbstract(p)
}

func (p TransferFromSwapEscrowToBuyBackProposal) String() string {
	return fmt.Sprintf(`Transfer from Swap Escrow to BuyBack Proposal:
  Title: 		%s
  Description: 	%s
  Amount: 		%s
`, p.Title, p.Description, p.Amount)
}