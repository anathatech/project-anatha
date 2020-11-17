package treasury

import (
	"github.com/anathatech/project-anatha/x/treasury/internal/keeper"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
)

const (
	ModuleName               = types.ModuleName
	RouterKey                = types.RouterKey
	StoreKey                 = types.StoreKey
	DefaultParamspace        = types.DefaultParamspace
	QuerierRoute             = types.QuerierRoute
	BuyBackLiquidityFundModuleName = types.BuyBackLiquidityFundModuleName
	BuyBackFundModuleName = types.BuyBackFundModuleName
	DistributionProfitsModuleName = types.DistributionProfitsModuleName
	TreasuryEscrowModuleName = types.TreasuryEscrowModuleName
	SwapEscrowModuleName	= types.SwapEscrowModuleName
)

var (
	// functions aliases
	NewKeeper                          = keeper.NewKeeper
	NewQuerier                         = keeper.NewQuerier
	RegisterCodec                      = types.RegisterCodec
	NewGenesisState                    = types.NewGenesisState
	DefaultGenesisState                = types.DefaultGenesisState
	ValidateGenesis                    = types.ValidateGenesis
	NewMsgAddOperator 					= types.NewMsgAddOperator
	NewMsgRemoveOperator 				= types.NewMsgRemoveOperator
	NewMsgDisburse 						= types.NewMsgDisburse
	NewMsgDisburseToEscrow              = types.NewMsgDisburseToEscrow
	NewMsgDisburseFromEscrow            = types.NewMsgDisburseFromEscrow
	NewMsgRevertFromEscrow              = types.NewMsgRevertFromEscrow
	NewMsgCancelDisbursement			= types.NewMsgCancelDisbursement
	NewMsgCreateSellOrder				= types.NewMsgCreateSellOrder
	NewMsgCreateBuyOrder				= types.NewMsgCreateBuyOrder
	NewMsgSwap							= types.NewMsgSwap

	NewAddBuyBackLiquidityProposal = types.NewAddBuyBackLiquidityProposal
	NewRemoveBuyBackLiquidityProposal = types.NewRemoveBuyBackLiquidityProposal
	NewBurnDistributionProfitsProposal = types.NewBurnDistributionProfitsProposal
	NewTransferFromDistributionProfitsToBuyBackLiquidityProposal = types.NewTransferFromDistributionProfitsToBuyBackLiquidityProposal
	NewTransferFromTreasuryToSwapEscrowProposal = types.NewTransferFromTreasuryToSwapEscrowProposal
	NewTransferFromSwapEscrowToBuyBackProposal = types.NewTransferFromSwapEscrowToBuyBackProposal

	// variable aliases
	ModuleCdc     = types.ModuleCdc
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	AddBuyBackLiquidityProposal = types.AddBuyBackLiquidityProposal
	RemoveBuyBackLiquidityProposal = types.RemoveBuyBackLiquidityProposal
	BurnDistributionProfitsProposal = types.BurnDistributionProfitsProposal
	TransferFromDistributionProfitsToBuyBackLiquidityProposal = types.TransferFromDistributionProfitsToBuyBackLiquidityProposal
	TransferFromTreasuryToSwapEscrowProposal = types.TransferFromTreasuryToSwapEscrowProposal
	TransferFromSwapEscrowToBuyBackProposal = types.TransferFromSwapEscrowToBuyBackProposal

	MsgAddOperator 					= types.MsgAddOperator
	MsgRemoveOperator 				= types.MsgRemoveOperator
	MsgDisburse 					= types.MsgDisburse
	MsgDisburseToEscrow             = types.MsgDisburseToEscrow
	MsgDisburseFromEscrow           = types.MsgDisburseFromEscrow
	MsgRevertFromEscrow             = types.MsgRevertFromEscrow
	MsgCancelDisbursement			= types.MsgCancelDisbursement
	MsgCreateSellOrder				= types.MsgCreateSellOrder
	MsgCreateBuyOrder				= types.MsgCreateBuyOrder
	MsgSwap							= types.MsgSwap
)
