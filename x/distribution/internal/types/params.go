package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"time"
)

const (
	DefaultNameDepositDelay              = time.Second * 1 // time.Second * 30
	DefaultRewardWithdrawalBlockedPeriod = time.Hour * 24 * 365 // time.Second * 30
)

var (
	DefaultRewardWithdrawalEnabledTime = time.Time{}
	DefaultSavingsSplitAdjustment = sdk.NewDecWithPrec(9, 1)

	DefaultDevelopmentFundShare = sdk.NewDecWithPrec(25, 2)
	DefaultSecurityTokenFundShare = sdk.NewDecWithPrec(25, 2)

	KeyNameDepositDelay              = []byte("NameDepositDelay")
	KeyRewardWithdrawalBlockedPeriod = []byte("RewardWithdrawalBlockedPeriod")
	KeyRewardWithdrawalEnabledTime   = []byte("RewardWithdrawalEnabledTime")

	KeySavingsSplitAdjustment        = []byte("SavingsSplitAdjustment")

	KeyDevelopmentFundShare          = []byte("DevelopmentFundShare")
	KeySecurityTokenFundShare        = []byte("SecurityTokenFundShare")
)

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

type Params struct {
	NameDepositDelay              time.Duration `json:"name_deposit_delay" yaml:"name_deposit_delay"`
	RewardWithdrawalBlockedPeriod time.Duration `json:"reward_withdrawal_blocked_period" yaml:"reward_withdrawal_blocked_period"`
	RewardWithdrawalEnabledTime   time.Time `json:"reward_withdrawal_enabled_time" yaml:"reward_withdrawal_enabled_time"`
	SavingsSplitAdjustment        sdk.Dec `json:"savings_split_adjustment" yaml:"savings_split_adjustment"`

	DevelopmentFundShare          sdk.Dec `json:"development_fund_share" yaml:"development_fund_share"`
	SecurityTokenFundShare        sdk.Dec `json:"security_token_fund_share" yaml:"security_token_fund_share"`
}

func NewParams(nameDepositDelay time.Duration) Params {
	return Params{
		NameDepositDelay: nameDepositDelay,
	}
}

func DefaultParams() Params {
	return Params{
		DefaultNameDepositDelay,
		DefaultRewardWithdrawalBlockedPeriod,
		DefaultRewardWithdrawalEnabledTime,
		DefaultSavingsSplitAdjustment,
		DefaultDevelopmentFundShare,
		DefaultSecurityTokenFundShare,
	}
}

func (p Params) String() string {
	return fmt.Sprintf(`
	Name Deposit Delay: %s
	Reward Withdrawal Blocked Period: %s
	Reward Withdrawal Enabled Time: %s
	Savings Split Adjustment: %s
	Default Development Fund Share: %s
	Default Security Token Fund Share: %s
	`, p.NameDepositDelay, p.RewardWithdrawalBlockedPeriod, p.RewardWithdrawalEnabledTime, p.SavingsSplitAdjustment, p.DevelopmentFundShare, p.SecurityTokenFundShare)
}

func (p Params) Validate() error {
	if err := validateDuration(p.NameDepositDelay); err != nil {
		return err
	}

	if err := validateDuration(p.RewardWithdrawalBlockedPeriod); err != nil {
		return err
	}

	if err := validateTime(p.RewardWithdrawalEnabledTime); err != nil {
		return err
	}

	if err := validatePercentage(p.SavingsSplitAdjustment); err != nil {
		return err
	}

	if err := validateDevelopmentFundShare(p.DevelopmentFundShare); err != nil {
		return err
	}

	if err := validateSecurityTokenFundShare(p.SecurityTokenFundShare); err != nil {
		return err
	}

	return nil
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyNameDepositDelay, &p.NameDepositDelay, validateDuration),
		params.NewParamSetPair(KeyRewardWithdrawalBlockedPeriod, &p.RewardWithdrawalBlockedPeriod, validateDuration),
		params.NewParamSetPair(KeyRewardWithdrawalEnabledTime, &p.RewardWithdrawalEnabledTime, validateTime),
		params.NewParamSetPair(KeySavingsSplitAdjustment, &p.SavingsSplitAdjustment, validatePercentage),
		params.NewParamSetPair(KeyDevelopmentFundShare, &p.DevelopmentFundShare, validateDevelopmentFundShare),
		params.NewParamSetPair(KeySecurityTokenFundShare, &p.SecurityTokenFundShare, validateSecurityTokenFundShare),
	}
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

func validateTime(i interface{}) error {
	_, ok := i.(time.Time)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validatePercentage(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("percentage must be less than 100%: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("percentage must be positive: %s", v)
	}

	return nil
}

func validateDevelopmentFundShare(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	max := sdk.NewDecWithPrec(25, 2)
	if v.GT(max) {
		return fmt.Errorf("percentage must be less than 25%: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("percentage must be positive: %s", v)
	}

	return nil
}

func validateSecurityTokenFundShare(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	max := sdk.NewDecWithPrec(25, 2)
	if v.GT(max) {
		return fmt.Errorf("percentage must be less than 25%: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("percentage must be positive: %s", v)
	}

	return nil
}