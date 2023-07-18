package request

type FtContract struct {
	GameId  int `validate:"required,gte=0"`
	ChainId int `validate:"required,gte=0"`
}

type QueryEquipment struct {
	GameID    uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	AssetName string `form:"asset_name" json:"asset_name" uri:"asset_name" binding:"-"`
	Owner     string `form:"owner" json:"owner" uri:"owner"  binding:"eth_addr,ne=0x0000000000000000000000000000000000000000"`
	ChainId   uint64 `form:"chain_id" json:"chain_id" uri:"chain_id" binding:"gt=0"`
	LastKey   string `form:"lastKey" json:"lastKey" uri:"lastKey"`
	PageSize  uint64 `form:"pageSize" json:"pageSize" uri:"pageSize" binding:"gt=0,lte=20"`
}

/*
   "lootAssetsEntities": [
     {
       "tokenId": "0",
       "contract": "0x36fe92a511d7024f8c1726b30dfb784cb043295f",
       "token": {
         "name": "ArchLoot"
       }
     },
*/
