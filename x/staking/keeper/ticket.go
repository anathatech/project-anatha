package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/staking/types"
)

func (k Keeper) GetLastTicket(ctx sdk.Context) (ticket uint64) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.TicketKey)
	if bz == nil {
		panic("Ticket should have been set")
	}

	ticket = types.GetTicketFromBytes(bz)

	return ticket
}

func (k Keeper) SetLastTicket(ctx sdk.Context, ticket uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TicketKey, types.GetTicketBytes(ticket))
}


// validator index
func (k Keeper) SetValidatorByTicket(ctx sdk.Context, validator types.Validator) {
	// jailed validators are not kept in the power index
	if validator.Jailed {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorsByTicketKey(validator), validator.OperatorAddress)

	k.Logger(ctx).Debug(
		fmt.Sprintf("Assigned ticket %d to %s", validator.Ticket, sdk.AccAddress(validator.OperatorAddress)),
	)
}

// validator index
func (k Keeper) DeleteValidatorByTicket(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorsByTicketKey(validator))

	k.Logger(ctx).Debug(
		fmt.Sprintf("Removed ticket %d from %s", validator.Ticket, sdk.AccAddress(validator.OperatorAddress)),
	)
}

// validator index
func (k Keeper) SetNewValidatorByTicket(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorsByTicketKey(validator), validator.OperatorAddress)

	k.Logger(ctx).Debug(
		fmt.Sprintf("Assigned ticket %d to %s", validator.Ticket, sdk.AccAddress(validator.OperatorAddress)),
	)
}

func (k Keeper) DumpTickets(ctx sdk.Context) {
	iterator := k.ValidatorTicketStoreIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		ticket, operator := types.ParseValidatorTicketKey(iterator.Key())

		validator, found := k.GetValidator(ctx, sdk.ValAddress(operator))
		if ! found {
			continue
		}

		k.Logger(ctx).Debug(
			fmt.Sprintf("%d (%s) %s", ticket, validator.GetMoniker() ,sdk.AccAddress(operator)),
		)
	}
}