package riskcontrol

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	PrefixFTCount = "ft.count"
	PrefixFTTotal = "ft.total"

	PrefixNFTCount = "nft.count"
)

func SetFTRiskCountCache(appID int, account, amount, coinName string, expire int64) error {
	now := time.Now().UTC().String()
	key := fmt.Sprintf("%s.%d.%s.%s:%s", PrefixFTCount, appID, coinName, account, now)
	return redis.SetString(key, amount, expire)
}

func GetFTRiskCountCache(appID int, account, coinName string) (int, error) {
	regkey := fmt.Sprintf("%s.%d.%s.%s*", PrefixFTCount, appID, coinName, account)
	return redis.GetAllLength(regkey)
}

func SetFTRiskTotalCache(appID int, account, amount, coinName string, expire int64) error {
	now := time.Now().UTC().String()
	key := fmt.Sprintf("%s.%d.%s.%s:%s", PrefixFTTotal, appID, coinName, account, now)
	return redis.SetString(key, amount, expire)
}

func GetFTRiskTotalCache(appID int, account, coinName string) (*big.Int, error) {
	regkey := fmt.Sprintf("%s.%d.%s.%s*", PrefixFTTotal, appID, coinName, account)
	values, err := redis.GetAllValues(regkey)
	if err != nil {
		return nil, fmt.Errorf("scan failed, %v", err)
	}

	sum := new(big.Int)
	for _, v := range values {
		val, err := decimal.NewFromString(v)
		if err != nil {
			return nil, err
		}

		sum = new(big.Int).Add(sum, val.BigInt())
	}

	logger.Logrus.WithFields(logrus.Fields{"GameCoinName": coinName, "Values": values, "Sum": sum.String()}).Info("GetFTRiskTotalCache info")

	return sum, nil
}

func SetNFTRiskCountCache(appID int, account, contractAddr, orderid string, expire int64) error {
	now := time.Now().UTC().String()
	key := fmt.Sprintf("%s.%d.%s.%s:%s", PrefixNFTCount, appID, contractAddr, account, now)
	value := fmt.Sprintf("%s.%s", orderid, now)
	return redis.SetString(key, value, expire)
}

func GetNFTRiskCountCache(appID int, account, contractAddr string) (int, error) {
	regkey := fmt.Sprintf("%s.%d.%s.%s*", PrefixNFTCount, appID, contractAddr, account)
	return redis.GetAllLength(regkey)
}
