package web

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAuth(c *gin.Context) {
	appIDStr := c.Query("appId")
	account := strings.ToLower(c.Query("account"))
	address := strings.ToLower(c.Query("address"))

	appid, err := strconv.ParseInt(appIDStr, 0, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.AuthFailed,
			"msg":  fmt.Sprintf("convert appid string to int, %v", err),
		})
		return
	}

	bindinfo, err := comminfo.GetBindInfo(int(appid), address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.AuthFailed,
			"msg":  fmt.Sprintf("get address bind info, %v", err),
		})
		return
	}

	if bindinfo.Account != account {
		c.JSON(http.StatusOK, gin.H{
			"code": common.AuthFailed,
			"msg":  fmt.Sprintf("input account %v is not match bind account %v", account, bindinfo.Account),
		})
		return
	}

	token, err := common.GenerateToken(account, address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.AuthFailed,
			"msg":  fmt.Sprintf("generate token, %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": common.SuccessCode,
		"msg":  "get token success",
		"data": token,
	})
}
