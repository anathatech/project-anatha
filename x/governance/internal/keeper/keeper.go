package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/anathatech/project-anatha/x/governance/internal/types"
	"time"

	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramSpace params.Subspace
	router gov.Router

	AccountKeeper 	auth.AccountKeeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, router gov.Router, accountKeeper auth.AccountKeeper,) Keeper {
	return Keeper{
		storeKey:     key,
		paramSpace:   paramSpace,
		cdc:          cdc,
		router:       router,
		AccountKeeper: accountKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Router() gov.Router {
	return k.router
}

// Proposal Queue

func (k Keeper) InsertActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := types.GetProposalIDBytes(proposalID)
	store.Set(types.ActiveProposalQueueKey(proposalID, endTime), bz)
}

func (k Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ActiveProposalQueueKey(proposalID, endTime))
}

func (k Keeper) InsertExpeditedProposalQueue(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := types.GetProposalIDBytes(proposalID)
	store.Set(types.ExpeditedProposalQueueKey(proposalID), bz)
}

func (k Keeper) RemoveFromExpeditedProposalQueue(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ExpeditedProposalQueueKey(proposalID))
}

// Iterators

func (k Keeper) IterateActiveProposalsQueue(ctx sdk.Context, endTime time.Time, cb func(proposal types.Proposal) (stop bool)) {
	iterator := k.ActiveProposalQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := types.SplitActiveProposalQueueKey(iterator.Key())
		proposal, err := k.GetProposal(ctx, proposalID)
		if err != nil {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

func (k Keeper) ActiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.ActiveProposalQueuePrefix, sdk.PrefixEndBytes(types.ActiveProposalByTimeKey(endTime)))
}

func (k Keeper) IterateExpeditedProposalsQueue(ctx sdk.Context, cb func(proposal types.Proposal) (stop bool)) {
	iterator := k.ExpeditedProposalQueueIterator(ctx)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID := types.SplitExpeditedProposalQueueKey(iterator.Key())
		proposal, err := k.GetProposal(ctx, proposalID)
		if err != nil {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

func (k Keeper) ExpeditedProposalQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.ExpeditedProposalQueuePrefix)
}