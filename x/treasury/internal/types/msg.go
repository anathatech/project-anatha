package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/config"
	"time"
)

// MsgAddOperator
type MsgAddOperator struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Operator sdk.AccAddress `json:"operator" yaml:"operator"`
}

func NewMsgAddOperator(sender sdk.AccAddress, operator sdk.AccAddress) MsgAddOperator {
	return MsgAddOperator{
		Sender: sender,
		Operator: operator,
	}
}

func (msg MsgAddOperator) Route() string { return RouterKey }

func (msg MsgAddOperator) Type() string { return "add_operator" }

func (msg MsgAddOperator) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}
	if msg.Operator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}
	return nil
}

func (msg MsgAddOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgAddOperator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgRemoveOperator
type MsgRemoveOperator struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Operator sdk.AccAddress `json:"operator" yaml:"operator"`
}

func NewMsgRemoveOperator(sender sdk.AccAddress, operator sdk.AccAddress) MsgRemoveOperator {
	return MsgRemoveOperator{
		Sender: sender,
		Operator: operator,
	}
}

func (msg MsgRemoveOperator) Route() string { return RouterKey }

func (msg MsgRemoveOperator) Type() string { return "remove_operator" }

func (msg MsgRemoveOperator) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}
	if msg.Operator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}
	return nil
}

func (msg MsgRemoveOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRemoveOperator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgDisburse
type MsgDisburse struct {
	Operator  sdk.AccAddress `json:"operator" yaml:"operator"`
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
	Amount    sdk.Coins      `json:"amount" yaml:"amount"`
	Reference string         `json:"reference" yaml:"reference"`
}

func NewMsgDisburse(sender sdk.AccAddress, recipient sdk.AccAddress, amount sdk.Coins, reference string) MsgDisburse {
	return MsgDisburse{
		Operator:  sender,
		Recipient: recipient,
		Amount:    amount,
		Reference: reference,
	}
}

func (msg MsgDisburse) Route() string { return RouterKey }

func (msg MsgDisburse) Type() string { return "disburse" }

func (msg MsgDisburse) ValidateBasic() error {
	if msg.Operator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Operator.String())
	}
	if msg.Recipient.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Operator.String())
	}
	if ! msg.Amount.AmountOf(config.DefaultStableDenom).IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	if len(msg.Reference) > 255 {
		return sdkerrors.Wrapf(sdkerrors.ErrMemoTooLarge, "Reference too long")
	}
	return nil
}

func (msg MsgDisburse) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDisburse) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Operator}
}

// MsgDisburseToEscrow
type MsgDisburseToEscrow struct {
	Operator  sdk.AccAddress `json:"operator" yaml:"operator"`
	Amount    sdk.Coins      `json:"amount" yaml:"amount"`
	Reference string         `json:"reference" yaml:"reference"`
}

func NewMsgDisburseToEscrow(sender sdk.AccAddress, amount sdk.Coins, reference string) MsgDisburseToEscrow {
	return MsgDisburseToEscrow{
		Operator:  sender,
		Amount:    amount,
		Reference: reference,
	}
}

func (msg MsgDisburseToEscrow) Route() string { return RouterKey }

func (msg MsgDisburseToEscrow) Type() string { return "disburse_to_escrow" }

func (msg MsgDisburseToEscrow) ValidateBasic() error {
	if msg.Operator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Operator.String())
	}
	if ! msg.Amount.AmountOf(config.DefaultStableDenom).IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	if len(msg.Reference) > 255 {
		return sdkerrors.Wrapf(sdkerrors.ErrMemoTooLarge, "Reference too long")
	}
	return nil
}

func (msg MsgDisburseToEscrow) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDisburseToEscrow) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Operator}
}

// MsgDisburseFromEscrow
type MsgDisburseFromEscrow struct {
	Operator  sdk.AccAddress `json:"operator" yaml:"operator"`
	Reference string         `json:"reference" yaml:"reference"`
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
}

func NewMsgDisburseFromEscrow(sender sdk.AccAddress, reference string, recipient sdk.AccAddress) MsgDisburseFromEscrow {
	return MsgDisburseFromEscrow{
		Operator:  sender,
		Reference: reference,
		Recipient: recipient,
	}
}

func (msg MsgDisburseFromEscrow) Route() string { return RouterKey }

func (msg MsgDisburseFromEscrow) Type() string { return "disburse_from_escrow" }

func (msg MsgDisburseFromEscrow) ValidateBasic() error {
	if msg.Operator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Operator.String())
	}
	if msg.Recipient.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Recipient.String())
	}
	if len(msg.Reference) > 255 {
		return sdkerrors.Wrapf(sdkerrors.ErrMemoTooLarge, "Reference too long")
	}
	return nil
}

func (msg MsgDisburseFromEscrow) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDisburseFromEscrow) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Operator}
}

// MsgRevertFromEscrow
type MsgRevertFromEscrow struct {
	Operator  sdk.AccAddress `json:"operator" yaml:"operator"`
	Amount    sdk.Coins      `json:"amount" yaml:"amount"`
	Reference string         `json:"reference" yaml:"reference"`
}

