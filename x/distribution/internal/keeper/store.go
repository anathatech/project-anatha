package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
)

// get accumulated rewards for a validator
func (k Keeper) GetValidatorAccumulatedRewards(ctx sdk.Context, val sdk.ValAddress) (rewards sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetValidatorAccumulatedRewardsKey(val))
	if b == nil {
		return sdk.NewDecCoins()
	}
	k.cdc.MustUnmarshalBinaryBare(b, &rewards)
	return
}

// set accumulated rewards for a validator
func (k Keeper) SetValidatorAccumulatedRewards(ctx sdk.Context, val sdk.ValAddress, rewards sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)

	if ! rewards.IsZero() {
		store.Set(types.GetValidatorAccumulatedRewardsKey(val), k.cdc.MustMarshalBinaryBare(rewards))
	} else {
		store.Delete(types.GetValidatorAccumulatedRewardsKey(val))
	}
}

// delete accumulated rewards for a validator
func (k Keeper) DeleteValidatorAccumulatedRewards(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorAccumulatedRewardsKey(val))
}

// iterate over validator accumulated rewards
func (k Keeper) IterateValidatorAccumulatedRewards(ctx sdk.Context, handler func(val sdk.ValAddress, rewards sdk.DecCoins) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorAccumulatedRewardsKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards sdk.DecCoins
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &rewards)
		addr := types.GetValidatorAccumulatedRewardsAddress(iter.Key())
		if handler(addr, rewards) {
			break
		}
	}
}

func (k Keeper) GetNvrpRemainder(ctx sdk.Context) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetNvrpdRemainderKey())
	if bz == nil {
		return sdk.NewDecCoins()
	}

	var remainder sdk.DecCoins
	k.cdc.MustUnmarshalBinaryBare(bz, &remainder)

	return remainder
}

func (k Keeper) SetNvrpRemainder(ctx sdk.Context, remainder sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)

	if ! remainder.IsZero() {
		store.Set(types.GetNvrpdRemainderKey(), k.cdc.MustMarshalBinaryBare(remainder))
	} else {
		store.Delete(types.GetNvrpdRemainderKey())
	}
}
