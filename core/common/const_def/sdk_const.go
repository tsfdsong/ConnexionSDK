package const_def

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// for parse logs
const (
	//ft
	Erc20DepositTopic      = "0xf0f5a6096b2a43fc10bc314148e9ac3851fd8455384682f3c3cd9d2f0bc71b17"
	Erc20WithdrawTopic     = "0xc603e3a56e4b1f9c71cf99ec62d824c6e4ff919af5a6dd3885c82b9d497b5cc7"
	OZ_Erc721TransferTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

	//nft loot
	NFTWithdrawLootTopic = "0x9561a8e4f5d28e8cc84227248daa23dbb2560e390340a9424a3447303a56e3f7"
	//nft standard
	NFTWithdrawTopic = "0x3478ef94fb18ff913d9966a0ce3c6168579d4c38f7c856e8ae9d40ab9720b717"

	NFTWithdrawUpdatTopic = "0x5560e2fe30c6e118818d212511a41d156e9da59318adc954739e60358780cec4"
	NFTDepositTopic       = "0xdf194dfab570e662896f8c8670a8a17541328d72a4d3fc02210becbf3d0be2f8"

	ParseSuiteBlockNum = 10

	GameServer_ERC20_Deposit         = "/ftToken/deposit"
	GameServer_ERC721_Deposit        = "/nftToken/deposit"
	GameServer_ERC721_Withdraw       = "/nftToken/withdrawTokenID"
	GameServer_ERC721_PreWithdraw    = "/nftToken/preWithdraw"
	GameServer_ERC721_CommitWithdraw = "/nftToken/withdraw"
)

// game server http status code return  --code
const (
	GAME_SERVER_SUCCESS_CODE    = 0
	GAME_SERVER_FAILED_CODE     = 1
	GAME_SERVER_SUCCESS_MESSAGE = "success"
)

// game server response for withdraw order status code --status
const (
	GAMESERVER_PASS   = 1
	GAMESERVER_REJECT = 2
)

// deposit & withdraw noti status.default=0 means not need noti at least now
const (

	//for t_ft_contracts/t_nft_contracts/t_user/t_game_equipment
	//default open 0 close 1

	SDK_TABLE_SWITCH_OPEN  = 0
	SDK_TABLE_SWITCH_CLOSE = 1

	SDK_GAMESERVER_NOTI_SUCCESS = 1
	SDK_GAMESERVER_NOTI_FAILED  = 2

	//game equipment status
	SDK_EQUIPMENT_STATUS_DEFAULT  = 0 //mystery box open
	SDK_EQUIPMENT_STATUS_DEPOSIT  = 1 //deposited
	SDK_EQUIPMENT_STATUS_WITHDRAW = 2 //withdrawed

	//contract selector
	CONTRACT_ERC20_DEPOSIT_SELECTOR  = "ERC20DepositSelector"
	CONTRACT_ERC20_WITHDRAW_SELECTOR = "ERC20WithdrawSelector"

	GameServer_FT_Query         = "/ftToken/getAssets"
	GameServer_NFT_Query        = "/nftToken/getAssets"
	GameServer_NFT_Detail_Query = "/nftToken/getAssetDetail"
	GameServer_FT_Pre_Withdraw  = "/ftToken/preWithdraw"
	GameServer_FT_Withdraw      = "/ftToken/withdraw"

	SDK_WITHDRAW_SIGN_URL = "/requestSignature"
	SDK_DEPOSIT_SIGN_URL  = "/requestSignatureAuto"
	SDK_QUERY_SIGN_URL    = "/querySignature"

	SDK_ATTR_SETTED  = 1
	SDK_ATTR_NOT_SET = 0
)

// status of second send withdraw request to game server
const (
	NOTI_GAMESERVER_DELETE  = 1
	NOTI_GAMESERVER_RECOVER = 2
)

const (
	SIGN_SUCCESS = "success"
	SIGN_FAILED  = "failed"
)

const (
	SDK_ERC20_NOTI_TYPE         = 1
	SDK_ERC721_NOTI_TYPE        = 2
	SDK_WITHDRAW_WITH_MINT_TYPE = 3
)

const (
	SDK_ERC20_DEPOSITSELECTOR  = "25DA3828"
	SDK_ERC20_WITHDRAWSELECTOR = "F1513F27"

	ABI_TYPE_UINT256 = "uint256"
	ABI_TYPE_ADDRESS = "address"
)

const (
	SIGN_PENDING = "pending"
)

const (
	ALERT_DISCORD = 1
)

const NEED_MANUAL_REPAIR_THIS = "NEED_MANUAL_REPAIR_THIS!!! "

const (
	DASHBOARD_WITHDRAW_PENDING   = 1
	DASHBOARD_WITHDRAW_SUCCESS   = 2
	DASHBOARD_WITHDRAW_FAILED    = 3
	DASHBOARD_WITHDRAW_REJECT    = 4
	DASHBOARD_WITHDRAW_CLAIMABLE = 5
)

const (
	GRAPH_ORDER_ASC  = 1
	GRAPH_ORDER_DESC = 2
)

const (
	MARKETPLACE_ORDER_DETAIL_LIST   = 1 //dashboard--order detail listing
	MARKETPLACE_ORDER_DETAIL_CLOSED = 2 //dashboard--order detail not listing

	MARKETPLACE_LIST     = 1 //graph and dashboard--activity  list
	MARKETPLACE_CANCEL   = 2 //graph anddashboard--activity  cancel
	MARKETPLACE_PURCHASE = 3 //graph anddashboard--activity purchase
	MARKETPLACE_SALE     = 4 //graph and dashboard--activity  sale
	MARKETPLACE_REDEEM   = 5
)

const FT_MAX_LIMIT = 200000

const (
	MARKETPLACE_ORDER_BUY      = 1
	MARKETPLACE_ORDER_CANCEL   = 2
	MARKETPLACE_ORDER_SELL     = 3
	MARKETPLACE_ORDER_REDEEM   = 4
	MARKETPLACE_ORDER_NOT_LIST = 5
)

// redis deposit&&withdraw prefix
const (
	FT_DEPOSIT     = "ft_deposit:"
	FT_WITHDRAW    = "ft_withdraw:"
	FT_HANDLE_NOW  = "handle_now"
	FT_SIGN_PREFIX = "ft_sign"
)

func GetFTSignKey(address string, signStr string) string {
	mdHash := fmt.Sprintf("%x", md5.Sum([]byte(signStr)))
	return FT_SIGN_PREFIX + strings.ToLower(address) + ":" + mdHash
}
