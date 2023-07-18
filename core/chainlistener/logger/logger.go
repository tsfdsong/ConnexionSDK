package logger

import (
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"

	"github.com/sirupsen/logrus"
)

func GameLogger(gameId int) *logrus.Entry {
	return logger.Logrus.WithFields(logrus.Fields{"AppId": gameId})
}
