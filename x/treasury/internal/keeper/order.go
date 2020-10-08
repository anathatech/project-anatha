package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
)

func (k Keeper) HandleCreateSellOrder(ctx sdk.Context, sender sdk.AccAddress, pinAmount sdk.Coins) error {
	buyBackPercentage := k.BuyBackPercentage(ctx)
	stagePrice := k.GetPriceForStage(
		ctx,
		k.GetStageFromDistribution(
			ctx,
			k.DistributedFromTreasury(ctx),
		),
	)

	dinAmountInt := pinAmount.AmountOf(config.DefaultDenom).Mul(stagePrice)
	dinAmountInt = dinAmountInt.ToDec().Mul(buyBackPercentage).TruncateInt()

	dinAmount := sdk.NewCoins(sdk.NewCoin(config.DefaultStableDenom, dinAmountInt))

	if ! dinAmount.IsZero() {
		err := k.TransferToBuyBackFund(ctx, sender, pinAmount)
		if err != nil {
			return err
		}

		err = k.TransferFromBuyBackFund(ctx, sender, dinAmount)
		if err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateSellOrder,
			sdk.NewAttribute(types.AttributeKeyPinAmount, pinAmount.String()),
			sdk.NewAttribute(types.AttributeKeyDinAmount, dinAmount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
		),
	})

	return nil
}

func (k Keeper) HandleCreateBuyOrder(ctx sdk.Context, buyer sdk.AccAddress, dinAmount sdk.Coins) error {
	if ! k.BankKeeper.HasCoins(ctx, buyer, dinAmount) {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Insufficient funds for ANATHA purchase.")
	}

	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, buyer, types.DistributionProfitsModuleName, dinAmount)
	if err != nil {
		return err
	}

	_, fromBuyBack, fromTreasury := k.CalculatePinAmountExtended(ctx, dinAmount)

	err = k.DisburseFunds(ctx, nil, buyer, dinAmount, fromBuyBack, fromTreasury)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateBuyOrder,
			sdk.NewAttribute(types.AttributeKeyPinAmount, fromBuyBack.Add(fromTreasury...).String()),
			sdk.NewAttribute(types.AttributeKeyDinAmount, dinAmount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, buyer.String()),
		),
	})

	return nil
}
