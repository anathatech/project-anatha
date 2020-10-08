package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
)

// Withdrawals are done in full with re-depositing the Name stake
func (k Keeper) HandleWithdrawNameReward(ctx sdk.Context, address sdk.AccAddress) error {
	if ctx.BlockTime().Before(k.RewardWithdrawalEnabledTime(ctx)) {
		return types.ErrRewardWithdrawalDisabled
	}

	_, receivingRewards := k.GetNameRewardRateByAddress(ctx, address)

	// If user is not signed up for rewards no state changes will happen and reward will be 0
	reward := k.withdrawNameReward(ctx, address)

	// Not checking if reward is zero because escrow might exist
	err := k.DistributeNameReward(ctx, address, reward)
	if err != nil {
		return err
	}

	// If the user was signed up for rewards the withdrawNameReward reset the reward state and we need to deposit the name again
	if receivingRewards {
		k.depositName(ctx, address)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawNameReward,
			sdk.NewAttribute(types.AttributeKeyReward, k.GetSavingsRewardEscrow(ctx, address).Add(reward...).String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, address.String()),
		),
	})

	return nil
}

// Algorithm

func (k Keeper) depositName(ctx sdk.Context, address sdk.AccAddress) {
	rate := k.GetNameRewardRate(ctx)
	k.SetNameRewardRateByAddress(ctx, address, rate)

	k.IncreaseNameStake(ctx)
}

func (k Keeper) distributeNameReward(ctx sdk.Context, amount sdk.Coins) bool {
	nameStake := k.GetNameStake(ctx)

	if ! nameStake.IsZero() && amount.AmountOf(config.DefaultDenom).IsPositive() {
		// S = S + r / T;

		rate := k.GetNameRewardRate(ctx)
		distribution := amount.AmountOf(config.DefaultDenom).ToDec()

		rate = rate.Add(
			distribution.QuoTruncate(nameStake),
		)

		k.SetNameRewardRate(ctx, rate)
		return true
	}

	return false
}

func (k Keeper) withdrawNameReward(ctx sdk.Context, address sdk.AccAddress) (amount sdk.DecCoins) {
	reward, found := k.calculateNameReward(ctx, address)

	if found {
		k.DecreaseNameStake(ctx)

		k.DeleteNameRewardRateByAddress(ctx, address)
	}

	return reward
}

func (k Keeper) calculateNameReward(ctx sdk.Context, address sdk.AccAddress) (sdk.DecCoins, bool) {
	deposit := sdk.OneDec() // one share per HRA holder
	rate := k.GetNameRewardRate(ctx)

	userRate, found := k.GetNameRewardRateByAddress(ctx, address)
	if ! found {
		return sdk.DecCoins{}, false
	}

	reward := deposit.MulTruncate(rate.Sub(userRate))

	return sdk.NewDecCoins(sdk.NewDecCoinFromDec(config.DefaultDenom, reward)), true
}

// Util

func (k Keeper) NameDeposited(ctx sdk.Context, address sdk.AccAddress) bool {
	_, found := k.GetNameRewardRateByAddress(ctx, address)

	return found
}

func (k Keeper) DepositName(ctx sdk.Context, address sdk.AccAddress) (ok bool) {
	if ! k.NameDeposited(ctx, address) {
		k.depositName(ctx, address)
		ok = true
	} else {
		ok = false
	}

	k.Logger(ctx).Debug(
		fmt.Sprintf("Name Deposited Status %s %s", address, ok),
	)

	return
}

// Storage

func (k Keeper) GetNameStake(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetNameStakeKey())
	if bz == nil {
		panic("Name stake should have been set")
	}

	var stake sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(bz, &stake)

	return stake
}

func (k Keeper) IncreaseNameStake(ctx sdk.Context) {
	nameStake := k.GetNameStake(ctx)
	k.SetNameStake(ctx, nameStake.Add(sdk.OneDec()))
}

func (k Keeper) DecreaseNameStake(ctx sdk.Context) {
	nameStake := k.GetNameStake(ctx)
	k.SetNameStake(ctx, nameStake.Sub(sdk.OneDec()))
}

