package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"regexp"
	"strconv"
)

const (
	validChar = `[a-z0-9\.,\+\-_]`
	blockchainAddressMaxLen = 128
)

var (
	validBlockchainId = regexp.MustCompile(`^[a-zA-Z0-9_.-]{3,32}$`).MatchString
	validName = regexp.MustCompile(`^` + validChar + `{3,64}$`).MatchString
)

func validateName(name string) error {
	if ! validName(name) {
		return ErrNameNotValid
	}

	return nil
}

func validateBlockchainId(blockchainId string) error {
	if ! validBlockchainId(blockchainId) {
		return ErrBlockchainIdNotValid
	}

	return nil
}

func validateIndex(index string) error {
	if _, err := strconv.Atoi(index); err != nil {
		return ErrAddressIndexNotValid
	}

	return nil
}

func validateBlockchainAddress(blockchainAddress string) error {
	length := len(blockchainAddress)

	if length == 0 {
		return sdkerrors.Wrap(ErrBlockchainAddressNotValid,"address is required")
	}
	if length > blockchainAddressMaxLen {
		return sdkerrors.Wrap(ErrBlockchainAddressNotValid, "address too long")
	}

	return nil
}
