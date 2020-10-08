package hra

import (
	"github.com/anathatech/project-anatha/x/hra/internal/keeper"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	QuerierRoute      = types.QuerierRoute
	NameConstraintBlock = types.NameConstraintBlock
)

var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	NewMultiNameHooks 	= types.NewMultiNameHooks
	NewMsgRegisterName  = types.NewMsgRegisterName
	NewMsgRenewName 	= types.NewMsgRenewName
	NewMsgSetPrice 		= types.NewMsgSetPrice
	NewMsgDeleteName 	= types.NewMsgDeleteName
	NewMsgBuyName		= types.NewMsgBuyName
	NewMsgTransferName	= types.NewMsgTransferName
	NewMsgRegisterAddress = types.NewMsgRegisterAddress
	NewMsgRemoveAddress = types.NewMsgRemoveAddress

	ModuleCdc     = types.ModuleCdc

	KeyNameInfoDuration = types.KeyNameInfoDuration

	ErrNameNotRegistered = types.ErrNameNotRegistered
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params
	MultiNameHooks = types.MultiNameHooks
	NameHooks = types.NameHooks

	MsgRegisterName = types.MsgRegisterName
	MsgRenewName	= types.MsgRenewName
	MsgSetPrice		= types.MsgSetPrice
	MsgDeleteName 	= types.MsgDeleteName
	MsgBuyName 		= types.MsgBuyName
	MsgTransferName = types.MsgTransferName
	MsgRegisterAddress = types.MsgRegisterAddress
	MsgRemoveAddress = types.MsgRemoveAddress
	MsgRemoveAllAddresses = types.MsgRemoveAllAddresses
)
