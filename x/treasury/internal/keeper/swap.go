package keeper

import (
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleSwap(ctx sdk.Context, operator sdk.AccAddress, recipient sdk.AccAddress, amount sdk.Coins, reference string) error {
	if ! k.IsOperator(ctx, operator) {
		return types.ErrNotOperator
	}

	if k.IsDisbursementReferenceSet(ctx, reference) {
		return types.ErrDuplicateReference
	}

	err := k.TransferFromSwapEscrow(ctx, recipient, amount)

	if err != nil {
		return err
	}

	k.SetDisbursementReferenceAmount(ctx, reference, sdk.ZeroInt())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
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
