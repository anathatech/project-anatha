package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"time"
)

type ValidatorAccumulatedRewardRecord struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	Accumulated      sdk.DecCoins 	`json:"accumulated" yaml:"accumulated"`
}

type AddressNameRewardRateRecord struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Rate sdk.Dec `json:"rate" yaml:"rate"`
}

type NameDepositQueueRecord struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Time time.Time `json:"time" yaml:"time"`
}

type NameRewardEscrowRecord struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Amount sdk.DecCoins `json:"amount" yaml:"amount"`
}

type NameRewardLeftoverRecord struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Amount sdk.Dec `json:"amount" yaml:"amount"`
}

type AddressSavingsRewardRateRecord struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Rate sdk.Dec `json:"rate" yaml:"rate"`
}

type AddressSavingsStakeRecord struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Stake   sdk.Int        `json:"stake" yaml:"stake"`
}

type SavingsRewardEscrowRecord struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Amount sdk.DecCoins `json:"amount" yaml:"amount"`
}

type SavingsRewardLeftoverRecord struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Amount sdk.Dec `json:"amount" yaml:"amount"`
}

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`

	NameStake sdk.Dec `json:"name_stake" yaml:"name_stake"`
	NameRewardRate sdk.Dec `json:"name_reward_rate" yaml:"name_reward_rate"`
	AddressNameRewardRates []AddressNameRewardRateRecord `json:"address_name_reward_rates" yaml:"address_name_reward_rates"`
	PendingNameDistribution sdk.Coins `json:"pending_name_distribution" yaml:"pending_name_distribution"`
	NameDepositQueue []NameDepositQueueRecord `json:"name_deposit_queue" yaml:"name_deposit_queue"`
	NameRewardEscrow []NameRewardEscrowRecord `json:"name_reward_escrow" yaml:"name_reward_escrow"`
	NameRewardLeftover []NameRewardLeftoverRecord `json:"name_reward_leftover" yaml:"name_reward_leftover"`

	SavingsStake sdk.Int `json:"savings_stake" yaml:"savings_stake"`
	SavingsRewardRate sdk.Dec `json:"savings_reward_rate" yaml:"savings_reward_rate"`
	AddressSavingsRewardRates []AddressSavingsRewardRateRecord `json:"address_savings_reward_rates" yaml:"address_savings_reward_rates"`
	AddressSavingsStake []AddressSavingsStakeRecord `json:"address_savings_stake" yaml:"address_savings_stake"`
	SavingsRewardEscrow []SavingsRewardEscrowRecord `json:"savings_reward_escrow" yaml:"savings_reward_escrow"`
	SavingsRewardLeftover []SavingsRewardLeftoverRecord `json:"savings_reward_leftover" yaml:"savings_reward_leftover"`

	ValidatorAccumulatedRewards []ValidatorAccumulatedRewardRecord `json:"validator_accumulated_rewards" yaml:"validator_accumulated_rewards"`
	NvrpRemainder sdk.DecCoins `json:"nvrp_remainder" yaml:"nvrp_remainder"`
}

func NewGenesisState(params Params, nameStake sdk.Dec, nameRewardRate sdk.Dec, addressNameRewardRates []AddressNameRewardRateRecord,
	pendingNameDistribution sdk.Coins, nameDepositQueue []NameDepositQueueRecord, nameRewardEscrow []NameRewardEscrowRecord, nameRewardLeftover []NameRewardLeftoverRecord,
	savingsStake sdk.Int, savingsRewardRate sdk.Dec, addressSavingsRewardRates []AddressSavingsRewardRateRecord, addressSavingsStake []AddressSavingsStakeRecord,
	savingsRewardEscrow []SavingsRewardEscrowRecord, savingsRewardLeftover []SavingsRewardLeftoverRecord,
	validatorAccumulatedRewards []ValidatorAccumulatedRewardRecord, nvrpRemainder sdk.DecCoins) GenesisState {

	return GenesisState{
		Params: params,

		NameStake: nameStake,
		NameRewardRate: nameRewardRate,
		AddressNameRewardRates: addressNameRewardRates,
		PendingNameDistribution: pendingNameDistribution,
		NameDepositQueue: nameDepositQueue,
		NameRewardEscrow: nameRewardEscrow,
		NameRewardLeftover: nameRewardLeftover,

		SavingsStake: savingsStake,
		SavingsRewardRate: savingsRewardRate,
		AddressSavingsRewardRates: addressSavingsRewardRates,
		AddressSavingsStake: addressSavingsStake,
		SavingsRewardEscrow: savingsRewardEscrow,
		SavingsRewardLeftover: savingsRewardLeftover,

		ValidatorAccumulatedRewards: validatorAccumulatedRewards,
		NvrpRemainder: nvrpRemainder,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),

		NameStake: sdk.ZeroDec(),
		NameRewardRate: sdk.ZeroDec(),
		AddressNameRewardRates: []AddressNameRewardRateRecord{},
		PendingNameDistribution: sdk.NewCoins(),
		NameDepositQueue: []NameDepositQueueRecord{},
		NameRewardEscrow: []NameRewardEscrowRecord{},
		NameRewardLeftover: []NameRewardLeftoverRecord{},

		SavingsStake: sdk.ZeroInt(),
		SavingsRewardRate: sdk.ZeroDec(),
		AddressSavingsRewardRates: []AddressSavingsRewardRateRecord{},
		AddressSavingsStake: []AddressSavingsStakeRecord{},
		SavingsRewardEscrow: []SavingsRewardEscrowRecord{},
		SavingsRewardLeftover: []SavingsRewardLeftoverRecord{},

		ValidatorAccumulatedRewards: []ValidatorAccumulatedRewardRecord{},
		NvrpRemainder: sdk.NewDecCoins(),
	}
}

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	for _, record := range data.AddressNameRewardRates {
		if record.Address.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, record.Address.String())
		}
		if record.Rate.IsZero() || record.Rate.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, record.Rate.String())
		}
	}

	for _, record := range data.NameDepositQueue {
		if record.Address.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, record.Address.String())
		}
	}

	for _, record := range data.NameRewardEscrow {
		if record.Address.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, record.Address.String())
		}
		if record.Amount.IsAnyNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, record.Amount.String())
		}
	}

	for _, record := range data.NameRewardLeftover {
		if record.Address.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, record.Address.String())
		}
		if record.Amount.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, record.Amount.String())
		}
	}

	for _, record := range data.AddressSavingsRewardRates {
		if record.Address.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, record.Address.String())
		}
		if record.Rate.IsZero() || record.Rate.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, record.Rate.String())
		}
	}

	for _, record := range data.AddressSavingsStake {
		if record.Address.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, record.Address.String())
		}
		if record.Stake.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, record.Stake.String())
		}
	}

	for _, record := range data.SavingsRewardEscrow {
		if record.Address.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, record.Address.String())
		}
		if record.Amount.IsAnyNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, record.Amount.String())
		}
	}

	for _, record := range data.SavingsRewardLeftover {
		if record.Address.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, record.Address.String())
		}
		if record.Amount.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, record.Amount.String())
		}
	}

	for _, record := range data.ValidatorAccumulatedRewards {
		if record.ValidatorAddress.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, record.ValidatorAddress.String())
		}
		if record.Accumulated.IsAnyNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, record.Accumulated.String())
		}

	}

	return nil
}