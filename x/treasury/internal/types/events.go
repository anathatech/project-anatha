package types

const (
	EventTypeDisburse    		= "disburse"
	EventTypeDisburseToEscrow   = "disburse_to_escrow"
	EventTypeDisburseFromEscrow = "disburse_from_escrow"
	EventTypeRevertFromEscrow   = "revert_from_escrow"
	EventTypeAddOperator 		= "add_operator"
	EventTypeRemoveOperator 	= "remove_operator"
	EventTypeCancelDisbursement	= "cancel_disbursement"
	EventTypeCreateSellOrder	= "create_sell_order"
	EventTypeCreateBuyOrder		= "create_buy_order"
	EventTypeTransfer			= "transfer_to_distribution_module"

	EventTypeAddBuyBackLiquidity = "AddBuyBackLiquidity"
	EventTypeRemoveBuyBackLiquidity = "RemoveBuyBackLiquidity"
	EventTypeBurnDistributionProfits = "BurnDistributionProfits"
	EventTypeTransferFromDistributionProfitsToBuyBackLiquidity = "TransferFromDistributionProfitsToBuyBackLiquidity"

	AttributeKeySender				= "sender"
	AttributeKeyRecipient  			= "recipient"
	AttributeKeyAmount				= "amount"
	AttributeKeyReference			= "reference"
	AttributeKeyScheduledFor 		= "scheduledFor"
	AttributeKeyOperator			= "operator"
	AttributeKeyPinAmount			= "pin_amount"
	AttributeKeyDinAmount			= "din_amount"
	AttributeKeyEscrowRemainder 	= "escrow_remainder"
	AttributeKeyTitle					= "title"
	AttributeKeyDescription				= "description"

	AttributeValueModule = ModuleName
)
