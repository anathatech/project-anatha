package distribution

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"time"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {

	// AMC Rewards

	toHhrm, _, _ := k.DistributeFromAmc(ctx)

	k.DistributeFromHhrm(ctx, toHhrm)

	k.IterateNameDepositQueueByTime(ctx, ctx.BlockTime(), func(address sdk.AccAddress, endTime time.Time) (stop bool) {
		k.Logger(ctx).Debug(
			fmt.Sprintf("Working on Scheduled Name Deposit: %s %s", address, endTime),
		)

		k.DepositName(ctx, address)
		k.RemoveFromNameDepositQueue(ctx, address, endTime)

		return false
	})

	// NVRP Rewards

	// Distribution from NVRP to Savers
	k.DistributeFromNvrp(ctx)

	// Distribution from NVRP to Validator Rewards
	if ctx.BlockHeight() > 1 {
		k.AllocateTokens(ctx, req.LastCommitInfo.GetVotes())
	}
}