package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/fee/internal/types"
)

// FeePercentage
func (k Keeper) FeePercentage(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, types.KeyFeePercentage, &res)
	return
}

// MinimumFee
func (k Keeper) MinimumFee(ctx sdk.Context) (res sdk.Coins) {
	k.paramspace.Get(ctx, types.KeyMinimumFee, &res)
	return
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramspace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramspace.SetParamSet(ctx, &params)
}
