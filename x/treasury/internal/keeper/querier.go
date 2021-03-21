package keeper

import (
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/utils"
	abci "github.com/tendermint/tendermint/abci/types"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
)

const (
	QueryParameters = "parameters"
	QueryTreasury = "treasury"
	QueryOperators = "operators"
	QueryDisbursements = "disbursements"
	QueryPrice = "price"
	QueryDisbursementEscrow = "disbursement-escrow"
)

// NewQuerier creates a new querier for treasury clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryParameters:
			return queryParams(ctx, k)
		case QueryTreasury:
			return queryTreasury(ctx, k)
		case QueryOperators:
			return queryOperators(ctx, k)
		case QueryDisbursements:
			return queryDisbursements(ctx, k)
		case QueryPrice:
			return queryPrice(ctx, path[1:], req, k)
		case QueryDisbursementEscrow:
			return queryDisbursementEscrow(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown treasury query endpoint")
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

func queryTreasury(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetTreasury(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryOperators(ctx sdk.Context, k Keeper) ([]byte, error) {
	operators := k.GetOperators(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, operators)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDisbursements(ctx sdk.Context, k Keeper) ([]byte, error) {
	disbursements := k.GetDisbursements(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, disbursements)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryPrice(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	coins, err := utils.ParseAndConvertCoins(path[0])
	if err != nil {
		return nil, err
	}

	var result sdk.Coins

	if coins.AmountOf(config.DefaultStableDenom).IsPositive() {

		result, _, _ = k.CalculatePinAmountExtended(ctx, coins)

	} else if coins.AmountOf(config.DefaultDenom).IsPositive() {

		result = k.CalculateDinAmountExtended(ctx, coins)

	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid denomination. Expected fiat or anatha.")
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, result)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDisbursementEscrow(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	reference := strings.ToLower(path[0])

	found, amount := k.GetDisbursementReferenceAmount(ctx, reference)

	if ! found {
		return nil, types.ErrInvalidReference
	}

	coins := sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, amount))

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, coins)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}