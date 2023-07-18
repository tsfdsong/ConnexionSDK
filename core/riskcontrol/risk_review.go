package riskcontrol

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/math"
	"math/big"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// FTRiskReview FT risk control and return risk status code and error
func FTRiskReview(appID int, account, amount, coinName string) (int, error) {
	cfg, err := comminfo.GetRiskControlConfig(appID, comminfo.NameFTType, coinName)
	if err == gorm.ErrRecordNotFound {
		logger.Logrus.WithFields(logrus.Fields{"RiskConfig": cfg, "Amount": amount, "Account": account, "GameCoinName": coinName}).Info("FTRiskReview no found config record")

		return CodeRiskSuccess, nil
	}

	if err != nil {
		return CodeRiskNone, err
	}

	//check risk control is closed or not
	if cfg.Status < CodeRiskControlOpen {
		logger.Logrus.WithFields(logrus.Fields{"RiskConfig": cfg, "Amount": amount, "Account": account, "GameCoinName": coinName}).Info("FTRiskReview is closed")

		return CodeRiskClosed, nil
	}

	ftContractInfo, err := comminfo.GetFirstFTContractByCoinName(appID, coinName)
	if err != nil {
		return CodeRiskNone, err
	}

	decimal := ftContractInfo.TokenDecimal

	//require withdraw amount <= cfg.AmountLimit
	amountLimit := math.ToWei(cfg.AmountLimit, decimal)

	curAmount, err := math.NewFromString(amount)
	if err != nil {
		return CodeRiskNone, fmt.Errorf("NewFromString %v", err)
	}

	if curAmount.Cmp(amountLimit) > 0 {
		logger.Logrus.WithFields(logrus.Fields{"RiskConfig": cfg, "Amount": amount, "Account": account, "GameCoinName": coinName}).Info("FTRiskReview withdraw amount is more than the amount limit")

		return CodeRiskWaiting, nil
	}

	//require count time and limit
	actCount, err := GetFTRiskCountCache(appID, account, coinName)
	if err != nil {
		return CodeRiskNone, fmt.Errorf("GetFTRiskCountCache %v", err)
	}

	if actCount >= cfg.CountLimit {
		logger.Logrus.WithFields(logrus.Fields{"RiskConfig": cfg, "Amount": amount, "CountLimit": actCount, "Account": account, "GameCoinName": coinName}).Info("FTRiskReview more than count limit")

		return CodeRiskWaiting, nil
	}

	//require total time and limit
	actSum, err := GetFTRiskTotalCache(appID, account, coinName)
	if err != nil {
		return CodeRiskNone, fmt.Errorf("GetFTRiskCountCache %v", err)
	}

	actTotal := new(big.Int).Add(actSum, curAmount)

	total := math.ToWei(cfg.TotalLimit, decimal)
	if actTotal.Cmp(total) > 0 {
		logger.Logrus.WithFields(logrus.Fields{"RiskConfig": cfg, "Amount": amount, "ActualTotalAmount": actTotal.String(), "TotalAmount": total.String(), "Decimal": decimal, "Account": account, "GameCoinName": coinName}).Info("FTRiskReview more than total limit")

		return CodeRiskWaiting, nil
	}

	err = SetFTRiskCountCache(appID, account, amount, coinName, int64(cfg.CountTime))
	if err != nil {
		return CodeRiskNone, fmt.Errorf("SetFTRiskCountCache %v", err)
	}

	err = SetFTRiskTotalCache(appID, account, amount, coinName, int64(cfg.TotalTime))
	if err != nil {
		return CodeRiskNone, fmt.Errorf("SetFTRiskCountCache %v", err)
	}

	logger.Logrus.WithFields(logrus.Fields{"RiskConfig": cfg, "Amount": amount, "CountLimit": actCount, "TotalLimit": actTotal.String(), "Account": account}).Info("FTRiskReview risk control review")

	return CodeRiskSuccess, nil
}

func NFTRiskReview(appID int, account, contractAddr, equipmentID string) (int, error) {
	cfg, err := comminfo.GetRiskControlConfig(appID, comminfo.NameNFTType, contractAddr)
	if err == gorm.ErrRecordNotFound {
		logger.Logrus.WithFields(logrus.Fields{"RiskConfig": cfg, "Account": account}).Info("NFTRiskReview risk control review pass when no config")

		return CodeRiskSuccess, nil
	}

	if err != nil {
		return CodeRiskNone, err
	}

	//check risk control is closed or not
	if cfg.Status < CodeRiskControlOpen {
		logger.Logrus.WithFields(logrus.Fields{"RiskConfig": cfg, "Account": account}).Info("NFTRiskReview risk control is closed")

		return CodeRiskClosed, nil
	}

	//require count time and limit
	actCount, err := GetNFTRiskCountCache(appID, account, contractAddr)
	if err != nil {
		return CodeRiskNone, fmt.Errorf("GetNFTRiskCountCache %v", err)
	}

	if actCount >= cfg.CountLimit {
		logger.Logrus.WithFields(logrus.Fields{"RiskConfig": cfg, "CountLimit": actCount, "Account": account}).Info("NFTRiskReview more than count limit")

		return CodeRiskWaiting, nil
	}

	err = SetNFTRiskCountCache(appID, account, contractAddr, equipmentID, int64(cfg.CountTime))
	if err != nil {
		return CodeRiskNone, fmt.Errorf("SetNFTRiskCountCache %v", err)
	}

	logger.Logrus.WithFields(logrus.Fields{"RiskConfig": cfg, "CountLimit": actCount, "Account": account}).Info("NFTRiskReview risk control review")

	return CodeRiskSuccess, nil
}
