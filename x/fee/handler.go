package fee

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/anathatech/project-anatha/x/fee/internal/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {

			default:
				errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName,  msg)
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func NewGovernanceProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
			case types.AddFeeExcludedMessageProposal:
				return handleProposalAddFeeExcludedMessage(ctx, k, c)

			case types.RemoveFeeExcludedMessageProposal:
				return handleProposalRemoveFeeExcludedMessage(ctx, k, c)

			default:
				return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized hra proposal content type: %T", c)
		}
	}
}

func handleProposalAddFeeExcludedMessage(ctx sdk.Context, k Keeper, proposal types.AddFeeExcludedMessageProposal) error {
	k.SetFeeExcludedMessage(ctx, proposal.MessageType)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddFeeExcludedMessage,
			sdk.NewAttribute(types.AttributeKeyTitle, proposal.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, proposal.Description),
			sdk.NewAttribute(types.AttributeKeyMessageType, proposal.MessageType),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}

func handleProposalRemoveFeeExcludedMessage(ctx sdk.Context, k Keeper, proposal types.RemoveFeeExcludedMessageProposal) error {
	k.RemoveFeeExcludedMessage(ctx, proposal.MessageType)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveFeeExcludedMessage,
			sdk.NewAttribute(types.AttributeKeyTitle, proposal.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, proposal.Description),
			sdk.NewAttribute(types.AttributeKeyMessageType, proposal.MessageType),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}