package ingame

import (
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
)

// WithdrawData input parames
type WithdrawData struct {
	AppID    int         `json:"appId"`
	Params   interface{} `json:"params"`
	SignHash string      `json:"sign"`
}

// PrewithdrawDataRes result
type PrewithdrawDataRes struct {
	GameAssetName string `json:"game_asset_name"`
	UID           uint64 `json:"uid"`
	EquipmentID   string `json:"equipment_id"`
	AppOrderID    string `json:"app_order_id"`
	Nonce         string `json:"nonce"`
	Status        int    `json:"status"`
}

// CommitWithdrawData input parames
type CommitWithdrawData struct {
	GameAssetName string `json:"game_asset_name"`
	UID           uint64 `json:"uid"`
	Nonce         string `json:"nonce"`
	AppOrderID    string `json:"app_order_id"`
	Operate       int    `json:"operation"`
}

// CommitwithdrawDataRes
type CommitwithdrawDataRes struct {
	GameAssetName string `json:"game_asset_name"`
	UID           uint64 `json:"uid"`
	Nonce         string `json:"nonce"`
	AppOrderID    string `json:"app_order_id"`
	Status        int    `json:"status"`
}

type NotifyFTDepositData struct {
	GameCoinName string `json:"game_coin_name"`
	Amount       string `json:"amount"`
	TxHash       string `json:"tx_hash"`
	Uid          int64  `json:"uid"`
}

type NotifyFTDeposits struct {
	Params   []NotifyFTDepositData `json:"params"`
	AppID    int                   `json:"appId"`
	SignHash string                `json:"sign"`
}

type GameFTDepositRes struct {
	GameCoinName string `json:"game_coin_name"`
	AppOrderID   string `json:"app_order_id"`
	TxHash       string `json:"tx_hash"`
	Status       int    `json:"status"`
}

// noti erc721 deposit data
type GameNFTDepositData struct {
	GameAssetName string                   `json:"game_asset_name"`
	TokenID       string                   `json:"token_id"`
	EquipmentID   string                   `json:"equipment_id"`
	TxHash        string                   `json:"tx_hash"`
	Uid           int64                    `json:"uid"`
	Attrs         []commdata.EquipmentAttr `json:"attrs"`
}

type NotifyNFTDeposits struct {
	Params   []GameNFTDepositData `json:"params"`
	AppID    int                  `json:"appId"`
	SignHash string               `json:"sign"`
}

type GameNFTDepositRes struct {
	GameAssetName string `json:"game_asset_name"`
	AppOrderID    string `json:"app_order_id"`
	TokenId       string `json:"token_id"`
	EquipmentId   string `json:"equipment_id"`
	TxHash        string `json:"tx_hash"`
	Status        int    `json:"status"`
}

type GameFTDepositFinalRet struct {
	Code    int64              `json:"code"`
	Message string             `json:"message"`
	Data    []GameFTDepositRes `json:"data"`
}

type GameNFTDepositFinalRet struct {
	Code    int64               `json:"code"`
	Message string              `json:"message"`
	Data    []GameNFTDepositRes `json:"data"`
}

type GameFTPreWithdrawData struct {
	AppOrderID string `json:"app_order_id"`
	Status     int    `json:"status"` //1 pass 2 reject
}

type GameCommonResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
type GameFTPreWithdraResp struct {
	GameCommonResponse
	Data GameFTPreWithdrawData `json:"data"`
}

type GameWithdrawComfirmData struct {
	AppOrderID   string `json:"app_order_id"`
	Nonce        string `json:"nonce"`
	GameCoinName string `json:"game_coin_name"`
	Status       int    `json:"status"` //1 pass 2 reject
}

type GameWithdrawComfirmResp struct {
	GameCommonResponse
	Data []GameWithdrawComfirmData `json:"data"`
}

func (rsp *GameFTPreWithdraResp) IsFailure() bool {
	return rsp.Data.Status != const_def.GAMESERVER_PASS || rsp.Data.AppOrderID == ""
}

func (rsp *GameWithdrawComfirmResp) IsFailure() bool {
	return len(rsp.Data) == 0 || rsp.Code != const_def.GAME_SERVER_SUCCESS_CODE || (rsp.Data[0].Status != const_def.GAMESERVER_PASS && rsp.Data[0].Status != const_def.GAMESERVER_REJECT)
}

func (rsp *GameWithdrawComfirmResp) IsReject() bool {
	return len(rsp.Data) > 0 && rsp.Data[0].Status == const_def.GAMESERVER_REJECT
}

func (rsp *GameWithdrawComfirmResp) IsSuccess() bool {
	return len(rsp.Data) > 0 && rsp.Data[0].Status == const_def.GAMESERVER_PASS
}
