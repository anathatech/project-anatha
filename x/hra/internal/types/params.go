package types

import (
	"fmt"
	"github.com/anathatech/project-anatha/config"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = ModuleName

	DefaultNameInfoDuration 	= time.Hour * 24 * 365

	DefaultNameInfoMaxDuration  = time.Hour * 24 * 365 * 3

	NameConstraintBlock = 750
)

var (
	DefaultAddressCredits		= sdk.NewInt(20)

	DefaultRegisteredBlockchainIds = []string {
		"omg",
		"tusd",
		"zrx",
		"btc",
		"bch",
		"eth",
		"dash",
		"ltc",
		"atom",
		"xrp",
		"bnb",
	}

	KeyNameInfoDuration			= []byte("NameInfoDuration")

	KeyNameInfoRegistrationFee 	= []byte("NameInfoRegistrationFee")
	KeyNameInfoRenewalFee 		= []byte("NameInfoRenewalFee")

	KeyAddressCredits 			= []byte("AddressCredits")
	KeyAddressRegistrationFee 	= []byte("AddressRegistrationFee")

	KeyNameInfoMaxDuration     = []byte("NameInfoMaxDuration")
)

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

type Params struct {
	NameInfoDuration		time.Duration 	`json:"nameinfo_duration" yaml:"nameinfo_duration"`
	NameInfoMaxDuration 	time.Duration 	`json:"nameinfo_max_duration" yaml:"nameinfo_max_duration"`
	NameInfoRegistrationFee sdk.Coins 		`json:"registration_fee" yaml:"registration_fee"`
	NameInfoRenewalFee 		sdk.Coins 		`json:"renewal_fee" yaml:"renewal_fee"`
	AddressCredits 			sdk.Int			`json:"address_credits" yaml:"address_credits"`
	AddressRegistrationFee	sdk.Coins		`json:"address_registration_fee" yaml:"address_registration_fee"`
}

func NewParams(nameInfoDuration time.Duration, nameInfoMaxDuration time.Duration, nameInfoRegistrationFee sdk.Coins,
	nameInfoRenewalFee sdk.Coins, addressCredits sdk.Int, addressRegistrationFee sdk.Coins) Params {
	return Params{
		NameInfoDuration: nameInfoDuration,
		NameInfoMaxDuration: nameInfoMaxDuration,
		NameInfoRegistrationFee: nameInfoRegistrationFee,
		NameInfoRenewalFee: nameInfoRenewalFee,
		AddressCredits: addressCredits,
		AddressRegistrationFee: addressRegistrationFee,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  NameInfo Duration:     %s
  NameInfo Max Duration: %s
  NameInfoRegistrationFee: %s
  NameInfoRenewalFee: %s
  AddressCredits: %s
  AddressRegistrationFee: %s`,
		p.NameInfoDuration,
		p.NameInfoMaxDuration,
		p.NameInfoRegistrationFee,
		p.NameInfoRenewalFee,
		p.AddressCredits,
		p.AddressRegistrationFee,
	)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyNameInfoDuration, &p.NameInfoDuration, validateNameInfoDuration),
		params.NewParamSetPair(KeyNameInfoMaxDuration, &p.NameInfoMaxDuration, validateNameInfoDuration),
		params.NewParamSetPair(KeyNameInfoRegistrationFee, &p.NameInfoRegistrationFee, validateFee),
		params.NewParamSetPair(KeyNameInfoRenewalFee, &p.NameInfoRenewalFee, validateFee),
		params.NewParamSetPair(KeyAddressCredits, &p.AddressCredits, validateAddressCredits),
		params.NewParamSetPair(KeyAddressRegistrationFee, &p.AddressRegistrationFee, validateFee),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	defaultNameInfoCoinFee, _ := 			sdk.ConvertCoin(sdk.NewInt64Coin("anatha", 1), "pin")
	defaultNameInfoCoinsFee := 				sdk.NewCoins(defaultNameInfoCoinFee)

	defaultAddressRegistrationCoinFee, _ := sdk.ConvertCoin(sdk.NewInt64Coin("anatha", 1), "pin")
	defaultAddressRegistrationCoinsFee := 	sdk.NewCoins(defaultAddressRegistrationCoinFee)

	return NewParams(
		DefaultNameInfoDuration,
		DefaultNameInfoMaxDuration,
		defaultNameInfoCoinsFee,
		defaultNameInfoCoinsFee,
		DefaultAddressCredits,
		defaultAddressRegistrationCoinsFee,
	)
}

func (p Params) Validate() error {
	if err := validateNameInfoDuration(p.NameInfoDuration); err != nil {
		return err
	}

	if err := validateNameInfoDuration(p.NameInfoMaxDuration); err != nil {
		return err
	}

	if err := validateFee(p.NameInfoRegistrationFee); err != nil {
		return err
	}

	if err := validateFee(p.NameInfoRenewalFee); err != nil {
		return err
	}

	if err := validateAddressCredits(p.AddressCredits); err != nil {
		return err
	}

	if err := validateFee(p.AddressRegistrationFee); err != nil {
		return err
	}

	return nil
}

func validateNameInfoDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("name info duration must be positive: %d", v)
	}

	return nil
}

func validateFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if ! v.IsValid() {
		return fmt.Errorf("invalid name info fee: %s", v)
	}
	if ! v.AmountOf(config.DefaultDenom).IsPositive() {
		return fmt.Errorf("invalid fee denomination. expected: %s", config.DefaultDenom)
	}
	return nil
}

func validateAddressCredits(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("address credits must be positive: %d", v)
	}

	return nil
}

