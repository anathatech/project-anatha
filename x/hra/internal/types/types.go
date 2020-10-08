package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/config"
	"time"
)

type NameInfo struct {
	Name string				`json:"name" yaml:"name"`
	Owner sdk.AccAddress	`json:"owner" yaml:"owner"`
	Price sdk.Coins			`json:"price" yaml:"price"`
	CreationTime time.Time	`json:"creation_time" yaml:"creation_time"`
	ExpiryTime time.Time	`json:"expiry_time" yaml:"expiry_time"`
}

func NewNameInfo(name string) NameInfo {
	return NameInfo{
		Name: name,
		Owner: nil,
		Price: sdk.NewCoins(sdk.NewInt64Coin(config.DefaultDenom, 0)),
		CreationTime: time.Time{},
		ExpiryTime: time.Time{},
	}
}

func (h NameInfo) String() string {
 	return fmt.Sprintf(`Name: %s
Owner: %s
Price: %s
Creation time: %s
Expiry time: %s`, h.Name, h.Owner, h.Price, h.CreationTime, h.ExpiryTime)
}


type BlockchainAddressRecordInfo struct {
	Address               sdk.AccAddress        `json:"address" yaml:"address"`
	BlockchainAddressInfo BlockchainAddressInfo `json:"blockchain_address_info" yaml:"blockchain_address_info"`
}

func NewBlockchainAddressRecordInfo(address sdk.AccAddress, blockchainAddress BlockchainAddressInfo) BlockchainAddressRecordInfo {
	return BlockchainAddressRecordInfo{
		Address:               address,
		BlockchainAddressInfo: blockchainAddress,
	}
}

func (a BlockchainAddressRecordInfo) String() string {
	return fmt.Sprintf(`Address: %s
Blockchain address: %s`, a.Address, a.BlockchainAddressInfo)
}

type BlockchainAddressInfo struct {
	BlockchainId string		`json:"blockchain_id" yaml:"blockchain_id"`
	Index string			`json:"index" yaml:"index"`
	BlockchainAddress string`json:"blockchain_address" yaml:"blockchain_address"`
}

func NewBlockchainAddressInfo(blockchainId string, index string, blockchainAddress string) BlockchainAddressInfo {
	return BlockchainAddressInfo{
		BlockchainId: blockchainId,
		Index: index,
		BlockchainAddress: blockchainAddress,
	}
}

func (a BlockchainAddressInfo) String() string {
	return fmt.Sprintf(`BlockchainId: %s
Index: %s
BlockchainAddress: %s`, a.BlockchainId, a.Index, a.BlockchainAddress)
}


type AddressCreditsInfo struct {
	Address sdk.AccAddress	`json:"address" yaml:"address"`
	Credits sdk.Int			`json:"credits" yaml:"credits"`
}

func NewAddressCreditsInfo(address sdk.AccAddress, credits sdk.Int) AddressCreditsInfo {
	return AddressCreditsInfo{Address:address,Credits:credits}
}

func (a AddressCreditsInfo) String() string {
	return fmt.Sprintf(`Address: %s
Credits: %s`, a.Address, a.Credits)
}