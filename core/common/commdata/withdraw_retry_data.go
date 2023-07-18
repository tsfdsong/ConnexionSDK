package commdata

type PrewithdrawRetryData struct {
	ID            uint64          `json:"id"`
	AppID         int             `json:"appId"`
	UID           uint64          `json:"uid"`
	EquipmentID   string          `json:"equipment_id"`
	GameAssetName string          `json:"game_asset_name"`
	Nonce         string          `json:"nonce"`
	Attrs         []EquipmentAttr `json:"attrs"`
	TraceID       string          `json:"trace_id"`
	SeqNumber     int             `json:"seq_number"`
}

type FTPrewithdrawRetryData struct {
	ContractAddress string `json:"contract"`
	UserAddress     string `json:"user_address"`
	Nonce           string `json:"nonce"`
	Account         string `json:"account'`

	AppID       int    `json:"appId"`
	Uid         int64  `json:"uid"`
	AppCoinName string `json:"game_coin_name"`
	Amount      string `json:"amount"`
	TraceID     string `json:"trace_id"`
	SeqNumber   int    `json:"seq_number"`
	ChainID     int64  `json:"chainId"`
}
