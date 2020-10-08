package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
)

func (k Keeper) HandleWithdrawValidatorReward(ctx sdk.Context, valAddr sdk.ValAddress) error {
	if ctx.BlockTime().Before(k.RewardWithdrawalEnabledTime(ctx)) {
		return types.ErrRewardWithdrawalDisabled
	}

	accumulatedRewards := k.GetValidatorAccumulatedRewards(ctx, valAddr)
	rewards, remainder := accumulatedRewards.TruncateDecimal()

	if rewards.IsZero() {
		return types.ErrNoValidatorRewards
	}

	k.SetValidatorAccumulatedRewards(ctx, valAddr, remainder)

	if ! rewards.IsZero() {
		accAddr := sdk.AccAddress(valAddr)
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.NvrpDistributionModuleName, accAddr, rewards)
		if err != nil {
			return err
		}
		k.Logger(ctx).Debug(
			fmt.Sprintf("nvrpd -> %s : %s", accAddr, rewards),
		)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawValidatorReward,
			sdk.NewAttribute(types.AttributeKeyReward, rewards.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, sdk.AccAddress(valAddr).String()),
		),
	})

	return nil
}
