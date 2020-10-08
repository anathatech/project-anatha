package types

const (
	EventTypeSubmitProposal 		= "submit_proposal"
	EventTypeVote					= "vote"
	EventTypeExpedite				= "expedite"
	EventTypeAddGovernor			= "add_governor"
	EventTypeRemoveGovernor			= "remove_governor"

	AttributeKeySender				= "sender"
	AttributeKeyProposalId  		= "proposal_id"
	AttributeKeyStatus				= "status"
	AttributeKeyOption				= "option"
	AttributeKeyVotingStartTime		= "voting_start_time"
	AttributeKeyVotingEndTime		= "voting_end_time"
	AttributeKeyContent				= "content"
	AttributeKeyGovernor			= "governor"
	AttributeKeyTitle				= "title"
	AttributeKeyDescription			= "description"

	AttributeValueModule = ModuleName
)

