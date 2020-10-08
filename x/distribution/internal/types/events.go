package types

const (
	EventTypeDevelopmentFundDistribution 	= "DevelopmentFundDistribution"
	EventTypeSecurityTokenFundDistribution 	= "SecurityTokenFundDistribution"
	EventTypeWithdrawNameReward				= "withdraw_name_reward"
	EventTypeWithdrawValidatorReward		= "withdraw_validator_rewards"
	EventTypeDepositSavings					= "deposit_savings"
	EventTypeWithdrawSavings				= "withdraw_savings"
	EventTypeWithdrawSavingsInterest		= "withdraw_savings_interest"

	AttributeKeyAmount					= "amount"
	AttributeKeyRecipient				= "recipient"
	AttributeKeySender					= "sender"
	AttributeKeyTitle					= "title"
	AttributeKeyDescription				= "description"
	AttributeKeyReward					= "reward"

	AttributeValueModule = ModuleName
)

