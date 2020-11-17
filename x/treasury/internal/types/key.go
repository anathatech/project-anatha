package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

var (
	TreasuryKey = []byte{0x00}

	OperatorKeyPrefix          = []byte{0x10}
	DisbursementQueueKeyPrefix = []byte{0x12}

	DisbursementReferenceKeyPrefix = []byte{0x15}

	StatusPresent = []byte{0x01}
)

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

const (
	// ModuleName is the name of the module
	ModuleName  					= "treasury"				// Module stores initial supply of Anatha

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName

	BuyBackLiquidityFundModuleName 	= "buyback_liquidity" 		// Module stores the supply of AST for executing BuyBacks
	BuyBackFundModuleName 			= "treasury_buyback" 		// Module stores the bought back supply of Anatha (participates in distribution)
	DistributionProfitsModuleName  	= "distribution_profits" 	// Module stores AST profits from Anatha purchases with AST
	TreasuryEscrowModuleName        = "treasury_escrow"         // Module stores distributions without supplied recipient addresses
	SwapEscrowModuleName			= "swap_escrow"             // Module stores ERC20 token balance
)


func GetOperatorKey(address sdk.AccAddress) []byte {
	return append(OperatorKeyPrefix, address...)
}

func GetOperatorIteratorKey() []byte {
	return OperatorKeyPrefix
}

func SplitOperatorKey(key []byte) sdk.AccAddress {
	return key[1:]
}

func DisbursementByTimeKey(endTime time.Time) []byte {
	return append(DisbursementQueueKeyPrefix, sdk.FormatTimeBytes(endTime)...)
}

func DisbursementQueueKey(address sdk.AccAddress, endTime time.Time) []byte {
	return append(DisbursementByTimeKey(endTime), address...)
}

func GetDisbursementReferenceKey(reference string) []byte {
	return append(DisbursementReferenceKeyPrefix, []byte(reference)...)
}

func GetDisbursementReferenceIteratorKey() []byte {
	return DisbursementReferenceKeyPrefix
}

func SplitDisbursementReferenceKey(key []byte) (string) {
	return string(key[1:])
}