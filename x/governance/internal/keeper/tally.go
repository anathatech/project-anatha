package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/x/governance/internal/types"
)

func (k Keeper) Tally(ctx sdk.Context, proposal types.Proposal) (passes bool, tallyResults types.TallyResult) {
	results := make(map[types.VoteOption]sdk.Dec)
	results[types.OptionYes] = sdk.ZeroDec()
	results[types.OptionNo] = sdk.ZeroDec()

	totalVotes := sdk.ZeroDec()

	k.IterateVotes(ctx, proposal.ProposalID, func(vote types.Vote) bool {
		results[vote.Option] = results[vote.Option].Add(sdk.OneDec())
		totalVotes = totalVotes.Add(sdk.OneDec())

		k.deleteVote(ctx, vote.ProposalID, vote.Voter)

		return false
	})

	tallyParams := k.GetTallyParams(ctx)
	tallyResults = types.NewTallyResultFromMap(results)

	if ! proposal.TallyResult.Equals(tallyResults) {
		panic(fmt.Sprintf("proposal %d tally does not match aggregate", proposal.ProposalID))
	}

	currentGovernorCount := k.GetGovernorCount(ctx)

	// Only pass adding a new governor when all of the current ones agree
	if proposal.Content.ProposalType() == types.ProposalTypeAddGovernor {
		if tallyResults.Yes.Equal(currentGovernorCount) {
			return true, tallyResults
		} else {
			return false, tallyResults
		}
	}

	// Only remove a governor when everyone excluding the governor in question agrees
	if proposal.Content.ProposalType() == types.ProposalTypeRemoveGovernor {
		if tallyResults.Yes.GTE(currentGovernorCount.Sub(sdk.OneInt())) {
			return true, tallyResults
		} else {
			return false, tallyResults
		}
	}

	// If there is not enough quorum of votes, the proposal fails
	percentVoting := totalVotes.Quo(currentGovernorCount.ToDec())
	if percentVoting.LT(tallyParams.Quorum) {
		return false, tallyResults
	}

	// If no one votes, proposal fails
	if totalVotes.Equal(sdk.ZeroDec()) {
		return false, tallyResults
	}

	// If more than Threshold of voters vote Yes, proposal passes
	if results[types.OptionYes].Quo(totalVotes).GT(tallyParams.Threshold) {
		return true, tallyResults
	}

	// If more than 1/2 of voters vote No, proposal fails
	return false, tallyResults
}

