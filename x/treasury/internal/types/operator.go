package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultOperatorsList = []string {
		"anatha18klyq6qxyemrgsuyapw80y7xz4gh0lhjk9tj7n",
		"anatha1lafvzw2f6nj4sv2k0cpsu4zdxur9l0e98ngzxz",
		"anatha13fsp5p3mjcxtr833vcln6kjwvyl8tz66tu2j8z",
	}
)

func DefaultOperators() []sdk.AccAddress {
	var operators []sdk.AccAddress

	for _, item := range DefaultOperatorsList {
		address, _ := sdk.AccAddressFromBech32(item)
		operators = append(operators, address)
	}

	return operators
}