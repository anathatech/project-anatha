package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/config"
)

var _, _ sdk.Msg = &MsgWithdrawNameReward{}, &MsgWithdrawValidatorReward{}

// MsgWithdrawNameReward
type MsgWithdrawNameReward struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
}

func NewMsgWithdrawNameReward(sender sdk.AccAddress) MsgWithdrawNameReward {
	return MsgWithdrawNameReward{
		Sender: sender,
	}
}

func (msg MsgWithdrawNameReward) Route() string { return RouterKey }

func (msg MsgWithdrawNameReward) Type() string { return "withdraw_name_reward" }

func (msg MsgWithdrawNameReward) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

func (msg MsgWithdrawNameReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgWithdrawNameReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgWithdrawValidatorReward
type MsgWithdrawValidatorReward struct {
	Validator sdk.ValAddress `json:"validator" yaml:"validator"`
}

func NewMsgWithdrawValidatorReward(validator sdk.ValAddress) MsgWithdrawValidatorReward {
	return MsgWithdrawValidatorReward{
		Validator: validator,
	}
}

func (msg MsgWithdrawValidatorReward) Route() string { return RouterKey }
func (msg MsgWithdrawValidatorReward) Type() string  { return "withdraw_validator_rewards" }

// quick validity check
func (msg MsgWithdrawValidatorReward) ValidateBasic() error {
	if msg.Validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator.String())
	}
	return nil
}

func (msg MsgWithdrawValidatorReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgWithdrawValidatorReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Validator.Bytes())}
}

// MsgDepositSavings
type MsgDepositSavings struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Amount sdk.Coins `json:"amount" yaml:"amount"`
}

func NewMsgDepositSavings(sender sdk.AccAddress, amount sdk.Coins) MsgDepositSavings {
	return MsgDepositSavings{
		Sender: sender,
		Amount: amount,
	}
}

func (msg MsgDepositSavings) Route() string { return RouterKey }

func (msg MsgDepositSavings) Type() string { return "deposit_savings" }

func (msg MsgDepositSavings) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}
	if msg.Amount.AmountOf(config.DefaultDenom).IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

func (msg MsgDepositSavings) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDepositSavings) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgWithdrawSavings
type MsgWithdrawSavings struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
}

func NewMsgWithdrawSavings(sender sdk.AccAddress) MsgWithdrawSavings {
	return MsgWithdrawSavings{
		Sender: sender,
	}
}

func (msg MsgWithdrawSavings) Route() string { return RouterKey }

func (msg MsgWithdrawSavings) Type() string { return "withdraw_savings" }

func (msg MsgWithdrawSavings) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

func (msg MsgWithdrawSavings) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgWithdrawSavings) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgWithdrawSavingsInterest
type MsgWithdrawSavingsInterest struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
}

func NewMsgWithdrawSavingsInterest(sender sdk.AccAddress) MsgWithdrawSavingsInterest {
	return MsgWithdrawSavingsInterest{
		Sender: sender,
	}
}

func (msg MsgWithdrawSavingsInterest) Route() string { return RouterKey }

func (msg MsgWithdrawSavingsInterest) Type() string { return "withdraw_savings_interest" }

func (msg MsgWithdrawSavingsInterest) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

func (msg MsgWithdrawSavingsInterest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgWithdrawSavingsInterest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}