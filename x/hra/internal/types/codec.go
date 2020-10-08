package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRegisterName{}, "hra/Register", nil)
	cdc.RegisterConcrete(MsgRenewName{}, "hra/Renew", nil)
	cdc.RegisterConcrete(MsgSetPrice{}, "hra/SetPrice", nil)
	cdc.RegisterConcrete(MsgDeleteName{}, "hra/Delete", nil)
	cdc.RegisterConcrete(MsgBuyName{}, "hra/Buy", nil)
	cdc.RegisterConcrete(MsgTransferName{}, "hra/Transfer", nil)
	cdc.RegisterConcrete(MsgRegisterAddress{}, "hra/RegisterAddress", nil)
	cdc.RegisterConcrete(MsgRemoveAddress{}, "hra/RemoveAddress", nil)
	cdc.RegisterConcrete(MsgRemoveAllAddresses{}, "hra/RemoveAllAddresses", nil)

	cdc.RegisterConcrete(RegisterBlockchainIdProposal{}, "hra/RegisterBlockchainIdProposal", nil)
	cdc.RegisterConcrete(RemoveBlockchainIdProposal{}, "hra/RemoveBlockchainIdProposal", nil)
}

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)

	codec.RegisterCrypto(ModuleCdc)

	ModuleCdc.Seal()
}

