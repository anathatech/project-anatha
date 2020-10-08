package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)


var (
	ErrRewardWithdrawalDisabled = sdkerrors.Register(ModuleName, 101, "Reward Withdrawal Disabled.")
	ErrNoValidatorRewards       = sdkerrors.Register(ModuleName, 102, "no validator commission to withdraw")

	ErrAlreadyHasSavings        = sdkerrors.Register(ModuleName, 103, "savings already active")
	ErrHasNoSavings             = sdkerrors.Register(ModuleName, 104, "user has no savings")
)