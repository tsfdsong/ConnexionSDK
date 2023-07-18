package skywalking

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SkyAPM/go2sky"
)

const (
	NFTRiskAutoCode   = 1
	NFTRiskReviewCode = 2
	NFTGameRetryCode  = 3

	FTRiskAutoCode   = 4
	FTRiskReviewCode = 5
	FTGameRetryCode  = 6
)

type SkyContext struct {
	SeqNumber  int    `json:"seq"`
	ServerPath string `json:"path"`
}

func SkyPutCorrelation(serverName string, seq int, serverPath string, ctx context.Context) error {
	cont := &SkyContext{
		SeqNumber:  seq,
		ServerPath: serverPath,
	}

	bytes, err := json.Marshal(cont)
	if err != nil {
		return err
	}

	isSucc := go2sky.PutCorrelation(ctx, serverName, string(bytes))
	if !isSucc {
		return fmt.Errorf("put correlation failed")
	}
	return nil
}
