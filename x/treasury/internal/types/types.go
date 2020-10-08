package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type Disbursement struct {
	Operator 		sdk.AccAddress 	`json:"operator" yaml:"operator"`
	Recipient 		sdk.AccAddress 	`json:"recipient" yaml:"recipient"`
	Amount 			sdk.Coins 		`json:"amount" yaml:"amount"`
	ScheduledFor 	time.Time 		`json:"scheduled_for" yaml:"scheduled_for"`
	Reference 		string 			`json:"reference" yaml:"reference"`
}

func NewDisbursement(operator sdk.AccAddress, recipient sdk.AccAddress, amount sdk.Coins, scheduledFor time.Time, reference string) Disbursement {
	return Disbursement{
		Operator: operator,
		Recipient: recipient,
		Amount: amount,
		ScheduledFor: scheduledFor,
		Reference: reference,
	}
}

func (d Disbursement) String() string {
	return fmt.Sprintf(`
	Operator: %s
	Recipient: %s
	Amount: %s
	ScheduledFor: %s
	Reference: %s
	`, d.Operator, d.Recipient, d.Amount, d.ScheduledFor, d.Reference)
}

type ReferenceAmountInfo struct {
	Reference string	`json:"reference" yaml:"reference"`
	Amount sdk.Int		`json:"amount" yaml:"amount"`
}

func NewReferenceAmountInfo(reference string, amount sdk.Int) ReferenceAmountInfo {
	return ReferenceAmountInfo{Reference:reference,Amount:amount}
}

func (a ReferenceAmountInfo) String() string {
	return fmt.Sprintf(`Reference: %s
Amount: %s`, a.Reference, a.Amount)
}