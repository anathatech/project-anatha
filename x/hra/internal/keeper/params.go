package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
	"time"
)


// NameInfoDuration
func (k Keeper) NameInfoDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, types.KeyNameInfoDuration, &res)
	return
}

// NameInfoDuration
func (k Keeper) NameInfoMaxDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, types.KeyNameInfoMaxDuration, &res)
	return
}

// NameInfoRegistrationFee
func (k Keeper) NameInfoRegistrationFee(ctx sdk.Context) (res sdk.Coins) {
	k.paramspace.Get(ctx, types.KeyNameInfoRegistrationFee, &res)
	return
}

// NameInfoRenewalFee
func (k Keeper) NameInfoRenewalFee(ctx sdk.Context) (res sdk.Coins) {
	k.paramspace.Get(ctx, types.KeyNameInfoRenewalFee, &res)
	return
}

// AddressCredits
func (k Keeper) AddressCredits(ctx sdk.Context) (res sdk.Int) {
	k.paramspace.Get(ctx, types.KeyAddressCredits, &res)
	return
}

// AddressRegistrationFee
func (k Keeper) AddressRegistrationFee(ctx sdk.Context) (res sdk.Coins) {
	k.paramspace.Get(ctx, types.KeyAddressRegistrationFee, &res)
	return
}

// GetParams returns the total set of hra parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramspace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the hra parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramspace.SetParamSet(ctx, &params)
}