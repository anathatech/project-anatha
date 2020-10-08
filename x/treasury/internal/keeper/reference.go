package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
	"strings"
)

func (k Keeper) SetDisbursementReferenceAmount(ctx sdk.Context, reference string, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)

	reference = strings.ToLower(reference)

	store.Set(types.GetDisbursementReferenceKey(reference), k.cdc.MustMarshalBinaryBare(amount))
}

func (k Keeper) RemoveDisbursementReferenceAmount(ctx sdk.Context, reference string) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetDisbursementReferenceKey(reference))
}

func (k Keeper) GetDisbursementReferenceAmount(ctx sdk.Context, reference string) (bool, sdk.Int) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetDisbursementReferenceKey(reference))
	if bz == nil {
		return false, sdk.ZeroInt()
	}

	var amount sdk.Int
	k.cdc.MustUnmarshalBinaryBare(bz, &amount)

	return true, amount
}

func (k Keeper) IsDisbursementReferenceSet(ctx sdk.Context, reference string) bool {
	store := ctx.KVStore(k.storeKey)

	return store.Has(types.GetDisbursementReferenceKey(reference))
}

func (k Keeper) GetDisbursementReferenceIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetDisbursementReferenceIteratorKey())
}

func (k Keeper) IterateDisbursementReferences(ctx sdk.Context, cb func(reference string, amount sdk.Int) (stop bool)) {
	iterator := k.GetDisbursementReferenceIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		reference := types.SplitDisbursementReferenceKey(iterator.Key())

		var amount sdk.Int
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &amount)

		if cb(reference, amount) {
			break
		}
	}
}
