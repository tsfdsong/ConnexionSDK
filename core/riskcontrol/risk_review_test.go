package riskcontrol

import (
	"flag"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/pool"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNFTRiskControl(t *testing.T) {
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

	db := mysql.GetDB()
	if db == nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("init db failed")
		return
	}

	err = redis.InitRedis()
	if err != nil {
		logger.Logrus.Error("init redis failed")
		return
	}
	//init http pool
	pool.InitClient(int(config.GetHttpPoolSize()))

	appID := int(2)
	account := "tsfdsong@163.com"
	contractAddress := "0x6adda3d65b65f8870f07a2a3e1e7b182fd259bea"
	statusCode, err := NFTRiskReview(appID, account, contractAddress, "447273373076029444")
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("NFTRiskReview  failed")

		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": statusCode}).Info("NFTRiskReview success")

}
