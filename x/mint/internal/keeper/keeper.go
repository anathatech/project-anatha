package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/anathatech/project-anatha/config"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/anathatech/project-anatha/x/mint/internal/types"
)

type Keeper struct {
	cdc              *codec.Codec
	storeKey         sdk.StoreKey
	paramSpace       params.Subspace
	supplyKeeper     supply.Keeper
	amcModuleName 	 string
	nvrpModuleName   string
}

func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	supplyKeeper supply.Keeper, amcModuleName string, nvrpModuleName string,
) Keeper {

	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}
	if addr := supplyKeeper.GetModuleAddress(amcModuleName); addr == nil {
		panic("the amc module account has not been set")
	}
	if addr := supplyKeeper.GetModuleAddress(nvrpModuleName); addr == nil {
		panic("the nvrp module account has not been set")
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper:     supplyKeeper,
		amcModuleName:    amcModuleName,
		nvrpModuleName:   nvrpModuleName,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshalBinaryBare(b, &minter)
	return
}

func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(minter)
	store.Set(types.MinterKey, b)
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		return nil
	}

	return k.supplyKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

func (k Keeper) TotalSupply(ctx sdk.Context) sdk.Coins {
	return k.supplyKeeper.GetSupply(ctx).GetTotal()
}

func (k Keeper) DistributeCollectedFees(ctx sdk.Context, toMintCoins sdk.Coins) error {
	toAmc := sdk.NewCoins(
		sdk.NewCoin(
			config.DefaultDenom,
			toMintCoins.AmountOf(config.DefaultDenom).QuoRaw(2),
		),
	)
	toNvrp := toMintCoins.Sub(toAmc) // intentional uneven division in favour of nvrp

	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.amcModuleName, toAmc)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("%s -> %s : %s", types.ModuleName, k.amcModuleName, toAmc),
	)

	err = k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.nvrpModuleName, toNvrp)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("%s -> %s : %s", types.ModuleName, k.nvrpModuleName, toNvrp),
	)

	return nil
}
