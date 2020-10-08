package client

import (
	"github.com/anathatech/project-anatha/x/distribution/client/cli"
	govclient "github.com/anathatech/project-anatha/x/governance/client"
)

var DevelopmentFundDistributionProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitDevelopmentFundDistributionProposal)
var SecurityTokenFundDistributionProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitSecurityTokenFundDistributionProposal)
