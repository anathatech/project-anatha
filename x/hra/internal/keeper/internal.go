package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
)


func (k Keeper) GetNameInfo(ctx sdk.Context, name string) (types.NameInfo, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetNameInfoByNameKey(name))
	if bz == nil {
		return types.NameInfo{}, false
	}

	var nameInfo types.NameInfo
	k.cdc.MustUnmarshalBinaryBare(bz, &nameInfo)

	return nameInfo, true
}

func (k Keeper) SetNameInfo(ctx sdk.Context, name string, nameInfo types.NameInfo) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetNameInfoByNameKey(name), k.cdc.MustMarshalBinaryBare(nameInfo))
}

func (k Keeper) DeleteNameInfo(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetNameInfoByNameKey(name))
}

func (k Keeper) DeleteNameInfoStatusMap(ctx sdk.Context, owner sdk.AccAddress, name string) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetStatusByAddressAndNameKey(owner, name))
}

func (k Keeper) SetNameInfoStatusMap(ctx sdk.Context, owner sdk.AccAddress, name string) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetStatusByAddressAndNameKey(owner, name), types.StatusPresent)
}
