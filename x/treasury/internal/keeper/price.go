package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
)

func (k Keeper) CoinsPerStage(ctx sdk.Context) sdk.Int {
	treasury := k.GetTreasury(ctx)
	return treasury.CoinsPerStage.AmountOf(config.DefaultDenom)
}

func (k Keeper) DistributedFromTreasury(ctx sdk.Context) sdk.Int {
	treasury := k.GetTreasury(ctx)
	return treasury.Distributed.AmountOf(config.DefaultDenom)
}

// Returns the Din amount for input Pin taking in to account Treasury BuyBack
func (k Keeper) CalculateDinAmountExtended(ctx sdk.Context, pinCoins sdk.Coins) sdk.Coins {
	stagePrice := k.GetPriceForStage(
		ctx,
		k.GetStageFromDistribution(
			ctx,
			k.DistributedFromTreasury(ctx),
		),
	)

	buyBackPinBalance := k.supplyKeeper.GetModuleAccount(ctx, types.BuyBackFundModuleName).GetCoins().AmountOf(config.DefaultDenom)

	dinFromBuyBack := sdk.NewCoins()
	dinFromTreasury := sdk.NewCoins()

	if buyBackPinBalance.GTE(pinCoins.AmountOf(config.DefaultDenom)) {
		// We are fulfilling the order only from the buyback account
		dinFromBuyBack = sdk.NewCoins(
			sdk.NewCoin(
				config.DefaultStableDenom,
				pinCoins.AmountOf(config.DefaultDenom).Mul(stagePrice),
			),
		)

		dinFromTreasury = sdk.NewCoins()
	} else {
		// We are partially fulfilling the order from the buyback and partially from the treasury
		dinFromBuyBack = sdk.NewCoins(
			sdk.NewCoin(
				config.DefaultStableDenom,
				buyBackPinBalance.Mul(stagePrice),
			),
		)

		pinCoins = pinCoins.Sub(
			sdk.NewCoins(
				sdk.NewCoin(
					config.DefaultDenom,
					buyBackPinBalance,
				),
			),
		)

		dinFromTreasury = k.CalculateDinAmount(ctx, pinCoins)
	}

	return dinFromBuyBack.Add(dinFromTreasury...)
}

func (k Keeper) CalculateDinAmount(ctx sdk.Context, pinCoins sdk.Coins) sdk.Coins {
	pinAmount := pinCoins.AmountOf(config.DefaultDenom)
	dinAmount := sdk.ZeroInt()

	if pinAmount.LTE(sdk.ZeroInt()) {
		return sdk.NewCoins(sdk.NewCoin(config.DefaultStableDenom, dinAmount))
	}

	stage := k.GetStageFromDistribution(ctx, k.DistributedFromTreasury(ctx))

	canonicalStage := stage

	for ! pinAmount.IsZero() {

		remaining := k.RemainingPinInStage(ctx, stage, canonicalStage)
		stagePrice := k.GetPriceForStage(ctx, stage)

		if remaining.LT(pinAmount) {
			// we are clearing out a stage and transitioning to the next

			dinAmount = dinAmount.Add(remaining.Mul(stagePrice))
			pinAmount = pinAmount.Sub(remaining)
		} else {
			// we have enough liquidity in the current stage

			dinAmount = dinAmount.Add(pinAmount.Mul(stagePrice))
			pinAmount = sdk.ZeroInt()
		}

		stage = stage.Add(sdk.OneInt())
	}

	return sdk.NewCoins(sdk.NewCoin(config.DefaultStableDenom, dinAmount))
}

// Returns the Pin amount for input Din taking in to account Treasury BuyBack
func (k Keeper) CalculatePinAmountExtended(ctx sdk.Context, dinCoins sdk.Coins) (total sdk.Coins, pinFromBuyBack sdk.Coins, pinFromTreasury sdk.Coins) {
	stagePrice := k.GetPriceForStage(
		ctx,
		k.GetStageFromDistribution(
			ctx,
			k.DistributedFromTreasury(ctx),
		),
	)

	buyBackPinBalance := k.supplyKeeper.GetModuleAccount(ctx, types.BuyBackFundModuleName).GetCoins().AmountOf(config.DefaultDenom)

	if stagePrice.Mul(buyBackPinBalance).GTE(dinCoins.AmountOf(config.DefaultStableDenom)) {
		// We are fulfilling the order only from the buyback account
		pinFromBuyBack = sdk.NewCoins(
			sdk.NewCoin(
				config.DefaultDenom,
				dinCoins.AmountOf(config.DefaultStableDenom).Quo(stagePrice),
			),
		)

		pinFromTreasury = sdk.NewCoins()
	} else {
		// We are partially fulfilling the order from the buyback and partially from the treasury
		pinFromBuyBack = sdk.NewCoins(
			sdk.NewCoin(
				config.DefaultDenom,
				buyBackPinBalance,
			),
		)

		dinCoins = dinCoins.Sub(
			sdk.NewCoins(
				sdk.NewCoin(
					config.DefaultStableDenom,
					pinFromBuyBack.AmountOf(config.DefaultDenom).Mul(stagePrice),
				),
			),
		)

		pinFromTreasury = k.CalculatePinAmount(ctx, dinCoins)
	}

	total = pinFromBuyBack.Add(pinFromTreasury...)

	return
}

func (k Keeper) CalculatePinAmount(ctx sdk.Context, dinCoins sdk.Coins) sdk.Coins {
	dinAmount := dinCoins.AmountOf(config.DefaultStableDenom)
	pinAmount := sdk.ZeroInt()

	if dinAmount.LTE(sdk.ZeroInt()) {
		return sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, pinAmount))
	}

	stage := k.GetStageFromDistribution(ctx, k.DistributedFromTreasury(ctx))

	canonicalStage := stage

	for ! dinAmount.IsZero() {

		remainingPin := k.RemainingPinInStage(ctx, stage, canonicalStage)
		stagePrice := k.GetPriceForStage(ctx, stage)

		remainingDinPrice := remainingPin.Mul(stagePrice)

		if remainingDinPrice.LT(dinAmount) {
			// we are clearing out a stage and transitioning to the next

			pinAmount = pinAmount.Add(remainingPin)
			dinAmount = dinAmount.Sub(remainingDinPrice)
		} else {
			// we have enough liquidity in the current stage

			pinAmount = pinAmount.Add(dinAmount.Quo(stagePrice))
			dinAmount = sdk.ZeroInt()
		}

		stage = stage.Add(sdk.OneInt())
	}

	return sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, pinAmount))
}

// returned amount is $0.0000000001 per PIN
func (k Keeper) GetPriceForStage(ctx sdk.Context, stage sdk.Int) sdk.Int {
	return stage.Add(sdk.OneInt())
}

func (k Keeper) GetStageFromDistribution(ctx sdk.Context, distribution sdk.Int) sdk.Int {
	return distribution.Quo(k.CoinsPerStage(ctx))
}

func (k Keeper) RemainingPinInStage(ctx sdk.Context, stage sdk.Int, cannonicalStage sdk.Int) sdk.Int {
	coinsPerStage := k.CoinsPerStage(ctx)

	if stage.LT(cannonicalStage) {
		// given stage is cleared

		return sdk.ZeroInt()
	} else if stage.Equal(cannonicalStage) {
		// we are in the current stage

		return coinsPerStage.Sub(k.DistributedFromTreasury(ctx).Mod(coinsPerStage))
	} else if stage.GT(cannonicalStage) {
		// we are in a future stage

		return coinsPerStage
	} else {
		panic("This should not happen.")
	}
}