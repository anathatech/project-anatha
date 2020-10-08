package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
	"time"
)

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) NameDepositDelay(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.KeyNameDepositDelay, &res)
	return
}

func (k Keeper) RewardWithdrawalBlockedPeriod(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.KeyRewardWithdrawalBlockedPeriod, &res)
	return
}

func (k Keeper) SetRewardWithdrawalEnabledTime(ctx sdk.Context, time time.Time) {
	k.paramSpace.Set(ctx, types.KeyRewardWithdrawalEnabledTime, time)
}

func (k Keeper) RewardWithdrawalEnabledTime(ctx sdk.Context) (res time.Time) {
	k.paramSpace.Get(ctx, types.KeyRewardWithdrawalEnabledTime, &res)
	return
}

func (k Keeper) SavingsSplitAdjustment(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeySavingsSplitAdjustment, &res)
	return
}

func (k Keeper) DevelopmentFundShare(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyDevelopmentFundShare, &res)
	return
}

func (k Keeper) SecurityTokenFundShare(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeySecurityTokenFundShare, &res)
	return
}