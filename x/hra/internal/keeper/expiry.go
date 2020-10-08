package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
	"time"
)

func (k Keeper) IterateExpiredNameInfoQueue(ctx sdk.Context, endTime time.Time, cb func(name string) (stop bool)) {
	iterator := k.ExpiredNameInfoQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		name, _ := types.SplitExpiredNameInfoQueueKey(iterator.Key())

		if cb(name) {
			break
		}
	}
}

func (keeper Keeper) InsertExpiredNameInfoQueue(ctx sdk.Context, name string, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := []byte(name)
	store.Set(types.ExpiredNameInfoQueueKey(name, endTime), bz)
}

func (keeper Keeper) RemoveFromExpiredNameInfoQueue(ctx sdk.Context, name string, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.ExpiredNameInfoQueueKey(name, endTime))
}

func (k Keeper) DeleteExpiredNameInfo(ctx sdk.Context, nameInfo types.NameInfo) error {
	k.DeleteNameInfo(ctx, nameInfo.Name)
	k.DeleteNameInfoStatusMap(ctx, nameInfo.Owner, nameInfo.Name)

	// if last HRA remove all associated addresses
	if ! k.OwnsAnyName(ctx, nameInfo.Owner) {
		k.RemoveAllAddresses(ctx, nameInfo.Owner)
		k.SetCredits(ctx, nameInfo.Owner, sdk.ZeroInt())
		err := k.AfterLastNameRemoved(ctx, nameInfo.Owner)
		if err != nil {
			return err
		}
	}

	k.RemoveFromExpiredNameInfoQueue(ctx, nameInfo.Name, nameInfo.ExpiryTime)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeExpiredName,
			sdk.NewAttribute(types.AttributeKeyName, nameInfo.Name),
			sdk.NewAttribute(types.AttributeKeySender, nameInfo.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}

func (keeper Keeper) ExpiredNameInfoQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types.ExpiredNameInfoQueueKeyPrefix, sdk.PrefixEndBytes(types.ExpiredNameInfoByTimeKey(endTime)))
}
