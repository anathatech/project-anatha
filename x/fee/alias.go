package fee

import (
	"github.com/anathatech/project-anatha/x/fee/internal/keeper"
	"github.com/anathatech/project-anatha/x/fee/internal/types"
)

const (
	ModuleName               = types.ModuleName
	RouterKey                = types.RouterKey
	StoreKey                 = types.StoreKey
	DefaultParamspace        = types.DefaultParamspace
	QuerierRoute             = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper                          = keeper.NewKeeper
	NewQuerier                         = keeper.NewQuerier
	RegisterCodec                      = types.RegisterCodec
	NewGenesisState                    = types.NewGenesisState
	DefaultGenesisState                = types.DefaultGenesisState
	ValidateGenesis                    = types.ValidateGenesis

	// variable aliases
	ModuleCdc     = types.ModuleCdc
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params
)