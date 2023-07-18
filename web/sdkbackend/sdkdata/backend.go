package sdkdata

import "github/Connector-Gamefi/ConnectorGoSDK/web/common"

type PRepairOrder struct {
	ID           int64  `json:"id" binding:"gt=0"`
	ContractType string `json:"type" binding:"gt=0"`
	//Status       int    `json:"status" binding:"gt=0"`
}

type PParseLogSwitch struct {
	Status int `json:"status" binding:"gte=0"`
	GameID int `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
}

type PParseSpecifiedBlockLog struct {
	Height int `json:"height" binding:"gt=0"`
	GameID int `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
}

type PGetUserAssets struct {
	GameID uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	Email  string `form:"email" json:"email" uri:"email" binding:"email"`
	UID    uint64 `form:"uid" json:"uid" uri:"uid" binding:"gt=0"`
}

type PGetUserNFTAssets struct {
	GameID    uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	AssetName string `form:"game_asset_name" json:"game_asset_name" uri:"game_asset_name" binding:"-"`
	Email     string `form:"email" json:"email" uri:"email" binding:"email"`
	Page      int    `form:"page" json:"page" uri:"page" binding:"gt=0"`
	PageSize  uint64 `form:"pageSize" json:"pageSize" uri:"pageSize" binding:"gt=0,lte=50"`
	UID       uint64 `form:"uid" json:"uid" uri:"uid" binding:"gt=0"`
}

type RGetUserNFTAssets struct {
	Erc721     []common.RBackendERC721Asset `json:"erc721"`
	UpdateTime string                       `json:"updateTime"`
}

type RGetUserAssets struct {
	Erc20      []common.RBackendERC20Asset `json:"erc20"`
	UpdateTime string                      `json:"updateTime"`
}

type PWithdrawExamineSet struct {
	ContractType string `json:"type" binding:"gt=0"`
	ID           int    `json:"id"  binding:"gt=0"`
	Status       int    `form:"status" json:"status" uri:"status" binding:"gt=0"`
	Reviewer     string `json:"reviewer" binding:"-"`
}
