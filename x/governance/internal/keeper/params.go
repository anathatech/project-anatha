package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/governance/internal/types"
)

func (k Keeper) GetVotingParams(ctx sdk.Context) types.VotingParams {
	var votingParams types.VotingParams
	k.paramSpace.Get(ctx, types.ParamStoreKeyVotingParams, &votingParams)
	return votingParams
}

func (k Keeper) GetTallyParams(ctx sdk.Context) types.TallyParams {
	var tallyParams types.TallyParams
	k.paramSpace.Get(ctx, types.ParamStoreKeyTallyParams, &tallyParams)
	return tallyParams
}

func (k Keeper) SetVotingParams(ctx sdk.Context, votingParams types.VotingParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyVotingParams, &votingParams)
}

func (k Keeper) SetTallyParams(ctx sdk.Context, tallyParams types.TallyParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyTallyParams, &tallyParams)
}