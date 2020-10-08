package keeper

import (
	"github.com/anathatech/project-anatha/config"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
)

func (k Keeper) HandleDisburse(ctx sdk.Context, operator sdk.AccAddress, recipient sdk.AccAddress, dinAmount sdk.Coins, reference string) error {
	if ! k.IsOperator(ctx, operator) {
		return types.ErrNotOperator
	}

	if k.IsDisbursementReferenceSet(ctx, reference) {
		return types.ErrDuplicateReference
	}

	scheduledFor := ctx.BlockTime()

	if dinAmount.IsAnyGTE(k.RiskAssessmentAmount(ctx)) {
		scheduledFor = scheduledFor.Add(k.RiskAssessmentDuration(ctx))
	}

	for k.HasDisbursementInQueue(ctx, recipient, scheduledFor) {
		scheduledFor = scheduledFor.Add(time.Millisecond)
	}

	k.InsertDisbursementQueue(ctx, types.NewDisbursement(
		operator,
		recipient,
		dinAmount,
		scheduledFor,
		reference,
	))

	k.SetDisbursementReferenceAmount(ctx, reference, sdk.ZeroInt())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDisburse,
			sdk.NewAttribute(types.AttributeKeyScheduledFor, scheduledFor.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, dinAmount.String()),
			sdk.NewAttribute(types.AttributeKeyReference, reference),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, operator.String()),
		),
	})

	return nil
}

func (k Keeper) HandleDisburseToEscrow(ctx sdk.Context, operator sdk.AccAddress, dinAmount sdk.Coins, reference string) error {
	if ! k.IsOperator(ctx, operator) {
		return types.ErrNotOperator
	}

	if k.IsDisbursementReferenceSet(ctx, reference) {
		return types.ErrDuplicateReference
	}

	if dinAmount.IsAnyGTE(k.RiskAssessmentAmount(ctx)) {
		return types.ErrEscrowDistributionAmountExceeded
	}

	totalPinAmount, fromBuyBack, fromTreasury := k.CalculatePinAmountExtended(ctx, dinAmount)

	err := k.DisburseFundsToEscrow(ctx, reference, dinAmount, fromBuyBack, fromTreasury)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDisburseToEscrow,
			sdk.NewAttribute(types.AttributeKeyPinAmount, totalPinAmount.String()),
			sdk.NewAttribute(types.AttributeKeyDinAmount, dinAmount.String()),
			sdk.NewAttribute(types.AttributeKeyReference, reference),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, operator.String()),
		),
	})

	return nil
}

func (k Keeper) HandleDisburseFromEscrow(ctx sdk.Context, operator sdk.AccAddress, reference string, recipient sdk.AccAddress) error {
	if ! k.IsOperator(ctx, operator) {
		return types.ErrNotOperator
	}

	found, amount := k.GetDisbursementReferenceAmount(ctx, reference)

	if ! found {
		return types.ErrInvalidReference
	}

	if amount.IsZero() {
		return types.ErrEscrowDisbursed
	}

	err := k.DisburseFundsFromEscrow(ctx, reference, amount, recipient)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDisburseFromEscrow,
			sdk.NewAttribute(types.AttributeKeyAmount, sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, amount)).String()),
			sdk.NewAttribute(types.AttributeKeyReference, reference),
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, operator.String()),
		),
	})

	return nil
}

func (k Keeper) HandleRevertFromEscrow(ctx sdk.Context, operator sdk.AccAddress, amount sdk.Coins, reference string) error {
	if ! k.IsOperator(ctx, operator) {
		return types.ErrNotOperator
	}

	found, escrowBalance := k.GetDisbursementReferenceAmount(ctx, reference)

	if ! found {
		return types.ErrInvalidReference
	}

	if escrowBalance.IsZero() {
		return types.ErrEscrowDisbursed
	}

	amountInt := amount.AmountOf(config.DefaultDenom)

	if amountInt.GT(escrowBalance) {
		return types.ErrEscrowRevertAmountTooBig
	}

	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.TreasuryEscrowModuleName, types.BuyBackFundModuleName, amount)
	if err != nil {
		return err
	}

	escrowRemainder := escrowBalance.Sub(amountInt)

	k.SetDisbursementReferenceAmount(ctx, reference, escrowRemainder)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRevertFromEscrow,
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyReference, reference),
			sdk.NewAttribute(types.AttributeKeyEscrowRemainder, escrowRemainder.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, operator.String()),
		),
	})

	return nil
}

