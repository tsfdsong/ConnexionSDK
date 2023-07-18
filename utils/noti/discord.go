package noti

import (
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/http_client"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"

	"github.com/sirupsen/logrus"
)

type discordRet struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

func DiscordNoti(content, url string) error {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("alert.Alert failed")
		}
	}()
	if url == "" {
		return errors.New("empty discord url")
	}
	v := discordRet{}
	body := map[string]string{"content": content}
	err := http_client.HttpClientReq(url, body, &v)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "content": content}).Error("discord noti failed")
		return err
	}

	if v.Code != 0 {
		logger.Logrus.WithFields(logrus.Fields{"content": content, "code": v.Code, "msg": v.Msg}).Error("discord noti failed")
		return errors.New(fmt.Sprintf("err is:%s", v.Msg))
	}
	return nil
}
