package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultGovernorList = []string {
		"anatha170mj6j6veall698u5xwkdhpf2nlza9j6n366v2",
		"anatha1e4ay99w2v3qqnl8wrmv3s85mexm6eum2vhucmf",
		"anatha1p470fhtytsym8fpp5j6rak2ec2qrvmwdvpg3ve",
	}
)

func DefaultGovernors() []sdk.AccAddress {
	var governors []sdk.AccAddress

	for _, item := range DefaultGovernorList {
		address, _ := sdk.AccAddressFromBech32(item)
		governors = append(governors, address)
	}

	return governors
}
