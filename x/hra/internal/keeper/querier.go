package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
)

const (
	QueryName = "resolve"
	QueryAddressNames = "names"
	QueryNameAddress = "address"
	QueryAddressNameInfos = "name-infos"
	QueryRegisteredBlockchainIds = "registered-blockchain-ids"
	QueryAddressCredits = "address-credits"
	QueryBlockchainAddresses = "blockchain-addresses"
	QueryParameters = "parameters"
	QueryModule = "module"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryName:
			return queryName(ctx, path[1:], req, k)
		case QueryAddressNames:
			return queryAddressNames(ctx, path[1:], req, k)
		case QueryNameAddress:
			return queryNameAddress(ctx, path[1:], req, k)
		case QueryAddressNameInfos:
			return queryAddressNameInfos(ctx, path[1:], req, k)
		case QueryRegisteredBlockchainIds:
			return queryRegisteredBlockchainIds(ctx, path[1:], req, k)
		case QueryAddressCredits:
			return queryAddressCredits(ctx, path[1:], req, k)
		case QueryBlockchainAddresses:
			return queryBlockchainAddresses(ctx, path[1:], req, k)
		case QueryParameters:
			return queryParams(ctx, k)
		case QueryModule:
			return queryModuleAccount(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unknown hra query endpoint: %s", path[0])
		}
	}
}

func queryName(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	nameInfo, found := k.GetNameInfo(ctx, path[0])
	if ! found {
		return nil, types.ErrNameNotRegistered
	}

	resNameInfo := types.QueryResNameInfo{
		NameInfo: nameInfo,
		Credits:  k.GetCredits(ctx, nameInfo.Owner),
	}

	k.IterateBlockchainAddressInfos(ctx, nameInfo.Owner, func (info types.BlockchainAddressInfo) (stop bool) {
		resNameInfo.Addresses = append(resNameInfo.Addresses, info)
		
		return false
	})

	res, marshalErr := codec.MarshalJSONIndent(types.ModuleCdc, resNameInfo)

	if marshalErr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, marshalErr.Error())
	}

	return res, nil
}

func queryAddressNames(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	iterator := k.GetNamesByAddressIterator(ctx, address)

	defer iterator.Close()

	var names types.QueryResNames

	for ; iterator.Valid(); iterator.Next() {
		key := string(iterator.Key())
		parts := strings.Split(key, types.Separator)
		name := strings.Join(parts[1:], "")
		names = append(names, name)
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, names)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryAddressNameInfos(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to parse the address.")
	}

	iterator := k.GetNamesByAddressIterator(ctx, address)

	defer iterator.Close()

	var names types.QueryResNameInfos

	for ; iterator.Valid(); iterator.Next() {
		key := string(iterator.Key())
		parts := strings.Split(key, types.Separator)
		name := strings.Join(parts[1:], "")

		nameInfo, found := k.GetNameInfo(ctx, name)
		if found {
			names = append(names, nameInfo)
		}

	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, names)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryNameAddress(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	nameInfo, found := k.GetNameInfo(ctx, path[0])
	if ! found {
		return nil, types.ErrNameNotRegistered
	}

	address, err := k.GetAddress(ctx, nameInfo.Owner, path[1], path[2])
	if err != nil {
		return nil, err
	}

	res, marshalErr := codec.MarshalJSONIndent(types.ModuleCdc, address)

	if marshalErr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, marshalErr.Error())
	}

	return res, nil
}

func queryRegisteredBlockchainIds(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {

	registeredBlockchainIds := k.GetRegisteredBlockchainIds(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, registeredBlockchainIds)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryAddressCredits(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	addressCredits := k.GetCredits(ctx, address)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, addressCredits.String())

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryBlockchainAddresses(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	var resBlockchainAddresses types.QueryResBlockchainAddresses

	k.IterateBlockchainAddressInfos(ctx, address, func (info types.BlockchainAddressInfo) (stop bool) {
		resBlockchainAddresses = append(resBlockchainAddresses, info)

		return false
	})

	res, marshalErr := codec.MarshalJSONIndent(types.ModuleCdc, resBlockchainAddresses)

	if marshalErr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, marshalErr.Error())
	}

	return res, nil
}

func queryModuleAccount(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	key := supply.NewModuleAddress(path[0])

	acc := k.AccountKeeper.GetAccount(ctx, key)

	res, marshalErr := codec.MarshalJSONIndent(types.ModuleCdc, acc)

	if marshalErr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, marshalErr.Error())
	}

	return res, nil
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}