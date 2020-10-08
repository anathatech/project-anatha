package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.Operators {
		k.AddOperator(ctx, record)
	}

	// devnet TODO: Temporary for testing purposes
	err := k.MintCoinsToBuyBackLiquidityFund(ctx, sdk.NewCoins(sdk.NewInt64Coin(config.DefaultStableDenom, 10000000000000)))
	if err != nil {
		panic(err)
	}

	if data.Treasury.MintGenesisSupply {
		targetSupply := data.Treasury.TargetSupply

		err := k.MintCoins(ctx, targetSupply)
		if err != nil {
			panic(err)
		}

		data.Treasury.MintGenesisSupply = false
	}

	k.SetTreasury(ctx, data.Treasury)
	k.SetParams(ctx, data.Params)

	for _, disbursement := range data.DisbursementQueue {
		k.InsertDisbursementQueue(ctx, disbursement)
	}

	for _, disbursementReference := range data.DisbursementReferences {
		k.SetDisbursementReferenceAmount(ctx, disbursementReference.Reference, disbursementReference.Amount)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	treasury := k.GetTreasury(ctx)
	params := k.GetParams(ctx)
	operators := k.GetOperators(ctx)
	disbursements := k.GetDisbursements(ctx)

	var disbursementReferences []types.ReferenceAmountInfo
	k.IterateDisbursementReferences(ctx, func(reference string, amount sdk.Int) (stop bool) {
		disbursementReferences = append(disbursementReferences, types.NewReferenceAmountInfo(reference, amount))

		return false
	})

	return NewGenesisState(treasury, params, operators, disbursements, disbursementReferences)
}
