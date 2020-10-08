package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
)

func (k Keeper) TransferFromDevelopmentFund(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.Coins) error {
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.DevelopmentFundModuleName, recipient , amount)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("development -> %s : %s", recipient, amount),
	)

	return nil
}

func (k Keeper) TransferFromSecurityTokenFund(ctx sdk.Context, recipients []types.Recipients) error {
	for i := 0; i < len(recipients); i++ {
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.SecurityTokenFundModuleName, recipients[i].Recipient , recipients[i].Amount)
		if err != nil {
			return err
		}
		k.Logger(ctx).Debug(
			fmt.Sprintf("security -> %s : %s", recipients[i].Recipient, recipients[i].Amount),
		)
	}

	return nil
}

func (k Keeper) DistributeNameReward(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.DecCoins) error {
	if ctx.BlockTime().Before(k.RewardWithdrawalEnabledTime(ctx)) {
		// If we are in the period where withdrawals are disabled, this function only gets executed as a result of a Hook call
		// In that case we move all the funds to escrow

		escrow := k.GetNameRewardEscrow(ctx, recipient)
		escrow = escrow.Add(amount...)
		k.SetNameRewardEscrow(ctx, recipient, escrow)

		k.Logger(ctx).Debug(
			fmt.Sprintf("[V] hhrm -> escrow[%s] : %s ", recipient, escrow),
		)

		return nil
	}

	// we add the balance of the escrow account to the current distribution and delete the escrow account
	inEscrow := k.GetNameRewardEscrow(ctx, recipient)
	amount = amount.Add(inEscrow...)
	k.DeleteNameRewardEscrow(ctx, recipient)

	k.Logger(ctx).Debug(
		fmt.Sprintf("Distributing Name rewards for %s. Total %s, in escrow %s", recipient, amount, inEscrow),
	)

	// we split the reward in to whole (rewardInt) and fraction (rewardLeftover)
	rewardDec := amount.AmountOf(config.DefaultDenom)
	rewardInt := rewardDec.TruncateInt()
	rewardLeftover := rewardDec.Sub(rewardInt.ToDec())

	if ! rewardLeftover.IsZero() {
		// Current reward leftover is greater than 0.
		// This means that we need to add it to the existing leftovers and send whole pin while persisting fractions

		previousLeftover := k.GetNameRewardLeftover(ctx, recipient)
		currentLeftover := rewardLeftover.Add(previousLeftover)

		// we split the leftover in to whole (currentLeftoverInt) and fraction (currentLeftover)
		currentLeftoverInt := currentLeftover.TruncateInt()
		currentLeftover = currentLeftover.Sub(currentLeftoverInt.ToDec())

		// we add whole pin to the rewardInt variable
		rewardInt = rewardInt.Add(currentLeftoverInt)

		k.Logger(ctx).Debug(
			fmt.Sprintf("Previous leftover: %s, Current leftover: %s", previousLeftover, currentLeftover),
		)

		// we update the remaining leftover
		k.SetNameRewardLeftover(ctx, recipient, currentLeftover)
	}

	toTransfer := sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, rewardInt))

	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.HRAHolderRewardModuleName, recipient , toTransfer)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("hhrm -> %s : %s", recipient, toTransfer),
	)

	return nil
}

func (k Keeper) DistributeSavingsReward(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.DecCoins) error {
	if ctx.BlockTime().Before(k.RewardWithdrawalEnabledTime(ctx)) {
		// If we are in the period where withdrawals are disabled, this function only gets executed as a result of
		// Withdraw Savings transaction where Funds being taken out of savings and interest is transferred to escrow
		// In that case we move all the funds to escrow

		escrow := k.GetSavingsRewardEscrow(ctx, recipient)
		escrow = escrow.Add(amount...)
		k.SetSavingsRewardEscrow(ctx, recipient, escrow)

		k.Logger(ctx).Debug(
			fmt.Sprintf("[V] savingsdistr -> escrow[%s] : %s ", recipient, escrow),
		)

		return nil
	}

	// we add the balance of the escrow account to the current distribution and delete the escrow account
	inEscrow := k.GetSavingsRewardEscrow(ctx, recipient)
	amount = amount.Add(inEscrow...)
	k.DeleteSavingsRewardEscrow(ctx, recipient)

	k.Logger(ctx).Debug(
		fmt.Sprintf("Distributing Savings rewards for %s. Total %s, in escrow %s", recipient, amount, inEscrow),
	)

	// we split the reward in to whole (rewardInt) and fraction (rewardLeftover)
	rewardDec := amount.AmountOf(config.DefaultDenom)
	rewardInt := rewardDec.TruncateInt()
	rewardLeftover := rewardDec.Sub(rewardInt.ToDec())

	if ! rewardLeftover.IsZero() {
		// Current reward leftover is greater than 0.
		// This means that we need to add it to the existing leftovers and send whole pin while persisting fractions

		previousLeftover := k.GetSavingsRewardLeftover(ctx, recipient)
		currentLeftover := rewardLeftover.Add(previousLeftover)

		// we split the leftover in to whole (currentLeftoverInt) and fraction (currentLeftover)
		currentLeftoverInt := currentLeftover.TruncateInt()
		currentLeftover = currentLeftover.Sub(currentLeftoverInt.ToDec())

		// we add whole pin to the rewardInt variable
		rewardInt = rewardInt.Add(currentLeftoverInt)

		k.Logger(ctx).Debug(
			fmt.Sprintf("Previous leftover: %s, Current leftover: %s", previousLeftover, currentLeftover),
		)

		// we update the remaining leftover
		k.SetSavingsRewardLeftover(ctx, recipient, currentLeftover)
	}

	if ! rewardInt.IsZero() {
		toTransfer := sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, rewardInt))

		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.SavingsDistributionModuleName, recipient , toTransfer)
		if err != nil {
			return err
		}
		k.Logger(ctx).Debug(
			fmt.Sprintf("savingsdistr -> %s : %s", recipient, toTransfer),
		)
	}

	return nil
}

func (k Keeper) RefundSavingsStake(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.Coins) error {
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.SavingsModuleName, recipient, amount)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("savings -> %s : %s", recipient, amount),
	)

	return nil
}


func (k Keeper) ClaimSavingsStake(ctx sdk.Context, sender sdk.AccAddress, amount sdk.Coins) error {
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.SavingsModuleName, amount)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug(
		fmt.Sprintf("%s -> savings : %s", sender, amount),
	)

	return nil
}