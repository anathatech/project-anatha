package keeper

import	(
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
	"time"
)

func (k Keeper) Managers(ctx sdk.Context) (res []sdk.AccAddress) {
	k.paramspace.Get(ctx, types.KeyManagers, &res)
	return
}

func (k Keeper) RiskAssessmentAmount(ctx sdk.Context) (res sdk.Coins) {
	k.paramspace.Get(ctx, types.KeyRiskAssessmentAmount, &res)
	return
}

func (k Keeper) RiskAssessmentDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, types.KeyRiskAssessmentDuration, &res)
	return
}

func (k Keeper) BuyBackPercentage(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, types.KeyBuyBackPercentage, &res)
	return
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramspace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramspace.SetParamSet(ctx, &params)
}
