package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/mint/internal/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	if ctx.BlockHeight() == 1 {
		minter.PreviousMintTime = ctx.BlockTime()
		k.SetMinter(ctx, minter)

		return
	}

	currentBlockTime := ctx.BlockTime()
	previousBlockTime := minter.PreviousMintTime

	if previousBlockTime.Equal(currentBlockTime) || previousBlockTime.After(currentBlockTime) {
		panic("Block Time Sequence Invalid")
	}

	elapsed := int(currentBlockTime.Sub(previousBlockTime).Seconds())

	totalSupply := k.TotalSupply(ctx)

	inflationPerSecond := params.PerSecondInflationRate

	totalSupplyDec := totalSupply.AmountOf(config.DefaultDenom).ToDec()

	toMintDec := minter.Leftover // initialize minting amount with previous iterations leftover

	for i := 0; i < elapsed; i++ {

		currentSecondInflationAmountDec := totalSupplyDec.MulTruncate(inflationPerSecond)

		toMintDec = toMintDec.Add(currentSecondInflationAmountDec)
		totalSupplyDec = totalSupplyDec.Add(currentSecondInflationAmountDec)

	}

	toMintInt := toMintDec.TruncateInt()
	leftover := toMintDec.Sub(toMintInt.ToDec())

	minter.PreviousMintTime = currentBlockTime
	minter.Leftover = leftover

	k.SetMinter(ctx, minter)

	toMintCoins := sdk.NewCoins(sdk.NewCoin(params.MintDenom, toMintInt))

	err := k.MintCoins(ctx, toMintCoins)
	if err != nil {
		panic(err)
	}

	err = k.DistributeCollectedFees(ctx, toMintCoins)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(sdk.AttributeKeyAmount, toMintCoins.String()),
		),
	)
}
