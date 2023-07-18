package comminfo

import (
	"flag"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDeleteAllContractCache(t *testing.T) {
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

	err = DeleteAllContractCache()
	if err != nil {
		fmt.Printf("DeleteAllContractCache batch del failed, %v", err)
		return
	}

	fmt.Printf("success")
}
