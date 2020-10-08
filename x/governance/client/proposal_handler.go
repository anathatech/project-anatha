package client

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
)


type CLIHandlerFn func(*codec.Codec) *cobra.Command

type ProposalHandler struct {
	CLIHandler  CLIHandlerFn
}

func NewProposalHandler(cliHandler CLIHandlerFn) ProposalHandler {
	return ProposalHandler{
		CLIHandler:  cliHandler,
	}
}
