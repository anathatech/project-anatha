package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	hraQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	hraQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdResolveName(queryRoute, cdc),
			GetCmdResolveAddress(queryRoute, cdc),
			GetCmdGetBlockchainAddress(queryRoute, cdc),
			GetCmdResolveAddressNameInfos(queryRoute, cdc),
			GetCmdGetRegisteredBlockchainIds(queryRoute, cdc),
			GetCmdGetAddressCredits(queryRoute, cdc),
			GetCmdGetBlockchainAddresses(queryRoute, cdc),
			GetCmdQueryParams(queryRoute, cdc),
		)...,
	)

	hraQueryCmd.AddCommand(GetModuleAccountCmd(cdc))

	return hraQueryCmd

}

func GetCmdResolveName(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve [name]",
		Short: "Query hra info of name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/resolve/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("Could not resolve hra - %s \n", name)
				return nil
			}

			var out types.QueryResNameInfo
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdResolveAddress(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "names [address]",
		Short: "Query address names",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/names/%s", queryRoute, address), nil)
			if err != nil {
				fmt.Printf("Could not resolve address - %s \n", address)
				return nil
			}

			var out types.QueryResNames
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdGetBlockchainAddress(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "address [name] [blockchain-id] [index]",
		Short: "Resolve a HRA to a blockchain address",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := args[0]
			blockchainId := args[1]

			index := "0"
			if len(args) == 3 {
				index = args[2]
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/address/%s/%s/%s", queryRoute, name, blockchainId, index), nil)
			if err != nil {
				fmt.Printf("Could not resolve address - %s/%s/%s \n", name, blockchainId, index)
				return nil
			}

			var out types.QueryResAddress
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdResolveAddressNameInfos(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "name-infos [address]",
		Short: "Query name infos for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/name-infos/%s", queryRoute, address), nil)
			if err != nil {
				fmt.Printf("Could not resolve address - %s \n", address)
				return nil
			}

			var out types.QueryResNameInfos
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdGetRegisteredBlockchainIds(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "registered-blockchain-ids",
		Short: "Query registered blockchain ids",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/registered-blockchain-ids", queryRoute), nil)
			if err != nil {
				fmt.Print("Could not get registered blockchain ids \n")
				return nil
			}

			var out types.QueryResRegisteredBlockchainIds
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdGetAddressCredits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "address-credits [address]",
		Short: "Query credits for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/address-credits/%s", queryRoute, address), nil)
			if err != nil {
				fmt.Printf("Could not resolve address - %s \n", address)
				return nil
			}

			var out types.QueryResAddressCredits
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdGetBlockchainAddresses(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "blockchain-addresses [address]",
		Short: "Query blockchain addresses for given address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/blockchain-addresses/%s", queryRoute, address), nil)
			if err != nil {
				fmt.Printf("Could not resolve address - %s \n", address)
				return nil
			}

			var out types.QueryResBlockchainAddresses
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetModuleAccountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module [name]",
		Short: "Query module account balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			accGetter := auth.NewAccountRetriever(cliCtx)

			key := supply.NewModuleAddress(args[0])

			acc, err := accGetter.GetAccount(key)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(acc)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current minting parameters",
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