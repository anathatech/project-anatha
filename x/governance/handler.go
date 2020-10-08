package governance

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/anathatech/project-anatha/x/governance/internal/types"
	"strconv"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
			case MsgSubmitProposal:
				return handleMsgSubmitProposal(ctx, keeper, msg)

			case MsgVote:
				return handleMsgVote(ctx, keeper, msg)

			case MsgExpedite:
				return handleMsgExpedite(ctx, keeper, msg)

			default:
				return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func NewGovernanceProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
			case TextProposal:
				return nil

			case types.AddGovernorProposal:
				return handleProposalAddGovernor(ctx, k, c)

			case types.RemoveGovernorProposal:
				return handleProposalRemoveGovernor(ctx, k, c)

			default:
				return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized governance proposal content type: %T", c)
		}
	}
}

func handleProposalAddGovernor(ctx sdk.Context, k Keeper, c types.AddGovernorProposal) error {
	err := k.HandleAddGovernor(ctx, c.Governor)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddGovernor,
			sdk.NewAttribute(types.AttributeKeyGovernor, c.Governor.String()),
			sdk.NewAttribute(types.AttributeKeyTitle, c.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, c.Description),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}

func handleProposalRemoveGovernor(ctx sdk.Context, k Keeper, c types.RemoveGovernorProposal) error {
	err := k.HandleRemoveGovernor(ctx, c.Governor)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveGovernor,
			sdk.NewAttribute(types.AttributeKeyGovernor, c.Governor.String()),
			sdk.NewAttribute(types.AttributeKeyTitle, c.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, c.Description),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitProposal) (*sdk.Result, error) {
	err := keeper.SubmitProposal(ctx, msg.Proposer, msg.Content)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgVote(ctx sdk.Context, keeper Keeper, msg MsgVote) (*sdk.Result, error) {
	err := keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgExpedite(ctx sdk.Context, keeper Keeper, msg MsgExpedite) (*sdk.Result, error) {
	err := keeper.HandleExpedite(ctx, msg.ProposalID)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeExpedite,
			sdk.NewAttribute(types.AttributeKeyProposalId, strconv.Itoa(int(msg.ProposalID))),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}