package distribution

import (
	"github.com/anathatech/project-anatha/x/distribution/internal/keeper"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
)

const (
	ModuleName = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey = types.StoreKey
	RouterKey = types.RouterKey
	QuerierRoute = types.QuerierRoute

	AmcModuleName               = types.AmcModuleName
	NvrpModuleName              = types.NvrpModuleName
	NvrpDistributionModuleName  = types.NvrpDistributionModuleName
	HRAHolderRewardModuleName   = types.HRAHolderRewardModuleName
	DevelopmentFundModuleName   = types.DevelopmentFundModuleName
	SecurityTokenFundModuleName = types.SecurityTokenFundModuleName
	SavingsModuleName           = types.SavingsModuleName
	SavingsDistributionModuleName = types.SavingsDistributionModuleName
)

var (
	NewKeeper                                = keeper.NewKeeper
	NewQuerier                               = keeper.NewQuerier
	NewGenesisState                          = types.NewGenesisState
	DefaultGenesisState                      = types.DefaultGenesisState
	ValidateGenesis                          = types.ValidateGenesis
	ParamKeyTable                            = types.ParamKeyTable
	NewParams                                = types.NewParams
	DefaultParams                            = types.DefaultParams
	NewDevelopmentFundDistributionProposal   = types.NewDevelopmentFundDistributionProposal
	NewSecurityTokenFundDistributionProposal = types.NewSecurityTokenFundDistributionProposal

	NewMsgWithdrawNameReward                 = types.NewMsgWithdrawNameReward
	ModuleCdc                                = types.ModuleCdc
	RegisterCodec                            = types.RegisterCodec
)

type (
	Keeper = keeper.Keeper
	Params = types.Params
	GenesisState = types.GenesisState

	DevelopmentFundDistributionProposal   = types.DevelopmentFundDistributionProposal
	SecurityTokenFundDistributionProposal = types.SecurityTokenFundDistributionProposal

	MsgWithdrawNameReward                 = types.MsgWithdrawNameReward
	MsgWithdrawValidatorReward            = types.MsgWithdrawValidatorReward
	MsgDepositSavings                     = types.MsgDepositSavings
	MsgWithdrawSavings                    = types.MsgWithdrawSavings
	MsgWithdrawSavingsInterest            = types.MsgWithdrawSavingsInterest
)
