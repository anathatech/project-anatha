package types

import (
	"errors"
	"fmt"
	"github.com/anathatech/project-anatha/config"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	KeyPerSecondInflationRate = []byte("PerSecondInflationRate")
	KeyMintDenom              = []byte("MintDenom")
)

type Params struct {
	PerSecondInflationRate sdk.Dec `json:"per_second_inflation_rate" yaml:"per_second_inflation_rate"` 	// per second inflation rate
	MintDenom              string  `json:"mint_denom" yaml:"mint_denom"`            		// type of coin to mint
}

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	perSecondInflationRate sdk.Dec, mintDenom string,
) Params {

	return Params{
		PerSecondInflationRate: perSecondInflationRate,
		MintDenom:              mintDenom,
	}
}

func DefaultParams() Params {
	return Params{
		PerSecondInflationRate: sdk.NewDecWithPrec(315306958, 18), // to target 1% annual
		MintDenom:              config.DefaultDenom,
	}
}

func (p Params) Validate() error {
	if err := validatePerSecondInflation(p.PerSecondInflationRate); err != nil {
		return err
	}
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}

	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Minting Params:
  Mint Denom:             %s
  Per Second Inflation Rate:  %s
`,
		p.PerSecondInflationRate, p.MintDenom,
	)
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyPerSecondInflationRate, &p.PerSecondInflationRate, validatePerSecondInflation),
		params.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validatePerSecondInflation(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("max inflation cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("max inflation too large: %s", v)
	}

	return nil
}
