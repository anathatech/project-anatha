package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(AddFeeExcludedMessageProposal{}, "fee/AddFeeExcludedMessageProposal", nil)
	cdc.RegisterConcrete(RemoveFeeExcludedMessageProposal{}, "fee/RemoveFeeExcludedMessageProposal", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)

	codec.RegisterCrypto(ModuleCdc)

	ModuleCdc.Seal()
}
