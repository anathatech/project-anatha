package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/anathatech/project-anatha/x/fee/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace params.Subspace
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramspace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Storage

func (k Keeper) SetFeeExcludedMessage(ctx sdk.Context, msgType string) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetFeeExcludedMessageKey(msgType), types.StatusPresent)
}

func (k Keeper) RemoveFeeExcludedMessage(ctx sdk.Context, msgType string) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetFeeExcludedMessageKey(msgType))
}

func (k Keeper) IsMessageFeeExcluded(ctx sdk.Context, msg sdk.Msg) bool {
	msgType := msg.Route() + "/" + msg.Type()

	store := ctx.KVStore(k.storeKey)

	return store.Has(types.GetFeeExcludedMessageKey(msgType))
}

func (k Keeper) GetFeeExcludedMessageIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetFeeExcludedMessageIteratorKey())
}

func (k Keeper) IterateFeeExcludedMessages(ctx sdk.Context, cb func(msgType string) (stop bool)) {
	iterator := k.GetFeeExcludedMessageIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msgType := types.SplitFeeExcludedMessageKey(iterator.Key())

		if cb(msgType) {
			break
		}
	}
}

func (k Keeper) GetFeeExcludedMessages(ctx sdk.Context) []string {
	var messageTypes []string
	k.IterateFeeExcludedMessages(ctx, func(msgType string) (stop bool) {
		messageTypes = append(messageTypes, msgType)
		return false
	})

	return messageTypes
}