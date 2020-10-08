package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
)

func (k Keeper) HandleSetPrice(ctx sdk.Context, name string, owner sdk.AccAddress, price sdk.Coins) error {
	nameInfo, found := k.GetNameInfo(ctx, name)
	if ! found {
		return types.ErrNameNotRegistered
	}

	if ! owner.Equals(nameInfo.Owner) {
		return types.ErrNotOwner
	}

	nameInfo.Price = price

	k.SetNameInfo(ctx, name, nameInfo)

	return nil
}

func (k Keeper) GetPrice(ctx sdk.Context, name string) (sdk.Coins, error) {
	nameInfo, found := k.GetNameInfo(ctx, name)
	if ! found {
		return nil, types.ErrNameNotRegistered
	}

	return nameInfo.Price, nil
}

func (k Keeper) HandleBuyName(ctx sdk.Context, name string, buyer sdk.AccAddress) error {
	nameInfo, found := k.GetNameInfo(ctx, name)
	if ! found {
		return types.ErrNameNotRegistered
	}

	if nameInfo.Owner.Equals(buyer) {
		return types.ErrAlreadyOwned
	}

	if nameInfo.Price.IsZero() {
		return types.ErrNotForSale
	}

	coins := nameInfo.Price

	if ! k.CoinKeeper.HasCoins(ctx, buyer, coins) {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Insufficient funds for HRA purchase.")
	}

	err := k.CoinKeeper.SendCoins(ctx, buyer, nameInfo.Owner, coins)
	if err != nil {
		return err
	}

	if ! k.OwnsAnyName(ctx, buyer) {
		k.SetCredits(ctx, buyer, k.AddressCredits(ctx))
		k.AfterFirstNameCreated(ctx, buyer)
	}

	// update the status mapping
	k.DeleteNameInfoStatusMap(ctx, nameInfo.Owner, name)
	k.SetNameInfoStatusMap(ctx,buyer, name)

	oldOwner := nameInfo.Owner

	// update the owner and reset the price
	nameInfo.Owner = buyer
	nameInfo.Price = sdk.NewCoins(sdk.NewInt64Coin(config.DefaultDenom, 0))

	k.SetNameInfo(ctx, name, nameInfo)

	if ! k.OwnsAnyName(ctx, oldOwner) {
		k.RemoveAllAddresses(ctx, oldOwner)

		k.SetCredits(ctx, oldOwner, sdk.ZeroInt())
		err = k.AfterLastNameRemoved(ctx, oldOwner)
		if err != nil {
			return err
		}
	}

	return nil
}