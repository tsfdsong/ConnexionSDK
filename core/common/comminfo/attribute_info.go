package comminfo

import (
	"encoding/json"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"strings"
)

const (
	CachePrefixAttributeName = "attribute.name"
)

type AttributeCache struct {
	AttributeTable map[uint64]string `json:"attr_table"`
}

func GetAttributeCache(appid int, address string) (map[uint64]string, error) {
	key := fmt.Sprintf("%s:%d.%s", CachePrefixAttributeName, appid, strings.ToLower(address))
	res, err := redis.GetString(key)
	if err == nil {
		var data AttributeCache
		errn := json.Unmarshal([]byte(res), &data)
		if errn != nil && data.AttributeTable != nil {
			return data.AttributeTable, nil
		}
	}

	var attrList []model.TAttribute
	err = mysql.WrapFindAllByCondition(model.TableAttrbute, map[string]interface{}{"app_id": appid, "contract_address": strings.ToLower(address)}, &attrList)
	if err != nil {
		return nil, fmt.Errorf("get all attr list {%d} {%s} failed,%v", appid, address, err)
	}

	attrNameList := make(map[uint64]string, 0)
	for _, v := range attrList {
		attrNameList[v.AttrID] = v.AttrName
	}

	data := &AttributeCache{
		AttributeTable: attrNameList,
	}

	resbytes, err := json.Marshal(&data)
	if err != nil {
		return nil, fmt.Errorf("marshal attrbute data failed, %v", err)
	}

	err = redis.SetString(key, string(resbytes), config.GetKeyNoExpireTime())
	if err != nil {
		return nil, fmt.Errorf("set attribute cache failed , %v", err)
	}

	return attrNameList, nil
}

func DeleteAllAttributeCache() error {
	matckkey := CachePrefixAttributeName + "*"
	err := redis.ClearByKey(matckkey)
	if err != nil {
		return fmt.Errorf("DeleteAllAttributeCache %v", err)
	}
	return nil
}
