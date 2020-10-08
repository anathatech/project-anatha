package governance

import (
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/anathatech/project-anatha/x/governance/internal/keeper"
	"github.com/anathatech/project-anatha/x/governance/internal/types"
)

const (
	MaxDescriptionLength  = types.MaxDescriptionLength
	MaxTitleLength        = types.MaxTitleLength
	DefaultPeriod         = types.DefaultPeriod
	ModuleName            = types.ModuleName
	StoreKey              = types.StoreKey
	RouterKey             = types.RouterKey
	QuerierRoute          = types.QuerierRoute
	DefaultParamspace     = types.DefaultParamspace
	TypeMsgVote           = types.TypeMsgVote
	TypeMsgSubmitProposal = types.TypeMsgSubmitProposal
	StatusNil             = types.StatusNil
	StatusVotingPeriod    = types.StatusVotingPeriod
	StatusPassed          = types.StatusPassed
	StatusRejected        = types.StatusRejected
	StatusFailed          = types.StatusFailed
	ProposalTypeText      = types.ProposalTypeText
	ProposalTypeAddGovernor = types.ProposalTypeAddGovernor
	ProposalTypeRemoveGovernor = types.ProposalTypeRemoveGovernor

	OptionEmpty           = types.OptionEmpty
	OptionYes             = types.OptionYes
	OptionNo              = types.OptionNo
)

var (
	// functions aliases
	//RegisterInvariants            = keeper.RegisterInvariants
	//AllInvariants                 = keeper.AllInvariants
	//ModuleAccountInvariant        = keeper.ModuleAccountInvariant
	NewKeeper                     = keeper.NewKeeper
	NewQuerier                    = keeper.NewQuerier
	RegisterCodec                 = types.RegisterCodec
	ValidateAbstract              = types.ValidateAbstract
	ErrUnknownProposal            = types.ErrUnknownProposal
	ErrUnknownVote            		= types.ErrUnknownVote
	ErrInactiveProposal           = types.ErrInactiveProposal
	ErrInvalidProposalContent     = types.ErrInvalidProposalContent
	ErrInvalidProposalType        = types.ErrInvalidProposalType
	ErrInvalidVote                = types.ErrInvalidVote
	ErrInvalidGenesis             = types.ErrInvalidGenesis
	ErrNoProposalHandlerExists    = types.ErrNoProposalHandlerExists
	NewGenesisState               = types.NewGenesisState
	DefaultGenesisState           = types.DefaultGenesisState
	ValidateGenesis               = types.ValidateGenesis
	GetProposalIDBytes            = types.GetProposalIDBytes
	GetProposalIDFromBytes        = types.GetProposalIDFromBytes
	ProposalKey                   = types.ProposalKey
	ActiveProposalByTimeKey       = types.ActiveProposalByTimeKey
	ActiveProposalQueueKey        = types.ActiveProposalQueueKey
	VotesKey                      = types.VotesKey
	VoteKey                       = types.VoteKey
	SplitProposalKey              = types.SplitProposalKey
	SplitActiveProposalQueueKey   = types.SplitActiveProposalQueueKey
	SplitKeyVote                  = types.SplitKeyVote
	NewMsgSubmitProposal          = types.NewMsgSubmitProposal
	NewMsgVote                    = types.NewMsgVote
	NewMsgExpedite				  = types.NewMsgExpedite
	ParamKeyTable                 = types.ParamKeyTable
	NewTallyParams                = types.NewTallyParams
	NewVotingParams               = types.NewVotingParams
	NewParams                     = types.NewParams
	NewProposal                   = types.NewProposal
	NewRouter                     = gov.NewRouter
	ProposalStatusFromString      = types.ProposalStatusFromString
	ValidProposalStatus           = types.ValidProposalStatus
	NewTextProposal               = types.NewTextProposal
	NewAddGovernorProposal		= types.NewAddGovernorProposal
	NewRemoveGovernorProposal	= types.NewRemoveGovernorProposal
	RegisterProposalType          = types.RegisterProposalType
	RegisterProposalTypeCodec 	= types.RegisterProposalTypeCodec
	ContentFromProposalType       = types.ContentFromProposalType
	IsValidProposalType           = types.IsValidProposalType
	NewTallyResult                = types.NewTallyResult
	NewTallyResultFromMap         = types.NewTallyResultFromMap
	EmptyTallyResult              = types.EmptyTallyResult
	NewVote                       = types.NewVote
	VoteOptionFromString          = types.VoteOptionFromString
	ValidVoteOption               = types.ValidVoteOption

	// variable aliases
	ModuleCdc                   = types.ModuleCdc
	ProposalsKeyPrefix          = types.ProposalsKeyPrefix
	ActiveProposalQueuePrefix   = types.ActiveProposalQueuePrefix
	ExpeditedProposalQueuePrefix = types.ExpeditedProposalQueuePrefix
	ProposalIDKey               = types.ProposalIDKey
	VotesKeyPrefix              = types.VotesKeyPrefix
	ParamStoreKeyVotingParams   = types.ParamStoreKeyVotingParams
	ParamStoreKeyTallyParams    = types.ParamStoreKeyTallyParams
)

type (
	Keeper               = keeper.Keeper
	Content              = gov.Content
	Handler              = gov.Handler
	GenesisState         = types.GenesisState
	MsgSubmitProposal    = types.MsgSubmitProposal
	MsgVote              = types.MsgVote
	MsgExpedite          = types.MsgExpedite
	TallyParams          = types.TallyParams
	VotingParams         = types.VotingParams
	Params               = types.Params
	Proposal             = types.Proposal
	Proposals            = types.Proposals
	ProposalStatus       = types.ProposalStatus
	TextProposal         = types.TextProposal
	TallyResult          = types.TallyResult
	Vote                 = types.Vote
	Votes                = types.Votes
	VoteOption           = types.VoteOption
)
