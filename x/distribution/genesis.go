package distribution

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
	"time"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	// We are adding the RewardWithdrawalBlockedPeriod to the genesis timestamp to set when withdrawals become enabled
	// devnet TODO: remove for production and set manually RewardWithdrawalEnabledTime
	keeper.SetRewardWithdrawalEnabledTime(ctx,
		ctx.BlockTime().Add(
			keeper.RewardWithdrawalBlockedPeriod(ctx),
		),
	)

	keeper.SetNameRewardRate(ctx, data.NameRewardRate)
	keeper.SetNameStake(ctx, data.NameStake)

	for _, anrr := range data.AddressNameRewardRates {
		keeper.SetNameRewardRateByAddress(ctx, anrr.Address, anrr.Rate)
	}

	keeper.SetPendingNameDistribution(ctx, data.PendingNameDistribution)

	for _, ndq := range data.NameDepositQueue {
		keeper.InsertNameDepositQueue(ctx, ndq.Address, ndq.Time)
	}

	for _, nre := range data.NameRewardEscrow {
		keeper.SetNameRewardEscrow(ctx, nre.Address, nre.Amount)
	}

	for _, nrl := range data.NameRewardLeftover {
		keeper.SetNameRewardLeftover(ctx, nrl.Address, nrl.Amount)
	}

	keeper.SetSavingsRewardRate(ctx, data.SavingsRewardRate)
	keeper.SetSavingsStake(ctx, data.SavingsStake)

	for _, asrr := range data.AddressSavingsRewardRates {
		keeper.SetSavingsRewardRateByAddress(ctx, asrr.Address, asrr.Rate)
	}

	for _, ass := range data.AddressSavingsStake {
		keeper.SetSavingsStakeByAddress(ctx, ass.Address, ass.Stake)
	}

	for _, sre := range data.SavingsRewardEscrow {
		keeper.SetSavingsRewardEscrow(ctx, sre.Address, sre.Amount)
	}

	for _, srl := range data.SavingsRewardLeftover {
		keeper.SetSavingsRewardLeftover(ctx, srl.Address, srl.Amount)
	}

	for _, rew := range data.ValidatorAccumulatedRewards {
		keeper.SetValidatorAccumulatedRewards(ctx, rew.ValidatorAddress, rew.Accumulated)
	}

	keeper.SetNvrpRemainder(ctx, data.NvrpRemainder)
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)

	validatorRewards := make([]types.ValidatorAccumulatedRewardRecord, 0)
	keeper.IterateValidatorAccumulatedRewards(ctx, func(address sdk.ValAddress, accumulated sdk.DecCoins) (stop bool) {
		validatorRewards = append(validatorRewards, types.ValidatorAccumulatedRewardRecord{
			ValidatorAddress: address,
			Accumulated: accumulated,
		})
		return false
	})

	addressNameRewardRates := make([]types.AddressNameRewardRateRecord, 0)
	keeper.IterateNameRewardRateByAddress(ctx, func(address sdk.AccAddress, rate sdk.Dec) (stop bool) {
		addressNameRewardRates = append(addressNameRewardRates, types.AddressNameRewardRateRecord{
			Address: address,
			Rate:    rate,
		})
		return false
	})

	nameDepositQueue := make([]types.NameDepositQueueRecord, 0)
	keeper.IterateNameDepositQueue(ctx, func(address sdk.AccAddress, endTime time.Time) (stop bool) {
		nameDepositQueue = append(nameDepositQueue, types.NameDepositQueueRecord{
			Address: address,
			Time:    endTime,
		})
		return false
	})

	nameRewardEscrow := make([]types.NameRewardEscrowRecord, 0)
	keeper.IterateNameRewardEscrow(ctx, func(address sdk.AccAddress, amount sdk.DecCoins) (stop bool) {
		nameRewardEscrow = append(nameRewardEscrow, types.NameRewardEscrowRecord{
			Address: address,
			Amount:  amount,
		})
		return false
	})

	nameRewardLeftover := make([]types.NameRewardLeftoverRecord, 0)
	keeper.IterateNameRewardLeftover(ctx, func(address sdk.AccAddress, amount sdk.Dec) (stop bool) {
		nameRewardLeftover = append(nameRewardLeftover, types.NameRewardLeftoverRecord{
			Address: address,
			Amount: amount,
		})
		return false
	})

	addressSavingsRewardRates := make([]types.AddressSavingsRewardRateRecord, 0)
	keeper.IterateSavingsRewardRateByAddress(ctx, func(address sdk.AccAddress, rate sdk.Dec) (stop bool) {
		addressSavingsRewardRates = append(addressSavingsRewardRates, types.AddressSavingsRewardRateRecord{
			Address: address,
			Rate:    rate,
		})
		return false
	})

	addressSavingsStake := make([]types.AddressSavingsStakeRecord, 0)
	keeper.IterateSavingsStakeByAddress(ctx, func(address sdk.AccAddress, rate sdk.Int) (stop bool) {
		addressSavingsStake = append(addressSavingsStake, types.AddressSavingsStakeRecord{
			Address: address,
			Stake:   rate,
		})
		return false
	})


	savingsRewardEscrow := make([]types.SavingsRewardEscrowRecord, 0)
	keeper.IterateSavingsRewardEscrow(ctx, func(address sdk.AccAddress, amount sdk.DecCoins) (stop bool) {
		savingsRewardEscrow = append(savingsRewardEscrow, types.SavingsRewardEscrowRecord{
			Address: address,
			Amount:  amount,
		})
		return false
	})

	savingsRewardLeftover := make([]types.SavingsRewardLeftoverRecord, 0)
	keeper.IterateSavingsRewardLeftover(ctx, func(address sdk.AccAddress, amount sdk.Dec) (stop bool) {
		savingsRewardLeftover = append(savingsRewardLeftover, types.SavingsRewardLeftoverRecord{
			Address: address,
			Amount: amount,
		})
		return false
	})

	return NewGenesisState(
		params,
		keeper.GetNameStake(ctx),
		keeper.GetNameRewardRate(ctx),
		addressNameRewardRates,
		keeper.GetPendingNameDistribution(ctx),
		nameDepositQueue,
		nameRewardEscrow,
		nameRewardLeftover,
		keeper.GetSavingsStake(ctx),
		keeper.GetSavingsRewardRate(ctx),
		addressSavingsRewardRates,
		addressSavingsStake,
		savingsRewardEscrow,
		savingsRewardLeftover,
		validatorRewards,
		keeper.GetNvrpRemainder(ctx),
	)
}
