package treasury

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
	"time"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {

		case MsgAddOperator:
			return handleMsgAddOperator(ctx, k, msg)

		case MsgRemoveOperator:
			return handleMsgRemoveOperator(ctx, k, msg)

		case MsgDisburse:
			return handleMsgDisburse(ctx, k, msg)

		case MsgDisburseToEscrow:
			return handleMsgDisburseToEscrow(ctx, k, msg)

		case MsgDisburseFromEscrow:
			return handleMsgDisburseFromEscrow(ctx, k, msg)

		case MsgRevertFromEscrow:
			return handleMsgRevertFromEscrow(ctx, k, msg)

		case MsgCancelDisbursement:
			return handleMsgCancelDisbursement(ctx, k, msg)

		case MsgCreateSellOrder:
			return handleMsgCreateSellOrder(ctx, k, msg)

		case MsgCreateBuyOrder:
			return handleMsgCreateBuyOrder(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName,  msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func NewTreasuryProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
			case AddBuyBackLiquidityProposal:
				return handleAddBuyBackLiquidityProposal(ctx, k, c)

			case RemoveBuyBackLiquidityProposal:
				return handleRemoveBuyBackLiquidityProposal(ctx, k, c)

			case BurnDistributionProfitsProposal:
				return handleBurnDistributionProfitsProposal(ctx, k, c)

			case TransferFromDistributionProfitsToBuyBackLiquidityProposal:
				return handleTransferFromDistributionProfitsToBuyBackLiquidityProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized treasury proposal content type: %T", c)
		}
	}
}

func handleAddBuyBackLiquidityProposal(ctx sdk.Context, k Keeper, p AddBuyBackLiquidityProposal) error {
	err := k.MintCoinsToBuyBackLiquidityFund(ctx, p.Amount)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddBuyBackLiquidity,
			sdk.NewAttribute(types.AttributeKeyTitle, p.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, p.Description),
			sdk.NewAttribute(types.AttributeKeyAmount, p.Amount.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}

func handleRemoveBuyBackLiquidityProposal(ctx sdk.Context, k Keeper, p RemoveBuyBackLiquidityProposal) error {
	err := k.BurnCoinsFromBuyBackLiquidityFund(ctx, p.Amount)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveBuyBackLiquidity,
			sdk.NewAttribute(types.AttributeKeyTitle, p.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, p.Description),
			sdk.NewAttribute(types.AttributeKeyAmount, p.Amount.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}

func handleBurnDistributionProfitsProposal(ctx sdk.Context, k Keeper, p BurnDistributionProfitsProposal) error {
	err := k.BurnFromDistributionProfits(ctx, p.Amount)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBurnDistributionProfits,
			sdk.NewAttribute(types.AttributeKeyTitle, p.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, p.Description),
			sdk.NewAttribute(types.AttributeKeyAmount, p.Amount.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}

func handleTransferFromDistributionProfitsToBuyBackLiquidityProposal(ctx sdk.Context, k Keeper, p TransferFromDistributionProfitsToBuyBackLiquidityProposal) error {
	err := k.TransferFromDistributionProfitsToBuyBackLiquidity(ctx, p.Amount)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransferFromDistributionProfitsToBuyBackLiquidity,
			sdk.NewAttribute(types.AttributeKeyTitle, p.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, p.Description),
			sdk.NewAttribute(types.AttributeKeyAmount, p.Amount.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}

func handleMsgAddOperator(ctx sdk.Context, k Keeper, msg MsgAddOperator) (*sdk.Result, error) {
	err := k.HandleAddOperator(ctx, msg.Sender, msg.Operator)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRemoveOperator(ctx sdk.Context, k Keeper, msg MsgRemoveOperator) (*sdk.Result, error) {
	err := k.HandleRemoveOperator(ctx, msg.Sender, msg.Operator)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDisburse(ctx sdk.Context, k Keeper, msg MsgDisburse) (*sdk.Result, error) {
	err := k.HandleDisburse(ctx, msg.Operator, msg.Recipient, msg.Amount, msg.Reference)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDisburseToEscrow(ctx sdk.Context, k Keeper, msg MsgDisburseToEscrow) (*sdk.Result, error) {
	err := k.HandleDisburseToEscrow(ctx, msg.Operator, msg.Amount, msg.Reference)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDisburseFromEscrow(ctx sdk.Context, k Keeper, msg MsgDisburseFromEscrow) (*sdk.Result, error) {
	err := k.HandleDisburseFromEscrow(ctx, msg.Operator, msg.Reference, msg.Recipient)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRevertFromEscrow(ctx sdk.Context, k Keeper, msg MsgRevertFromEscrow) (*sdk.Result, error) {
	err := k.HandleRevertFromEscrow(ctx, msg.Operator, msg.Amount, msg.Reference)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCancelDisbursement(ctx sdk.Context, k Keeper, msg MsgCancelDisbursement) (*sdk.Result, error) {
	scheduledFor, err := time.Parse("2006-01-02T15:04:05.99999999999Z", msg.ScheduledFor)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidTime, msg.ScheduledFor)
	}

	err = k.HandleCancelDisbursement(ctx, msg.Manager, msg.Recipient, scheduledFor)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCreateSellOrder(ctx sdk.Context, k Keeper, msg MsgCreateSellOrder) (*sdk.Result, error) {
	err := k.HandleCreateSellOrder(ctx, msg.Seller, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCreateBuyOrder(ctx sdk.Context, k Keeper, msg MsgCreateBuyOrder) (*sdk.Result, error) {
	err := k.HandleCreateBuyOrder(ctx, msg.Buyer, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
