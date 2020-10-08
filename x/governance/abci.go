package governance

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, keeper Keeper) {
	// Iterate expedited proposals that can't be overturned
	keeper.IterateExpeditedProposalsQueue(ctx, func(proposal Proposal) (stop bool) {
		keeper.HandleProposal(ctx, proposal, true)

		return false
	})

	// iterate proposals that reached the end of their voting period
	keeper.IterateActiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal Proposal) bool {
		keeper.HandleProposal(ctx, proposal, false)

		return false
	})
}
