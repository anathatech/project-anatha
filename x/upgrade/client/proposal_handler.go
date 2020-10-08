package client

import (
	govclient "github.com/anathatech/project-anatha/x/governance/client"
	"github.com/anathatech/project-anatha/x/upgrade/client/cli"
)

var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitSoftwareUpgradeProposal)
