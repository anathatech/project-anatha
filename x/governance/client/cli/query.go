package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/anathatech/project-anatha/x/governance/internal/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	govQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Aliases: 					[]string{"gov"},
		Short:                      "Querying commands for the governance module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	govQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryProposal(queryRoute, cdc),
			GetCmdQueryProposals(queryRoute, cdc),
			GetCmdQueryVote(queryRoute, cdc),
			GetCmdQueryVotes(queryRoute, cdc),
			GetCmdQueryParams(queryRoute, cdc),
			GetCmdQueryGovernors(queryRoute, cdc),
		)...
	)

	return govQueryCmd
}

func GetCmdQueryProposal(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proposal [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid uint, please input a valid proposal-id", args[0])
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/proposal/%d", queryRoute, proposalID), nil)
			if err != nil {
				fmt.Printf("Could not find Proposal - %d \n", proposalID)
				return nil
			}

			var proposal types.Proposal
			cdc.MustUnmarshalJSON(res, &proposal)
			return cliCtx.PrintOutput(proposal)
		},
	}
}

func GetCmdQueryProposals(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proposals",
		Short: "Query all proposals",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/proposals", queryRoute), nil)
			if err != nil {
				fmt.Print("Could not get proposals \n")
				return nil
			}

			var out []types.Proposal
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdQueryVote(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "vote [proposal-id] [voter-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a single vote",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			_, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			_, err = sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/vote/%s/%s", queryRoute, args[0], args[1]), nil)
			if err != nil {
				fmt.Printf("Could not resolve vote - %s \n", args[1])
				return nil
			}

			var vote types.Vote
			cdc.MustUnmarshalJSON(res, &vote)
			return cliCtx.PrintOutput(vote)
		},
	}
}

func GetCmdQueryVotes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "votes [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query votes on a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			_, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/votes/%s", queryRoute, args[0]), nil)
			if err != nil {
				fmt.Printf("Could not resolve votes - %s \n", args[1])
				return nil
			}

			var votes []types.Vote
			cdc.MustUnmarshalJSON(res, &votes)
			return cliCtx.PrintOutput(votes)
		},
	}
}

func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the governance process.

Example:
$ %s query gov params
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params/tallying", queryRoute), nil)
			if err != nil {
				return err
			}
			vp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params/voting", queryRoute), nil)
			if err != nil {
				return err
			}

			var tallyParams types.TallyParams
			cdc.MustUnmarshalJSON(tp, &tallyParams)
			var votingParams types.VotingParams
			cdc.MustUnmarshalJSON(vp, &votingParams)

			return cliCtx.PrintOutput(types.NewParams(votingParams, tallyParams))
		},
	}
}

func GetCmdQueryGovernors(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "governors",
		Short: "Query governors",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/governors", queryRoute), nil)
			if err != nil {
				fmt.Printf("Could not resolve governors\n")
				return nil
			}

			var governors []sdk.AccAddress
			cdc.MustUnmarshalJSON(res, &governors)
			return cliCtx.PrintOutput(governors)
		},
	}
}

