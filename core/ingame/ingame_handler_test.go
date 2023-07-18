package ingame

import (
	"flag"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/pool"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSign(t *testing.T) {
	appId := 10
	uid := 1001
	orderId := "1642043144"
	op := 0
	body := &WithdrawData{
		AppID: appId,
		Params: []CommitWithdrawData{
			{
				UID:        uint64(uid),
				AppOrderID: orderId,
				Operate:    op,
			},
		},
	}

	in := make(map[string]interface{}, 0)
	in["appId"] = body.AppID
	in["params"] = body.Params
	sidata, err := InGameSign(body.AppID, in)
	if err != nil {
		t.Fatalf("game sign: %v", err)
		return
	}
	body.SignHash = sidata

	fmt.Printf("game sign data: %v\n", sidata)
	fmt.Printf("game input: %v\n", body)
}

func TestNFTWithdrawRecoverAsset(t *testing.T) {
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

	appid := int(2)
	uid := uint64(0)
	gameAssetName := "MonsterWeapon"
	nonce := ""
	appOrderID := ""

	err = RequestCommitWithdrawToGame(appid, uid, gameAssetName, nonce, appOrderID, const_def.NOTI_GAMESERVER_RECOVER)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("RequestCommitWithdrawToGame")
		return
	}

	logger.Logrus.WithFields(logrus.Fields{}).Info("RequestCommitWithdrawToGame success")
}

func TestNFTDeposit(t *testing.T) {
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

	appid := int(2)

	para := &GameNFTDepositData{
		GameAssetName: "MonsterWeapon",
		TokenID:       "2",
		EquipmentID:   "",
		TxHash:        "",
		Uid:           int64(0),
		Attrs:         []commdata.EquipmentAttr{},
	}

	baseUrl, err := comminfo.GetBaseUrl(appid)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("GetBaseUrl failed")
		return
	}

	deporder, err := RequestNFTDepositToGame(appid, para, baseUrl)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("RequestNFTDepositToGame failed")
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"DepositOrderData": deporder}).Info("RequestNFTDepositToGame success")
}
