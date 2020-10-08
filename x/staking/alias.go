package staking

// nolint

import (
	"github.com/anathatech/project-anatha/x/staking/exported"
	"github.com/anathatech/project-anatha/x/staking/keeper"
	"github.com/anathatech/project-anatha/x/staking/types"
)

const (
	DefaultParamspace                  = keeper.DefaultParamspace
	ModuleName                         = types.ModuleName
	StoreKey                           = types.StoreKey
	TStoreKey                          = types.TStoreKey
	QuerierRoute                       = types.QuerierRoute
	RouterKey                          = types.RouterKey
	DefaultUnbondingTime               = types.DefaultUnbondingTime
	DefaultMaxValidators               = types.DefaultMaxValidators
	DefaultMaxEntries                  = types.DefaultMaxEntries
	NotBondedPoolName                  = types.NotBondedPoolName
	BondedPoolName                     = types.BondedPoolName
	QueryValidators                    = types.QueryValidators
	QueryValidator                     = types.QueryValidator
	QueryDelegatorDelegations          = types.QueryDelegatorDelegations
	QueryDelegatorUnbondingDelegations = types.QueryDelegatorUnbondingDelegations
	QueryValidatorDelegations          = types.QueryValidatorDelegations
	QueryValidatorUnbondingDelegations = types.QueryValidatorUnbondingDelegations
	QueryDelegation                    = types.QueryDelegation
	QueryUnbondingDelegation           = types.QueryUnbondingDelegation
	QueryDelegatorValidators           = types.QueryDelegatorValidators
	QueryDelegatorValidator            = types.QueryDelegatorValidator
	QueryPool                          = types.QueryPool
	QueryParameters                    = types.QueryParameters
	QueryHistoricalInfo                = types.QueryHistoricalInfo
	MaxMonikerLength                   = types.MaxMonikerLength
	MaxIdentityLength                  = types.MaxIdentityLength
	MaxWebsiteLength                   = types.MaxWebsiteLength
	MaxDetailsLength                   = types.MaxDetailsLength
	DoNotModifyDesc                    = types.DoNotModifyDesc
)

