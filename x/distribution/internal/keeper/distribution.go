package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
)

func (k Keeper) DistributeFromAmc(ctx sdk.Context) (sdk.Coins, sdk.Coins, sdk.Coins) {
	// Distributin from AMC to:
	// Anatha Development Fund Module - DevelopmentFundShare%
	// Anata Security Token Holder Fund Module - SecurityTokenFundShare%
	// Anatha HRA Holder Reward Module - rest
	amcCoins := k.supplyKeeper.GetModuleAccount(ctx, types.AmcModuleName).GetCoins()

	toDfm := sdk.NewCoins(
		sdk.NewCoin(
			config.DefaultDenom,
			amcCoins.AmountOf(config.DefaultDenom).ToDec().MulTruncate(k.DevelopmentFundShare(ctx)).TruncateInt(),
		),
	)

	toStfm := sdk.NewCoins(
		sdk.NewCoin(
			config.DefaultDenom,
			amcCoins.AmountOf(config.DefaultDenom).ToDec().MulTruncate(k.SecurityTokenFundShare(ctx)).TruncateInt(),
		),
	)

	toHhrm := amcCoins.Sub(toDfm).Sub(toStfm)

	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.AmcModuleName, types.HRAHolderRewardModuleName, toHhrm)
	if err != nil {
		panic(fmt.Sprintf("failed transfer from %s to %s, amount: %s, %s", types.AmcModuleName, types.HRAHolderRewardModuleName, toHhrm, err))
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("%s -> %s : %s", types.AmcModuleName, types.HRAHolderRewardModuleName, toHhrm),
	)

	err = k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.AmcModuleName, types.DevelopmentFundModuleName, toDfm)
	if err != nil {
		panic(fmt.Sprintf("failed transfer from %s to %s, amount: %s, %s", types.AmcModuleName, types.DevelopmentFundModuleName, toDfm, err))
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("%s -> %s : %s", types.AmcModuleName, types.DevelopmentFundModuleName, toDfm),
	)

	err = k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.AmcModuleName, types.SecurityTokenFundModuleName, toStfm)
	if err != nil {
		panic(fmt.Sprintf("failed transfer from %s to %s, amount: %s, %s", types.AmcModuleName, types.SecurityTokenFundModuleName, toStfm, err))
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("%s -> %s : %s", types.AmcModuleName, types.SecurityTokenFundModuleName, toStfm),
	)

	return toHhrm, toDfm, toStfm
}

func (k Keeper) DistributeFromHhrm(ctx sdk.Context, amount sdk.Coins) {
	pending := k.GetPendingNameDistribution(ctx)
	total := pending.Add(amount...)

	distributed := k.distributeNameReward(ctx, total)

	if distributed {
		// After successful distribution we remove pending distribution
		k.DeletePendingNameDistribution(ctx)
	} else {
		// If there was no distibution due to the lack of HRA holders, we accumulate the pending distribution
		k.SetPendingNameDistribution(ctx, total)
	}
}

func (k Keeper) DistributeFromNvrp(ctx sdk.Context){
	// Distribute From Nvrp to Savings. The rest of the balance will be allocated for the Validators
	nvrpBalanceInt := k.supplyKeeper.GetModuleAccount(ctx, types.NvrpModuleName).GetCoins().AmountOf(config.DefaultDenom)
	nvrpBalanceDec := nvrpBalanceInt.ToDec()

	bonded := k.stakingKeeper.TotalBondedTokens(ctx).ToDec()
	saved := k.GetSavingsStake(ctx).ToDec()

	adjustment := k.SavingsSplitAdjustment(ctx)

	forSavingsDec := nvrpBalanceDec.Mul(
		saved.QuoTruncate(
			saved.Add(bonded),
		),
	).MulTruncate(adjustment)

	forSavingsInt := forSavingsDec.TruncateInt()

	if ! forSavingsInt.IsZero() {
		forSavings := sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, forSavingsInt))

		err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.NvrpModuleName, types.SavingsDistributionModuleName, forSavings)
		if err != nil {
			panic(fmt.Sprintf("failed transfer from %s to %s, amount: %s, %s", types.NvrpModuleName, types.SavingsDistributionModuleName, forSavings, err))
		}
		k.Logger(ctx).Debug(
			fmt.Sprintf("%s -> %s : %s", types.NvrpModuleName, types.SavingsDistributionModuleName, forSavings),
		)

		k.distributeSavingsReward(ctx, forSavings)
	}
}