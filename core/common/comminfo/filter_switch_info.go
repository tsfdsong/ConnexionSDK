package comminfo

import (
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
)

const (
	CacheFilterLogSwitch = "filter.log.switch"
)

func GetFilterLogSwitch() (int, error) {
	res, err := redis.GetInt(CacheFilterLogSwitch)
	if err == nil {
		return res, nil
	}

	heightInfo := []model.TBlockHeight{}
	err = mysql.WrapFindAll(model.TableBlockHeight, &heightInfo)
	if err != nil {
		return 0, err
	}

	if len(heightInfo) < 1 {
		return 0, errors.New("height info not exist")
	}

	err = redis.SetString(CacheFilterLogSwitch, heightInfo[0].Switch, config.GetKeyExpireTime())
	if err != nil {
		return 0, err
	}

	return int(heightInfo[0].Switch), nil
}

func DeleteFilterLogSwitch() error {
	return redis.DeleteString(CacheFilterLogSwitch)
}

func GetFilterLogSwitchByApp(appId int) (int, error) {
	redisKey := fmt.Sprintf(CacheFilterLogSwitch+".%d", appId)
	res, err := redis.GetInt(redisKey)
	if err == nil {
		return res, nil
	}

	blockHeight := model.TBlockHeight{}
	err = mysql.GetDB().Where("app_id = ?", appId).First(&blockHeight).Error
	if err != nil {
		return 0, err
	}

	err = redis.SetString(redisKey, blockHeight.Switch, config.GetKeyExpireTime())
	if err != nil {
		return 0, err
	}

	return int(blockHeight.Switch), nil
}

func DeleteFilterLogSwitchByApp(appId int) error {
	redisKey := fmt.Sprintf(CacheFilterLogSwitch+".%d", appId)
	return redis.DeleteString(redisKey)
}
