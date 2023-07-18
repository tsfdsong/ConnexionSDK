package comminfo

import (
	"flag"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestAttributeCache(t *testing.T) {
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

	attrTable, err := GetAttributeCache(1, "0x53a50c33506b2d23598d229b94798de83a86386e")
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("GetAttributeCache failed")
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"AttrTable": attrTable}).Info("GetAttributeCache")
}

func TestAttributeCacheDelete(t *testing.T) {
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

	err = DeleteAllAttributeCache()
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DeleteAllAttributeCache failed")
		return
	}

	logger.Logrus.WithFields(logrus.Fields{}).Info("DeleteAllAttributeCache success")
}
