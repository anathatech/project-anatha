package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	app.RegisterCoinDenoms()
}

func TestParseAndConvertCoins(t *testing.T) {
	cases := []struct {
		input    string
		valid    bool  // if false, we expect an error on parse
		expected sdk.Coins // if valid is true, make sure this is returned
	}{
		{"", true, nil},

		{"100pin", true, sdk.Coins{{"pin", sdk.NewInt(100)}}},
		{"1anatha", true, sdk.Coins{{"pin", sdk.NewInt(100000000)}}},
		{"0.25anatha", true, sdk.Coins{{"pin", sdk.NewInt(25000000)}}},
		{"0.00000001anatha", true, sdk.Coins{{"pin", sdk.NewInt(1)}}},

		{"1sense", true, sdk.Coins{{"pin", sdk.NewInt(100000)}}},
		{"0.25sense", true, sdk.Coins{{"pin", sdk.NewInt(25000)}}},

		{"100din", true, sdk.Coins{{"din", sdk.NewInt(100)}}},
		{"1fiat", true, sdk.Coins{{"din", sdk.NewInt(10000000000)}}},
		{"0.25fiat", true, sdk.Coins{{"din", sdk.NewInt(2500000000)}}},
		{"0.0000000001fiat", true, sdk.Coins{{"din", sdk.NewInt(1)}}},

		{"1anatha,1fiat", true, sdk.Coins{{"din", sdk.NewInt(10000000000)}, {"pin", sdk.NewInt(100000000)}}},
		{"0.5anatha,1fiat", true, sdk.Coins{{"din", sdk.NewInt(10000000000)}, {"pin", sdk.NewInt(50000000)}}},
	}

	for tcIndex, tc := range cases {
		res, err := ParseAndConvertCoins(tc.input)
		if !tc.valid {
			require.NotNil(t, err, "%s: %#v. tc #%d", tc.input, res, tcIndex)
		} else if assert.Nil(t, err, "%s: %+v", tc.input, err) {
			require.Equal(t, tc.expected, res, "coin parsing was incorrect, tc #%d", tcIndex)
		}
	}
}