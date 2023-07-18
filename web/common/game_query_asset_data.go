package common

import "github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"

// request from dashboard
type PQueryGameAssets struct {
	GameID uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	Email  string `form:"email" json:"email" uri:"email" binding:"email"`
}

type PQueryGameFTAssets struct {
	GameID uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	Email  string `form:"email" json:"email" uri:"email" binding:"email"`
}

type PQueryGameNFTAssets struct {
	GameID    uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	Email     string `form:"email" json:"email" uri:"email" binding:"email"`
	AssetName string `form:"asset_name" json:"asset_name" uri:"asset_name"`
	Page      int    `form:"page" json:"page" uri:"page" binding:"gt=0"`
	PageSize  uint64 `form:"pageSize" json:"pageSize" uri:"pageSize" binding:"gt=0,lte=50"`
}

type RQueryGameNFTAssets struct {
	GameERC721Assets []RDashboardGameERC721Asset `json:"gameERC721Assets`
	Total            int64                       `json:"total"`
}

type RNFTAssetResponse struct {
	Code    int64                    `json:"code"`
	Message string                   `json:"message"`
	Total   int64                    `json:"total"`
	Data    []RGameServerERC721Asset `json:"data"`
}
type PQueryGameAssetDetail struct {
	GameID        uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	TokenID       string `json:"token_id"`
	EquipmentID   string `json:"equipment_id"`
	GameAssetName string `json:"game_asset_name"`
}

// response from gameserver
type RGameServerERC20Asset struct {
	AppCoinName   string `json:"game_coin_name"`
	Balance       string `json:"coin_balance"`
	FrozenBalance string `json:"coin_frozen_balance"`
}

type SortRGameServerERC20Assets []RGameServerERC20Asset

func (s SortRGameServerERC20Assets) Len() int           { return len(s) }
func (s SortRGameServerERC20Assets) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortRGameServerERC20Assets) Less(i, j int) bool { return s[i].AppCoinName < s[j].AppCoinName }

// response from gameserver
type RGameServerERC721Asset struct {
	GameAssetName string                   `json:"game_asset_name"`
	TokenID       string                   `json:"token_id"`
	EquipmentID   string                   `json:"equipment_id"`
	Frozen        bool                     `json:"frozen"`
	Image         string                   `json:"image"`
	Attrs         []commdata.EquipmentAttr `json:"attrs"`
}

type RGameServerERC721AssetDetail struct {
	GameAssetName string                   `json:"game_asset_name"`
	EquipmentID   string                   `json:"equipment_id"`
	Frozen        bool                     `json:"frozen"`
	Image         string                   `json:"image"`
	Attrs         []commdata.EquipmentAttr `json:"attrs"`
}

type RGameAbility struct {
	Image string `json:"image"`
}

//response to dashboard

type RDashboardGameERC20Asset struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Contract    string `json:"contract"`
	Decimal     uint64 `json:"decimal"`
	Balance     string `json:"balance"` //8
	GameDecimal uint64 `json:"gameDecimal"`
}

type RDashboardGameERC721Asset struct {
	GameAssetName string                   `json:"game_asset_name"`
	Name          string                   `json:"name"`
	Symbol        string                   `json:"symbol"`
	Contract      string                   `json:"contract"`
	TokenID       string                   `json:"token_id"`
	EquipmentID   string                   `json:"equipmentId"`
	Image         string                   `json:"image"`
	Attrs         []commdata.EquipmentAttr `json:"attrs"`
}

type RDashboardGameQueryGameAssets struct {
	GameERC20Assets  []RDashboardGameERC20Asset  `json:"GameERC20Assets"`
	GameERC721Assets []RDashboardGameERC721Asset `json:"GameERC721Assets`
}

type RDashboardGameQueryAssetDetail struct {
	GameAssetName string                   `json:"game_asset_name"`
	Name          string                   `json:"name"`
	Symbol        string                   `json:"symbol"`
	Contract      string                   `json:"contract"`
	Trease        string                   `json:"treasure"`
	TokenID       string                   `json:"token_id"`
	EquipmentID   string                   `json:"equipmentId"`
	Image         string                   `json:"image"`
	ChainImage    string                   `json:"chainImage"`
	Attrs         []commdata.EquipmentAttr `json:"attrs"`
	Network       string                   `json:"chain"`
	Standard      string                   `json:"standard"`
}

// response to backend
type RBackendERC20Asset struct {
	Symbol        string `json:"symbol"`
	Contract      string `json:"address"`
	Balance       string `json:"balance"`
	FrozenBalance string `json:"frozenBalance"` //8
}

type RBackendERC721Asset struct {
	Symbol        string `json:"symbol"`
	TokenContract string `json:"address"`
	TokenID       string `json:"token_id"`
	EquipmentID   string `json:"equipmentId"`
	Balance       int    `json:"balance"`       //8
	FrozenBalance int    `json:"frozenBalance"` //8

}

type PChainFTAssets struct {
	GameID  uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	Address string `form:"address" json:"address" uri:"address" binding:"eth_addr,ne=0x0000000000000000000000000000000000000000"`
}
type RChainFTAssets struct {
	Contract        string `json:"contract"`
	Symbol          string `json:"symbol"`
	Treasure        string `json:"treasure"`
	DepositTreasure string `json:"depositTreasure"`
	GameDecimal     int    `json:"gameDecimal"`
	Decimal         int    `json:"tokenDecimal"`
	Balance         string `json:"balance"`
	OriginBalance   string `json:"originBalance"`
	ChainID         int    `json:"chainId"`
}
