package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/distribution/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryParameters = "parameters"
	QueryNameReward = "name-reward"
	QueryValidatorReward = "validator-reward"
	QuerySavingsReward = "savings-reward"
	QuerySavings = "savings"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
			case QueryParameters:
				return queryParams(ctx, k)
			case QueryNameReward:
				return queryNameReward(ctx, path[1:], req, k)
			case QueryValidatorReward:
				return queryValidatorReward(ctx, path[1:], req, k)
			case QuerySavingsReward:
				return querySavingsReward(ctx,path[1:], req, k)
			case QuerySavings:
				return querySavings(ctx, path[1:], req, k)

			default:
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown distribution query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryNameReward(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	reward, _ := k.calculateNameReward(ctx, address)
	escrow := k.GetNameRewardEscrow(ctx, address)

	reward = reward.Add(escrow...)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, reward)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func querySavingsReward(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	reward, _ := k.calculateSavingsReward(ctx, address)
	escrow := k.GetSavingsRewardEscrow(ctx, address)

	reward = reward.Add(escrow...)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, reward)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func querySavings(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	savings, _ := k.GetSavingsStakeByAddress(ctx, address)

	coins := sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, savings))

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, coins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryValidatorReward(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.ValAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	reward := k.GetValidatorAccumulatedRewards(ctx, address)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, reward)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
