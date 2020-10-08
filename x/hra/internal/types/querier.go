package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)


type QueryResNameInfo struct {
	NameInfo 	NameInfo 					`json:"name_info" yaml:"name_info"`
	Credits 	sdk.Int 					`json:"credits" yaml:"credits"`
	Addresses 	[]BlockchainAddressInfo 	`json:"addresses" yaml:"addresses"`
}

func (n QueryResNameInfo) String() string {
	return fmt.Sprintf(`%s
%s
%s`, n.NameInfo, n.Credits, n.Addresses)
}

type QueryResNameInfos []NameInfo

func (n QueryResNameInfos) String() string {
	var nameInfos []string

	for _, nameInfo := range n {
		nameInfos = append(nameInfos, nameInfo.String())
	}

	return strings.Join(nameInfos, "\n")
}

type QueryResNames []string

func (n QueryResNames) String() string {
	return strings.Join(n[:], "\n")
}

type QueryResAddress string

func (n QueryResAddress) String() string {
	return string(n)
}

type QueryResRegisteredBlockchainIds []string

func (n QueryResRegisteredBlockchainIds) String() string {
	return strings.Join(n[:], "\n")
}

type QueryResAddressCredits string

func (n QueryResAddressCredits) String() string {
	return string(n)
}

type QueryResBlockchainAddresses []BlockchainAddressInfo

func (n QueryResBlockchainAddresses) String() string {
	return fmt.Sprintf(`%s`, n)
}
