package ingame

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"reflect"
	"sort"
	"strings"
)

//InGameSign get game input parameter sign data
func InGameSign(appid int, paras map[string]interface{}) (string, error) {
	keys := make([]string, 0)
	maps := make(map[string]string, 0)

	for k, v := range paras {
		keys = append(keys, k)

		switch reflect.TypeOf(v).Kind() {
		case reflect.String, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			maps[k] = fmt.Sprintf("%v", v)
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			maps[k] = fmt.Sprintf("%v", v)
		case reflect.Struct, reflect.Slice:
			bytes, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			maps[k] = string(bytes)
		default:
			return "", fmt.Errorf("%v cannot convert to string", v)
		}
	}

	sort.Strings(keys)

	appsecret, err := comminfo.GetAppSecret(appid)
	if err != nil {
		return "", err
	}

	var rawStr string
	for _, val := range keys {
		if value, ok := maps[val]; ok {
			rawStr = rawStr + fmt.Sprintf("%s=%s&", val, value)
		}
	}

	rawStr = rawStr + fmt.Sprintf("secretKey=%s", appsecret)

	res := md5.Sum([]byte(rawStr))
	mdHash := fmt.Sprintf("%x", res)

	return strings.ToUpper(mdHash), nil
}
