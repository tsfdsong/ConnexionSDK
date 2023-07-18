package rpc

type FtDepositReq struct {
	GameId   int    `json:"gameId"`
	ChainId  uint64 `json:"chainId"`
	From     string `json:"from"`
	Amount   string `json:"amount"`
	Nonce    string `json:"nonce"`
	Treasure string `json:"treasure"`
	Height   uint64 `json:"height"`
	Tx       string `json:"tx"`
}

type TxConfirmedReq struct {
	GameId int    `json:"gameId"`
	Tx     string `json:"tx"`
	Height uint64 `json:"height"`
	Status uint64 `json:"status"`
}

type FtWithdrawReq struct {
	GameId   int    `json:"gameId"`
	ChainId  uint64 `json:"chainId"`
	From     string `json:"from"`
	Amount   string `json:"amount"`
	Nonce    string `json:"nonce"`
	Treasure string `json:"treasure"`
	Height   uint64 `json:"height"`
	Tx       string `json:"tx"`
}

type NftWithdrawReq struct {
	GameId       int    `json:"gameId"`
	From         string `json:"from"`
	Nonce        string `json:"nonce"`
	Treasure     string `json:"treasure"`
	Height       uint64 `json:"height"`
	TokenID      string `json:"tokenId"`
	EquipID      string `json:"equipId"`
	ContractAddr string `json:"contractAddr"`
	MinterAddr   string `json:"minterAddr"`
	Tx           string `json:"tx"`
}

type NftDepositReq struct {
	GameId       int    `json:"gameId"`
	From         string `json:"from"`
	Nonce        string `json:"nonce"`
	Height       uint64 `json:"height"`
	TokenID      string `json:"tokenId"`
	ContractAddr string `json:"contractAddr"`
	TargetAddr   string `json:"targetAddr"`
	Tx           string `json:"tx"`
}
