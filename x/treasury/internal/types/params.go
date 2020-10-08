package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = ModuleName

	DefaultRiskAssesmentDuration 	= time.Hour * 24 * 3 // devnet: time.Second * 60 * 2
	DefaultRiskAssesmentAmount		= 10000 // usd
)

var (
	KeyManagers               = []byte("Managers")
	KeyRiskAssessmentAmount   = []byte("RiskAssesmentAmount")
	KeyRiskAssessmentDuration = []byte("RiskAssesmentDuration")
	KeyBuyBackPercentage      = []byte("BuyBackPercentage")

	DefaultManagerAddress = "anatha1qaf2gssp652s6np00a5cxdwytdf3vutdumwc0q"

	DefaultBuyBackPercentage = sdk.NewDecWithPrec(24, 2)
)

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

type Params struct {
	Managers               []sdk.AccAddress `json:"managers" yaml:"managers"`
	RiskAssessmentAmount   sdk.Coins      `json:"risk_assesment_amount" yaml:"risk_assesment_amount"`
	RiskAssessmentDuration time.Duration  `json:"risk_assesment_duration" yaml:"risk_assesment_duration"`
	BuyBackPercentage      sdk.Dec        `json:"buyback_percentage" yaml:"buyback_percentage"`
}

func NewParams(managers []sdk.AccAddress, amount sdk.Coins, riskAssessmentDuration time.Duration, buybackPercentage sdk.Dec) Params {
	return Params{
		Managers:               managers,
		RiskAssessmentAmount:   amount,
		RiskAssessmentDuration: riskAssessmentDuration,
		BuyBackPercentage:      buybackPercentage,
	}
}

func (p Params) String() string {
	return fmt.Sprintf(`
	Managers: %s
	RiskAssesmentAmount: %s
	RiskAssesmentDuration: %s
	`, p.Managers, p.RiskAssessmentAmount, p.RiskAssessmentDuration)
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyManagers, &p.Managers, validateManager),
		params.NewParamSetPair(KeyRiskAssessmentAmount, &p.RiskAssessmentAmount, validateCoins),
		params.NewParamSetPair(KeyRiskAssessmentDuration, &p.RiskAssessmentDuration, validateDuration),
		params.NewParamSetPair(KeyBuyBackPercentage, &p.BuyBackPercentage, validateBuyBackPercentage),
	}
}

func DefaultParams() Params {
	manager, _ := sdk.AccAddressFromBech32(DefaultManagerAddress)
	amount, _ := sdk.ConvertCoin(sdk.NewInt64Coin("usd", DefaultRiskAssesmentAmount), "din")

	return NewParams(
		[]sdk.AccAddress{manager},
		sdk.NewCoins(amount),
		DefaultRiskAssesmentDuration,
		DefaultBuyBackPercentage,
	)
}

func (p Params) Validate() error {
	if err := validateManager(p.Managers); err != nil {
		return err
	}

	if err := validateCoins(p.RiskAssessmentAmount); err != nil {
		return err
	}

	if err := validateDuration(p.RiskAssessmentDuration); err != nil {
		return err
	}

	if err := validateBuyBackPercentage(p.BuyBackPercentage); err != nil {
		return err
	}

	return nil
}

func validateManager(i interface{}) error {
	addresses, ok := i.([]sdk.AccAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(addresses) <= 0 || len(addresses) > 4 {
		return fmt.Errorf("only 1-4 Managers can be specified: %d", len(addresses))
	}

	return nil
}

func validateDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("duration must be positive: %d", v)
	}

	return nil
}

func validateCoins(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsValid() {
		return fmt.Errorf("invalid amount: %s", v)
	}

	return nil
}

func validateBuyBackPercentage(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("buyback percentage must be less than 100%: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("buyback percentage must be positive: %s", v)
	}

	return nil
}