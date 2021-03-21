package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/anathatech/project-anatha/x/treasury/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group treasury queries under a subcommand
	treasuryQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	treasuryQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryParams(queryRoute, cdc),
			GetCmdQueryTreasury(queryRoute, cdc),
			GetCmdOperators(queryRoute, cdc),
			GetCmdDisbursements(queryRoute, cdc),
			GetCmdQueryPrice(queryRoute, cdc),
			GetCmdQueryDisbursementEscrow(queryRoute, cdc),
		)...,
	)

	return treasuryQueryCmd
}

func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current treasury module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/parameters", queryRoute)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
}

func GetCmdQueryTreasury(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "treasury",
		Short: "Query the current treasury",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/treasury", queryRoute)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Treasury
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
}

func GetCmdOperators(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "operators",
		Short: "Query Treasury Distribution Operators",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/operators", queryRoute), nil)
			if err != nil {
				fmt.Printf("Could not resolve operators\n")
				return nil
			}

			var out types.QueryResOperators
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdDisbursements(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "disbursements",
		Short: "Query Treasury Scheduled Disbursements",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/disbursements", queryRoute), nil)
			if err != nil {
				fmt.Printf("Could not resolve disbursements\n")
				return nil
			}

			var out types.QueryResDisbursements
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdQueryPrice(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "price [amount]",
		Short: "Query the current price (fiat, anatha)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/price/%s", queryRoute, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var price types.QueryResPrice
			if err := cdc.UnmarshalJSON(res, &price); err != nil {
				return err
			}

			return cliCtx.PrintOutput(price)
		},
	}
}

func GetCmdQueryDisbursementEscrow(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "disbursement-escrow [reference]",
		Short: "Query escrow of a distribution by its reference. 0 = distributed, Int = pending, error = reference doesn't exist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/disbursement-escrow/%s", queryRoute, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var amount sdk.Coins
			if err := cdc.UnmarshalJSON(res, &amount); err != nil {
				return err
			}

			return cliCtx.PrintOutput(amount)
		},
	}
}