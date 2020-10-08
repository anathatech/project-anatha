package types

import (
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

const (
	ModuleName = "governance"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName

	DefaultParamspace = ModuleName
)


// - 0x00<proposalID_Bytes>: Proposal
// - 0x01<endTime_Bytes><proposalID_Bytes>: activeProposalID
// - 0x02<endTime_Bytes><proposalID_Bytes>: inactiveProposalID
// - 0x03: nextProposalID
// - 0x10<proposalID_Bytes><depositorAddr_Bytes>: Deposit
// - 0x20<proposalID_Bytes><voterAddr_Bytes>: Voter
var (
	ProposalsKeyPrefix          = []byte{0x00}
	ActiveProposalQueuePrefix   = []byte{0x01}
	ExpeditedProposalQueuePrefix = []byte{0x04}
	ProposalIDKey               = []byte{0x02}
	GovernorCountKey			= []byte{0x03}

	VotesKeyPrefix = []byte{0x10}
	GovernorKeyPrefix = []byte{0x21}

	StatusPresent = []byte{0x01}
)

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

func GetProposalIDBytes(proposalID uint64) (proposalIDBz []byte) {
	proposalIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(proposalIDBz, proposalID)
	return
}

func GetProposalIDFromBytes(bz []byte) (proposalID uint64) {
	return binary.BigEndian.Uint64(bz)
}

func ProposalKey(proposalID uint64) []byte {
	return append(ProposalsKeyPrefix, GetProposalIDBytes(proposalID)...)
}

func ActiveProposalByTimeKey(endTime time.Time) []byte {
	return append(ActiveProposalQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

func ActiveProposalQueueKey(proposalID uint64, endTime time.Time) []byte {
	return append(ActiveProposalByTimeKey(endTime), GetProposalIDBytes(proposalID)...)
}

func ExpeditedProposalQueueKey(proposalID uint64) []byte {
	return append(ExpeditedProposalQueuePrefix, GetProposalIDBytes(proposalID)...)
}

func VotesKey(proposalID uint64) []byte {
	return append(VotesKeyPrefix, GetProposalIDBytes(proposalID)...)
}

func VoteKey(proposalID uint64, voterAddr sdk.AccAddress) []byte {
	return append(VotesKey(proposalID), voterAddr.Bytes()...)
}

func SplitProposalKey(key []byte) (proposalID uint64) {
	if len(key[1:]) != 8 {
		panic(fmt.Sprintf("unexpected key length (%d ≠ 8)", len(key[1:])))
	}

	return GetProposalIDFromBytes(key[1:])
}

func SplitActiveProposalQueueKey(key []byte) (proposalID uint64, endTime time.Time) {
	return splitKeyWithTime(key)
}

func SplitExpeditedProposalQueueKey(key []byte) (proposalID uint64) {
	return splitKey(key)
}

func SplitKeyVote(key []byte) (proposalID uint64, voterAddr sdk.AccAddress) {
	return splitKeyWithAddress(key)
}

// Governor

func GetGovernorKey(address sdk.AccAddress) []byte {
	return append(GovernorKeyPrefix, address...)
}

func GetGovernorIteratorKey() []byte {
	return GovernorKeyPrefix
}

func SplitGovernorKey(key []byte) sdk.AccAddress {
	return key[1:]
}

func GetGovernorCountBytes(count uint64) (countBz []byte) {
	countBz = make([]byte, 8)
	binary.BigEndian.PutUint64(countBz, count)
	return
}

func GetGovernorCountFromBytes(bz []byte) (count uint64) {
	return binary.BigEndian.Uint64(bz)
}

// Internal

func splitKeyWithTime(key []byte) (proposalID uint64, endTime time.Time) {
	if len(key[1:]) != 8+lenTime {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), lenTime+8))
	}

	endTime, err := sdk.ParseTimeBytes(key[1 : 1+lenTime])
	if err != nil {
		panic(err)
	}

	proposalID = GetProposalIDFromBytes(key[1+lenTime:])
	return
}

func splitKeyWithAddress(key []byte) (proposalID uint64, addr sdk.AccAddress) {
	if len(key[1:]) != 8+sdk.AddrLen {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key), 8+sdk.AddrLen))
	}

	proposalID = GetProposalIDFromBytes(key[1:9])
	addr = sdk.AccAddress(key[9:])
	return
}

func splitKey(key []byte) (proposalID uint64) {
	if len(key[1:]) != 8 {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key), 8))
	}

	proposalID = GetProposalIDFromBytes(key[1:])
	return
}