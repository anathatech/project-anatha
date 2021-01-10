package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
)

func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		return nil
	}

	return k.supplyKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

func (k Keeper) MintCoinsToBuyBackLiquidityFund(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		return nil
	}

	return k.supplyKeeper.MintCoins(ctx, types.BuyBackLiquidityFundModuleName, newCoins)
}

func (k Keeper) BurnCoinsFromBuyBackLiquidityFund(ctx sdk.Context, amount sdk.Coins) error {
	if amount.Empty() {
		return nil
	}

	return k.supplyKeeper.BurnCoins(ctx, types.BuyBackLiquidityFundModuleName, amount)
}


func (k Keeper) TotalSupply(ctx sdk.Context) sdk.Coins {
	return k.supplyKeeper.GetSupply(ctx).GetTotal()
}

func (k Keeper) TransferFromBuyBackFund(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.Coins) error {
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.BuyBackLiquidityFundModuleName, recipient , amount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) TransferToBuyBackFund(ctx sdk.Context, sender sdk.AccAddress, amount sdk.Coins) error {
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.BuyBackFundModuleName, amount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) BurnFromDistributionProfits(ctx sdk.Context, amount sdk.Coins) error {
	err := k.supplyKeeper.BurnCoins(ctx, types.DistributionProfitsModuleName, amount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) TransferFromDistributionProfitsToBuyBackLiquidity(ctx sdk.Context, amount sdk.Coins) error {
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.DistributionProfitsModuleName, types.BuyBackLiquidityFundModuleName, amount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) TransferFromTreasuryToSwapEscrow(ctx sdk.Context, amount sdk.Coins) error {
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.SwapEscrowModuleName, amount)
	if err != nil {
		return err
	}

	treasury := k.GetTreasury(ctx)
	treasury.Distributed = treasury.Distributed.Add(amount...)
	k.SetTreasury(ctx, treasury)

	return nil
}

func (k Keeper) TransferFromSwapEscrowToBuyBack(ctx sdk.Context, amount sdk.Coins) error {
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.SwapEscrowModuleName, types.BuyBackFundModuleName, amount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) TransferFromSwapEscrow(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.Coins) error {
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.SwapEscrowModuleName, recipient, amount)
	if err != nil {
		return err
	}

	return nil
}