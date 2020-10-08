package utils

import "github.com/anathatech/project-anatha/x/governance/internal/types"

func NormalizeVoteOption(option string) string {
	switch option {
	case "Yes", "yes":
		return types.OptionYes.String()
	case "No", "no":
		return types.OptionNo.String()
	default:
		return ""
	}
}

func NormalizeProposalType(proposalType string) string {
	switch proposalType {
	case "Text", "text":
		return types.ProposalTypeText

	default:
		return ""
	}
}

func NormalizeProposalStatus(status string) string {
	switch status {
	case "VotingPeriod", "voting_period":
		return "VotingPeriod"
	case "Passed", "passed":
		return "Passed"
	case "Rejected", "rejected":
		return "Rejected"
	}
	return ""
}
