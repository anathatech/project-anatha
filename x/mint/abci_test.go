package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func TestWithLeftover(t *testing.T) {
	t.Log("With leftover")

	initialSupply := sdk.NewDec(770000000000000000)
	totalSupplyDec := sdk.NewDec(770000000000000000)
	inflationPerSecond := sdk.NewDecWithPrec(315306958, 18)

	elapsed := int(365.25 * 24 * 60 * 60)

	leftover := sdk.ZeroDec()
	leftoverAccummulator := sdk.ZeroDec()

	for i := 0; i < elapsed; i++ {

		toMintDec := sdk.ZeroDec().Add(leftover) // initialize minting amount with previous iterations leftover

		currentSecondInflationAmountDec := totalSupplyDec.Mul(inflationPerSecond)

		toMintDec = toMintDec.Add(currentSecondInflationAmountDec)

		toMintInt := toMintDec.TruncateInt()
		leftover = toMintDec.Sub(toMintInt.ToDec())

		leftoverAccummulator = leftoverAccummulator.Add(leftover)

		totalSupplyDec = totalSupplyDec.Add(toMintInt.ToDec())

	}

	t.Log("Got")
	got := totalSupplyDec
	t.Log(got)

	t.Log("Expected")
	expected := initialSupply.Mul(sdk.OneDec().Add(sdk.NewDecWithPrec(1, 2)))
	t.Log(expected)

	diff := got.Sub(expected).Abs()
	t.Log("Diff")
	t.Log(diff)

	t.Log("Total Leftover accumulated")
	t.Log(leftoverAccummulator)
}

func TestWithoutLeftover(t *testing.T) {
	t.Log("Without leftover")

	initialSupply := sdk.NewDec(770000000000000000)
	totalSupplyDec := sdk.NewDec(770000000000000000)
	inflationPerSecond := sdk.NewDecWithPrec(315306958, 18)

	elapsed := int(365.25 * 24 * 60 * 60)

	for i := 0; i < elapsed; i++ {

		toMintDec := sdk.ZeroDec()

		currentSecondInflationAmountDec := totalSupplyDec.Mul(inflationPerSecond)

		toMintDec = toMintDec.Add(currentSecondInflationAmountDec)

		toMintInt := toMintDec.TruncateInt()

		totalSupplyDec = totalSupplyDec.Add(toMintInt.ToDec())

	}

	t.Log("Got")
	got := totalSupplyDec
	t.Log(got)

	t.Log("Expected")
	expected := initialSupply.Mul(sdk.OneDec().Add(sdk.NewDecWithPrec(1, 2)))
	t.Log(expected)

	diff := got.Sub(expected).Abs()
	t.Log("Diff")
	t.Log(diff)
}
