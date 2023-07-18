package comminfo

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
)

//GetBaseUrl get base server url
func GetBaseUrl(appid int) (string, error) {
	key := fmt.Sprintf("game.info:%d.baseurl", appid)
	res, err := redis.GetString(key)
	if err == nil {
		return res, nil
	}

	var data model.TGameInfo
	err = mysql.GetDB().Model(&model.TGameInfo{}).Where("app_id=?", appid).First(&data).Error
	if err != nil {

		return "", fmt.Errorf("get base url of %d, %v", appid, err)
	}

	redis.SetString(key, data.BaseServerURL, config.GetKeyNoExpireTime())

	return data.BaseServerURL, nil
}

func GetBaseUrlOriginError(appid int) (string, error) {
	key := fmt.Sprintf("game.info:%d.baseurl", appid)
	res, err := redis.GetString(key)
	if err == nil {
		return res, nil
	}

	var data model.TGameInfo
	err = mysql.GetDB().Model(&model.TGameInfo{}).Where("app_id=?", appid).First(&data).Error
	if err != nil {

		return "", err
	}

	redis.SetString(key, data.BaseServerURL, config.GetKeyNoExpireTime())

	return data.BaseServerURL, nil
}

//GetAppSecret get app secret
func GetAppSecret(appid int) (string, error) {
	key := fmt.Sprintf("app.secret:%d", appid)
	res, err := redis.GetString(key)
	if err == nil {
		return res, nil
	}

	var data model.TGameInfo

	err = mysql.GetDB().Model(&model.TGameInfo{}).Where("app_id=?", appid).First(&data).Error
	if err != nil {
		return "", err
	}

	appsecret := data.AppSecret
	redis.SetString(key, appsecret, config.GetKeyNoExpireTime())

	return appsecret, nil
}

func DeleteGameInfoCache(appid int) error {
	key := fmt.Sprintf("game.info:%d.baseurl", appid)
	err := redis.DeleteString(key)
	if err != nil {
		return err
	}

	keys := fmt.Sprintf("app.secret:%d", appid)
	err = redis.DeleteString(keys)
	if err != nil {
		return err
	}
	return nil
}
