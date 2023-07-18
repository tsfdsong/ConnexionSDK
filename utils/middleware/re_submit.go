package middleware

import (
	"bytes"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ReSubmitMiddleware(prefixKey string, paramName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		jsonData, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": common.InnerError,
				"msg":  "get request body fail",
			})
			c.Abort()
			return
		}
		//replace the body with a reader that reads from the buffer
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))

		redisKey := prefixKey

		if len(paramName) > 0 {
			paramKey, err := jsonparser.GetString(jsonData, paramName...)

			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": common.InnerError,
					"msg":  "field not found",
				})
				c.Abort()
				return
			}
			redisKey = redisKey + strings.ToLower(paramKey)
		}

		luaScript := `
			if redis.call("EXISTS", KEYS[1] ) > 0 then
				return 0
			else
			redis.call("SET", KEYS[1], ARGV[1])
			redis.call("EXPIRE", KEYS[1], ARGV[2])
			return 1
			end`
		err, redisResult := redis.LuaRunWithValuefunc(luaScript, []string{redisKey}, const_def.FT_HANDLE_NOW, config.GetFTDepositWithdrawRedisLockTime())
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": common.InnerError,
				"msg":  common.ErrorMap[int(common.InnerError)],
			})
			c.Abort()
			return
		}

		if redisResult == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": common.FTRepeatClick,
				"msg":  common.ErrorMap[int(common.FTRepeatClick)],
			})
			c.Abort()
			return
		}

		c.Next()

		go func() {
			err := redis.DeleteString(redisKey)
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("redis delete failed")
			}
		}()
	}
}
