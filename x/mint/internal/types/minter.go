package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type Minter struct {
	PreviousMintTime 	time.Time `json:"previous_mint_time" yaml:"previous_mint_time"`
	Leftover 			sdk.Dec `json:"leftover" yaml:"leftover"`
}

func NewMinter(previousMintTime time.Time, leftover sdk.Dec) Minter {
	return Minter{
		PreviousMintTime: previousMintTime,
		Leftover: leftover,
	}
}

func InitialMinter(previousMintTime time.Time, leftover sdk.Dec) Minter {
	return NewMinter(
		previousMintTime,
		leftover,
	)
}

func DefaultInitialMinter() Minter {
	return InitialMinter(
		time.Time{},
		sdk.ZeroDec(),
	)
}

func (m Minter) Validate() error {
	if err := validatePreviousMintTime(m.PreviousMintTime); err != nil {
		return err
	}
	if err := validateLeftover(m.Leftover); err != nil {
		return err
	}
	return nil
}

func validatePreviousMintTime(i interface{}) error {
	v, ok := i.(time.Time)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.After(time.Now()) {
		return fmt.Errorf("Previous mint time must be in the past")
	}

	return nil
}

func validateLeftover(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("leftover must be positive: %s", v)
	}

	return nil
}

