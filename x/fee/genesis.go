package fee

import sdk "github.com/cosmos/cosmos-sdk/types"

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetParams(ctx, data.Params)

	for _, record := range data.FeeExcludedMessages {
		k.SetFeeExcludedMessage(ctx, record)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)

	return GenesisState{
		Params: params,
		FeeExcludedMessages: k.GetFeeExcludedMessages(ctx),
	}
}