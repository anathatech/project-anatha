package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type MultiNameHooks []NameHooks

func NewMultiNameHooks(hooks ...NameHooks) MultiNameHooks {
	return hooks
}

func (h MultiNameHooks) AfterFirstNameCreated(ctx sdk.Context, address sdk.AccAddress) error {
	for i := range h {
		err := h[i].AfterFirstNameCreated(ctx, address)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiNameHooks) AfterLastNameRemoved(ctx sdk.Context, address sdk.AccAddress) error {
	for i := range h {
		err := h[i].AfterLastNameRemoved(ctx, address)
		if err != nil {
			return err
		}
	}
	return nil
}
