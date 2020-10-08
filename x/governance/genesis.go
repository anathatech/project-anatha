package governance

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetProposalID(ctx, data.StartingProposalID)
	k.SetVotingParams(ctx, data.VotingParams)
	k.SetTallyParams(ctx, data.TallyParams)

	for _, vote := range data.Votes {
		k.SetVote(ctx, vote)
	}

	for _, proposal := range data.Proposals {
		if proposal.Status == StatusVotingPeriod {
			k.InsertActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)
		}

		k.SetProposal(ctx, proposal)
	}

	for _, governor := range data.Governors {
		k.AddGovernor(ctx, governor)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingProposalID, _ := k.GetProposalID(ctx)
	votingParams := k.GetVotingParams(ctx)
	tallyParams := k.GetTallyParams(ctx)
	proposals := k.GetProposals(ctx)
	governors := k.GetGovernors(ctx)

	var proposalsVotes Votes
	for _, proposal := range proposals {
		votes := k.GetVotes(ctx, proposal.ProposalID)
		proposalsVotes = append(proposalsVotes, votes...)
	}

	return GenesisState{
		StartingProposalID: startingProposalID,
		Votes:              proposalsVotes,
		Proposals:          proposals,
		VotingParams:       votingParams,
		TallyParams:        tallyParams,
		Governors: 			governors,
	}
}
