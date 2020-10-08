package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/config"
)

// MsgRegisterName
type MsgRegisterName struct {
	Name  string         `json:"name" yaml:"name"`
	Owner sdk.AccAddress `json:"owner" yaml:"owner"`
}

func NewMsgRegisterName(name string, owner sdk.AccAddress) MsgRegisterName {
	return MsgRegisterName{
		Name:  name,
		Owner: owner,
	}
}

func (msg MsgRegisterName) Route() string { return RouterKey }

func (msg MsgRegisterName) Type() string { return "register_name" }

func (msg MsgRegisterName) ValidateBasic() error {
	err := validateName(msg.Name)
	if err != nil {
		return err
	}
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	return nil
}

func (msg MsgRegisterName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRegisterName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgRenewName
type MsgRenewName struct {
	Name  string         `json:"name" yaml:"name"`
	Owner sdk.AccAddress `json:"owner" yaml:"owner"`
}

func NewMsgRenewName(name string, owner sdk.AccAddress) MsgRenewName {
	return MsgRenewName{
		Name:  name,
		Owner: owner,
	}
}

func (msg MsgRenewName) Route() string { return RouterKey }

func (msg MsgRenewName) Type() string { return "renew_name" }

func (msg MsgRenewName) ValidateBasic() error {
	err := validateName(msg.Name)
	if err != nil {
		return err
	}
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	return nil
}

func (msg MsgRenewName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRenewName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgSetPrice
type MsgSetPrice struct {
	Name  string         	`json:"name" yaml:"name"`
	Owner sdk.AccAddress 	`json:"owner" yaml:"owner"`
	Price sdk.Coins			`json:"price" yaml:"price"`
}

func NewMsgSetPrice(name string, owner sdk.AccAddress, price sdk.Coins) MsgSetPrice {
	return MsgSetPrice{
		Name:  name,
		Owner: owner,
		Price: price,
	}
}

func (msg MsgSetPrice) Route() string { return RouterKey }

func (msg MsgSetPrice) Type() string { return "set_price" }

func (msg MsgSetPrice) ValidateBasic() error {
	err := validateName(msg.Name)
	if err != nil {
		return err
	}
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if msg.Price.AmountOf(config.DefaultDenom).IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid denomination.")
	}
	return nil
}

func (msg MsgSetPrice) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSetPrice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgDeleteName
type MsgDeleteName struct {
	Name  string         	`json:"name" yaml:"name"`
	Owner sdk.AccAddress 	`json:"owner" yaml:"owner"`
}

func NewMsgDeleteName(name string, owner sdk.AccAddress) MsgDeleteName {
	return MsgDeleteName{
		Name:  name,
		Owner: owner,
	}
}

func (msg MsgDeleteName) Route() string { return RouterKey }

func (msg MsgDeleteName) Type() string { return "delete_name" }

func (msg MsgDeleteName) ValidateBasic() error {
	err := validateName(msg.Name)
	if err != nil {
		return err
	}
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	return nil
}

func (msg MsgDeleteName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDeleteName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgBuyName
type MsgBuyName struct {
	Name  string         	`json:"name" yaml:"name"`
	Buyer sdk.AccAddress 	`json:"buyer" yaml:"buyer"`
}

func NewMsgBuyName(name string, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:  name,
		Buyer: buyer,
	}
}

func (msg MsgBuyName) Route() string { return RouterKey }

func (msg MsgBuyName) Type() string { return "buy_name" }

func (msg MsgBuyName) ValidateBasic() error {
	err := validateName(msg.Name)
	if err != nil {
		return err
	}
	if msg.Buyer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Buyer.String())
	}

	return nil
}

func (msg MsgBuyName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

// MsgTransferName
type MsgTransferName struct {
	Name  string         	`json:"name" yaml:"name"`
	Owner sdk.AccAddress 	`json:"owner" yaml:"owner"`
	NewOwner sdk.AccAddress	`json:"new_owner" yaml:"new_owner"`
}

func NewMsgTransferName(name string, owner sdk.AccAddress, newOwner sdk.AccAddress) MsgTransferName {
	return MsgTransferName{
		Name:  name,
		Owner: owner,
		NewOwner: newOwner,
	}
}

func (msg MsgTransferName) Route() string { return RouterKey }

func (msg MsgTransferName) Type() string { return "transfer_name" }

func (msg MsgTransferName) ValidateBasic() error {
	err := validateName(msg.Name)
	if err != nil {
		return err
	}
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if msg.NewOwner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.NewOwner.String())
	}

	return nil
}

func (msg MsgTransferName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgRegisterAddress
type MsgRegisterAddress struct {
	Owner sdk.AccAddress 	`json:"owner" yaml:"owner"`
	BlockchainId string		`json:"blockchain_id" yaml:"blockchain_id"`
	Index string			`json:"index" yaml:"index"`
	BlockchainAddress string`json:"blockchain_address" yaml:"blockchain_address"`
}

func NewMsgRegisterAddress(owner sdk.AccAddress, blockchainId string, index string, blockchainAddress string) MsgRegisterAddress {
	return MsgRegisterAddress{
		Owner: owner,
		BlockchainId: blockchainId,
		Index: index,
		BlockchainAddress: blockchainAddress,
	}
}

func (msg MsgRegisterAddress) Route() string { return RouterKey }

func (msg MsgRegisterAddress) Type() string { return "register_address" }

func (msg MsgRegisterAddress) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	err := validateBlockchainId(msg.BlockchainId)
	if err != nil {
		return err
	}

	err = validateIndex(msg.Index)
	if err != nil {
		return err
	}

	err = validateBlockchainAddress(msg.BlockchainAddress)
	if err != nil {
		return err
	}

	return nil
}

func (msg MsgRegisterAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRegisterAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgRemoveAddress
type MsgRemoveAddress struct {
	Owner sdk.AccAddress 	`json:"owner" yaml:"owner"`
	BlockchainId string		`json:"blockchain_id" yaml:"blockchain_id"`
	Index string			`json:"index" yaml:"index"`
}

func NewMsgRemoveAddress(owner sdk.AccAddress, blockchainId string, index string) MsgRemoveAddress {
	return MsgRemoveAddress{
		Owner: owner,
		BlockchainId: blockchainId,
		Index: index,
	}
}

func (msg MsgRemoveAddress) Route() string { return RouterKey }

func (msg MsgRemoveAddress) Type() string { return "remove_address" }

func (msg MsgRemoveAddress) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	err := validateBlockchainId(msg.BlockchainId)
	if err != nil {
		return err
	}

	err = validateIndex(msg.Index)
	if err != nil {
		return err
	}

	return nil
}

func (msg MsgRemoveAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRemoveAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgRemoveAllAddresses
type MsgRemoveAllAddresses struct {
	Owner sdk.AccAddress 	`json:"owner" yaml:"owner"`
}

func NewMsgRemoveAllAddresses(owner sdk.AccAddress) MsgRemoveAllAddresses {
	return MsgRemoveAllAddresses{
		Owner: owner,
	}
}

func (msg MsgRemoveAllAddresses) Route() string { return RouterKey }

func (msg MsgRemoveAllAddresses) Type() string { return "remove_all_addresses" }

func (msg MsgRemoveAllAddresses) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	return nil
}

func (msg MsgRemoveAllAddresses) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRemoveAllAddresses) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}