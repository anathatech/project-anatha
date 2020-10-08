package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)


var (
	ErrNameRegistered 				= sdkerrors.Register(ModuleName, 101, "Name already registered.")
	ErrNameNotRegistered 			= sdkerrors.Register(ModuleName, 102, "Name not registered.")
	ErrNotOwner 					= sdkerrors.Register(ModuleName, 103, "Initiator is not the owner.")
	ErrExpiredNameRenewal        	= sdkerrors.Register(ModuleName, 104, "Expired name renewal.")
	ErrNotForSale                	= sdkerrors.Register(ModuleName, 105, "Name is not for sale.")
	ErrAlreadyOwned              	= sdkerrors.Register(ModuleName, 106, "Name is already owned by the buyer.")
	ErrBlockchainIdNotValid      	= sdkerrors.Register(ModuleName, 107, "Blockchain ID not valid.")
	ErrAddressIndexNotValid      	= sdkerrors.Register(ModuleName, 108, "Address Index not valid.")
	ErrBlockchainAddressNotValid 	= sdkerrors.Register(ModuleName, 109, "Blockchain Address not valid.")
	ErrNameNotValid              	= sdkerrors.Register(ModuleName, 110, "Invalid name provided.")
	ErrBlockchainAddressNotFound 	= sdkerrors.Register(ModuleName, 111, "Blockchain Address not found.")
	ErrNoNamesRegistered         	= sdkerrors.Register(ModuleName, 112, "No names registered.")
	ErrMaximumDurationExceeded 		= sdkerrors.Register(ModuleName, 113, "Maximum Name Info duration exceeded.")
)
