package marketplace

import (
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
)

type ROrderList struct {
	LastKey   string            `json:"lastKey"`
	OrderList []SingleOrderList `json:"orderList"`
}

type SingleOrderList struct {
	Index         uint64                   `json:"index"`
	ChainImage    string                   `json:"chainImage"`
	Image         string                   `json:"image"`
	TokenName     string                   `json:"tokenName"`
	TokenID       string                   `json:"token_id"`
	Price         string                   `json:"price"`
	PurchaseToken string                   `json:"purchaseToken"`
	AssetName     string                   `json:"game_asset_name"`
	Attrs         []commdata.EquipmentAttr `json:"attrs"`
	SortIndex     int
}

type ROrderDetail struct {
	Index               uint64                   `json:"index"`
	TokenID             string                   `json:"token_id"`
	Price               string                   `json:"price"`
	OriginPrice         string                   `json:"originPrice"`
	PurchaseToken       string                   `json:"purchaseToken"`
	Status              int                      `json:"status"`
	ContractAddress     string                   `json:"contract"`
	Attrs               []commdata.EquipmentAttr `json:"attrs"`
	Image               string                   `json:"image"`
	ChainImage          string                   `json:"chainImage"`
	Name                string                   `json:"name"`
	Symbol              string                   `json:"symbol"`
	Standard            string                   `json:"standard"`
	Network             string                   `json:"chain"`
	Approved            bool                     `json:"approved"`
	MarketplaceContract string                   `json:"marketplaceContract"`
	Creator             string                   `json:"creator"`
}

type RActivity struct {
	LastKey      string           `json:"lastKey"`
	ActivityList []SingleActivity `json:"activityList"`
}

type SingleActivity struct {
	ActivityType int    `json:"activityType"`
	TokenName    string `json:"tokenName"`
	TokenSymbol  string `json:"tokenSymbol"`
	Price        string `json:"price"`
	TokenID      string `json:"token_id"`
	From         string `json:"from"`
	To           string `json:"to"`
	Time         string `json:"time"`
	Hash         string `json:"hash"`
	GameImageURL string `json:"nftUrl"`
}

// {
// 	"data": {
// 	  "pools": [
// 		{
// 		  "id": "0x1",
// 		  "index": "1",
// 		  "tokenId": "39",
// 		  "amountTotal1": "3000000000000000000"
// 		},
// 		{
// 		  "id": "0x0",
// 		  "index": "0",
// 		  "tokenId": "10",
// 		  "amountTotal1": "2000000000000000000"

type RGraphOrderListPools struct {
	Pools []RGraphOrderListPool `json:"pools"`
}
type RGraphOrderListPool struct {
	ID              string `json:"id"`
	Index           string `json:"index"`
	NFTContract     string `json:"token0"`
	TokenId         string `json:"tokenId"`
	AmountTotal1    string `json:"amountTotal1"`
	AmountTotal1Key string `json:"amountTotal1Key"`
	CreatedTimeKey  string `json:"createdTimeKey"`
}

/*
	{
	  "data": {
	    "lootAssetsEntities": [
	      {
	        "id": "0x36fe92a511d7024f8c1726b30dfb784cb043295f_0xc1",
	        "tokenId": "193",
	        "tokenURI": "ipfs://QmWb1yYPB6SexCVmHg3KDkES3zR9h398Z7trscZL6DkatN",
	        "contract": "0x36fe92a511d7024f8c1726b30dfb784cb043295f"
	      }
	    ]
	  }
	}
*/
type RGraphTokenURI struct {
	LootAssetsEntities []RLootAssetsEntities `json:"lootAssetsEntities"`
}

type RLootAssetsEntities struct {
	ID       string `json:"id"`
	TokenId  string `json:"tokenId"`
	TokenURI string `json:"tokenURI"`
	Contract string `json:"contract"`
}

/*
{
  "data": {
    "pool": {
      "id": "0x1",
      "index": "1",
      "tokenId": "39",
      "amountTotal1": "3000000000000000000",
      "closeAt": "1652371901",
      "swapped": false,
      "closed": false
    }
  }
}
*/

type RGraphOrderDetailPool struct {
	Pool RGraphOrderDetail `json:"pool"`
}
type RGraphOrderDetail struct {
	ID           string `json:"id"`
	Index        string `json:"index"`
	Token0       string `json:"token0"`
	TokenId      string `json:"tokenId"`
	AmountTotal1 string `json:"amountTotal1"`
	CloseAt      string `json:"closeAt"`
	Swapped      bool   `json:"swapped"`
	Redeemed     bool   `json:"redeemed"`
	Closed       bool   `json:"closed"`
	Creator      string `json:"creator"`
}

/*
{
  "data": {
    "activities": [
      {
        "id": "0x3ce48a094141f22b280795118ae7896e16129e8e10fbd67a9535aa9b11f327e3_81_0x5ece8b7038d83ddbb7375090a4e81e48ce062815",
        "activityType": "3",
        "amountTotal1": "22",
        "tokenId": "43",
        "activityer": "0xaeaff03acb124bf77915c0b7f2426b03dca1ce8c",
        "creator": "0x5ece8b7038d83ddbb7375090a4e81e48ce062815",
        "buyer": "0xaeaff03acb124bf77915c0b7f2426b03dca1ce8c",
        "time": "1652286372",
        "hash": "0x3ce48a094141f22b280795118ae7896e16129e8e10fbd67a9535aa9b11f327e3"
      }
    ]
  }
}
*/

type RGraphActivities struct {
	Activities []RGraphSingleActivity `json:"activities"`
}
type RGraphSingleActivity struct {
	ID           string `json:"id"`
	ActivityType string `json:"activityType"`
	Token0       string `json:"token0"`
	AmountTotal1 string `json:"amountTotal1"`
	TokenId      string `json:"tokenId"`
	Activityer   string `json:"activityer"`
	Creator      string `json:"creator"`
	Buyer        string `json:"buyer"`
	Time         string `json:"time"`
	Hash         string `json:"hash"`
	TimeKey      string `json:"timeKey"`
}
