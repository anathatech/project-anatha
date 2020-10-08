package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
)

type Keeper struct {
	CoinKeeper bank.Keeper
	AccountKeeper auth.AccountKeeper
	SupplyKeeper supply.Keeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace types.ParamSubspace
	feeCollectorName string
	hooks types.NameHooks
}

func NewKeeper(coinKeeper bank.Keeper, accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper, cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace, feeCollectorName string) Keeper {

	keeper := Keeper{
		CoinKeeper:	coinKeeper,
		AccountKeeper: accountKeeper,
		SupplyKeeper: supplyKeeper,
		storeKey:   key,
		cdc:        cdc,
		paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
		feeCollectorName: feeCollectorName,
	}
	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) SetHooks(sh types.NameHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set name hooks twice")
	}
	k.hooks = sh
	return k
}