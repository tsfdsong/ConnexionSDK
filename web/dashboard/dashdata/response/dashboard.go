package response

import "github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"

type FTContract struct {
	Address string `json:"address"`
	Symbol  string `json:"symbol"`
	Decimal int    `json:"decimal"`
}

type NFTAsset struct {
	Name     string      `json:"name"`
	Contract string      `json:"contract"`
	TokenID  int         `json:"token_id"`
	TokenURI string      `json:"tokenURI"`
	Attrs    []Attribute `json:"attrs"`
}

type NFTChainAssets struct {
	ChainImage    string                   `json:"chainImage"`
	Image         string                   `json:"image"`
	Name          string                   `json:"name"`
	Contract      string                   `json:"contract"`
	GameAssetName string                   `json:"game_asset_name"`
	Trease        string                   `json:"treasure"`
	TokenID       string                   `json:"token_id"`
	Attrs         []commdata.EquipmentAttr `json:"attrs"`
	SortIndex     int
}

type QueryEquipmentResp struct {
	LastKey   string           `json:"lastKey"`
	NFTAssets []NFTChainAssets `json:"nftAssets"`
}

type Attribute struct {
	AttrID    int    `json:"attr_id"`
	AttrValue string `json:"attr_value"`
}

type Response struct {
	Code    int64       `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type U8ListResponse struct {
	Code    int64       `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
}

type MResponse struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