func (k Keeper) HandleCancelDisbursement(ctx sdk.Context, manager sdk.AccAddress, recipient sdk.AccAddress, scheduledFor time.Time) error {
	if ! k.IsManager(ctx, manager) {
		return types.ErrNotManager
	}

	if ! k.HasDisbursementInQueue(ctx, recipient, scheduledFor) {
		return types.ErrDisbursementNotScheduled
	}

	k.RemoveFromDisbursementQueue(ctx, recipient, scheduledFor)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelDisbursement,
			sdk.NewAttribute(types.AttributeKeyScheduledFor, scheduledFor.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, manager.String()),
		),
	})

	return nil
}

func (k Keeper) DisburseFunds(ctx sdk.Context, operator sdk.AccAddress, recipient sdk.AccAddress, dinAmount sdk.Coins, fromBuyBack sdk.Coins, fromTreasury sdk.Coins) error {
	if ! operator.Empty() && ! k.IsOperator(ctx, operator) {
		return types.ErrNotOperator
	}

	if ! fromBuyBack.IsZero() {
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.BuyBackFundModuleName, recipient, fromBuyBack)
		if err != nil {
			return err
		}
	}

	if ! fromTreasury.IsZero() {
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, fromTreasury)
		if err != nil {
			return err
		}
	}

	treasury := k.GetTreasury(ctx)
	treasury.Distributed = treasury.Distributed.Add(fromTreasury...)
	treasury.InboundDin = treasury.InboundDin.Add(dinAmount...)
	k.SetTreasury(ctx, treasury)

	return nil
}

func (k Keeper) DisburseFundsToEscrow(ctx sdk.Context, reference string, dinAmount sdk.Coins, fromBuyBack sdk.Coins, fromTreasury sdk.Coins) error {
	if ! fromBuyBack.IsZero() {
		err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.BuyBackFundModuleName, types.TreasuryEscrowModuleName, fromBuyBack)
		if err != nil {
			return err
		}
	}

	if ! fromTreasury.IsZero() {
		err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.TreasuryEscrowModuleName, fromTreasury)
		if err != nil {
			return err
		}
	}

	treasury := k.GetTreasury(ctx)
	treasury.Distributed = treasury.Distributed.Add(fromTreasury...)
	treasury.InboundDin = treasury.InboundDin.Add(dinAmount...)
	k.SetTreasury(ctx, treasury)

	k.SetDisbursementReferenceAmount(ctx, reference, fromBuyBack.Add(fromTreasury...).AmountOf(config.DefaultDenom))

	return nil
}

func (k Keeper) DisburseFundsFromEscrow(ctx sdk.Context, reference string, amount sdk.Int, recipient sdk.AccAddress) error {
	coins := sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, amount))

	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.TreasuryEscrowModuleName, recipient, coins)
	if err != nil {
		return err
	}

	k.SetDisbursementReferenceAmount(ctx, reference, sdk.ZeroInt())

	return nil
}

func (k Keeper) IterateScheduledDisbursementQueue(ctx sdk.Context, endTime time.Time, cb func(disbursement types.Disbursement) (stop bool)) {
	iterator := k.ScheduledDisbursementQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		var disbursement types.Disbursement
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &disbursement)

		if cb(disbursement) {
			break
		}
	}
}

func (keeper Keeper) InsertDisbursementQueue(ctx sdk.Context, disbursement types.Disbursement) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryBare(disbursement)
	store.Set(types.DisbursementQueueKey(disbursement.Recipient, disbursement.ScheduledFor), bz)
}

func (keeper Keeper) RemoveFromDisbursementQueue(ctx sdk.Context, address sdk.AccAddress, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.DisbursementQueueKey(address, endTime))
}

func (keeper Keeper) HasDisbursementInQueue(ctx sdk.Context, recipient sdk.AccAddress, scheduledFor time.Time) bool {
	store := ctx.KVStore(keeper.storeKey)
	return store.Has(types.DisbursementQueueKey(recipient, scheduledFor))
}

func (keeper Keeper) ScheduledDisbursementQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types.DisbursementQueueKeyPrefix, sdk.PrefixEndBytes(types.DisbursementByTimeKey(endTime)))
}

func (k Keeper) IterateDisbursementQueue(ctx sdk.Context, cb func(disbursement types.Disbursement) (stop bool)) {
	iterator := k.DisbursementQueueIterator(ctx)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		var disbursement types.Disbursement
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &disbursement)

		if cb(disbursement) {
			break
		}
	}
}

func (keeper Keeper) DisbursementQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, types.DisbursementQueueKeyPrefix)

}

func (k Keeper) GetDisbursements(ctx sdk.Context) []types.Disbursement {
	var disbursements []types.Disbursement
	k.IterateDisbursementQueue(ctx, func(disbursement types.Disbursement) (stop bool) {
		disbursements = append(disbursements, disbursement)
		return false
	})

	return disbursements
}
