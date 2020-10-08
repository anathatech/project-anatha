package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
)

func (k Keeper) IsManager(ctx sdk.Context, address sdk.AccAddress) bool {
	managers := k.Managers(ctx)

	for _, manager := range managers {
		if manager.Equals(address) {
			return true
		}
	}

	return false
}

func (k Keeper) HandleAddOperator(ctx sdk.Context, sender sdk.AccAddress, operator sdk.AccAddress) error {
	if ! k.IsManager(ctx, sender) {
		return types.ErrNotManager
	}

	k.AddOperator(ctx, operator)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddOperator,
			sdk.NewAttribute(types.AttributeKeyOperator, operator.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
		),
	})

	return nil
}

func (k Keeper) HandleRemoveOperator(ctx sdk.Context, sender sdk.AccAddress, operator sdk.AccAddress) error {
	if ! k.IsManager(ctx, sender) {
		return types.ErrNotManager
	}

	k.RemoveOperator(ctx, operator)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveOperator,
			sdk.NewAttribute(types.AttributeKeyOperator, operator.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
		),
	})

	return nil
}

func (k Keeper) AddOperator(ctx sdk.Context, address sdk.AccAddress) {
	account := k.AccountKeeper.GetAccount(ctx, address)
	if account == nil {
		account = k.AccountKeeper.NewAccountWithAddress(ctx, address)
		k.AccountKeeper.SetAccount(ctx, account)
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOperatorKey(address), types.StatusPresent)
}

func (k Keeper) RemoveOperator(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetOperatorKey(address))
}

func (k Keeper) IsOperator(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)

	return store.Has(types.GetOperatorKey(address))
}

func (k Keeper) GetOperatorIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetOperatorIteratorKey())
}

func (k Keeper) IterateOperators(ctx sdk.Context, cb func(address sdk.AccAddress) (stop bool)) {
	iterator := k.GetOperatorIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		address := types.SplitOperatorKey(iterator.Key())

		if cb(address) {
			break
		}
	}
}

func (k Keeper) GetOperators(ctx sdk.Context) []sdk.AccAddress {
	var operators []sdk.AccAddress
	k.IterateOperators(ctx, func(address sdk.AccAddress) (stop bool) {
		operators = append(operators, address)
		return false
	})

	return operators
}