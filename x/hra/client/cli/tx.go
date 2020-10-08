package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	denom "github.com/anathatech/project-anatha/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/anathatech/project-anatha/x/hra/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	hraTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%S transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	hraTxCmd.AddCommand(flags.PostCommands(
		GetCmdRegisterName(cdc),
		GetCmdRenewName(cdc),
		GetCmdSetPrice(cdc),
		GetCmdDeleteName(cdc),
		GetCmdBuyName(cdc),
		GetCmdTransferName(cdc),
		GetCmdRegisterAddress(cdc),
		GetCmdRegisterAddressBatch(cdc),
		GetCmdRemoveAddress(cdc),
		GetCmdRemoveAllAddresses(cdc),
	)...)

	return hraTxCmd
}

func GetCmdRegisterName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{

		Use:   "register [name] [addresses-file]",
		Short: "register a new hra",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var addresses []types.BlockchainAddressInfo
			var msgs []sdk.Msg

			if len(args) > 1 {

				file, err := os.Open(args[1])
				if err != nil {
					return err
				}

				defer closeFile(file)

				bytes, _ := ioutil.ReadAll(file)

				err = json.Unmarshal(bytes, &addresses)
				if err != nil {
					return err
				}

			}

			msg := types.NewMsgRegisterName(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			msgs = append(msgs, msg)

			for _, address := range addresses {
				registerAddressMsg := types.NewMsgRegisterAddress(
					cliCtx.GetFromAddress(),
					address.BlockchainId,
					address.Index,
					address.BlockchainAddress,
				)

				err := registerAddressMsg.ValidateBasic()
				if err != nil {
					return err
				}

				msgs = append(msgs, registerAddressMsg)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
		},
	}
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		panic(fmt.Sprintf("error closing file: %v\n", err))
	}
}

func GetCmdRenewName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "renew [name]",
		Short: "renew a hra",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgRenewName(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdSetPrice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-price [name] [price]",
		Short: "set price for selling a hra",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			coins, err := denom.ParseAndConvertCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetPrice(args[0], cliCtx.GetFromAddress(), coins)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdDeleteName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [name]",
		Short: "delete a hra",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgDeleteName(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdBuyName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "buy [name]",
		Short: "buy a hra",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgBuyName(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdTransferName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "transfer [name] [new_owner]",
		Short: "transfer a hra",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			newOwner, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferName(args[0], cliCtx.GetFromAddress(), newOwner)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdRegisterAddress(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "register-address [blockchanId] [index] [blockchainAddress]",
		Short: "register a new blockchain address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgRegisterAddress(cliCtx.GetFromAddress(), args[0], args[1], args[2])
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdRegisterAddressBatch(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "register-address-batch [addresses-file]",
		Short: "batch blockchain address register",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			var addresses []types.BlockchainAddressInfo
			var msgs []sdk.Msg

			file, err := os.Open(args[0])
			if err != nil {
				return err
			}

			defer closeFile(file)

			bytes, _ := ioutil.ReadAll(file)

			err = json.Unmarshal(bytes, &addresses)
			if err != nil {
				return err
			}

			for _, address := range addresses {
				registerAddressMsg := types.NewMsgRegisterAddress(
					cliCtx.GetFromAddress(),
					address.BlockchainId,
					address.Index,
					address.BlockchainAddress,
				)

				err := registerAddressMsg.ValidateBasic()
				if err != nil {
					return err
				}

				msgs = append(msgs, registerAddressMsg)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
		},
	}
}

func GetCmdRemoveAddress(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "remove-address [blockchanId] [index]",
		Short: "remove a blockchain address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgRemoveAddress(cliCtx.GetFromAddress(), args[0], args[1])
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdRemoveAllAddresses(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "remove-all-addresses",
		Short: "remove all associated blockchain addresses",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgRemoveAllAddresses(cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}