func NewMsgRevertFromEscrow(sender sdk.AccAddress, amount sdk.Coins, reference string) MsgRevertFromEscrow {
	return MsgRevertFromEscrow{
		Operator:  sender,
		Amount:    amount,
		Reference: reference,
	}
}

func (msg MsgRevertFromEscrow) Route() string { return RouterKey }

func (msg MsgRevertFromEscrow) Type() string { return "revert_from_escrow" }

func (msg MsgRevertFromEscrow) ValidateBasic() error {
	if msg.Operator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Operator.String())
	}
	if ! msg.Amount.AmountOf(config.DefaultDenom).IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	if len(msg.Reference) > 255 {
		return sdkerrors.Wrapf(sdkerrors.ErrMemoTooLarge, "Reference too long")
	}
	return nil
}

func (msg MsgRevertFromEscrow) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRevertFromEscrow) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Operator}
}

// MsgCancelDisbursement
type MsgCancelDisbursement struct {
	Manager sdk.AccAddress `json:"manager" yaml:"manager"`
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
	ScheduledFor string `json:"scheduled_for" yaml:"scheduled_for"` // has to be string and parsed on server because of broken amino decoding
}

func NewMsgCancelDisbursement(manager sdk.AccAddress, recipient sdk.AccAddress, scheduledFor string) MsgCancelDisbursement {
	return MsgCancelDisbursement{
		Manager: manager,
		Recipient: recipient,
		ScheduledFor: scheduledFor,
	}
}

func (msg MsgCancelDisbursement) Route() string { return RouterKey }

func (msg MsgCancelDisbursement) Type() string { return "cancel_disbursement" }

func (msg MsgCancelDisbursement) ValidateBasic() error {
	if msg.Manager.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Manager.String())
	}
	if msg.Recipient.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Recipient.String())
	}
	scheduledFor, err := time.Parse("2006-01-02T15:04:05.99999999999Z", msg.ScheduledFor)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidTime, msg.ScheduledFor)
	}
	if scheduledFor.IsZero() {
		return sdkerrors.Wrap(ErrInvalidTime, msg.ScheduledFor)
	}
	return nil
}

func (msg MsgCancelDisbursement) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCancelDisbursement) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Manager}
}

// MsgCreateSellOrder
type MsgCreateSellOrder struct {
	Seller          sdk.AccAddress `json:"seller" yaml:"seller"`
	Amount          sdk.Coins      `json:"amount" yaml:"amount"`
}

func NewMsgCreateSellOrder(seller sdk.AccAddress, amount sdk.Coins) MsgCreateSellOrder {
	return MsgCreateSellOrder{
		Seller:          seller,
		Amount:          amount,
	}
}

func (msg MsgCreateSellOrder) Route() string { return RouterKey }

func (msg MsgCreateSellOrder) Type() string { return "create_sell_order" }

func (msg MsgCreateSellOrder) ValidateBasic() error {
	if msg.Seller.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Seller.String())
	}
	if ! msg.Amount.AmountOf(config.DefaultDenom).IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	return nil
}

func (msg MsgCreateSellOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateSellOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Seller}
}

// MsgCreateBuyOrder
type MsgCreateBuyOrder struct {
	Buyer 	sdk.AccAddress 	`json:"buyer" yaml:"buyer"`
	Amount 	sdk.Coins 		`json:"amount" yaml:"amount"`
}

func NewMsgCreateBuyOrder(buyer sdk.AccAddress, amount sdk.Coins) MsgCreateBuyOrder {
	return MsgCreateBuyOrder{
		Buyer: buyer,
		Amount: amount,
	}
}

func (msg MsgCreateBuyOrder) Route() string { return RouterKey }

func (msg MsgCreateBuyOrder) Type() string { return "create_buy_order" }

func (msg MsgCreateBuyOrder) ValidateBasic() error {
	if msg.Buyer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Buyer.String())
	}
	if ! msg.Amount.AmountOf(config.DefaultStableDenom).IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	return nil
}

func (msg MsgCreateBuyOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateBuyOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

// MsgSwap
type MsgSwap struct {
	Operator  sdk.AccAddress `json:"operator" yaml:"operator"`
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
	Amount    sdk.Coins      `json:"amount" yaml:"amount"`
	Reference string         `json:"reference" yaml:"reference"`
}

func NewMsgSwap(sender sdk.AccAddress, recipient sdk.AccAddress, amount sdk.Coins, reference string) MsgSwap {
	return MsgSwap{
		Operator:  sender,
		Recipient: recipient,
		Amount:    amount,
		Reference: reference,
	}
}

func (msg MsgSwap) Route() string { return RouterKey }

func (msg MsgSwap) Type() string { return "swap" }

func (msg MsgSwap) ValidateBasic() error {
	if msg.Operator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Operator.String())
	}
	if msg.Recipient.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Operator.String())
	}
	if ! msg.Amount.AmountOf(config.DefaultDenom).IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	if len(msg.Reference) > 255 {
		return sdkerrors.Wrapf(sdkerrors.ErrMemoTooLarge, "Reference too long")
	}
	return nil
}

func (msg MsgSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSwap) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Operator}
}