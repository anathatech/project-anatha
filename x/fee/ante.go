package fee

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/anathatech/project-anatha/config"
	types2 "github.com/anathatech/project-anatha/x/fee/internal/types"
	"github.com/anathatech/project-anatha/x/hra"
	"github.com/anathatech/project-anatha/x/treasury"
)


type FeeDecorator struct {
	feeKeeper			Keeper
	bankKeeper   		bank.Keeper
	hraKeeper    		hra.Keeper
	supplyKeeper 		supply.Keeper
	feeCollectorModule 	string
}

func NewFeeDecorator(fk Keeper, bk bank.Keeper, hk hra.Keeper, sk supply.Keeper, feeCollectorModule string) FeeDecorator {
	if addr := sk.GetModuleAddress(feeCollectorModule); addr == nil {
		panic("the fee collector module account has not been set")
	}

	return FeeDecorator{
		feeKeeper: fk,
		bankKeeper: bk,
		hraKeeper: hk,
		supplyKeeper: sk,
		feeCollectorModule: feeCollectorModule,
	}
}

func (d FeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	stdTx, ok := tx.(types.StdTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a StdTx")
	}

	msgs := stdTx.GetMsgs()
	feePayer := stdTx.FeePayer()

	txFees := sdk.NewCoins()
	systemFees := sdk.NewCoins()

	namesOwnedCount := d.hraKeeper.GetNamesByAddressCount(ctx, feePayer)
	credits := d.hraKeeper.GetCredits(ctx, feePayer)

	for _, msg := range msgs {
		msgFee := sdk.NewCoins()

		switch msg := msg.(type) {
		case hra.MsgRegisterName:
			if namesOwnedCount == 0 {
				namesOwnedCount++
				credits = d.hraKeeper.AddressCredits(ctx)
			}

			msgFee = msgFee.Add(d.hraKeeper.NameInfoRegistrationFee(ctx)...)

		case hra.MsgRenewName:
			msgFee = msgFee.Add(d.hraKeeper.NameInfoRenewalFee(ctx)...)

		case hra.MsgBuyName:
			fee, err := d.hraKeeper.GetPrice(ctx, msg.Name)
			if err != nil {
				return ctx, err
			}

			if namesOwnedCount == 0 {
				namesOwnedCount++
				credits = d.hraKeeper.AddressCredits(ctx)
			}

			msgFee = msgFee.Add(fee...)

		case hra.MsgRegisterAddress:
			if credits.LTE(sdk.ZeroInt()) {
				msgFee = msgFee.Add(d.hraKeeper.AddressRegistrationFee(ctx)...)
			}

			credits = credits.Sub(sdk.OneInt())

		case hra.MsgTransferName:
			if namesOwnedCount == 1 {
				namesOwnedCount--
				credits = sdk.ZeroInt()
			}

		case hra.MsgDeleteName:
			if namesOwnedCount == 1 {
				namesOwnedCount--
				credits = sdk.ZeroInt()
			}

		case bank.MsgSend:
			msgFee = msgFee.Add(msg.Amount...)

		case bank.MsgMultiSend:
			for _, input := range msg.Inputs {
				msgFee = msgFee.Add(input.Coins...)
			}

		case treasury.MsgCreateSellOrder:
			msgFee = msgFee.Add(msg.Amount...)
		}

		txFees = txFees.Add(msgFee...)

		if ! d.feeKeeper.IsMessageFeeExcluded(ctx, msg) {
			systemFees = systemFees.Add(d.CalculateSystemFee(ctx, msgFee)...)
		}
	}

	totalFees := txFees.Add(systemFees...)

	if ! d.bankKeeper.HasCoins(ctx, feePayer, totalFees) {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; required: %s", totalFees)
	}

	err = d.supplyKeeper.SendCoinsFromAccountToModule(
		ctx,
		feePayer,
		d.feeCollectorModule,
		systemFees,			// we deduct only system fees in the ante handler!
	)

	if err != nil {
		return ctx, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types2.AttributeValueModule,
			sdk.NewAttribute(types2.AttributeKeyTotalFee, totalFees.String()),
			sdk.NewAttribute(types2.AttributeKeySystemFee, systemFees.String()),
		),
	)

	return next(ctx, tx, simulate)
}

func (d FeeDecorator) CalculateSystemFee(ctx sdk.Context, txFee sdk.Coins) sdk.Coins {
	feePercentage := d.feeKeeper.FeePercentage(ctx)
	systemFeeInt := txFee.AmountOf(config.DefaultDenom).ToDec().Mul(feePercentage).TruncateInt()

	systemFee := sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, systemFeeInt))
	minimumFee := d.feeKeeper.MinimumFee(ctx)
	maximumFee := d.feeKeeper.MaximumFee(ctx)

	if systemFee.IsAllLT(minimumFee) {
		return minimumFee
	}

	if systemFee.IsAllGT(maximumFee) {
		return maximumFee
	}

	return systemFee
}