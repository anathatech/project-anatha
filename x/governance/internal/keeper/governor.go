package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/governance/internal/types"
)

func (k Keeper) HandleRemoveGovernor(ctx sdk.Context, address sdk.AccAddress) error {
	k.RemoveGovernor(ctx, address)

	return nil
}

func (k Keeper) HandleAddGovernor(ctx sdk.Context, address sdk.AccAddress) error {
	k.AddGovernor(ctx, address)

	return nil
}

func (k Keeper) AddGovernor(ctx sdk.Context, address sdk.AccAddress) {
	account := k.AccountKeeper.GetAccount(ctx, address)
	if account == nil {
		account = k.AccountKeeper.NewAccountWithAddress(ctx, address)
		k.AccountKeeper.SetAccount(ctx, account)
	}

	store := ctx.KVStore(k.storeKey)

	if ! store.Has(types.GetGovernorKey(address)) {
		store.Set(types.GetGovernorKey(address), types.StatusPresent)

		count := k.GetGovernorCount(ctx)
		k.SetGovernorCount(ctx, count.Add(sdk.OneInt()))
	}
}

func (k Keeper) RemoveGovernor(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	if store.Has(types.GetGovernorKey(address)) {
		store.Delete(types.GetGovernorKey(address))

		count := k.GetGovernorCount(ctx)

		k.SetGovernorCount(ctx, count.Sub(sdk.OneInt()))
	}
}

func (k Keeper) IsGovernor(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)

	return store.Has(types.GetGovernorKey(address))
}

func (k Keeper) GetGovernorIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetGovernorIteratorKey())
}

func (k Keeper) IterateGovernors(ctx sdk.Context, cb func(address sdk.AccAddress) (stop bool)) {
	iterator := k.GetGovernorIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		address := types.SplitGovernorKey(iterator.Key())

		if cb(address) {
			break
		}
	}
}

func (k Keeper) GetGovernors(ctx sdk.Context) []sdk.AccAddress {
	var governors []sdk.AccAddress
	k.IterateGovernors(ctx, func(address sdk.AccAddress) (stop bool) {
		governors = append(governors, address)
		return false
	})

	return governors
}

func (k Keeper) GetGovernorCount(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GovernorCountKey)
	if bz == nil {
		return sdk.ZeroInt()
	}

	var count sdk.Int
	k.cdc.MustUnmarshalBinaryBare(bz, &count)

	return count
}

func (k Keeper) SetGovernorCount(ctx sdk.Context, count sdk.Int) {
	store := ctx.KVStore(k.storeKey)

	if count.IsZero() {
		store.Delete(types.GovernorCountKey)
	} else {
		store.Set(types.GovernorCountKey, k.cdc.MustMarshalBinaryBare(count))
	}
}