var (
	// functions aliases
	RegisterInvariants                 = keeper.RegisterInvariants
	AllInvariants                      = keeper.AllInvariants
	ModuleAccountInvariants            = keeper.ModuleAccountInvariants
	NonNegativePowerInvariant          = keeper.NonNegativePowerInvariant
	PositiveDelegationInvariant        = keeper.PositiveDelegationInvariant
	DelegatorSharesInvariant           = keeper.DelegatorSharesInvariant
	NewKeeper                          = keeper.NewKeeper
	ParamKeyTable                      = keeper.ParamKeyTable
	NewQuerier                         = keeper.NewQuerier
	RegisterCodec                      = types.RegisterCodec
	NewDelegation                      = types.NewDelegation
	MustMarshalDelegation              = types.MustMarshalDelegation
	MustUnmarshalDelegation            = types.MustUnmarshalDelegation
	UnmarshalDelegation                = types.UnmarshalDelegation
	NewUnbondingDelegation             = types.NewUnbondingDelegation
	NewUnbondingDelegationEntry        = types.NewUnbondingDelegationEntry
	MustMarshalUBD                     = types.MustMarshalUBD
	MustUnmarshalUBD                   = types.MustUnmarshalUBD
	UnmarshalUBD                       = types.UnmarshalUBD
	NewDelegationResp                  = types.NewDelegationResp
	NewHistoricalInfo                  = types.NewHistoricalInfo
	MustMarshalHistoricalInfo          = types.MustMarshalHistoricalInfo
	MustUnmarshalHistoricalInfo        = types.MustUnmarshalHistoricalInfo
	UnmarshalHistoricalInfo            = types.UnmarshalHistoricalInfo
	ErrEmptyValidatorAddr              = types.ErrEmptyValidatorAddr
	ErrBadValidatorAddr                = types.ErrBadValidatorAddr
	ErrNoValidatorFound                = types.ErrNoValidatorFound
	ErrValidatorOwnerExists            = types.ErrValidatorOwnerExists
	ErrValidatorPubKeyExists           = types.ErrValidatorPubKeyExists
	ErrValidatorPubKeyTypeNotSupported = types.ErrValidatorPubKeyTypeNotSupported
	ErrValidatorJailed                 = types.ErrValidatorJailed
	ErrBadRemoveValidator              = types.ErrBadRemoveValidator
	ErrSelfDelegationBelowMinimum      = types.ErrSelfDelegationBelowMinimum
	ErrMinSelfDelegationInvalid        = types.ErrMinSelfDelegationInvalid
	ErrMinSelfDelegationDecreased      = types.ErrMinSelfDelegationDecreased
	ErrEmptyDelegatorAddr              = types.ErrEmptyDelegatorAddr
	ErrBadDenom                        = types.ErrBadDenom
	ErrBadDelegationAddr               = types.ErrBadDelegationAddr
	ErrBadDelegationAmount             = types.ErrBadDelegationAmount
	ErrNoDelegation                    = types.ErrNoDelegation
	ErrBadDelegatorAddr                = types.ErrBadDelegatorAddr
	ErrNoDelegatorForAddress           = types.ErrNoDelegatorForAddress
	ErrInsufficientShares              = types.ErrInsufficientShares
	ErrDelegationValidatorEmpty        = types.ErrDelegationValidatorEmpty
	ErrNotEnoughDelegationShares       = types.ErrNotEnoughDelegationShares
	ErrBadSharesAmount                 = types.ErrBadSharesAmount
	ErrBadSharesPercent                = types.ErrBadSharesPercent
	ErrNotMature                       = types.ErrNotMature
	ErrNoUnbondingDelegation           = types.ErrNoUnbondingDelegation
	ErrMaxUnbondingDelegationEntries   = types.ErrMaxUnbondingDelegationEntries
	ErrDelegatorShareExRateInvalid     = types.ErrDelegatorShareExRateInvalid
	ErrBothShareMsgsGiven              = types.ErrBothShareMsgsGiven
	ErrNeitherShareMsgsGiven           = types.ErrNeitherShareMsgsGiven
	ErrInvalidHistoricalInfo           = types.ErrInvalidHistoricalInfo
	ErrNoHistoricalInfo                = types.ErrNoHistoricalInfo
	NewGenesisState                    = types.NewGenesisState
	DefaultGenesisState                = types.DefaultGenesisState
	NewMultiStakingHooks               = types.NewMultiStakingHooks
	GetValidatorKey                    = types.GetValidatorKey
	GetValidatorByConsAddrKey          = types.GetValidatorByConsAddrKey
	AddressFromLastValidatorPowerKey   = types.AddressFromLastValidatorPowerKey
	GetValidatorsByPowerIndexKey       = types.GetValidatorsByPowerIndexKey
	GetLastValidatorPowerKey           = types.GetLastValidatorPowerKey
	ParseValidatorPowerRankKey         = types.ParseValidatorPowerRankKey
	GetValidatorQueueTimeKey           = types.GetValidatorQueueTimeKey
	GetDelegationKey                   = types.GetDelegationKey
	GetDelegationsKey                  = types.GetDelegationsKey
	GetUBDKey                          = types.GetUBDKey
	GetUBDByValIndexKey                = types.GetUBDByValIndexKey
	GetUBDKeyFromValIndexKey           = types.GetUBDKeyFromValIndexKey
	GetUBDsKey                         = types.GetUBDsKey
	GetUBDsByValIndexKey               = types.GetUBDsByValIndexKey
	GetUnbondingDelegationTimeKey      = types.GetUnbondingDelegationTimeKey
	GetHistoricalInfoKey               = types.GetHistoricalInfoKey
	NewMsgCreateValidator              = types.NewMsgCreateValidator
	NewMsgEditValidator                = types.NewMsgEditValidator
	NewMsgDelegate                     = types.NewMsgDelegate
	NewMsgUndelegate                   = types.NewMsgUndelegate
	NewParams                          = types.NewParams
	DefaultParams                      = types.DefaultParams
	MustUnmarshalParams                = types.MustUnmarshalParams
	UnmarshalParams                    = types.UnmarshalParams
	NewPool                            = types.NewPool
	NewQueryDelegatorParams            = types.NewQueryDelegatorParams
	NewQueryValidatorParams            = types.NewQueryValidatorParams
	NewQueryBondsParams                = types.NewQueryBondsParams
	NewQueryValidatorsParams           = types.NewQueryValidatorsParams
	NewQueryHistoricalInfoParams       = types.NewQueryHistoricalInfoParams
	NewValidator                       = types.NewValidator
	MustMarshalValidator               = types.MustMarshalValidator
	MustUnmarshalValidator             = types.MustUnmarshalValidator
	UnmarshalValidator                 = types.UnmarshalValidator
	NewDescription                     = types.NewDescription

	// variable aliases
	ModuleCdc                        = types.ModuleCdc
	LastValidatorPowerKey            = types.LastValidatorPowerKey
	LastTotalPowerKey                = types.LastTotalPowerKey
	ValidatorsKey                    = types.ValidatorsKey
	ValidatorsByConsAddrKey          = types.ValidatorsByConsAddrKey
	ValidatorsByPowerIndexKey        = types.ValidatorsByPowerIndexKey
	DelegationKey                    = types.DelegationKey
	UnbondingDelegationKey           = types.UnbondingDelegationKey
	UnbondingDelegationByValIndexKey = types.UnbondingDelegationByValIndexKey
	UnbondingQueueKey                = types.UnbondingQueueKey
	ValidatorQueueKey                = types.ValidatorQueueKey
	HistoricalInfoKey                = types.HistoricalInfoKey
	KeyUnbondingTime                 = types.KeyUnbondingTime
	KeyMaxValidators                 = types.KeyMaxValidators
	KeyMaxEntries                    = types.KeyMaxEntries
	KeyBondDenom                     = types.KeyBondDenom
)

type (
	Keeper                    = keeper.Keeper
	DVPair                    = types.DVPair
	Delegation                = types.Delegation
	Delegations               = types.Delegations
	UnbondingDelegation       = types.UnbondingDelegation
	UnbondingDelegationEntry  = types.UnbondingDelegationEntry
	UnbondingDelegations      = types.UnbondingDelegations
	HistoricalInfo            = types.HistoricalInfo
	DelegationResponse        = types.DelegationResponse
	DelegationResponses       = types.DelegationResponses
	GenesisState              = types.GenesisState
	LastValidatorPower        = types.LastValidatorPower
	MultiStakingHooks         = types.MultiStakingHooks
	MsgCreateValidator        = types.MsgCreateValidator
	MsgEditValidator          = types.MsgEditValidator
	MsgDelegate               = types.MsgDelegate
	MsgUndelegate             = types.MsgUndelegate
	Params                    = types.Params
	Pool                      = types.Pool
	QueryDelegatorParams      = types.QueryDelegatorParams
	QueryValidatorParams      = types.QueryValidatorParams
	QueryBondsParams          = types.QueryBondsParams
	QueryValidatorsParams     = types.QueryValidatorsParams
	QueryHistoricalInfoParams = types.QueryHistoricalInfoParams
	Validator                 = types.Validator
	Validators                = types.Validators
	Description               = types.Description
	DelegationI               = exported.DelegationI
	ValidatorI                = exported.ValidatorI
)
