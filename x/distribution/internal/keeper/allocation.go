package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/staking/exported"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// AllocateTokens handles distribution of the collected fees
func (k Keeper) AllocateTokens(ctx sdk.Context, previousVotes []abci.VoteInfo) {
	var totalPreviousPower int64
	for _, voteInfo := range previousVotes {
		totalPreviousPower += voteInfo.Validator.Power
	}

	if totalPreviousPower == 0 {
		k.Logger(ctx).Debug(
			fmt.Sprintf("Skipping token allocation because of totalPreviousPower = 0"),
		)
		return
	}

	feeCollector := k.supplyKeeper.GetModuleAccount(ctx, types.NvrpModuleName)
	feesCollectedInt := feeCollector.GetCoins()
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// transfer collected fees to the distribution module account
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.NvrpModuleName, types.NvrpDistributionModuleName, feesCollectedInt)
	if err != nil {
		panic(err)
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("nvrp -> nvrpd : %s", feesCollectedInt),
	)

	remaining := feesCollected

	// allocate tokens proportionally to voting power
	for _, vote := range previousVotes {
		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)

		powerFraction := sdk.NewDec(vote.Validator.Power).QuoTruncate(sdk.NewDec(totalPreviousPower))
		reward := feesCollected.MulDecTruncate(powerFraction)

		k.AllocateTokensToValidator(ctx, validator, reward)

		remaining = remaining.Sub(reward)
	}

	accumulatedRemainder := remaining.Add(k.GetNvrpRemainder(ctx)...)

	if ! accumulatedRemainder.IsZero() {
		k.Logger(ctx).Debug(
			fmt.Sprintf("NVRPD accumulated remainder : %s", accumulatedRemainder),
		)

		whole, remainder := accumulatedRemainder.TruncateDecimal()
		k.SetNvrpRemainder(ctx, remainder)

		if ! whole.IsZero() {
			err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.NvrpDistributionModuleName, types.NvrpModuleName, whole)
			if err != nil {
				panic(err)
			}
			k.Logger(ctx).Debug(
				fmt.Sprintf("nvrpd -> nvrp : %s", whole),
			)
		}
	}
}

func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val exported.ValidatorI, tokens sdk.DecCoins) {
	currentRewards := k.GetValidatorAccumulatedRewards(ctx, val.GetOperator())
	currentRewards = currentRewards.Add(tokens...)
	k.SetValidatorAccumulatedRewards(ctx, val.GetOperator(), currentRewards)

	k.Logger(ctx).Debug(
		fmt.Sprintf("nvrpd -> (%s) %s : %s", val.GetMoniker(), sdk.AccAddress(val.GetOperator()), tokens),
	)
}

