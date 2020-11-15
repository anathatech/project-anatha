package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"time"
)

const (
	ModuleName = "hra"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName

	Separator = ":"
)

// Keys for HRA store
// Items are stored with the following key: values
//
// - 0x10<Name_Bytes>: NameInfo
// - 0x11<Addr_Bytes><Separator><Name_Bytes>: boolean
// - 0x12<endTime_Bytes><Name_Bytes>: Name
// - 0x13<Addr_Bytes><Separator><BlockchainId_Bytes><Separator><AddressIndex_Bytes>: BlockchainAddress
// - 0x14<Addr_Bytes>: Int
// - 0x15<BlockchainId_Bytes>: boolean
var (
	NameInfoByNameKeyPrefix         = []byte{0x10}
	StatusByAddressAndNameKeyPrefix = []byte{0x11}
	ExpiredNameInfoQueueKeyPrefix   = []byte{0x12}
	AddressKeyPrefix                = []byte{0x13}
	CreditsKeyPrefix                = []byte{0x14}
	RegisteredBlockchainIdKeyPrefix = []byte{0x15}

	StatusPresent = []byte{0x01}
	StatusAbsent = []byte{0x00}
)

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

func GetNameInfoByNameKey(name string) []byte {
	return append(NameInfoByNameKeyPrefix, []byte(name)...)
}

func GetStatusByAddressAndNameKey(address sdk.AccAddress, name string) []byte {
	key := append(StatusByAddressAndNameKeyPrefix, address...)
	key = append(key, []byte(Separator)...)
	key = append(key, []byte(name)...)
	return key
}

// ExpiredNameInfoByTimeKey gets the inactive proposal queue key by endTime
func ExpiredNameInfoByTimeKey(endTime time.Time) []byte {
	return append(ExpiredNameInfoQueueKeyPrefix, sdk.FormatTimeBytes(endTime)...)
}

func SplitExpiredNameInfoQueueKey(key []byte) (name string, endTime time.Time) {
	return splitKeyWithTime(key)
}

func ExpiredNameInfoQueueKey(name string, endTime time.Time) []byte {
	return append(ExpiredNameInfoByTimeKey(endTime), []byte(name)...)
}

func GetStatusByAddressAndNameIteratorKey(address sdk.AccAddress) []byte {
	return append(StatusByAddressAndNameKeyPrefix, address...)
}

// Address
func GetAddressKey(address sdk.AccAddress, blockchainId string, index string) []byte {
	// Craft KVStore key in format: AnathaAddress:BlockchainId:AddressIndex
	key := append(AddressKeyPrefix,[]byte(address)...)
	key = append(key, []byte(Separator)...)
	key = append(key, []byte(blockchainId)...)
	key = append(key, []byte(Separator)...)
	key = append(key, []byte(index)...)

	return key
}

func GetBlockchainAddressIteratorKey(address sdk.AccAddress) []byte {
	return append(AddressKeyPrefix, address...)
}

func SplitBlockchainAddressKey(key []byte) (blockchainAddress BlockchainAddressInfo)  {
	parts := strings.Split(string(key[22:]), Separator) // prefix + address + first separator

	return NewBlockchainAddressInfo(parts[0], parts[1], "")
}

func SplitBlockchainAddressRecordKey(key []byte) (blockchainAddressRecord BlockchainAddressRecordInfo)  {
	parts := strings.Split(string(key[22:]), Separator) // prefix + address + first separator

	blockchainAddressInfo := NewBlockchainAddressInfo(parts[0], parts[1], "")

	return NewBlockchainAddressRecordInfo(sdk.AccAddress(key[1:21]), blockchainAddressInfo)
}

func GetCreditsKey(address sdk.AccAddress) []byte {
	return append(CreditsKeyPrefix, address...)
}

func GetRegisteredBlockchainIdKey(blockchainId string) []byte {
	return append(RegisteredBlockchainIdKeyPrefix, []byte(blockchainId)...)
}

func GetRegisteredBlockchainIdIteratorKey() []byte {
	return RegisteredBlockchainIdKeyPrefix
}

func SplitRegisteredBlockchainIdKey(key []byte) (string) {
	return string(key[1:])
}
// private functions

func splitKeyWithTime(key []byte) (name string, endTime time.Time) {
	/*
		if len(key[1:]) != 8 + lenTime {
			panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), lenTime + 8))
		}
	*/

	endTime, err := sdk.ParseTimeBytes(key[1 : 1 + lenTime])
	if err != nil {
		panic(err)
	}

	name = string(key[1 + lenTime:])
	return
}