func (k Keeper) SetNameStake(ctx sdk.Context, stake sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetNameStakeKey(), k.cdc.MustMarshalBinaryBare(stake))
}

func (k Keeper) GetNameRewardRate(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetNameRewardRateKey())
	if bz == nil {
		panic("Name reward rate should have been set")
	}

	var rate sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(bz, &rate)

	return rate
}

func (k Keeper) SetNameRewardRate(ctx sdk.Context, rate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetNameRewardRateKey(), k.cdc.MustMarshalBinaryBare(rate))
}

func (k Keeper) GetNameRewardRateByAddress(ctx sdk.Context, address sdk.AccAddress) (sdk.Dec, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetNameRewardByAddressKey(address))
	if bz == nil {
		return sdk.ZeroDec(), false
	}

	var rate sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(bz, &rate)

	return rate, true
}

func (k Keeper) SetNameRewardRateByAddress(ctx sdk.Context, address sdk.AccAddress, rate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetNameRewardByAddressKey(address), k.cdc.MustMarshalBinaryBare(rate))
}

func (k Keeper) DeleteNameRewardRateByAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetNameRewardByAddressKey(address))
}

func (k Keeper) IterateNameRewardRateByAddress(ctx sdk.Context, handler func(address sdk.AccAddress, rate sdk.Dec) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.NameRewardRateByAddressKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rate sdk.Dec
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &rate)
		address := types.GetNameRewardAddress(iter.Key())
		if handler(address, rate) {
			break
		}
	}
}

func (k Keeper) GetPendingNameDistribution(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetPendingNameDistributionKey())
	if bz == nil {
		return sdk.NewCoins()
	}

	var distribution sdk.Coins
	k.cdc.MustUnmarshalBinaryBare(bz, &distribution)

	return distribution
}

func (k Keeper) SetPendingNameDistribution(ctx sdk.Context, distribution sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	if ! distribution.IsZero() {
		store.Set(types.GetPendingNameDistributionKey(), k.cdc.MustMarshalBinaryBare(distribution))
	} else {
		store.Delete(types.GetPendingNameDistributionKey())
	}
}

func (k Keeper) DeletePendingNameDistribution(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetPendingNameDistributionKey())
}

func (k Keeper) GetNameRewardEscrow(ctx sdk.Context, address sdk.AccAddress) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetNameRewardEscrowKey(address))
	if bz == nil {
		return sdk.NewDecCoins()
	}

	var amount sdk.DecCoins
	k.cdc.MustUnmarshalBinaryBare(bz, &amount)

	return amount
}

func (k Keeper) SetNameRewardEscrow(ctx sdk.Context, address sdk.AccAddress, amount sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetNameRewardEscrowKey(address), k.cdc.MustMarshalBinaryBare(amount))
}

func (k Keeper) DeleteNameRewardEscrow(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetNameRewardEscrowKey(address))
}

func (k Keeper) IterateNameRewardEscrow(ctx sdk.Context, handler func(address sdk.AccAddress, amount sdk.DecCoins) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.NameRewardEscrowKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var amount sdk.DecCoins
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &amount)
		address := types.GetNameRewardEscrowAddress(iter.Key())
		if handler(address, amount) {
			break
		}
	}
}

func (k Keeper) GetNameRewardLeftover(ctx sdk.Context, address sdk.AccAddress) (sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetNameRewardLeftoverKey(address))
	if bz == nil {
		return sdk.ZeroDec()
	}

	var leftover sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(bz, &leftover)

	return leftover
}

func (k Keeper) SetNameRewardLeftover(ctx sdk.Context, address sdk.AccAddress, leftover sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetNameRewardLeftoverKey(address), k.cdc.MustMarshalBinaryBare(leftover))
}

func (k Keeper) IterateNameRewardLeftover(ctx sdk.Context, handler func(address sdk.AccAddress, amount sdk.Dec) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.NameRewardLeftoverKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var amount sdk.Dec
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &amount)
		address := types.GetNameRewardLeftoverAddress(iter.Key())
		if handler(address, amount) {
			break
		}
	}
}