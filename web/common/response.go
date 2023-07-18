package common

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

type GameResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
