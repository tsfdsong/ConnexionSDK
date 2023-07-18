package dashdata

type RGraphERC20Info struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Decimal string `json:"decimal"`
}

type RGraphERC20 struct {
	Contract  RGraphERC20Info `json:"contract"`
	NumTokens string          `json:"numTokens"`
}

type RGraphERC20s struct {
	OwnerPerTokenContracts []RGraphERC20 `json:"ownerPerTokenContracts"`
}

type RGraph20 struct {
	Data RGraphERC20s `json:"data"`
}

type RGraphERC721Info struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type RGameERC721 struct {
	AttrIDs    []string         `json:"attrID"`
	AttrValues []string         `json:"attrValue"`
	Contract   RGraphERC721Info `json:"contract"`
	TokenID    string           `json:"tokenID"`
	TokenURI   string           `json:"tokenURI"`
}

type RGraphERC721s struct {
	Tokens []RGameERC721 `json:"tokens"`
}

type RGraph721 struct {
	Data RGraphERC721s `json:"data"`
}

type RGraphLootAssets struct {
	Assets []RGraphSingleLootAssets `json:"lootAssetsEntities"`
}

type RGraphSingleLootAssets struct {
	TokenId  string `json:"tokenId"`
	TokenURI string `json:"tokenURI"`
	Contract string `json:"contract"`
}
