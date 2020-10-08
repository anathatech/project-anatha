package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNotManager = sdkerrors.Register(ModuleName, 101, "No Manager permissions to execute the action.")
	ErrNotOperator = sdkerrors.Register(ModuleName, 102, "No Operator permissions to execute the action.")
	ErrInvalidTime = sdkerrors.Register(ModuleName, 103, "Invalid time")
	ErrEscrowDistributionAmountExceeded = sdkerrors.Register(ModuleName, 109, "escrow distribution amount exceeded")
	ErrInvalidReference = sdkerrors.Register(ModuleName, 110, "Invalid reference")
	ErrEscrowDisbursed = sdkerrors.Register(ModuleName, 111, "Escrow already disbursed")
	ErrDuplicateReference = sdkerrors.Register(ModuleName, 112, "Reference already used")
	ErrEscrowRevertAmountTooBig = sdkerrors.Register(ModuleName, 113, "Escrow revert amount too big")
	ErrDisbursementNotScheduled = sdkerrors.Register(ModuleName, 114, "Disbursement not scheduled")
)
