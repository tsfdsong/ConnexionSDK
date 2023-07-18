package commdata

type RiskCacheData struct {
	AmountLimit string `json:"amount_limit"`
	CountLimit  int    `json:"count_limit"`
	CountTime   int    `json:"count_time"`
	TotalLimit  string `json:"total_limit"`
	TotalTime   int    `json:"total_time"`
	Status      int    `json:"status"`
}
