package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
)

// Handler

func (k Keeper) HandleDepositSavings(ctx sdk.Context, sender sdk.AccAddress, amount sdk.Coins) error {
	if k.HasSavings(ctx, sender) {
		return types.ErrAlreadyHasSavings
	}

	err := k.depositSavings(ctx, sender, amount, true)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDepositSavings,
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
		),
	})

	return nil
}

func (k Keeper) HandleWithdrawSavings(ctx sdk.Context, sender sdk.AccAddress) error {
	if ! k.HasSavings(ctx, sender) {
		return types.ErrHasNoSavings
	}

	reward, err := k.withdrawSavingsReward(ctx, sender, true)
	if err != nil {
		return err
	}

	if ! reward.IsZero() {
		err := k.DistributeSavingsReward(ctx, sender, reward)
		if err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawSavings,
			sdk.NewAttribute(types.AttributeKeyReward, k.GetSavingsRewardEscrow(ctx, sender).Add(reward...).String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
		),
	})

	return nil
}

func (k Keeper) HandleWithdrawSavingsInterest(ctx sdk.Context, sender sdk.AccAddress) error {
	if ctx.BlockTime().Before(k.RewardWithdrawalEnabledTime(ctx)) {
		return types.ErrRewardWithdrawalDisabled
	}

	amount, receivingRewards := k.GetSavingsStakeByAddress(ctx, sender)

	reward, _ := k.withdrawSavingsReward(ctx, sender, false) // Error can be ignored due to supply keeper not being invoked

	err := k.DistributeSavingsReward(ctx, sender, reward)
	if err != nil {
		return err
	}

	// If the user was signed up for interest the withdrawSavingsReward reset the reward state and we need to deposit the amount again
	if receivingRewards {
		coins := sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, amount))

		err := k.depositSavings(ctx, sender, coins, false)
		if err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawSavingsInterest,
			sdk.NewAttribute(types.AttributeKeyReward, k.GetSavingsRewardEscrow(ctx, sender).Add(reward...).String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
		),
	})

	return nil
}

// Util

func (k Keeper) HasSavings(ctx sdk.Context, address sdk.AccAddress) bool {
	_, found := k.GetSavingsStakeByAddress(ctx, address)

	return found
}

// Algorithm

func (k Keeper) depositSavings(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coins, shouldMoveCoins bool) error {
	rate := k.GetSavingsRewardRate(ctx)

	k.SetSavingsRewardRateByAddress(ctx, address, rate)

	amountInt := amount.AmountOf(config.DefaultDenom)

	k.SetSavingsStakeByAddress(ctx, address, amountInt)

	if shouldMoveCoins {
		err := k.ClaimSavingsStake(ctx, address, amount)
		if err != nil {
			return err
		}
	}

	stake := k.GetSavingsStake(ctx)
	k.SetSavingsStake(ctx, stake.Add(amountInt))

	return nil
}

func (k Keeper) distributeSavingsReward(ctx sdk.Context, amount sdk.Coins) bool {
	savingsStake := k.GetSavingsStake(ctx)

	if ! savingsStake.IsZero() && amount.AmountOf(config.DefaultDenom).IsPositive() {
		// S = S + r / T;

		rate := k.GetSavingsRewardRate(ctx)
		distribution := amount.AmountOf(config.DefaultDenom).ToDec()

		rate = rate.Add(
			distribution.QuoTruncate(
				savingsStake.ToDec(),
			),
		)

		k.SetSavingsRewardRate(ctx, rate)
		return true
	}

	return false
}

func (k Keeper) withdrawSavingsReward(ctx sdk.Context, address sdk.AccAddress, shouldMoveCoins bool) (sdk.DecCoins, error) {
	reward, found := k.calculateSavingsReward(ctx, address)
	if ! found {
		return reward, nil
	}

	userStake, _ := k.GetSavingsStakeByAddress(ctx, address)
	totalStake := k.GetSavingsStake(ctx)

	k.SetSavingsStake(ctx, totalStake.Sub(userStake))

	k.DeleteSavingsStakeByAddress(ctx, address)
	k.DeleteSavingsRewardRateByAddress(ctx, address)

	if shouldMoveCoins {
		coins := sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, userStake))
		err := k.RefundSavingsStake(ctx, address, coins)
		if err != nil {
			return sdk.DecCoins{}, err
		}
	}

	return reward, nil
}

func (k Keeper) calculateSavingsReward(ctx sdk.Context, address sdk.AccAddress) (sdk.DecCoins, bool) {
	deposit, found := k.GetSavingsStakeByAddress(ctx, address)
	if ! found {
		return sdk.DecCoins{}, false
	}

	rate := k.GetSavingsRewardRate(ctx)

	userRate, found := k.GetSavingsRewardRateByAddress(ctx, address)
	if ! found {
		return sdk.DecCoins{}, false
	}

	reward := deposit.ToDec().MulTruncate(
		rate.Sub(userRate),
	)

	return sdk.NewDecCoins(sdk.NewDecCoinFromDec(config.DefaultDenom, reward)), true
}

