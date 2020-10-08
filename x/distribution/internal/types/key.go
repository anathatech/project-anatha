package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

const (
	ModuleName = "distribution"

	AmcModuleName = "amc"
	NvrpModuleName = "nvrp"
	NvrpDistributionModuleName = "nvrpd"

	HRAHolderRewardModuleName   = "hhrm"
	DevelopmentFundModuleName   = "dfm"
	SecurityTokenFundModuleName = "stfm"

	SavingsModuleName = "savings"
	SavingsDistributionModuleName = "savingsdistr"

	DefaultParamspace = ModuleName
	StoreKey = ModuleName
	QuerierRoute = StoreKey
	RouterKey = ModuleName
)

var (
	NameStakeKey                         = []byte{0x10}
	NameRewardRateKey                    = []byte{0x11}
	NameRewardRateByAddressKeyPrefix     = []byte{0x12}
	PendingNameDistributionKey           = []byte{0x13}
	NameDepositQueueKeyPrefix            = []byte{0x14}
	NameRewardEscrowKeyPrefix            = []byte{0x15}
	NameRewardLeftoverKeyPrefix          = []byte{0x16}

	SavingsStakeKey                      = []byte{0x20}
	SavingsStakeByAddressKeyPrefix       = []byte{0x21}
	SavingsRewardRateKey                 = []byte{0x22}
	SavingsRewardRateByAddressKeyPrefix  = []byte{0x23}
	SavingsRewardEscrowKeyPrefix         = []byte{0x24}
	SavingsRewardLeftoverKeyPrefix       = []byte{0x25}

	ValidatorAccumulatedRewardsKeyPrefix = []byte{0x30} // key for accumulated validator rewards
	NvrpdRemainderKey                    = []byte{0x31}
)

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

func GetNameStakeKey() []byte {
	return NameStakeKey
}

func GetNameRewardRateKey() []byte {
	return NameRewardRateKey
}

func GetNameRewardByAddressKey(address sdk.AccAddress) []byte {
	return append(NameRewardRateByAddressKeyPrefix, address...)
}

func GetNameRewardAddress(key []byte) (address sdk.AccAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.AccAddress(addr)
}

func GetPendingNameDistributionKey() []byte {
	return PendingNameDistributionKey
}

func NameDepositByTimeKey(endTime time.Time) []byte {
	return append(NameDepositQueueKeyPrefix, sdk.FormatTimeBytes(endTime)...)
}

func NameDepositQueueKey(address sdk.AccAddress, endTime time.Time) []byte {
	return append(NameDepositByTimeKey(endTime), address...)
}

func SplitNameDepositQueueKey(key []byte) (address sdk.AccAddress, endTime time.Time) {
	return splitKeyWithTime(key)
}

func GetNameRewardEscrowKey(address sdk.AccAddress) []byte {
	return append(NameRewardEscrowKeyPrefix, address...)
}

func GetNameRewardEscrowAddress(key []byte) (address sdk.AccAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.AccAddress(addr)
}

func GetNameRewardLeftoverKey(address sdk.AccAddress) []byte {
	return append(NameRewardLeftoverKeyPrefix, address...)
}

func GetNameRewardLeftoverAddress(key []byte) (address sdk.AccAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.AccAddress(addr)
}

// gets the key for a validator's current rewards
func GetValidatorAccumulatedRewardsKey(v sdk.ValAddress) []byte {
	return append(ValidatorAccumulatedRewardsKeyPrefix, v.Bytes()...)
}

// gets the address from a validator's accumulated rewards key
func GetValidatorAccumulatedRewardsAddress(key []byte) (valAddr sdk.ValAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.ValAddress(addr)
}

func GetNvrpdRemainderKey() []byte {
	return NvrpdRemainderKey
}

// Savings

func GetSavingsStakeKey() []byte {
	return SavingsStakeKey
}

func GetSavingsRewardRateKey() []byte {
	return SavingsRewardRateKey
}

func GetSavingsRewardByAddressKey(address sdk.AccAddress) []byte {
	return append(SavingsRewardRateByAddressKeyPrefix, address...)
}

func GetSavingsRewardAddress(key []byte) (address sdk.AccAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.AccAddress(addr)
}

func GetSavingsStakeByAddressKey(address sdk.AccAddress) []byte {
	return append(SavingsStakeByAddressKeyPrefix, address...)
}

func GetSavingsStakeAddress(key []byte) (address sdk.AccAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.AccAddress(addr)
}

func GetSavingsRewardEscrowKey(address sdk.AccAddress) []byte {
	return append(SavingsRewardEscrowKeyPrefix, address...)
}

func GetSavingsRewardEscrowAddress(key []byte) (address sdk.AccAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.AccAddress(addr)
}

func GetSavingsRewardLeftoverKey(address sdk.AccAddress) []byte {
	return append(SavingsRewardLeftoverKeyPrefix, address...)
}

func GetSavingsRewardLeftoverAddress(key []byte) (address sdk.AccAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.AccAddress(addr)
}

// Internal

func splitKeyWithTime(key []byte) (address sdk.AccAddress, endTime time.Time) {
	if len(key[1:]) != sdk.AddrLen + lenTime {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), sdk.AddrLen + lenTime ))
	}


	endTime, err := sdk.ParseTimeBytes(key[1 : 1 + lenTime])
	if err != nil {
		panic(err)
	}

	address = key[1 + lenTime:]

	return
}