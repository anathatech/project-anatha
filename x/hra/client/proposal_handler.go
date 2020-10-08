package client


import (
	govclient "github.com/anathatech/project-anatha/x/governance/client"
	"github.com/anathatech/project-anatha/x/hra/client/cli"
)

var (
	RegisterBlockchainIdProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitRegisterBlockchainIdProposal)
	RemoveBlockchainIdProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveBlockchainIdProposal)
)