package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/x/governance/internal/types"
	"strconv"
)

func (k Keeper) AddVote(ctx sdk.Context, proposalID uint64, voter sdk.AccAddress, option types.VoteOption) error {
	if ! k.IsGovernor(ctx, voter) {
		return sdkerrors.Wrap(types.ErrNotGovernor, voter.String())
	}

	proposal, err := k.GetProposal(ctx, proposalID)

	if err != nil {
		return sdkerrors.Wrapf(types.ErrUnknownProposal, "%d", proposalID)
	}
	if proposal.Status != types.StatusVotingPeriod {
		return sdkerrors.Wrapf(types.ErrInactiveProposal, "%d", proposalID)
	}
	if k.HasVoted(ctx, proposalID, voter) {
		return sdkerrors.Wrapf(types.ErrAlreadyVoted, "%d - %s", proposalID, voter)
	}

	switch option {
		case types.OptionYes:
			proposal.TallyResult.Yes = proposal.TallyResult.Yes.Add(sdk.OneInt())
		case types.OptionNo:
			proposal.TallyResult.No = proposal.TallyResult.No.Add(sdk.OneInt())
		default:
			return sdkerrors.Wrap(types.ErrInvalidVote, option.String())
	}

	vote := types.NewVote(proposalID, voter, option)

	k.SetVote(ctx, vote)

	if k.CanBeExpedited(ctx, proposal) {
		proposal.Status = types.StatusExpediting
		proposal.Expedited = true

		k.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)
		k.InsertExpeditedProposalQueue(ctx, proposal.ProposalID)
		k.Logger(ctx).Info(fmt.Sprintf("Expediting proposal ID: %d", proposalID))
	}

	k.SetProposal(ctx, proposal)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeVote,
			sdk.NewAttribute(types.AttributeKeyProposalId, strconv.Itoa(int(proposalID))),
			sdk.NewAttribute(types.AttributeKeyOption, option.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, voter.String()),
		),
	})

	return nil
}

func (k Keeper) GetAllVotes(ctx sdk.Context) (votes types.Votes) {
	k.IterateAllVotes(ctx, func(vote types.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

func (k Keeper) GetVotes(ctx sdk.Context, proposalID uint64) (votes types.Votes) {
	k.IterateVotes(ctx, proposalID, func(vote types.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

func (k Keeper) GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (types.Vote, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VoteKey(proposalID, voterAddr))
	var vote types.Vote
	if bz == nil {
		return vote, types.ErrUnknownVote
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &vote)
	return vote, nil
}

func (k Keeper) HasVoted(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.VoteKey(proposalID, voterAddr))
}

func (k Keeper) SetVote(ctx sdk.Context, vote types.Vote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(vote)
	store.Set(types.VoteKey(vote.ProposalID, vote.Voter), bz)
}

func (k Keeper) IterateAllVotes(ctx sdk.Context, cb func(vote types.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.VotesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

func (k Keeper) IterateVotes(ctx sdk.Context, proposalID uint64, cb func(vote types.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.VotesKey(proposalID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

func (k Keeper) deleteVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.VoteKey(proposalID, voterAddr))
}