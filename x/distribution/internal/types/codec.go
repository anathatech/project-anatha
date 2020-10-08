package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(DevelopmentFundDistributionProposal{}, "distribution/DevelopmentFundDistributionProposal", nil)
	cdc.RegisterConcrete(SecurityTokenFundDistributionProposal{}, "distribution/SecurityTokenFundDistributionProposal", nil)

	cdc.RegisterConcrete(MsgWithdrawNameReward{}, "distribution/WithdrawNameReward", nil)
	cdc.RegisterConcrete(MsgWithdrawValidatorReward{}, "distribution/WithdrawValidatorReward", nil)

	cdc.RegisterConcrete(MsgDepositSavings{}, "distribution/DepositSavings", nil)
	cdc.RegisterConcrete(MsgWithdrawSavings{}, "distribution/WithdrawSavings", nil)
	cdc.RegisterConcrete(MsgWithdrawSavingsInterest{}, "distribution/WithdrawSavingsInterest", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
