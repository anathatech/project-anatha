package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrUnknownProposal         = sdkerrors.Register(ModuleName, 1, "unknown proposal")
	ErrUnknownVote				= sdkerrors.Register(ModuleName, 9, "unknown vote")
	ErrInactiveProposal        = sdkerrors.Register(ModuleName, 2, "inactive proposal")
	//ErrAlreadyActiveProposal   = sdkerrors.Register(ModuleName, 3, "proposal already active")
	ErrInvalidProposalContent  = sdkerrors.Register(ModuleName, 4, "invalid proposal content")
	ErrInvalidProposalType     = sdkerrors.Register(ModuleName, 5, "invalid proposal type")
	ErrInvalidVote             = sdkerrors.Register(ModuleName, 6, "invalid vote option")
	ErrInvalidGenesis          = sdkerrors.Register(ModuleName, 7, "invalid genesis state")
	ErrNoProposalHandlerExists = sdkerrors.Register(ModuleName, 8, "no handler exists for proposal type")

	ErrNotGovernor 				= sdkerrors.Register(ModuleName, 10, "proposer/voter is not a governor")
	ErrAlreadyVoted				= sdkerrors.Register(ModuleName, 11, "already voted")
)