package const_def

const (
	//state code in deposit for ft or nft
	CodeDepositNone    = 0 //init state in deposit
	CodeDepositSign    = 1 // deposit signed and wait for send to blockchain
	CodeDepositOnChain = 2 // deposit tx on chain
	CodeDepositSuccess = 3 // parse log comfirmed and wait for noti success directly

	CodeDepositTxFailed    = 4 //deposit tx execute failed in deposit // erc20 no use
	CodeDepositGameFailed  = 5 //send deposit request to game server and execute failed //  -- parselog deposit success and noti faild
	CodeDepositNotiSuccess = 6 //-- parselog deposit success and noti success when customers check

	//state code in withdraw for ft or nft
	CodeWithdrawNone      = 0 //init state in withdraw
	CodeWithdrawRisking   = 1 //wait for risk control in withdraw
	CodeWithdrawSign      = 2 //wait signature result in withdraw
	CodeWithdrawWaitClaim = 3 //sign success and wait for claim in withdraw --- for ft means noti game server delete asset success and wait for claim
	CodeWithdrawOnChain   = 4 //find withdraw tx on chain in withdraw
	CodeWithdrawSuccess   = 5 //withdraw success
	CodeWithdrawTimeout   = 6 //ft prewithdraw timeout

	CodeWithdrawRiskFailed        = 7  //risk control failed in withdraw //----wait risk control reject or failed
	CodeWithdrawPreWithdrawFailed = 8  //prepare withdraw failed //---- game server first reject or failed
	CodeWithdrawSignFailed        = 9  // sign failed in withdraw //---- sign reject or  failed
	CodeWithdrawCommitFailed      = 10 //game server second  failed ..need repair
	CodeWithdrawTxFailed          = 11 //withdraw tx execute failed

	CodeInnerError = 12 //etc new request failed,db error/request error..  please notice ！！！ only ft use now. this code means repair!!!

	CodeNotiRecorverSuccess = 13

	//equipmen withdraw switch
	CodeEquipmentWithdrawOpen   = 0
	CodeEquipmentWithdrawClosed = 1
)
