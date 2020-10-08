package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
	"strings"
)

func (k Keeper) HandleRegisterAddress(ctx sdk.Context, address sdk.AccAddress, blockchainId string, index string, blockchainAddress string) error {
	blockchainAddress = strings.TrimSpace(blockchainAddress)

	if ! k.OwnsAnyName(ctx, address) {
		return types.ErrNoNamesRegistered
	}

	if ! k.IsBlockchainIdRegistered(ctx, blockchainId) {
		return types.ErrBlockchainIdNotValid
	}

	credits := k.GetCredits(ctx, address)

	if credits.LTE(sdk.ZeroInt()) {
		fee := k.AddressRegistrationFee(ctx)
		err := k.SupplyKeeper.SendCoinsFromAccountToModule(
			ctx,
			address,
			k.feeCollectorName,
			fee,
		)
		if err != nil {
			return err
		}
	} else {
		k.SetCredits(ctx, address, credits.Sub(sdk.OneInt()))
	}

	k.SetAddress(ctx, address, blockchainId,index, blockchainAddress)

	return nil
}

func (k Keeper) HandleRemoveAddress(ctx sdk.Context, address sdk.AccAddress, blockchainId string, index string) error {
	k.RemoveAddress(ctx, address, blockchainId, index)

	return nil
}

func (k Keeper) HandleRemoveAllAddresses(ctx sdk.Context, address sdk.AccAddress) error {
	k.RemoveAllAddresses(ctx, address)

	return nil
}

func (k Keeper) SetAddress(ctx sdk.Context, address sdk.AccAddress, blockchainId string, index string, blockchainAddress string) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetAddressKey(address, blockchainId, index), []byte(blockchainAddress))
}

func (k Keeper) RemoveAddress(ctx sdk.Context, address sdk.AccAddress, blockchainId string, index string) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetAddressKey(address, blockchainId, index))
}

func (k Keeper) RemoveAllAddresses(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	iterator := k.GetAddressIterator(ctx, address)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) GetAddress(ctx sdk.Context, address sdk.AccAddress, blockchainId string, index string) (string, error) {
	store := ctx.KVStore(k.storeKey)

	blockchainAddress := store.Get(types.GetAddressKey(address, blockchainId, index))

	if blockchainAddress == nil {
		return "", types.ErrBlockchainAddressNotFound
	}

	return string(blockchainAddress), nil
}

// Returns an iterator through users blockchain addresses
func (k Keeper) GetAddressIterator(ctx sdk.Context, address sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetBlockchainAddressIteratorKey(address))
}

func (k Keeper) IterateBlockchainAddressInfos(ctx sdk.Context, address sdk.AccAddress, cb func(blockchainAddress types.BlockchainAddressInfo) (stop bool)) {
	iterator := k.GetAddressIterator(ctx, address)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		blockchainAddress := types.SplitBlockchainAddressKey(iterator.Key())
		blockchainAddress.BlockchainAddress = string(iterator.Value())

		if cb(blockchainAddress) {
			break
		}

	}
}

func (k Keeper) GetAllAddressesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.AddressKeyPrefix)
}

func (k Keeper) IterateAllBlockchainAddressInfos(ctx sdk.Context, cb func(blockchainAddressRecord types.BlockchainAddressRecordInfo) (stop bool)) {
	iterator := k.GetAllAddressesIterator(ctx)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		blockchainAddressRecord := types.SplitBlockchainAddressRecordKey(iterator.Key())
		blockchainAddressRecord.BlockchainAddressInfo.BlockchainAddress = string(iterator.Value())

		if cb(blockchainAddressRecord) {
			break
		}

	}
}

func (k Keeper) GetNamesByAddressIterator(ctx sdk.Context, address sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)

	return sdk.KVStorePrefixIterator(store, types.GetStatusByAddressAndNameIteratorKey(address))
}

func (k Keeper) GetNamesByAddressCount(ctx sdk.Context, address sdk.AccAddress) int {
	iterator := k.GetNamesByAddressIterator(ctx, address)
	defer iterator.Close()

	count := 0

	for ; iterator.Valid(); iterator.Next() {
		count++
	}

	return count
}

func (k Keeper) GetCredits(ctx sdk.Context, address sdk.AccAddress) sdk.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetCreditsKey(address))
	if bz == nil {
		return sdk.ZeroInt()
	}

	var credits sdk.Int
	k.cdc.MustUnmarshalBinaryBare(bz, &credits)

	return credits
}

func (k Keeper) SetCredits(ctx sdk.Context, address sdk.AccAddress, credits sdk.Int) {
	store := ctx.KVStore(k.storeKey)

	if credits.IsZero(){
		store.Delete(types.GetCreditsKey(address))
	} else {
		store.Set(types.GetCreditsKey(address), k.cdc.MustMarshalBinaryBare(credits))
	}
}

func (k Keeper) GetAddressCreditsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.CreditsKeyPrefix)
}

func (k Keeper) IterateAddressCredits(ctx sdk.Context, cb func(address sdk.AccAddress, credits sdk.Int) (stop bool)) {
	iterator := k.GetAddressCreditsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		address := sdk.AccAddress(iterator.Key()[1:])
		var credits sdk.Int

		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &credits)

		if cb(address, credits) {
			break
		}
	}
}

func (k Keeper) SetRegisteredBlockchainId(ctx sdk.Context, blockchainId string) {
	store := ctx.KVStore(k.storeKey)

	blockchainId = strings.ToLower(blockchainId)

	store.Set(types.GetRegisteredBlockchainIdKey(blockchainId), types.StatusPresent)
}

func (k Keeper) RemoveRegisteredBlockchainId(ctx sdk.Context, blockchainId string) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetRegisteredBlockchainIdKey(blockchainId))
}

func (k Keeper) IsBlockchainIdRegistered(ctx sdk.Context, blockchainId string) bool {
	store := ctx.KVStore(k.storeKey)

	return store.Has(types.GetRegisteredBlockchainIdKey(blockchainId))
}

func (k Keeper) GetRegisteredBlockchainIdIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetRegisteredBlockchainIdIteratorKey())
}

func (k Keeper) IterateRegisteredBlockchainIds(ctx sdk.Context, cb func(blockchainId string) (stop bool)) {
 	iterator := k.GetRegisteredBlockchainIdIterator(ctx)
 	defer iterator.Close()

 	for ; iterator.Valid(); iterator.Next() {
 		blockchainId := types.SplitRegisteredBlockchainIdKey(iterator.Key())

 		if cb(blockchainId) {
 			break
		}
	}
 }

func (k Keeper) GetRegisteredBlockchainIds(ctx sdk.Context) []string {
	var registeredBlockchainIds []string
	k.IterateRegisteredBlockchainIds(ctx, func(blockchainId string) (stop bool) {
		registeredBlockchainIds = append(registeredBlockchainIds, blockchainId)
		return false
	})

	return registeredBlockchainIds
}