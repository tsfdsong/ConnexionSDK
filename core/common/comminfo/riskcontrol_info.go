package comminfo

import (
	"encoding/json"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
)

const (
	NameFTType  = "erc20"
	NameNFTType = "erc721"
)

// GetRiskControlConfig get all config of risk control by appid
func GetRiskControlConfig(appid int, tType, coinName string) (*commdata.RiskCacheData, error) {
	key := fmt.Sprintf("risk.cfg:%d.%s.%s", appid, tType, coinName)
	cfginfo, err := redis.GetString(key)
	if err == nil {
		var res commdata.RiskCacheData
		errn := json.Unmarshal([]byte(cfginfo), &res)
		if errn == nil {
			return &res, nil
		}
	}

	var data model.TRiskControl
	res := mysql.GetDB().Model(&model.TRiskControl{}).Where(map[string]interface{}{"app_id": appid, "token_type": tType, "game_coin_name": coinName}).First(&data)
	if res == nil || res.Error != nil {
		if res == nil {
			return nil, fmt.Errorf("res is nil")
		} else {
			return nil, res.Error
		}
	}

	result := &commdata.RiskCacheData{
		AmountLimit: data.AmountLimit,
		CountLimit:  data.CountLimit,
		CountTime:   data.CountTime,
		TotalLimit:  data.TotalLimit,
		TotalTime:   data.TotalTime,
		Status:      int(data.Status),
	}

	cfgbytes, err := json.Marshal(&result)
	if err == nil {
		redis.SetString(key, string(cfgbytes), config.GetKeyNoExpireTime())
	}

	return result, nil
}

func SetRiskControlCache(cfg *model.TRiskControl) error {
	key := fmt.Sprintf("risk.cfg:%d.%s.%s", cfg.AppID, cfg.TokenType, cfg.GameCoinName)

	result := &commdata.RiskCacheData{
		AmountLimit: cfg.AmountLimit,
		CountLimit:  cfg.CountLimit,
		CountTime:   cfg.CountTime,
		TotalLimit:  cfg.TotalLimit,
		TotalTime:   cfg.TotalTime,
		Status:      int(cfg.Status),
	}

	cfgbytes, err := json.Marshal(&result)
	if err != nil {
		return err
	}

	return redis.SetString(key, string(cfgbytes), config.GetKeyNoExpireTime())
}

func DeleteRiskControlCache(cfg *model.TRiskControl) error {
	key := fmt.Sprintf("risk.cfg:%d.%s.%s", cfg.AppID, cfg.TokenType, cfg.GameCoinName)
	return redis.DeleteString(key)
}
