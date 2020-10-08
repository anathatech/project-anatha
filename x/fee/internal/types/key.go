package types

const (
	ModuleName = "fee"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName

	DefaultParamspace = ModuleName
)

var (
	FeeExcludedMessageKeyPrefix = []byte{0x10}

	StatusPresent = []byte{0x01}
)

func GetFeeExcludedMessageKey(msgType string) []byte {
	return append(FeeExcludedMessageKeyPrefix, []byte(msgType)...)
}

func GetFeeExcludedMessageIteratorKey() []byte {
	return FeeExcludedMessageKeyPrefix
}

func SplitFeeExcludedMessageKey(key []byte) (string) {
	return string(key[1:])
}