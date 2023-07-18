package sdkcontroller

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/chainlistener"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"github/Connector-Gamefi/ConnectorGoSDK/web/sdkbackend/sdkdata"
	"net/http"

	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"

	webcommon "github/Connector-Gamefi/ConnectorGoSDK/web/common"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//use for deposit && withdraw

// getLogs return []
// finaly will be parse 1 block logs
func ParseSpecifiedBlockLog(c *gin.Context) {
	r := &webcommon.Response{
		Code: webcommon.SuccessCode,
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = sdkdata.PParseSpecifiedBlockLog{}
	err := c.ShouldBind(&p)
	logger.Logrus.WithFields(logrus.Fields{"p": p}).Info("ParseSpecifiedBlockLog info")
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("validator reject")

		r.Code = webcommon.IncorrectParams
		r.Message = webcommon.ErrorMap[int(r.Code)]
		return
	}

	gameConf := config.GetGameConfig()
	if _, ok := gameConf[p.GameID]; !ok {
		r.Code = webcommon.IncorrectParams
		r.Message = webcommon.ErrorMap[int(r.Code)]
		return
	}
	listener, err := chainlistener.NewListener(gameConf[p.GameID].FilterConfig, p.GameID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("new listener error")
		r.Code = webcommon.IncorrectParams
		r.Message = webcommon.ErrorMap[int(r.Code)]
		return
	}

	err = listener.HandleSpecHeight(int64(p.Height))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Height": p.Height}).Error("parse specify height")

		r.Code = webcommon.InnerError
		r.Message = err.Error()
		return
	}

	logger.Logrus.Info(fmt.Sprintf("parse specify height success %v", p.Height))

	r.Message = webcommon.ErrorMap[int(r.Code)]
}
