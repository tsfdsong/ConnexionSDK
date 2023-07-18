package commdata

type NFTContractCacheData struct {
	ID              uint64 `json:"id"`
	AppID           int    `json:"appId"`
	ChainID         int    `json:"chain_id"`
	Treasure        string `json:"trease"`
	ContractAddress string `json:"contract"`
	MinterAddress   string `json:"minter"`

	TokenName   string `json:"token_name"`
	TokenSymbol string `json:"token_symbol"`
	TokenSupply string `json:"token_supply"`

	GameAssetName string `json:"game_asset_name"`
	Decimal       int    `json:"decimal"`

	DepositSwitch  int8   `json:"deposit_switch"`
	WithdrawSwitch int8   `json:"withdraw_switch"`
	BaseURL        string `json:"base_url"`
}

type FTContractCacheData struct {
	ID              uint64 `json:"id"`
	AppID           int    `json:"appId"`
	ChainID         int    `json:"chain_id"`
	Treasure        string `json:"trease"`
	DepositTreasure string `json:"deposit_treasure"`
	ContractAddress string `json:"contract"`

	TokenName    string `json:"token_name"`
	TokenSymbol  string `json:"token_symbol"`
	TokenSupply  string `json:"token_supply"`
	TokenDecimal int    `json:"token_decimal"`

	GameCoinName string `json:"game_coin_name"`
	GameDecimal  int    `json:"game_decimal"`

	DepositSwitch  int8 `json:"deposit_switch"`
	WithdrawSwitch int8 `json:"withdraw_switch"`
}

type NFTContractTreaseCacheData struct {
	Treasure string `json:"trease"`
}
