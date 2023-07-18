package marketplace

// use for marketplace orderlist and profile/onsale
type POrderList struct {
	GameID    uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	AssetName string `form:"asset_name" json:"asset_name" uri:"asset_name" binding:"-"`
	Creator   string `form:"creator" json:"creator" uri:"creator" binding:"-"`
	PriceSort int    `form:"priceSort" json:"priceSort" uri:"priceSort" binding:"-"`
	PageSize  uint64 `form:"pageSize" json:"pageSize" uri:"pageSize" binding:"gt=0,lte=20"`
	LastKey   string `form:"lastKey" json:"lastKey" uri:"lastKey"`
}

type POrderDetail struct {
	GameID               uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	AssetName            string `form:"game_asset_name" json:"game_asset_name" uri:"asset_name" binding:"-"`
	UserAddress          string `form:"userAddress" json:"userAddress" uri:"userAddress" binding:"eth_addr,ne=0x0000000000000000000000000000000000000000"`
	Index                uint64 `form:"index" json:"index" uri:"index" binding:"-"`
	TokenID              string `form:"token_id" json:"token_id" uri:"token_id" binding:"gte=0"`
	OwnedContractAddress string `form:"ownedContractAddr" json:"ownedContractAddr" uri:"ownedContractAddr" binding:"-"`
	Owned                bool   `form:"owned" json:"owned" uri:"owned" binding:"-"`
}

type PActivity struct {
	GameID    uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	AssetName string `form:"asset_name" json:"asset_name" uri:"asset_name" binding:"-"`
	Address   string `form:"address" json:"address" uri:"address" binding:"eth_addr,ne=0x0000000000000000000000000000000000000000"`
	List      bool   `form:"list" json:"list" uri:"list" binding:"-"`
	Canceled  bool   `form:"canceled" json:"canceled" uri:"canceled" binding:"-"`
	Purchase  bool   `form:"purchase" json:"purchase" uri:"purchase" binding:"-"`
	Redeem    bool   `form:"redeem" json:"redeem" uri:"redeem" binding:"-"`
	Sale      bool   `form:"sale" json:"sale" uri:"sale" binding:"-"`
	LastKey   string `form:"lastKey" json:"lastKey" uri:"lastKey"`
	PageSize  uint64 `form:"pageSize" json:"pageSize" uri:"pageSize" binding:"gt=0,lte=20"`
}
