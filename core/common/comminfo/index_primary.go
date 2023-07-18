package comminfo

import (
	"encoding/json"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"sort"

	"github.com/sirupsen/logrus"
)

//all is weak consitency
//HGetAll return map[string]string. so use json serialized string

/*
query with unique index. and fill to val

step1 query cache primarykey by index
1.1 if not found. query from db and set index-primarykey-val. return val
1.2 if found. go on

step2 query cache val by primarykey
2.1 if not found. query from db and set index-primarykey-val. return val
2.2 if found. return val
*/

//MUST--NOTICE !!!! val is a model point(has ID value)
func C_QueryByUniqueIndex(val interface{}, tableName string, condition map[string]interface{}, expire int64) (error, bool) {
	if len(condition) <= 0 || tableName == "" {
		return nil, false
	}

	//build indexKey
	indexKey := tableName + ":"
	keys := []string{}
	for k := range condition {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		indexKey += fmt.Sprintf("_%+v=%+v", k, condition[k])
	}

	//query primaryKey by indexKey
	primaryKey, err := redis.GetStringAcceptable(indexKey)
	//logger.Logrus.WithFields(logrus.Fields{"indexKey": indexKey, "primarKey": primaryKey, "err": err}).Info("GetPrimaryKey")
	if err != nil || primaryKey == "" {
		if err == redis.Nil || primaryKey == "" {
			//primaryKey not found
			err, found := mysql.WrapFindFirst(tableName, val, condition)
			if err != nil {
				return err, false
			}
			if !found {
				return nil, false
			}
			//found at db. must return correct. and lua set redis

			//get id && json string
			id := tools.GetTableID(val)
			if id == "" {
				logger.Logrus.WithFields(logrus.Fields{"tablename": tableName}).Error("has no ids")

				return nil, true
			}

			primaryKey = BuildPrimaryKeyString(tableName, id)
			data, err := json.Marshal(val)
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"error": err.Error()}).Error("json marshal faied")

				return nil, true
			}
			value := string(data)

			//lua set redis
			luaScript := `
			redis.call("SET",KEYS[1],ARGV[1])
			redis.call("SET",KEYS[2],ARGV[2])
			redis.call("EXPIRE",KEYS[1],ARGV[3])
			redis.call("EXPIRE",KEYS[2],ARGV[4])
			return 1
			`

			err = redis.LuaRun(luaScript, []string{indexKey, primaryKey}, primaryKey, value, expire, expire)
			if err != nil {
				return nil, true
			}

			return nil, true
		} else {
			return err, false
		}
	}
	//query value by primary key
	value, err := redis.GetStringAcceptable(primaryKey)
	//logger.Logrus.WithFields(logrus.Fields{"parimaryKey": primaryKey, "value": value, "err": err}).Info("GetValue")
	if err != nil || value == "" {
		if err == redis.Nil || value == "" {
			//value not found
			err, found := mysql.WrapFindFirst(tableName, val, condition)
			if err != nil {
				return err, false
			}
			if !found {
				return nil, false
			}
			//found at db. must return correct. and may lua set redis

			//get id && json string
			id := tools.GetTableID(val)
			if id == "" {
				logger.Logrus.WithFields(logrus.Fields{"tablename": tableName}).Error("has no ids")

				return nil, true
			}

			newPrimaryKey := BuildPrimaryKeyString(tableName, id)
			data, err := json.Marshal(val)
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"error": err.Error()}).Error("json marshal failed")

				return nil, true
			}

			if newPrimaryKey == primaryKey {
				err := redis.SetString(primaryKey, string(data), expire)
				if err != nil {
					return nil, true
				}
			} else {
				//lua set it
				luaScript := `
				redis.call("SET",KEYS[1],ARGV[1])
				redis.call("SET",KEYS[2],ARGV[2])
				redis.call("EXPIRE",KEYS[1],ARGV[3])
				redis.call("EXPIRE",KEYS[2],ARGV[4])
				return 1
				`
				err = redis.LuaRun(luaScript, []string{indexKey, newPrimaryKey}, newPrimaryKey, string(data), expire, expire)
				if err != nil {
					return nil, true
				}
			}

			return nil, true
		} else {
			return err, false
		}
	} else {
		//json unmarshal
		err := json.Unmarshal([]byte(value), val)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"error": err.Error()}).Error("json unmarshal failed")

			return err, false
		}
		return nil, true
	}
}

func BuildPrimaryKeyString(tableName, id string) string {
	return fmt.Sprintf("%+v_id:%+v", tableName, id)
}

func BuildPrimaryKeyNumber(tableName string, id uint64) string {
	return fmt.Sprintf("%+v_id:%+v", tableName, id)
}
