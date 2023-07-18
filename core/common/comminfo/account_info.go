package comminfo

import (
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"strings"
)

// GetBindInfo get bind info
func GetBindInfo(appID int, address string) (*commdata.BindCacheData, error) {

	var data = model.TEmailBind{}
	con := map[string]interface{}{"app_id": appID, "address": strings.ToLower(address)}
	expire := config.GetKeyExpireTime()
	err, found := C_QueryByUniqueIndex(&data, model.TableEmailBind, con, expire)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errors.New("can't find data by condition")
	}

	res := commdata.BindCacheData{
		UID:     data.UID,
		Account: data.Account,
	}

	return &res, nil
}

// GetBindInfoWithChainID get bind info
func GetBindInfoWithChainID(appID int, address string, chainID int64) (*commdata.BindCacheData, error) {
	var data = model.TEmailBind{}
	if chainID == 56 || chainID == 97 {
		con := map[string]interface{}{"app_id": appID, "address": strings.ToLower(address)}
		expire := config.GetKeyExpireTime()
		err, found := C_QueryByUniqueIndex(&data, model.TableEmailBind, con, expire)
		if err != nil {
			return nil, err
		}

		if !found {
			return nil, errors.New("can't find bsc data by condition")
		}
	} else {
		con := map[string]interface{}{"app_id": appID, "zks_address": strings.ToLower(address)}
		expire := config.GetKeyExpireTime()
		err, found := C_QueryByUniqueIndex(&data, model.TableEmailBind, con, expire)
		if err != nil {
			return nil, err
		}

		if !found {
			return nil, errors.New("can't find other data by condition")
		}
	}

	res := commdata.BindCacheData{
		UID:     data.UID,
		Account: data.Account,
	}

	return &res, nil
}

// GetUserWithdrawSwitch
func GetUserWithdrawSwitch(appID int, account string) (int, error) {

	var data = model.TUser{}
	con := map[string]interface{}{"app_id": appID, "account": account}
	expire := config.GetKeyExpireTime()
	err, found := C_QueryByUniqueIndex(&data, model.TableUser, con, expire)
	if err != nil {
		return 0, err
	}

	if !found {
		return 0, errors.New("can't find data by condition")
	}

	return int(data.WithdrawSwitch), nil

}

func GetEmailBind(condition map[string]interface{}, val *model.TEmailBind) (error, bool) {
	expire := config.GetKeyExpireTime()
	err, found := C_QueryByUniqueIndex(val, model.TableEmailBind, condition, expire)
	return err, found
}

func GetEmailBindInfo(appID, chainID uint, userAddr string) (*model.TEmailBind, error) {
	var bindCondition map[string]interface{}
	if chainID != 97 && chainID != 56 {
		bindCondition = map[string]interface{}{"app_id": appID, "zks_address": strings.ToLower(userAddr)}
	} else {
		bindCondition = map[string]interface{}{"app_id": appID, "address": strings.ToLower(userAddr)}
	}

	//TODO should modify if multi chain
	bindEmail := &model.TEmailBind{}
	expire := config.GetKeyExpireTime()
	err, found := C_QueryByUniqueIndex(bindEmail, model.TableEmailBind, bindCondition, expire)
	if err != nil {
		return nil, err
	}

	if !found || bindEmail.UID == 0 {
		return nil, fmt.Errorf("%s not bind email", userAddr)
	}

	return bindEmail, nil
}

func GetUIDByAccount(appID uint, account string) (uint64, error) {
	bindCondition := map[string]interface{}{"app_id": appID, "account": account}

	//TODO should modify if multi chain
	bindEmail := &model.TEmailBind{}
	expire := config.GetKeyExpireTime()
	err, found := C_QueryByUniqueIndex(bindEmail, model.TableEmailBind, bindCondition, expire)
	if err != nil {
		return 0, err
	}

	if !found || bindEmail.UID == 0 {
		return 0, fmt.Errorf("%s not bind email", account)
	}

	return bindEmail.UID, nil
}

func GetUserInfo(condition map[string]interface{}, val *model.TUser) (error, bool) {
	expire := config.GetKeyExpireTime()
	err, found := C_QueryByUniqueIndex(val, model.TableUser, condition, expire)
	return err, found
}
