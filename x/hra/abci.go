package hra

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
}

func EndBlocker(ctx sdk.Context, k Keeper) {
	k.IterateExpiredNameInfoQueue(ctx, ctx.BlockTime(), func(name string) (stop bool) {
		nameInfo, found := k.GetNameInfo(ctx, name)
		if ! found {
			panic(fmt.Sprintf("name info %s does not exist", name))
		}

		err := k.DeleteExpiredNameInfo(ctx, nameInfo)
		if err != nil {
			panic("error deleting expired name info")
		}

		return false
	})
}