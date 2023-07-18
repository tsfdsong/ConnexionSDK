package state

import (
	"flag"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDelete(t *testing.T) {
	configPath := flag.String("config_path", "./", "config file")
	logicLogFile := flag.String("logic_log_file", "./log/sdk.log", "logic log file")
	flag.Parse()

	//init logic logger
	logger.Init(*logicLogFile)

	//load config
	err := config.LoadConf(*configPath)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("load config failed")
		return
	}

	err = redis.InitRedis()
	if err != nil {
		logger.Logrus.Error("init redis failed")
		return
	}

	c, err := GetPrewithdrawEngine().DeleteDeadLetter(NAMESPACE, NFTTASKQUEUE, 1)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DeleteDeadLetter failed")
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Count": c}).Error("DeleteDeadLetter success")

	dc, err := GetPrewithdrawEngine().Destroy(NAMESPACE, NFTTASKQUEUE)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("Destroy failed")
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Count": dc}).Error("Destroy success")
}
