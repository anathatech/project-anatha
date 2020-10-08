package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type ParamSubspace interface {
	WithKeyTable(table params.KeyTable) params.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

type NameHooks interface {
	AfterFirstNameCreated(ctx sdk.Context, address sdk.AccAddress) error
	AfterLastNameRemoved(ctx sdk.Context, address sdk.AccAddress) error
}