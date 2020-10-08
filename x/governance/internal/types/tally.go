package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TallyResult struct {
	Yes        sdk.Int `json:"yes" yaml:"yes"`
	No         sdk.Int `json:"no" yaml:"no"`
}

func NewTallyResult(yes, no sdk.Int) TallyResult {
	return TallyResult{
		Yes:        yes,
		No:         no,
	}
}

func NewTallyResultFromMap(results map[VoteOption]sdk.Dec) TallyResult {
	return NewTallyResult(
		results[OptionYes].TruncateInt(),
		results[OptionNo].TruncateInt(),
	)
}

func EmptyTallyResult() TallyResult {
	return NewTallyResult(sdk.ZeroInt(), sdk.ZeroInt())
}

func (tr TallyResult) Equals(comp TallyResult) bool {
	return tr.Yes.Equal(comp.Yes) &&
		tr.No.Equal(comp.No)
}

func (tr TallyResult) String() string {
	return fmt.Sprintf(`Tally Result:
  Yes:        %s
  No:         %s`, tr.Yes, tr.No)
}