package state

import (
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/core/sign"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetPrintStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

func Loop() {
	go ftLoop()
	go nftLoop()

	TaskLoop()
}

func SingleSignRequest(param sign.PQuerySignature, ftWithdrawRecord model.TFtWithdrawRecord, wg *sync.WaitGroup) {
	defer wg.Done()

	defer func() {
		err := recover()
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("SingleSignRequest failed")
		}
	}()

	condition := map[string]interface{}{"app_order_id": ftWithdrawRecord.AppOrderID, "nonce": ftWithdrawRecord.Nonce, "app_id": ftWithdrawRecord.AppID, "contract_address": strings.ToLower(ftWithdrawRecord.ContractAddress)}

	// gameInfocondition := map[string]interface{}{"app_id": gameId}
	baseURL, err := comminfo.GetBaseUrlOriginError(ftWithdrawRecord.AppID)
	if err == gorm.ErrRecordNotFound || baseURL == "" {
		logger.Logrus.WithFields(logrus.Fields{"condition": condition}).Error("GameInfo Not Found")
		return
	}

	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"condition": condition, "ErrMsg": err.Error()}).Error("GetBaseUrlOriginError Error")
		return
	}

	//sign.FtQueryDepositSign(param.ReqHash)
	r, err := FTWithdrawSignResultRequest(param)
	logger.Logrus.WithFields(logrus.Fields{"signResult": r, "err": err, "p": ftWithdrawRecord}).Info("SingleSignRequest")
	if err != nil || r.Sign == const_def.SIGN_PENDING {

		return
	}

	updateMap := map[string]interface{}{}
	status := const_def.NOTI_GAMESERVER_DELETE
	if r.Code == const_def.SIGN_FAILED || r.Sign == "" {
		//sign reject.recover
		status = const_def.NOTI_GAMESERVER_RECOVER
		updateMap["order_status"] = const_def.CodeWithdrawSignFailed
	} else {
		s, err := tools.EthSignFix(r.Sign)
		if err != nil {
			status = const_def.NOTI_GAMESERVER_RECOVER
			updateMap["order_status"] = const_def.CodeWithdrawSignFailed
		} else {
			updateMap["signature"] = s
		}
	}

	secondResult, err := ingame.RequestFTOperationToGame(ftWithdrawRecord.AppID, ftWithdrawRecord.UID, ftWithdrawRecord.AppOrderID, ftWithdrawRecord.GameCoinName, ftWithdrawRecord.Nonce, status)

	logger.Logrus.WithFields(logrus.Fields{"secondResult": secondResult, "err": err}).Info("FTCommitWithdrawHandler")

	if err != nil || secondResult.IsFailure() {
		updateMap["order_status"] = const_def.CodeWithdrawCommitFailed
		logger.Logrus.WithFields(logrus.Fields{"p": ftWithdrawRecord, "err": err.Error()}).Error("commit failed")
	} else if secondResult.IsSuccess() {
		//be care. sign has value. but it does not  matter
		if status == const_def.NOTI_GAMESERVER_DELETE {
			updateMap["order_status"] = const_def.CodeWithdrawWaitClaim
			logger.Logrus.WithFields(logrus.Fields{"p": ftWithdrawRecord}).Info("noti game server delete asset success")
		} else if status == const_def.NOTI_GAMESERVER_RECOVER {
			updateMap["order_status"] = const_def.CodeNotiRecorverSuccess
			logger.Logrus.WithFields(logrus.Fields{"p": ftWithdrawRecord}).Info("noti game server recorver asset success")
		}
	} else if secondResult.IsReject() {
		logger.Logrus.WithFields(logrus.Fields{"p": ftWithdrawRecord}).Error("commit reject")
		updateMap["order_status"] = const_def.CodeWithdrawCommitFailed
	} else {
		logger.Logrus.WithFields(logrus.Fields{"secondresult": secondResult}).Error("RequestFTOperationToGame Unknow Status")
		return
	}

	err = mysql.WrapUpdateByCondition(model.TableFtWithdrawRecord, condition, updateMap)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"condition": condition, "updateMap": updateMap, "error": err}).Error(const_def.NEED_MANUAL_REPAIR_THIS + "WrapUpdateByCondition failed")
	}
}

//case1 withdraw  sign result check ->result to recover or delete
//case2 withdraw  recover noti
//case3 withdraw  delete noti

func SingleFTLoop() {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "Stack": GetPrintStack()}).Fatalf("SingleFTLoop panic")
		}
	}()

	signList, err := FTFindAllSendSignRecords()
	if err != nil {
		logger.Logrus.Error("FTFindAllSendSignRecords Error")
		return
	}
	if len(signList) <= 0 {
		return
	}

	//coroutine
	wg := new(sync.WaitGroup)
	wg.Add(len(signList))
	for _, e := range signList {
		ftWithdraw := e
		param := sign.PQuerySignature{
			ReqHash: ftWithdraw.SignatureHash,
		}

		if ftWithdraw.UID > 0 {
			go SingleSignRequest(param, ftWithdraw, wg)
		}

	}
	wg.Wait()
}
func ftLoop() {
	tick := time.Tick(time.Duration(10) * time.Second)
	for range tick {
		SingleFTLoop()
	}
}

func nftLoop() {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "Stack": GetPrintStack()}).Fatalf("nftLoop panic")
		}
	}()

	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for range t.C {
		datalist, err := FindAllNFTRecordSign()
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("NFTLoop FindAllNFTRecordSign failed")

			continue
		}

		for _, data := range datalist {
			NftCommitWithdrawHandler(&data)
		}
	}
}
