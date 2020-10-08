package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type GenesisState struct {
	Params 					Params 		`json:"params" yaml:"params"`
	NameRecords 			[]NameInfo 	`json:"name_records" yaml:"name_records"`
	AddressRecords			[]BlockchainAddressRecordInfo `json:"address_records" yaml:"address_records"`
	AddressCredits          []AddressCreditsInfo    `json:"address_credits" yaml:"address_credits"`
	RegisteredBlockchainIds []string	`json:"registered_blockchain_ids" yaml:"registered_blockchain_ids"`
}


func NewGenesisState(params Params, nameRecords []NameInfo, addressRecords []BlockchainAddressRecordInfo, addressCredits []AddressCreditsInfo, registeredBlockchainIds []string) GenesisState {
	return GenesisState{
		Params: params,
		NameRecords: nameRecords,
		AddressRecords: addressRecords,
		AddressCredits: addressCredits,
		RegisteredBlockchainIds: registeredBlockchainIds,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
		NameRecords: []NameInfo{},
		AddressRecords: []BlockchainAddressRecordInfo{},
		AddressCredits: []AddressCreditsInfo{},
		RegisteredBlockchainIds: DefaultRegisteredBlockchainIds,
	}
}

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	for _, record := range data.NameRecords {
		if record.Owner.Empty() {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid NameRecord: Value: %s. Error: Missing owner", record.Owner)
		}
		err := validateName(record.Name)
		if err != nil {
			return sdkerrors.Wrapf(err, "invalid NameRecord: Name: %s. Error: Missing name", record.Name)
		}
	}
	for _, record := range data.AddressRecords {
		if record.Address.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, record.Address.String())
		}
		err := validateBlockchainId(record.BlockchainAddressInfo.BlockchainId)
		if err != nil {
			return err
		}

		err = validateIndex(record.BlockchainAddressInfo.Index)
		if err != nil {
			return err
		}

		err = validateBlockchainAddress(record.BlockchainAddressInfo.BlockchainAddress)
		if err != nil {
			return err
		}
	}
	for _, record := range data.AddressCredits {
		if record.Address.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, record.Address.String())
		}
	}
	for _, record := range data.RegisteredBlockchainIds {
		err := validateBlockchainId(record)
		if err != nil {
			return err
		}
	}
	return nil
}