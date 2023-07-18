package sign

type PRequestSignature struct {
	Timestamp  string `json:"timeStamp"`
	EncodeData string `json:"encodeData"`
	AssetsData string `json:"assetsData"`
	Key        string `json:"key"`
}

type RRequestSignature struct {
	Code    string `json:"code"`
	ReqHash string `json:"reqHash"`
}

type PQuerySignature struct {
	ReqHash string `json:"reqHash"`
}

type RQuerySignature struct {
	Code string `json:"code"`
	Sign string `json:"sig"`
}
