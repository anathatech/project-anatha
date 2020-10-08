package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gov "github.com/cosmos/cosmos-sdk/x/gov"
)

const (
	TypeMsgVote           = "vote"
	TypeMsgSubmitProposal = "submit_proposal"
)

var _, _ sdk.Msg = MsgSubmitProposal{}, MsgVote{}


type MsgSubmitProposal struct {
	Content        gov.Content        `json:"content" yaml:"content"`
	Proposer       sdk.AccAddress `json:"proposer" yaml:"proposer"`
}

func NewMsgSubmitProposal(content gov.Content, proposer sdk.AccAddress) MsgSubmitProposal {
	return MsgSubmitProposal{content, proposer}
}

func (msg MsgSubmitProposal) Route() string { return RouterKey }

func (msg MsgSubmitProposal) Type() string { return TypeMsgSubmitProposal }

func (msg MsgSubmitProposal) ValidateBasic() error {
	if msg.Content == nil {
		return sdkerrors.Wrap(ErrInvalidProposalContent, "missing content")
	}
	if msg.Proposer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Proposer.String())
	}
	if !IsValidProposalType(msg.Content.ProposalType()) {
		return sdkerrors.Wrap(ErrInvalidProposalType, msg.Content.ProposalType())
	}

	return msg.Content.ValidateBasic()
}

func (msg MsgSubmitProposal) String() string {
	return fmt.Sprintf(`Submit Proposal Message:
  Content:         %s
`, msg.Content.String())
}

func (msg MsgSubmitProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}

// MsgVote
type MsgVote struct {
	ProposalID uint64         `json:"proposal_id" yaml:"proposal_id"`
	Voter      sdk.AccAddress `json:"voter" yaml:"voter"`
	Option     VoteOption     `json:"option" yaml:"option"`
}

func NewMsgVote(voter sdk.AccAddress, proposalID uint64, option VoteOption) MsgVote {
	return MsgVote{proposalID, voter, option}
}

func (msg MsgVote) Route() string { return RouterKey }

func (msg MsgVote) Type() string { return TypeMsgVote }

// ValidateBasic implements Msg
func (msg MsgVote) ValidateBasic() error {
	if msg.Voter.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Voter.String())
	}
	if !ValidVoteOption(msg.Option) {
		return sdkerrors.Wrap(ErrInvalidVote, msg.Option.String())
	}

	return nil
}

func (msg MsgVote) String() string {
	return fmt.Sprintf(`Vote Message:
  Proposal ID: %d
  Option:      %s
`, msg.ProposalID, msg.Option)
}

func (msg MsgVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

// MsgExpedite
type MsgExpedite struct {
	ProposalID 	uint64         `json:"proposal_id" yaml:"proposal_id"`
	Sender 		sdk.AccAddress `json:"sender" yaml:"sender"`
}

func NewMsgExpedite(sender sdk.AccAddress, proposalID uint64) MsgExpedite {
	return MsgExpedite{proposalID, sender}
}

func (msg MsgExpedite) Route() string { return RouterKey }

func (msg MsgExpedite) Type() string { return "expedite" }

func (msg MsgExpedite) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	return nil
}

func (msg MsgExpedite) String() string {
	return fmt.Sprintf(`Expedite proposal:
  Proposal ID: %d
`, msg.ProposalID)
}

func (msg MsgExpedite) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgExpedite) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}