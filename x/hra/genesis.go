package hra

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)

	for _, record := range data.RegisteredBlockchainIds {
		keeper.SetRegisteredBlockchainId(ctx, record)
	}

	for _, record := range data.NameRecords {
		keeper.SetNameInfo(ctx, record.Name, record)

		keeper.InsertExpiredNameInfoQueue(ctx, record.Name, record.ExpiryTime)
	}

	for _, record := range data.AddressCredits {
		keeper.SetCredits(ctx, record.Address, record.Credits)
	}

	for _, record := range data.AddressRecords {
		keeper.SetAddress(ctx, record.Address, record.BlockchainAddressInfo.BlockchainId, record.BlockchainAddressInfo.Index, record.BlockchainAddressInfo.BlockchainAddress)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var nameInfos []types.NameInfo
	iterator := k.GetNamesIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		nameInfo, found := k.GetNameInfo(ctx, string(iterator.Key()[1:]))
		if found {
			nameInfos = append(nameInfos, nameInfo)
		}
	}

	params := k.GetParams(ctx)

	var addressCredits []types.AddressCreditsInfo
	k.IterateAddressCredits(ctx, func (address sdk.AccAddress, credits sdk.Int) (stop bool) {
		addressCredits = append(addressCredits, types.NewAddressCreditsInfo(address, credits))

		return false
	})

	var addressRecords []types.BlockchainAddressRecordInfo
	k.IterateAllBlockchainAddressInfos(ctx, func (addressRecordInfo types.BlockchainAddressRecordInfo) (stop bool) {
		addressRecords = append(addressRecords, addressRecordInfo)

		return false
	})

	return GenesisState{
		Params:      params,
		NameRecords: nameInfos,
		AddressCredits: addressCredits,
		AddressRecords: addressRecords,
		RegisteredBlockchainIds: k.GetRegisteredBlockchainIds(ctx),
	}
}