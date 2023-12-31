package const_def

const (
	U8_SET_WITHDRAW_STATUS_OPEN  = 1
	U8_SET_WITHDRAW_STATUS_CLOSE = 0

	U8_CONTRACT_OP_OPEN  = 1
	U8_CONTRACT_OP_CLOSE = 0

	U8_CONTRACT_ERC20  = "erc20"
	U8_CONTRACT_ERC721 = "erc721"

	U8_EQUIPMENT_WITHDRAW_OPEN  = 1
	U8_EQUIPMENT_WITHDRAW_CLOSE = 2

	U8_DEPOSIT_SUCCESS      = 1
	U8_DEPOSIT_WAIT_HANDLE  = 2
	U8_WITHDRAW_SUCCESS     = 1
	U8_WITHDRAW_WAIT_HANDLE = 2
	U8_WITHDRAW_PASS        = 3
	U8_WITHDRAW_REJECT      = 4
	U8_WITHDRAW_FAILED      = 5

	U8_WITHDRAW_SET_PASS   = 2
	U8_WITHDRAW_SET_REJECT = 3

	U8_REPAIR_UN_FINISHED = 1
	U8_REPAIR_FINISHED    = 2

	U8_CONTRCT_OP_DEPOSIT   = "deposit"
	U8_CONTRACT_OP_WITHDRAW = "withdraw"
)
