package hra

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgRegisterName:
			return handleMsgRegisterName(ctx, msg, k)
		case MsgRenewName:
			return handleMsgRenewName(ctx, msg, k)
		case MsgSetPrice:
			return handleMsgSetPrice(ctx, msg, k)
		case MsgDeleteName:
			return handleMsgDeleteName(ctx, msg, k)
		case MsgBuyName:
			return handleMsgBuyName(ctx, msg, k)
		case MsgTransferName:
			return handleMsgTransferName(ctx, msg, k)
		case MsgRegisterAddress:
			return handleMsgRegisterAddress(ctx, msg, k)
		case MsgRemoveAddress:
			return handleMsgRemoveAddress(ctx, msg, k)
		case MsgRemoveAllAddresses:
			return handleMsgRemoveAllAddresses(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func NewGovernanceProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
			case types.RegisterBlockchainIdProposal:
				return handleProposalRegisterBlockchainId(ctx, k, c)

			case types.RemoveBlockchainIdProposal:
				return handleProposalRemoveBlockchainId(ctx, k, c)

			default:
					return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized hra proposal content type: %T", c)
		}
	}
}

func handleProposalRegisterBlockchainId(ctx sdk.Context, k Keeper, proposal types.RegisterBlockchainIdProposal) error {
	k.SetRegisteredBlockchainId(ctx, proposal.BlockchainId)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterBlockchainId,
			sdk.NewAttribute(types.AttributeKeyTitle, proposal.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, proposal.Description),
			sdk.NewAttribute(types.AttributeKeyBlockchainId, proposal.BlockchainId),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}

func handleProposalRemoveBlockchainId(ctx sdk.Context, k Keeper, proposal types.RemoveBlockchainIdProposal) error {
	k.RemoveRegisteredBlockchainId(ctx, proposal.BlockchainId)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveBlockchainId,
			sdk.NewAttribute(types.AttributeKeyTitle, proposal.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, proposal.Description),
			sdk.NewAttribute(types.AttributeKeyBlockchainId, proposal.BlockchainId),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}

func handleMsgRegisterName(ctx sdk.Context, msg MsgRegisterName, k Keeper) (*sdk.Result, error) {
	err := k.HandleRegisterName(ctx, msg.Name, msg.Owner)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRenewName(ctx sdk.Context, msg MsgRenewName, k Keeper) (*sdk.Result, error) {
	err := k.HandleRenewName(ctx, msg.Name, msg.Owner)

	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSetPrice(ctx sdk.Context, msg MsgSetPrice, k Keeper) (*sdk.Result, error) {
	err := k.HandleSetPrice(ctx, msg.Name, msg.Owner, msg.Price)

	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetPrice,
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyPrice, msg.Price.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDeleteName(ctx sdk.Context, msg MsgDeleteName, k Keeper) (*sdk.Result, error) {
	err := k.HandleDeleteName(ctx, msg.Name, msg.Owner)

	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelete,
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgBuyName(ctx sdk.Context, msg MsgBuyName, k Keeper) (*sdk.Result, error) {
	err := k.HandleBuyName(ctx, msg.Name, msg.Buyer)

	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBuy,
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, msg.Buyer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTransferName(ctx sdk.Context, msg MsgTransferName, k Keeper) (*sdk.Result, error) {
	err := k.HandleTransferName(ctx, msg.Name, msg.Owner, msg.NewOwner)

	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferName,
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyNewOwner, msg.NewOwner.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRegisterAddress(ctx sdk.Context, msg MsgRegisterAddress, k Keeper) (*sdk.Result, error) {
	err := k.HandleRegisterAddress(ctx, msg.Owner, msg.BlockchainId, msg.Index, msg.BlockchainAddress)

	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegisterAddress,
			sdk.NewAttribute(types.AttributeKeyBlockchainId, msg.BlockchainId),
			sdk.NewAttribute(types.AttributeKeyIndex, msg.Index),
			sdk.NewAttribute(types.AttributeKeyBlockchainAddress, msg.BlockchainAddress),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRemoveAddress(ctx sdk.Context, msg MsgRemoveAddress, k Keeper) (*sdk.Result, error) {
	err := k.HandleRemoveAddress(ctx, msg.Owner, msg.BlockchainId, msg.Index)

	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveAddress,
			sdk.NewAttribute(types.AttributeKeyBlockchainId, msg.BlockchainId),
			sdk.NewAttribute(types.AttributeKeyIndex, msg.Index),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRemoveAllAddresses(ctx sdk.Context, msg MsgRemoveAllAddresses, k Keeper) (*sdk.Result, error) {
	err := k.HandleRemoveAllAddresses(ctx, msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		sdk.NewAttribute(types.AttributeKeySender, msg.Owner.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}