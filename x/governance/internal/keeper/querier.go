package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/anathatech/project-anatha/x/governance/internal/types"
)

const (
	QueryParams    = "params"
	QueryProposals = "proposals"
	QueryProposal  = "proposal"
	QueryVotes     = "votes"
	QueryVote      = "vote"

	QueryGovernors = "governors"

	ParamVoting   = "voting"
	ParamTallying = "tallying"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryParams:
			return queryParams(ctx, path[1:], req, keeper)

		case QueryProposals:
			return queryProposals(ctx, path[1:], req, keeper)

		case QueryProposal:
			return queryProposal(ctx, path[1:], req, keeper)

		case QueryVotes:
			return queryVotes(ctx, path[1:], req, keeper)

		case QueryVote:
			return queryVote(ctx, path[1:], req, keeper)

		case QueryGovernors:
			return queryGovernors(ctx, path[1:], req, keeper)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	switch path[0] {
		case ParamVoting:
			bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetVotingParams(ctx))
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}
			return bz, nil

		case ParamTallying:
			bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetTallyParams(ctx))
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
			}
			return bz, nil

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "%s is not a valid query request path", req.Path)
	}
}

func queryProposals(ctx sdk.Context, _ []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	proposals := keeper.GetProposals(ctx)
	if proposals == nil {
		proposals = types.Proposals{}
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, proposals)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryProposal(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	proposalID, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "proposal-id %s not a valid int, please input a valid proposal-id", path[0])
	}

	proposal, err := keeper.GetProposal(ctx, proposalID)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownProposal, "%d", proposalID)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, proposal)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryVotes(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	proposalID, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "proposal-id %s not a valid int, please input a valid proposal-id", path[0])
	}

	votes := keeper.GetVotes(ctx, proposalID)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, votes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryVote(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	proposalID, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "proposal-id %s not a valid int, please input a valid proposal-id", path[0])
	}

	voter, err := sdk.AccAddressFromBech32(path[1])
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to parse the voter address.")
	}

	vote, err := keeper.GetVote(ctx, proposalID, voter)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to find specified vote")
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, vote)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryGovernors(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	governors := keeper.GetGovernors(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, governors)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
