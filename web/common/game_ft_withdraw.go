package common

type PERC20PreWithdraw struct {
	Uid         int64  `json:"uid"`
	AppCoinName string `json:"game_coin_name"`
	Amount      string `json:"amount"`
	Nonce       string `json:"nonce"`
}

type PGameERC20PreWithdraw struct {
	Params   PERC20PreWithdraw `json:"params"`
	AppID    int               `json:"appId"`
	SignHash string            `json:"sign"`
}

type RERC20PreWithdraw struct {
	Code    int64                 `json:"code"`
	Message string                `json:"msg"`
	Data    RERC20PreWithdrawData `json:"data"`
}

type RERC20PreWithdrawData struct {
	AppOrderID string `json:"app_order_id"`
	Status     int    `json:"status"` //1 pass 2 reject
}

type PERC20WithdrawComfirm struct {
	AppOrderID   string `json:"app_order_id"`
	Nonce        string `json:"nonce"`
	GameCoinName string `json:"game_coin_name"`
	Operation    int    `json:"operation"`
	Uid          int64  `json:"uid"`
}

type PGameERC20WithdrawComfirm struct {
	Params   []PERC20WithdrawComfirm `json:"params"`
	AppID    int                     `json:"appId"`
	SignHash string                  `json:"sign"`
}

type RERC20WithdrawComfirm struct {
	Code    int64                         `json:"code"`
	Message string                        `json:"msg"`
	Data    []RERC20WithdrawComfirmStatus `json:"data"`
}

type RERC20WithdrawComfirmStatus struct {
	AppOrderID   string `json:"app_order_id"`
	Nonce        string `json:"nonce"`
	GameCoinName string `json:"game_coin_name"`
	Status       int    `json:"status"` //1 pass 2 reject
}
