package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/hra"
)

type NameHooks struct {
	k Keeper
}

var _ hra.NameHooks = NameHooks{}

func (k Keeper) NameHooks() NameHooks { return NameHooks{k} }

func (h NameHooks) AfterFirstNameCreated(ctx sdk.Context, address sdk.AccAddress) error {
	// sing up for HRA rewards

	h.k.InsertNameDepositQueue(
		ctx,
		address,
		ctx.BlockTime().Add(h.k.NameDepositDelay(ctx)),
	)

	return nil
}

func (h NameHooks) AfterLastNameRemoved(ctx sdk.Context, address sdk.AccAddress) error {
	// withdraw reward and stake
	reward := h.k.withdrawNameReward(ctx, address)

	if ! reward.IsZero() {
		err := h.k.DistributeNameReward(ctx, address, reward)
		if err != nil {
			return err
		}
	}

	return nil
}