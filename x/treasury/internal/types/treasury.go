package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Treasury struct {
	MintGenesisSupply 	bool 		`json:"mint_genesis_supply" yaml:"mint_genesis_supply"`
	TargetSupply 		sdk.Coins 	`json:"target_supply" yaml:"target_supply"`
	Distributed			sdk.Coins	`json:"distributed" yaml:"distributed"`
	CoinsPerStage		sdk.Coins 	`json:"coins_per_stage" yaml:"coins_per_stage"`
	InboundDin			sdk.Coins	`json:"inbound_din" yaml:"inbound_din"`
}

func NewTreasury(mintGenesisSupply bool, targetSupply sdk.Coins, distributed sdk.Coins, coinsPerStage sdk.Coins, inboundDin sdk.Coins) Treasury {
	return Treasury{
		MintGenesisSupply: 	mintGenesisSupply,
		TargetSupply: 		targetSupply,
		Distributed:		distributed,
		CoinsPerStage: 		coinsPerStage,
		InboundDin: 		inboundDin,
	}
}

func InitialTreasury(mintGenesisSupply bool, targetSupply sdk.Coins, distributed sdk.Coins, coinsPerStage sdk.Coins, inboundDin sdk.Coins) Treasury {
	return NewTreasury(
		mintGenesisSupply,
		targetSupply,
		distributed,
		coinsPerStage,
		inboundDin,
	)
}

func DefaultInitialTreasury() Treasury {
	targetSupply, _ := sdk.ConvertCoin(sdk.NewInt64Coin("anatha", 6909275134), "pin")
	coinsPerStage, _ := sdk.ConvertCoin(sdk.NewInt64Coin("anatha", 10000000), "pin")

	return InitialTreasury(
		true,
		sdk.NewCoins(targetSupply),
		sdk.NewCoins(),
		sdk.NewCoins(coinsPerStage),
		sdk.NewCoins(),
	)
}
