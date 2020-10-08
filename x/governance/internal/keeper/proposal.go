package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gov "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/anathatech/project-anatha/x/governance/internal/types"
	"strconv"
)

func (k Keeper) SubmitProposal(ctx sdk.Context, proposer sdk.AccAddress, content gov.Content) error {
	if ! k.IsGovernor(ctx, proposer) {
		return sdkerrors.Wrap(types.ErrNotGovernor, proposer.String())
	}
	if ! k.router.HasRoute(content.ProposalRoute()) {
		return sdkerrors.Wrap(types.ErrNoProposalHandlerExists, content.ProposalRoute())
	}

	cacheCtx, _ := ctx.CacheContext()
	handler := k.router.GetRoute(content.ProposalRoute())
	if err := handler(cacheCtx, content); err != nil {
		return sdkerrors.Wrap(types.ErrInvalidProposalContent, err.Error())
	}

	proposalID, err := k.GetProposalID(ctx)
	if err != nil {
		return err
	}

	votingPeriod := k.GetVotingParams(ctx).VotingPeriod

	proposal := types.NewProposal(content, proposalID)

	proposal.Status = types.StatusVotingPeriod
	proposal.VotingStartTime = ctx.BlockHeader().Time
	proposal.VotingEndTime = proposal.VotingStartTime.Add(votingPeriod)

	k.SetProposal(ctx, proposal)
	k.InsertActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)

	k.SetProposalID(ctx, proposalID + 1)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalId, strconv.Itoa(int(proposalID))),
			sdk.NewAttribute(types.AttributeKeyContent, content.String()),
			sdk.NewAttribute(types.AttributeKeyStatus, proposal.Status.String()),
			sdk.NewAttribute(types.AttributeKeyVotingStartTime, proposal.VotingStartTime.String()),
			sdk.NewAttribute(types.AttributeKeyVotingEndTime, proposal.VotingEndTime.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, proposer.String()),
		),
	})

	return nil
}

func (k Keeper) ExecuteProposal(ctx sdk.Context, proposal types.Proposal) (status types.ProposalStatus) {
	handler := k.Router().GetRoute(proposal.ProposalRoute())
	cacheCtx, writeCache := ctx.CacheContext()

	err := handler(cacheCtx, proposal.Content)
	if err == nil {
		status = types.StatusPassed

		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

		writeCache()
	} else {
		status = types.StatusFailed
	}

	return
}

func (k Keeper) HandleProposal(ctx sdk.Context, proposal types.Proposal, expedited bool) {
	passes, _ := k.Tally(ctx, proposal)

	if passes {
		proposal.Status = k.ExecuteProposal(ctx, proposal)
	} else {
		proposal.Status = types.StatusRejected
	}

	proposal.ExecutionTime = ctx.BlockTime()

	k.SetProposal(ctx, proposal)

	if expedited {
		k.RemoveFromExpeditedProposalQueue(ctx, proposal.ProposalID)
	} else {
		k.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)
	}

	// TODO: events
}

func (k Keeper) CanBeExpedited(ctx sdk.Context, proposal types.Proposal) bool {
	governorCount := k.GetGovernorCount(ctx)
	totalVotes := proposal.TallyResult.Yes.Add(proposal.TallyResult.No)

	tallyParams := k.GetTallyParams(ctx)

	if totalVotes.ToDec().Quo(governorCount.ToDec()).GTE(tallyParams.Quorum) {
		// quorum condition fulfilled
		if proposal.TallyResult.Yes.ToDec().Quo(governorCount.ToDec()).GTE(tallyParams.Threshold) {
			/*
				Proportion of Yes votes relative to the total amount of possible voters surpassed the threshold
				meaning that the vote can not be overturned.

				Calculation based on the following inequality:

				x - number of No votes
				y - number of Yes votes
				n - total voters (governors)
				t - vote threshold

				not_voted = n - (x + y)

				(y + not_voted) / (x + n) >= t

				which equals to

				y / n >= t

				https://desmos.com/calculator/kfhbqao95m
			*/
			return true
		}
	}

	return false
}

func (k Keeper) HandleExpedite(ctx sdk.Context, proposalId uint64) error {
	proposal, err := k.GetProposal(ctx, proposalId)

	if err != nil {
		return sdkerrors.Wrapf(types.ErrUnknownProposal, "%d", proposalId)
	}
	if proposal.Status != types.StatusVotingPeriod {
		return sdkerrors.Wrapf(types.ErrInactiveProposal, "%d", proposalId)
	}

	if k.CanBeExpedited(ctx, proposal) {
		proposal.Status = types.StatusExpediting
		proposal.Expedited = true

		k.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)
		k.InsertExpeditedProposalQueue(ctx, proposal.ProposalID)
	}

	k.SetProposal(ctx, proposal)

	return nil
}

func (k Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (types.Proposal, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposalKey(proposalID))
	if bz == nil {
		return types.Proposal{}, types.ErrUnknownProposal
	}

	var proposal types.Proposal
	k.cdc.MustUnmarshalBinaryBare(bz, &proposal)
	return proposal, nil
}

func (k Keeper) SetProposal(ctx sdk.Context, proposal types.Proposal) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(proposal)
	store.Set(types.ProposalKey(proposal.ProposalID), bz)
}

func (k Keeper) DeleteProposal(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	proposal, err := k.GetProposal(ctx, proposalID)
	if err != nil {
		panic(fmt.Sprintf("couldn't find proposal with id#%d", proposalID))
	}
	k.RemoveFromActiveProposalQueue(ctx, proposalID, proposal.VotingEndTime)
	store.Delete(types.ProposalKey(proposalID))
}

func (k Keeper) IterateProposals(ctx sdk.Context, cb func(proposal types.Proposal) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ProposalsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proposal types.Proposal
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &proposal)

		if cb(proposal) {
			break
		}
	}
}

func (k Keeper) GetProposals(ctx sdk.Context) (proposals types.Proposals) {
	k.IterateProposals(ctx, func(proposal types.Proposal) bool {
		proposals = append(proposals, proposal)
		return false
	})
	return
}

func (k Keeper) GetProposalID(ctx sdk.Context) (proposalID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposalIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial proposal ID hasn't been set")
	}

	proposalID = types.GetProposalIDFromBytes(bz)
	return proposalID, nil
}

func (k Keeper) SetProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProposalIDKey, types.GetProposalIDBytes(proposalID))
}

