package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
	"time"
)

func (k Keeper) IterateNameDepositQueueByTime(ctx sdk.Context, endTime time.Time, cb func(address sdk.AccAddress, endTime time.Time) (stop bool)) {
	iterator := k.ScheduledNameDepositQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		address, endTime := types.SplitNameDepositQueueKey(iterator.Key())

		if cb(address, endTime) {
			break
		}
	}
}

func (k Keeper) IterateNameDepositQueue(ctx sdk.Context, handler func(address sdk.AccAddress, endTime time.Time) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.NameDepositQueueKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		address, endTime := types.SplitNameDepositQueueKey(iter.Key())

		if handler(address, endTime) {
			break
		}
	}
}

func (k Keeper) InsertNameDepositQueue(ctx sdk.Context, recipient sdk.AccAddress, scheduledFor time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(recipient)
	store.Set(types.NameDepositQueueKey(recipient, scheduledFor), bz)

	k.Logger(ctx).Debug(
		fmt.Sprintf("Scheduled Name Deposit: %s %s", recipient, scheduledFor),
	)
}

func (k Keeper) RemoveFromNameDepositQueue(ctx sdk.Context, address sdk.AccAddress, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.NameDepositQueueKey(address, endTime))

	k.Logger(ctx).Debug(
		fmt.Sprintf("Removed from Name Deposit Queue: %s %s", address, endTime),
	)
}

func (k Keeper) ScheduledNameDepositQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.NameDepositQueueKeyPrefix, sdk.PrefixEndBytes(types.NameDepositByTimeKey(endTime)))
}
