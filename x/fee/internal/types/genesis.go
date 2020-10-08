package types

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
	FeeExcludedMessages []string `json:"fee_excluded_messages" yaml:"fee_excluded_messages"`
}

func NewGenesisState(params Params, feeExcludedMessages []string) GenesisState {
	return GenesisState{
		Params: params,
		FeeExcludedMessages: feeExcludedMessages,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
		FeeExcludedMessages: DefaultFeeExcludedMessages,
	}
}

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return nil
}