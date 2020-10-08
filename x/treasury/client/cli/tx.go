package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	denom "github.com/anathatech/project-anatha/utils"
	"github.com/anathatech/project-anatha/x/treasury/internal/types"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	treasuryTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	treasuryTxCmd.AddCommand(flags.PostCommands(
		GetCmdOperator(cdc),
		GetCmdOrder(cdc),
		GetCmdDisburse(cdc),
		GetCmdCancelDisbursement(cdc),
		GetCmdDisburseToEscrow(cdc),
		GetCmdDisburseFromEscrow(cdc),
		GetCmdRevertFromEscrow(cdc),
	)...)

	return treasuryTxCmd
}

func GetCmdDisburse(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "disburse [recipient] [amount] [reference]",
		Short: "Disburse treasury funds",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			amount, err := denom.ParseAndConvertCoins(args[1])
			if err != nil {
				return err
			}

			reference := args[2]

			msg := types.NewMsgDisburse(cliCtx.GetFromAddress(), recipient, amount, reference)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdCancelDisbursement(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "cancel-disbursement [recipient] [scheduled]",
		Short: "Cancel scheduled distribution",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelDisbursement(cliCtx.GetFromAddress(), recipient, args[1])
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdDisburseToEscrow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "disburse-to-escrow [amount] [reference]",
		Short: "Disburse treasury funds to escrow",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			amount, err := denom.ParseAndConvertCoins(args[0])
			if err != nil {
				return err
			}

			reference := args[1]

			msg := types.NewMsgDisburseToEscrow(cliCtx.GetFromAddress(), amount, reference)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdDisburseFromEscrow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "disburse-from-escrow [reference] [recipient]",
		Short: "Disburse funds from escrow",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			reference := args[0]

			recipient, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDisburseFromEscrow(cliCtx.GetFromAddress(), reference, recipient)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdRevertFromEscrow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "revert-from-escrow [amount] [reference]",
		Short: "Revert frunds from the escrow to the buyback fund",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			amount, err := denom.ParseAndConvertCoins(args[0])
			if err != nil {
				return err
			}

			reference := args[1]

			msg := types.NewMsgRevertFromEscrow(cliCtx.GetFromAddress(), amount, reference)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

