package utils

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/config"
	"regexp"
	"strings"
)

var (
	// Denominations can be 3 ~ 16 characters long.
	reDnmString = `[a-z][a-z0-9]{2,15}`
	reAmt       = `[[:digit:]]+`
	reDecAmt    = `[[:digit:]]*\.[[:digit:]]+`
	reSpc       = `[[:space:]]*`
	reCoin      = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnmString))
	reDecCoin   = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reDecAmt, reSpc, reDnmString))
)

// usd -> din
// anatha, sense -> pin
func ParseAndConvertCoins(input string) (sdk.Coins, error) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil, nil
	}

	parts := strings.Split(input, ",")

	coins := make(sdk.Coins, len(parts))

	for i, coinStr := range parts {

		var currentCoin sdk.Coin
		var err error

		matches := reCoin.FindStringSubmatch(coinStr)

		if matches != nil {
			// found int coin
			amountStr, denomStr := matches[1], matches[2]

			amount, ok := sdk.NewIntFromString(amountStr)
			if ! ok {
				return sdk.Coins{}, fmt.Errorf("failed to parse coin amount: %s", amountStr)
			}

			if err := sdk.ValidateDenom(denomStr); err != nil {
				return sdk.Coins{}, fmt.Errorf("invalid denom cannot contain upper case characters or spaces: %s", err)
			}

			currentCoin, err = ToBaseDenom(sdk.NewCoin(denomStr, amount))
			if err != nil {
				return sdk.Coins{}, err
			}
		} else {
			matches = reDecCoin.FindStringSubmatch(coinStr)

			if matches != nil {
				// found dec coin
				amountStr, denomStr := matches[1], matches[2]

				amount, err := sdk.NewDecFromStr(amountStr)
				if err != nil {
					return sdk.Coins{}, fmt.Errorf("failed to parse decimal coin amount: %s", amountStr)
				}

				if err := sdk.ValidateDenom(denomStr); err != nil {
					return sdk.Coins{}, fmt.Errorf("invalid denom cannot contain upper case characters or spaces: %s", err)
				}

				currentCoin, err = DecToBaseDenom(sdk.NewDecCoinFromDec(denomStr, amount))
				if err != nil {
					return sdk.Coins{}, err
				}
			} else {
				// failed to parse coins
				return sdk.Coins{}, fmt.Errorf("failed to parse coin amount: %s", coinStr)
			}
		}
		
		coins[i] = currentCoin
	}

	// sort coins for determinism
	coins.Sort()

	// validate coins before returning
	if !coins.IsValid() {
		return nil, fmt.Errorf("parseCoins invalid: %#v", coins)
	}

	return coins, nil
}


func ToBaseDenom(coin sdk.Coin) (sdk.Coin, error) {
	if IsBaseDenom(coin.Denom) {
		return coin, nil
	}

	var err error
	baseCoin := coin

	if coin.Denom == "usd" {

		baseCoin, err = sdk.ConvertCoin(coin, "din")

	} else if coin.Denom == "anatha" || coin.Denom == "sense" {

		baseCoin, err = sdk.ConvertCoin(coin, "pin")

	} else {
		return sdk.Coin{}, fmt.Errorf("failed to parse denom. Not registered: %s", coin.Denom)
	}

	if err != nil {
		return sdk.Coin{}, err
	}

	return baseCoin, nil
}

func IsBaseDenom(denom string) bool {
	return denom == config.DefaultDenom || denom == config.DefaultStableDenom
}

func DecToBaseDenom(decCoin sdk.DecCoin) (sdk.Coin, error) {
	if IsBaseDenom(decCoin.Denom) {
		result, _ := decCoin.TruncateDecimal()
		return result, nil
	}

	var err error
	baseDecCoin := decCoin

	if decCoin.Denom == "usd" {

		unit, _ := sdk.GetDenomUnit("din")

		multiplier := sdk.OneDec().Quo(unit)

		baseDecCoin = sdk.NewDecCoinFromDec("din", decCoin.Amount.Mul(multiplier))

	} else if decCoin.Denom == "anatha" || decCoin.Denom == "sense" {

		targetUnit, _ := sdk.GetDenomUnit("pin")

		sourceUnit, _ := sdk.GetDenomUnit(decCoin.Denom)

		multiplier := sourceUnit.Quo(targetUnit)

		baseDecCoin = sdk.NewDecCoinFromDec("pin", decCoin.Amount.Mul(multiplier))

	} else {
		return sdk.Coin{}, fmt.Errorf("failed to parse denom. Not registered: %s", decCoin.Denom)
	}

	if err != nil {
		return sdk.Coin{}, err
	}

	baseCoin, _ := baseDecCoin.TruncateDecimal()

	return baseCoin, nil
}
