package client

import (
	"github.com/anathatech/project-anatha/x/fee/client/cli"
	govclient "github.com/anathatech/project-anatha/x/governance/client"
)

var (
	AddFeeExcludedMessageProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitAddFeeExcludedMessageProposal)
	RemoveFeeExcludedMessageProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveFeeExcludedMessageProposal)
)