// Storage

func (k Keeper) GetSavingsStake(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetSavingsStakeKey())
	if bz == nil {
		panic("Savings stake should have been set")
	}

	var stake sdk.Int
	k.cdc.MustUnmarshalBinaryBare(bz, &stake)

	return stake
}

func (k Keeper) SetSavingsStake(ctx sdk.Context, stake sdk.Int) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetSavingsStakeKey(), k.cdc.MustMarshalBinaryBare(stake))
}

func (k Keeper) GetSavingsRewardRate(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetSavingsRewardRateKey())
	if bz == nil {
		panic("Savings reward rate should have been set")
	}

	var rate sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(bz, &rate)

	return rate
}

func (k Keeper) SetSavingsRewardRate(ctx sdk.Context, rate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetSavingsRewardRateKey(), k.cdc.MustMarshalBinaryBare(rate))
}

func (k Keeper) GetSavingsRewardRateByAddress(ctx sdk.Context, address sdk.AccAddress) (sdk.Dec, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetSavingsRewardByAddressKey(address))
	if bz == nil {
		return sdk.ZeroDec(), false
	}

	var rate sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(bz, &rate)

	return rate, true
}

func (k Keeper) SetSavingsRewardRateByAddress(ctx sdk.Context, address sdk.AccAddress, rate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetSavingsRewardByAddressKey(address), k.cdc.MustMarshalBinaryBare(rate))
}

func (k Keeper) DeleteSavingsRewardRateByAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetSavingsRewardByAddressKey(address))
}

func (k Keeper) IterateSavingsRewardRateByAddress(ctx sdk.Context, handler func(address sdk.AccAddress, rate sdk.Dec) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.SavingsRewardRateByAddressKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rate sdk.Dec
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &rate)
		address := types.GetSavingsRewardAddress(iter.Key())
		if handler(address, rate) {
			break
		}
	}
}

func (k Keeper) GetSavingsStakeByAddress(ctx sdk.Context, address sdk.AccAddress) (sdk.Int, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetSavingsStakeByAddressKey(address))
	if bz == nil {
		return sdk.ZeroInt(), false
	}

	var rate sdk.Int
	k.cdc.MustUnmarshalBinaryBare(bz, &rate)

	return rate, true
}

func (k Keeper) SetSavingsStakeByAddress(ctx sdk.Context, address sdk.AccAddress, rate sdk.Int) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetSavingsStakeByAddressKey(address), k.cdc.MustMarshalBinaryBare(rate))
}

func (k Keeper) DeleteSavingsStakeByAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetSavingsStakeByAddressKey(address))
}

func (k Keeper) IterateSavingsStakeByAddress(ctx sdk.Context, handler func(address sdk.AccAddress, rate sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.SavingsStakeByAddressKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rate sdk.Int
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &rate)
		address := types.GetSavingsStakeAddress(iter.Key())
		if handler(address, rate) {
			break
		}
	}
}

func (k Keeper) GetSavingsRewardEscrow(ctx sdk.Context, address sdk.AccAddress) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetSavingsRewardEscrowKey(address))
	if bz == nil {
		return sdk.NewDecCoins()
	}

	var amount sdk.DecCoins
	k.cdc.MustUnmarshalBinaryBare(bz, &amount)

	return amount
}

func (k Keeper) SetSavingsRewardEscrow(ctx sdk.Context, address sdk.AccAddress, amount sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetSavingsRewardEscrowKey(address), k.cdc.MustMarshalBinaryBare(amount))
}

func (k Keeper) DeleteSavingsRewardEscrow(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetSavingsRewardEscrowKey(address))
}

func (k Keeper) IterateSavingsRewardEscrow(ctx sdk.Context, handler func(address sdk.AccAddress, amount sdk.DecCoins) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.SavingsRewardEscrowKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var amount sdk.DecCoins
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &amount)
		address := types.GetSavingsRewardEscrowAddress(iter.Key())
		if handler(address, amount) {
			break
		}
	}
}

func (k Keeper) GetSavingsRewardLeftover(ctx sdk.Context, address sdk.AccAddress) (sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetSavingsRewardLeftoverKey(address))
	if bz == nil {
		return sdk.ZeroDec()
	}

	var leftover sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(bz, &leftover)

	return leftover
}

func (k Keeper) SetSavingsRewardLeftover(ctx sdk.Context, address sdk.AccAddress, leftover sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetNameRewardLeftoverKey(address), k.cdc.MustMarshalBinaryBare(leftover))
}

func (k Keeper) IterateSavingsRewardLeftover(ctx sdk.Context, handler func(address sdk.AccAddress, amount sdk.Dec) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.SavingsRewardLeftoverKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var amount sdk.Dec
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &amount)
		address := types.GetSavingsRewardLeftoverAddress(iter.Key())
		if handler(address, amount) {
			break
		}
	}
}