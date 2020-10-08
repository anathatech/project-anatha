package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	params "github.com/cosmos/cosmos-sdk/x/params/subspace"
)

const (
	DefaultPeriod time.Duration = time.Hour * 24 * 3 // devnet: time.Second * 120
)

var (
	DefaultQuorum           = sdk.NewDecWithPrec(334, 3)
	DefaultThreshold        = sdk.NewDecWithPrec(5, 1)
)

var (
	ParamStoreKeyVotingParams  = []byte("votingparams")
	ParamStoreKeyTallyParams   = []byte("tallyparams")
)

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(ParamStoreKeyVotingParams, VotingParams{}, validateVotingParams),
		params.NewParamSetPair(ParamStoreKeyTallyParams, TallyParams{}, validateTallyParams),
	)
}

type TallyParams struct {
	Quorum    sdk.Dec `json:"quorum" yaml:"quorum"`
	Threshold sdk.Dec `json:"threshold" yaml:"threshold"`
}

func NewTallyParams(quorum sdk.Dec, threshold sdk.Dec) TallyParams {
	return TallyParams{
		Quorum:    quorum,
		Threshold: threshold,
	}
}

func DefaultTallyParams() TallyParams {
	return NewTallyParams(DefaultQuorum, DefaultThreshold)
}

func (tp TallyParams) String() string {
	return fmt.Sprintf(`Tally Params:
  Quorum:             %s
  Threshold:          %s`,
		tp.Quorum, tp.Threshold)
}

func validateTallyParams(i interface{}) error {
	v, ok := i.(TallyParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Quorum.IsNegative() {
		return fmt.Errorf("quorom cannot be negative: %s", v.Quorum)
	}
	if v.Quorum.GT(sdk.OneDec()) {
		return fmt.Errorf("quorom too large: %s", v)
	}
	if !v.Threshold.IsPositive() {
		return fmt.Errorf("vote threshold must be positive: %s", v.Threshold)
	}
	if v.Threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("vote threshold too large: %s", v)
	}

	return nil
}

type VotingParams struct {
	VotingPeriod time.Duration `json:"voting_period" yaml:"voting_period"`
}

func NewVotingParams(votingPeriod time.Duration) VotingParams {
	return VotingParams{
		VotingPeriod: votingPeriod,
	}
}

func DefaultVotingParams() VotingParams {
	return NewVotingParams(DefaultPeriod)
}

func (vp VotingParams) String() string {
	return fmt.Sprintf(`Voting Params:
  Voting Period:      %s`, vp.VotingPeriod)
}

func validateVotingParams(i interface{}) error {
	v, ok := i.(VotingParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.VotingPeriod <= 0 {
		return fmt.Errorf("voting period must be positive: %s", v.VotingPeriod)
	}

	return nil
}

type Params struct {
	VotingParams  VotingParams  `json:"voting_params" yaml:"voting_params"`
	TallyParams   TallyParams   `json:"tally_params" yaml:"tally_params"`
}

func (gp Params) String() string {
	return gp.VotingParams.String() + "\n" +
		gp.TallyParams.String()
}

func NewParams(vp VotingParams, tp TallyParams) Params {
	return Params{
		VotingParams:  vp,
		TallyParams:   tp,
	}
}

func DefaultParams() Params {
	return NewParams(DefaultVotingParams(), DefaultTallyParams())
}

func (p Params) Validate() error {
	if err := validateVotingParams(p.VotingParams); err != nil {
		return err
	}

	if err := validateTallyParams(p.TallyParams); err != nil {
		return err
	}

	return nil
}