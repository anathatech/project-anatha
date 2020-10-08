package staking

import (
	"github.com/anathatech/project-anatha/x/hra"
	"time"

	tmstrings "github.com/tendermint/tendermint/libs/strings"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/x/staking/keeper"
	"github.com/anathatech/project-anatha/x/staking/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateValidator:
			return handleMsgCreateValidator(ctx, msg, k)

		case types.MsgEditValidator:
			return handleMsgEditValidator(ctx, msg, k)

		case types.MsgDelegate:
			return handleMsgDelegate(ctx, msg, k)

		case types.MsgUndelegate:
			return handleMsgUndelegate(ctx, msg, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// These functions assume everything has been authenticated,
// now we just perform action and save

func handleMsgCreateValidator(ctx sdk.Context, msg types.MsgCreateValidator, k keeper.Keeper) (*sdk.Result, error) {
	// check if sender has a HRA
	if ctx.BlockHeight() > hra.NameConstraintBlock && ! k.HraKeeper.OwnsAnyName(ctx, msg.DelegatorAddress) {
		return nil, hra.ErrNameNotRegistered
	}

	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetValidator(ctx, msg.ValidatorAddress); found {
		return nil, ErrValidatorOwnerExists
	}

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(msg.PubKey)); found {
		return nil, ErrValidatorPubKeyExists
	}

	if msg.Value.Denom != k.BondDenom(ctx) {
		return nil, ErrBadDenom
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}

	if ctx.ConsensusParams() != nil {
		tmPubKey := tmtypes.TM2PB.PubKey(msg.PubKey)
		if !tmstrings.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
			return nil, sdkerrors.Wrapf(
				ErrValidatorPubKeyTypeNotSupported,
				"got: %s, valid: %s", tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes,
			)
		}
	}

	validator := NewValidator(msg.ValidatorAddress, msg.PubKey, msg.Description)

	validator.MinSelfDelegation = types.DefaultStake

	validator.Ticket = k.GetLastTicket(ctx)
	k.SetLastTicket(ctx, validator.Ticket + 1)

	k.SetValidator(ctx, validator)
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)
	k.SetNewValidatorByTicket(ctx, validator)

	// call the after-creation hook
	k.AfterValidatorCreated(ctx, validator.OperatorAddress)

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	// NOTE source will always be from a wallet which are unbonded
	_, err := k.Delegate(ctx, msg.DelegatorAddress, types.DefaultStake, sdk.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateValidator,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgEditValidator(ctx sdk.Context, msg types.MsgEditValidator, k keeper.Keeper) (*sdk.Result, error) {
	// check if sender has a HRA
	if ctx.BlockHeight() > hra.NameConstraintBlock && ! k.HraKeeper.OwnsAnyName(ctx, msg.GetSigners()[0]) {
		return nil, hra.ErrNameNotRegistered
	}

	// validator must already be registered
	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return nil, ErrNoValidatorFound
	}

	// replace all editable fields (clients should autofill existing values)
	description, err := validator.Description.UpdateDescription(msg.Description)
	if err != nil {
		return nil, err
	}

	validator.Description = description

	k.SetValidator(ctx, validator)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDelegate(ctx sdk.Context, msg types.MsgDelegate, k keeper.Keeper) (*sdk.Result, error) {
	// check if sender has a HRA
	if ctx.BlockHeight() > hra.NameConstraintBlock && ! k.HraKeeper.OwnsAnyName(ctx, msg.GetSigners()[0]) {
		return nil, hra.ErrNameNotRegistered
	}

	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return nil, ErrNoValidatorFound
	}

	toDelegate := types.DefaultStake.Sub(validator.Tokens)
	if ! toDelegate.IsPositive() {
		return nil, ErrBadDelegationAmount
	}

	// NOTE: source funds are always unbonded
	_, err := k.Delegate(ctx, msg.DelegatorAddress, toDelegate, sdk.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, toDelegate.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUndelegate(ctx sdk.Context, msg types.MsgUndelegate, k keeper.Keeper) (*sdk.Result, error) {
	// check if sender has a HRA
	if ctx.BlockHeight() > hra.NameConstraintBlock && ! k.HraKeeper.OwnsAnyName(ctx, msg.DelegatorAddress) {
		return nil, hra.ErrNameNotRegistered
	}

	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return nil, ErrNoValidatorFound
	}

	amount := validator.Tokens

	shares, err := k.ValidateUnbondAmount(
		ctx, msg.DelegatorAddress, msg.ValidatorAddress, amount,
	)
	if err != nil {
		return nil, err
	}

	completionTime, err := k.Undelegate(ctx, msg.DelegatorAddress, msg.ValidatorAddress, shares)
	if err != nil {
		return nil, err
	}

	completionTimeBz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(completionTime)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbond,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	})

	return &sdk.Result{Data: completionTimeBz, Events: ctx.EventManager().Events()}, nil
}