package distribution

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
	"github.com/anathatech/project-anatha/x/hra"
)


func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
			case MsgWithdrawNameReward:
				return handleMsgWithdrawNameReward(ctx, k, msg)

			case MsgWithdrawValidatorReward:
				return handleMsgWithdrawValidatorReward(ctx, k, msg)

			case MsgDepositSavings:
				return handleMsgDepositSavings(ctx, k, msg)

			case MsgWithdrawSavings:
				return handleMsgWithdrawSavings(ctx, k, msg)

			case MsgWithdrawSavingsInterest:
				return handleMsgWithdrawSavingsInterest(ctx, k, msg)

			default:
				errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName,  msg)
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgWithdrawNameReward(ctx sdk.Context, k Keeper, msg MsgWithdrawNameReward)(*sdk.Result, error) {
	// check if sender has a HRA
	if ctx.BlockHeight() > hra.NameConstraintBlock && ! k.HraKeeper.OwnsAnyName(ctx, msg.Sender) {
		return nil, hra.ErrNameNotRegistered
	}

	err := k.HandleWithdrawNameReward(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdrawValidatorReward(ctx sdk.Context, k Keeper, msg MsgWithdrawValidatorReward) (*sdk.Result, error) {
	// check if sender has a HRA
	if ctx.BlockHeight() > hra.NameConstraintBlock && ! k.HraKeeper.OwnsAnyName(ctx, msg.GetSigners()[0]) {
		return nil, hra.ErrNameNotRegistered
	}

	err := k.HandleWithdrawValidatorReward(ctx, msg.Validator)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDepositSavings(ctx sdk.Context, k Keeper, msg MsgDepositSavings) (*sdk.Result, error) {
	// check if sender has a HRA
	if ctx.BlockHeight() > hra.NameConstraintBlock && ! k.HraKeeper.OwnsAnyName(ctx, msg.Sender) {
		return nil, hra.ErrNameNotRegistered
	}

	err := k.HandleDepositSavings(ctx, msg.Sender, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdrawSavings(ctx sdk.Context, k Keeper, msg MsgWithdrawSavings) (*sdk.Result, error) {
	// check if sender has a HRA
	if ctx.BlockHeight() > hra.NameConstraintBlock && ! k.HraKeeper.OwnsAnyName(ctx, msg.Sender) {
		return nil, hra.ErrNameNotRegistered
	}

	err := k.HandleWithdrawSavings(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdrawSavingsInterest(ctx sdk.Context, k Keeper, msg MsgWithdrawSavingsInterest) (*sdk.Result, error) {
	// check if sender has a HRA
	if ctx.BlockHeight() > hra.NameConstraintBlock && ! k.HraKeeper.OwnsAnyName(ctx, msg.Sender) {
		return nil, hra.ErrNameNotRegistered
	}

	err := k.HandleWithdrawSavingsInterest(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func NewDistributionProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
			case DevelopmentFundDistributionProposal:
				return handleDevelopmentFundDistributionProposal(ctx, k, c)

			case SecurityTokenFundDistributionProposal:
				return handleSecurityTokenFundDistributionProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distribution proposal content type: %T", c)
		}
	}
}

func handleDevelopmentFundDistributionProposal(ctx sdk.Context, k Keeper, p DevelopmentFundDistributionProposal) error {
	err := k.TransferFromDevelopmentFund(ctx, p.Recipient, p.Amount)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDevelopmentFundDistribution,
			sdk.NewAttribute(types.AttributeKeyTitle, p.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, p.Description),
			sdk.NewAttribute(types.AttributeKeyAmount, p.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, p.Recipient.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}

func handleSecurityTokenFundDistributionProposal(ctx sdk.Context, k Keeper, p SecurityTokenFundDistributionProposal) error {
	err := k.TransferFromSecurityTokenFund(ctx, p.Recipients)
	if err != nil {
		return err
	}

	var attributes []sdk.Attribute

	for _, recipient := range p.Recipients {
		attributes = append(attributes, sdk.NewAttribute(types.AttributeKeyAmount, recipient.Amount.String()))
		attributes = append(attributes, sdk.NewAttribute(types.AttributeKeyRecipient, recipient.Recipient.String()))
	}

	event := sdk.NewEvent(
		types.EventTypeSecurityTokenFundDistribution,
		sdk.NewAttribute(types.AttributeKeyTitle, p.Title),
		sdk.NewAttribute(types.AttributeKeyDescription, p.Description),
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
	)

	event = event.AppendAttributes(attributes...)

	ctx.EventManager().EmitEvent(event)

	return nil
}