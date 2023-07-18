package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/core/sign"
	"github/Connector-Gamefi/ConnectorGoSDK/distmiddleware/skywalking"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/delayqueue"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"strings"
	"sync"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/sirupsen/logrus"
)

const (
	REDISNAME = "redis_delay_queue"

	NAMESPACE = "prewithdraw"

	NFTTASKQUEUE = "nft_queue"
	FTTASKQUEUE  = "ft_queue"

	BatchCount  = 5
	TTRSecond   = 5 //less than loop ticker interval --hades
	TTLSecond   = 2592000
	DelySecond  = 10
	RetryCounts = 3
)

var engine *delayqueue.Engine
var once sync.Once

func GetPrewithdrawEngine() *delayqueue.Engine {
	once.Do(func() {
		redisInst := redis.GetRedisInst()
		if redisInst == nil {
			panic(fmt.Errorf("redis instance is empty"))
		}
		inst, err := delayqueue.NewEngine(REDISNAME, redisInst)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Fatal("GetPrewithdrawEngine failed")
			panic(err)
		}

		if inst == nil {
			panic(fmt.Errorf("engine instance is empty"))
		}

		engine = inst
	})

	return engine
}

func TaskLoop() {
	//time.Sleep(time.Second * 3)
	err := GetPrewithdrawEngine().RegisterQueue(NAMESPACE, NFTTASKQUEUE)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("TaskLoop RegisterQueue failed")
		panic(err)
	}
	err = GetPrewithdrawEngine().RegisterQueue(NAMESPACE, FTTASKQUEUE)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("TaskLoop RegisterQueue failed")
		panic(err)
	}

	//NFT loop task
	go func() {
		t := time.NewTicker(10 * time.Second)
		defer t.Stop()

		defer func() {
			err := recover()
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("NFT prewithdraw task failed")
			}
		}()

		for range t.C {
			jobs, err := GetPrewithdrawEngine().BatchConsume(NAMESPACE, []string{NFTTASKQUEUE}, BatchCount, TTRSecond, 0)
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("TaskLoop BatchConsume failed")
				continue
			}

			for _, job := range jobs {
				jobid := job.ID()
				err := NFTConsumePrewithdraw(job.Body(), jobid)
				if err != nil {
					continue
				} else {
					jerr := GetPrewithdrawEngine().Delete(NAMESPACE, NFTTASKQUEUE, jobid)
					if jerr != nil {
						logger.Logrus.WithFields(logrus.Fields{"ErrMsg": jerr.Error(), "Job": job}).Error("NFTConsumePrewithdraw delete task failed")
					} else {
						logger.Logrus.WithFields(logrus.Fields{"Job": job}).Info("NFTConsumePrewithdraw delete task success")
					}
				}
			}
		}

	}()

	//hand nft dead queue
	go NFTDeadPrewithdrawTask()

	//ft loop task
	go FTPrewithdrawTask()
	go FTDeadPrewithdrawTask()
}

func NFTConsumePrewithdraw(job []byte, jobid string) error {
	var data commdata.PrewithdrawRetryData
	err := json.Unmarshal(job, &data)
	if err != nil {
		return fmt.Errorf("unmarshal nft data failed, %v", err)
	}

	logger.Logrus.WithFields(logrus.Fields{"Order": data, "JobID": jobid}).Info("NFTConsumePrewithdraw consume nft withdraw order task queue info")

	order := &model.TNftWithdrawRecord{
		ID:            data.ID,
		AppID:         data.AppID,
		UID:           data.UID,
		EquipmentID:   data.EquipmentID,
		GameAssetName: data.GameAssetName,
		Nonce:         data.Nonce,
	}

	//create risk control sub span
	ctx := context.Background()
	span, subctx, err := skywalking.CreateLocalSpan(ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID}).Error("NFTConsumePrewithdraw CreateLocalSpan failed")
		return err
	}
	span.SetOperationName("NFTGamePrewithdraw")
	ctx = subctx
	defer span.End()

	go2sky.PutCorrelation(ctx, "type", fmt.Sprintf("%d", skywalking.NFTGameRetryCode))
	go2sky.PutCorrelation(ctx, "pre_trace_id", data.TraceID)

	err = skywalking.SkyPutCorrelation(config.GetSkyWalkingConfig().KeyGameWithdraw, data.SeqNumber, config.GetSkyWalkingConfig().ValueGameWithdraw, ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID, "ErrMsg": err.Error()}).Error("GamePrewithdraw SkyPutCorrelation failed")

		return err
	}

	gameorder, err := ingame.RequestPreWithdrawToGame(order.AppID, order.UID, order.EquipmentID, order.GameAssetName, order.Nonce)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID, "JobID": jobid}).Error("NFTConsumePrewithdraw request prewithdraw to game server failed")

		return fmt.Errorf("RequestPreWithdrawToGame failed, %v", err)
	}

	defer func() {
		con := map[string]interface{}{"id": order.ID}

		values := map[string]interface{}{"app_order_id": order.AppOrderID, "order_status": order.OrderStatus, "signature_hash": order.SignatureHash, "signature_source": order.SignatureSrc, "timestamp": time.Now().UnixMilli()}
		err = mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(con).Updates(values).Error
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderInfo": order, "JobID": jobid}).Error("NFTConsumePrewithdraw defer update app order id and order status failed")
		} else {
			logger.Logrus.WithFields(logrus.Fields{"OrderInfo": order, "JobID": jobid}).Info("NFTConsumePrewithdraw defer update app order id and order status success")
		}
	}()

	order.AppOrderID = gameorder.AppOrderID
	order.EquipmentID = gameorder.EquipmentID

	neworder, err := sign.NFTWithdrawSign(order, data.Attrs, ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Order": order, "JobID": jobid}).Error("NFTConsumePrewithdraw NFTWithdrawSign failed")

		//auto recover asset
		gameerr := ingame.RequestCommitWithdrawToGame(order.AppID, order.UID, order.GameAssetName, order.Nonce, order.AppOrderID, const_def.NOTI_GAMESERVER_RECOVER)
		if gameerr != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": gameerr.Error(), "OrderID": order.ID}).Error("NFTConsumePrewithdraw auto recover asset failed")

			order.OrderStatus = const_def.CodeWithdrawSignFailed
		} else {
			logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID}).Info("NFTConsumePrewithdraw auto recover asset success")

			order.OrderStatus = const_def.CodeNotiRecorverSuccess
		}

		return err
	}

	order.SignatureHash = neworder.SignatureHash
	order.SignatureSrc = neworder.SignatureSrc

	order.OrderStatus = const_def.CodeWithdrawSign

	return nil
}

func NFTPublishPrewithdraw(data *commdata.PrewithdrawRetryData) (string, error) {
	if data == nil {
		return "", fmt.Errorf("input data is empty")
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshal input data failed, %v", err)
	}

	ttl := TTLSecond
	delaytime := DelySecond
	jobid, err := GetPrewithdrawEngine().Publish(NAMESPACE, NFTTASKQUEUE, body, uint32(ttl), uint32(delaytime), uint16(RetryCounts))
	if err != nil {
		return "", fmt.Errorf("publish job failed, %v", err)
	}

	return jobid, nil
}

func FTPrewithdrawTask() {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("FTPrewithdrawTask failed")
		}
	}()

	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for range t.C {
		jobs, err := GetPrewithdrawEngine().BatchConsume(NAMESPACE, []string{FTTASKQUEUE}, BatchCount, TTRSecond, 0)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("FTPrewithdrawTask BatchConsume failed")

			continue
		}

		for _, job := range jobs {
			jobid := job.ID()
			err := FTConsumePrewithdraw(job.Body())
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "JobID": jobid}).Error("FTPrewithdrawTask FTConsumePrewithdraw failed")
				continue
			} else {
				err := GetPrewithdrawEngine().Delete(NAMESPACE, FTTASKQUEUE, jobid)
				if err != nil {
					logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Job": job}).Error("FTPrewithdrawTask delete task failed")
				} else {
					logger.Logrus.WithFields(logrus.Fields{"Job": job}).Info("FTPrewithdrawTask delete task success")
				}
			}
		}
	}
}

func FTConsumePrewithdraw(job []byte) error {
	data := commdata.FTPrewithdrawRetryData{}
	err := json.Unmarshal(job, &data)
	if err != nil {
		return fmt.Errorf("unmarshal nft data failed, %v", err)
	}

	logger.Logrus.WithFields(logrus.Fields{"data": data}).Info("FTConsumePrewithdraw:data")

	//create risk control sub span
	ctx := context.Background()
	span, subctx, err := skywalking.CreateLocalSpan(ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderNonce": data.Nonce}).Error("FTConsumePrewithdraw CreateLocalSpan failed")
		return err
	}
	span.SetOperationName("FTGamePrewithdraw")
	ctx = subctx
	defer span.End()

	go2sky.PutCorrelation(ctx, "type", fmt.Sprintf("%d", skywalking.FTGameRetryCode))
	go2sky.PutCorrelation(ctx, "pre_trace_id", data.TraceID)

	err = skywalking.SkyPutCorrelation(config.GetSkyWalkingConfig().KeyGameWithdraw, data.SeqNumber, config.GetSkyWalkingConfig().ValueGameWithdraw, ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"OrderNonce": data.Nonce, "ErrMsg": err.Error()}).Error("FTConsumePrewithdraw SkyPutCorrelation failed")

		return err
	}

	//send request. if http200 success/failed then flow go on / else return

	//step1 get base infos

	baseURL, err := comminfo.GetBaseUrl(int(data.AppID))
	if baseURL == "" || err != nil {
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "data": data}).Error("FTConsumePrewithdraw::GetBaseUrl failed")
		}

		return fmt.Errorf("get ft contract baseurl failed:: url:%+v,err:%+v", baseURL, err)
	}

	contractInfo, err := comminfo.GetFTContractByAddress(int(data.AppID), strings.ToLower(data.ContractAddress))
	if err != nil || contractInfo == nil {
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "data": data}).Error("FTConsumePrewithdraw::GetFTContractByAddress failed")
		}

		return fmt.Errorf("get ft  contract data failed")
	}

	if contractInfo.Treasure == "" {
		logger.Logrus.WithFields(logrus.Fields{"data": data}).Error("FTConsumePrewithdraw:: treasure is emtpy")
		return fmt.Errorf("get ft  contract data failed")
	}

	balance, err := tools.GetTokenExactAmount(data.Amount, int32(contractInfo.TokenDecimal))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "data": data}).Error("FTConsumePrewithdraw::GetTokenExactAmount failed")
		return fmt.Errorf("get ft  contract data failed")
	}
	condition := map[string]interface{}{"app_id": data.AppID, "contract_address": strings.ToLower(data.ContractAddress), "nonce": data.Nonce}
	ftWithdrawData := model.TFtWithdrawRecord{}
	err, found := mysql.WrapFindFirst(model.TableFtWithdrawRecord, &ftWithdrawData, condition)
	if err != nil || !found {
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "found": found}).Error("FTConsumePrewithdraw::not found this order or find error")
		}

		return fmt.Errorf("get ft  contract data failed")
	}

	//step2 retry
	rsp, err := ingame.RequestFTPreWithdrawToGame(int(data.AppID), uint64(data.Uid), balance, data.AppCoinName, data.Nonce)

	logger.Logrus.WithFields(logrus.Fields{"err": err, "result": rsp, "data": data}).Info("FTConsumePrewithdraw result")
	//if not 200. return

	if err != nil {

		logger.Logrus.Error("FTConsumePrewithdraw::FTPreWithdrawHandler failed")
		return fmt.Errorf("FTConsumePrewithdraw::FTPreWithdrawHandler failed")
	}

	if rsp.IsFailure() {
		//invalid order.

		//update CodeWithdrawPreWithdrawFailed
		erc20UpdateMap := map[string]interface{}{"order_status": const_def.CodeWithdrawPreWithdrawFailed}
		err := mysql.WrapUpdateByCondition(model.TableFtWithdrawRecord, condition, erc20UpdateMap)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "data": data}).Error(const_def.NEED_MANUAL_REPAIR_THIS + "FTConsumePrewithdraw::WrapUpdateByCondition failed")
		}

		return nil
	}

	//step3 flow go on(if need repair or sign success need save db && always return nil)

	// sign
	//1 --send sign failed.
	//2 --send sign success but code is failure
	//3 --send sign success and code is success and save db
	logger.Logrus.Info(fmt.Sprintf("FTConsumePrewithdraw::begin sign:%+v", data))
	var dbError error
	signResult, signErr := FTSignHandler(data.Amount, data.UserAddress, data.ContractAddress, contractInfo.Treasure, data.Nonce, data.ChainID, ctx)

	logger.Logrus.WithFields(logrus.Fields{"signErr": signErr, "result": signResult, "data": data}).Info("FTConsumePrewithdraw::sign result")

	if signErr == nil {

		erc20UpdateMap := map[string]interface{}{"order_status": const_def.CodeWithdrawSign, "signature_hash": signResult.ReqHash, "app_order_id": rsp.Data.AppOrderID}
		err := mysql.WrapUpdateByCondition(model.TableFtWithdrawRecord, condition, erc20UpdateMap)
		if err != nil {
			dbError = err
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "data": data}).Error("FTConsumePrewithdraw::WrapUpdateByCondition failed")
		} else {
			return nil
		}

	}

	//recorver  asset
	logger.Logrus.Info(fmt.Sprintf("FTConsumePrewithdraw::begin second gameserver commit:%+v", data))
	recorverUpdateMap := map[string]interface{}{}

	if dbError != nil {
		recorverUpdateMap["order_status"] = const_def.CodeInnerError
	} else {
		recorverUpdateMap["order_status"] = const_def.CodeWithdrawSignFailed
	}

	//noti gameserver recover assets
	secondResult, err := ingame.RequestFTRecoverToGame(int(data.AppID), uint64(data.Uid), rsp.Data.AppOrderID, data.AppCoinName, data.Nonce)

	logger.Logrus.WithFields(logrus.Fields{"err": err, "data": data, "secondresult": secondResult}).Info("RequestFTRecoverToGame result")

	if err != nil || secondResult.IsFailure() || secondResult.IsReject() {
		recorverUpdateMap["order_status"] = const_def.CodeWithdrawCommitFailed
		//need repair
		logger.Logrus.WithFields(logrus.Fields{"p": ftWithdrawData}).Error("recorver failed")
	} else if secondResult.IsSuccess() {
		//not need save db
		recorverUpdateMap["order_status"] = const_def.CodeNotiRecorverSuccess
		logger.Logrus.WithFields(logrus.Fields{"p": ftWithdrawData}).Info("recorver success")
	} else {
		logger.Logrus.WithFields(logrus.Fields{"secondResult": secondResult}).Error("RequestFTRecoverToGame Unknow Status")

	}

	err = mysql.WrapUpdateByCondition(model.TableFtWithdrawRecord, condition, recorverUpdateMap)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "data": data}).Error(const_def.NEED_MANUAL_REPAIR_THIS + "FTConsumePrewithdraw::WrapUpdateByCondition failed")
	}

	return nil

}

func FTPublishPrewithdraw(data *commdata.FTPrewithdrawRetryData) (string, error) {
	if data == nil {
		return "", fmt.Errorf("input ft data is empty")
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshal ft input data failed, %v", err)
	}

	ttl := TTLSecond
	delaytime := DelySecond
	jobid, err := GetPrewithdrawEngine().Publish(NAMESPACE, FTTASKQUEUE, body, uint32(ttl), uint32(delaytime), uint16(RetryCounts))
	if err != nil {
		return "", fmt.Errorf("publish ft job failed, %v", err)
	}

	return jobid, nil
}

func FTDeadPrewithdrawTask() {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("FTDeadPrewithdrawTask failed")
		}
	}()

	//time.Sleep(time.Second * 4)
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for range t.C {
		size, err := GetPrewithdrawEngine().SizeOfDeadLetter(NAMESPACE, FTTASKQUEUE)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("FTDeadPrewithdrawTask::SizeOfDeadLetter failed")
			continue
		}

		if size == 0 {
			continue
		}

		logger.Logrus.WithFields(logrus.Fields{"size": size}).Info("FTDeadPrewithdrawTask Dead queue Size")

		count := size
		if size > 5 {
			count = 5
		}

		for count > 0 {
			//can't skip so go on. and max retry 5 times
			HandleDeadFTPrewithdrawTask()
			count--
		}
	}
}

func HandleDeadFTPrewithdrawTask() error {
	_, headJobID, err := GetPrewithdrawEngine().PeekDeadLetter(NAMESPACE, FTTASKQUEUE)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("HandleDeadFTPrewithdrawTask::PeekDeadLetter failed")

		return err
	}

	if headJobID == "" {
		return nil
	}

	job, err := GetPrewithdrawEngine().Peek(NAMESPACE, FTTASKQUEUE, headJobID)
	if err != nil || job == nil {
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("HandleDeadFTPrewithdrawTask::Peek failed")
		}

		return err
	}
	var data commdata.FTPrewithdrawRetryData
	err = json.Unmarshal(job.Body(), &data)
	if err != nil {
		return err
	}

	logger.Logrus.WithFields(logrus.Fields{"data": data}).Info("HandleDeadFTPrewithdrawTask:data")

	//update status
	condition := map[string]interface{}{"app_id": data.AppID, "contract_address": strings.ToLower(data.ContractAddress), "nonce": data.Nonce}
	recorverUpdateMap := map[string]interface{}{"order_status": const_def.CodeWithdrawPreWithdrawFailed}
	err = mysql.WrapUpdateByCondition(model.TableFtWithdrawRecord, condition, recorverUpdateMap)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "data": data}).Error("HandleDeadFTPrewithdrawTask::WrapUpdateByCondition failed")
		return err
	}

	//delete is last
	delCount, err := GetPrewithdrawEngine().DeleteDeadLetter(NAMESPACE, FTTASKQUEUE, 1)
	if delCount != 1 || err != nil {
		if err != delayqueue.ErrNotFound && err != delayqueue.ErrEmptyQueue && err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "data": data}).Error("HandleDeadFTPrewithdrawTask::DeleteDeadLetter failed")

		}
		return errors.New("FTDeadPrewithdrawTask::DeleteDeadLetter failed")
	}

	return nil
}

func NFTDeadPrewithdrawTask() {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("NFTDeadPrewithdrawTask failed")
		}
	}()

	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for range t.C {
		size, err := GetPrewithdrawEngine().SizeOfDeadLetter(NAMESPACE, NFTTASKQUEUE)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NFTDeadPrewithdrawTask::SizeOfDeadLetter failed")
			continue
		}

		if size == 0 {
			continue
		}

		logger.Logrus.WithFields(logrus.Fields{"size": size}).Info("NFTDeadPrewithdrawTask Dead queue Size")

		count := size
		if size > 5 {
			count = 5
		}

		for count > 0 {
			//can't skip so go on. and max retry 5 times
			HandleDeadNFTPrewithdrawTask()
			count--
		}
	}
}

func HandleDeadNFTPrewithdrawTask() error {
	_, headJobID, err := GetPrewithdrawEngine().PeekDeadLetter(NAMESPACE, NFTTASKQUEUE)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("HandleDeadNFTPrewithdrawTask::PeekDeadLetter failed")

		return err
	}

	if headJobID == "" {
		return nil
	}

	job, err := GetPrewithdrawEngine().Peek(NAMESPACE, NFTTASKQUEUE, headJobID)
	if err != nil || job == nil {
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "JobID": headJobID}).Error("HandleDeadNFTPrewithdrawTask::Peek failed")
		}

		return err
	}
	var data commdata.PrewithdrawRetryData
	err = json.Unmarshal(job.Body(), &data)
	if err != nil {
		return err
	}

	logger.Logrus.WithFields(logrus.Fields{"data": data, "JobID": headJobID}).Info("HandleDeadNFTPrewithdrawTask:data")

	//update status
	con := map[string]interface{}{"id": data.ID}
	values := map[string]interface{}{"order_status": const_def.CodeWithdrawPreWithdrawFailed, "timestamp": time.Now().UnixMilli()}

	err = mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(con).Updates(values).Error
	if err != nil {

		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "ID": data.ID, "JobID": headJobID}).Error("HandleDeadFTPrewithdrawTask update order status failed")
		return err
	}

	//delete is last
	delCount, err := GetPrewithdrawEngine().DeleteDeadLetter(NAMESPACE, NFTTASKQUEUE, 1)
	if delCount != 1 || err != nil {
		if err != delayqueue.ErrNotFound && err != delayqueue.ErrEmptyQueue && err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "data": data, "JobID": headJobID}).Error("HandleDeadNFTPrewithdrawTask::DeleteDeadLetter failed")

		}

		return errors.New("NFTDeadPrewithdrawTask::DeleteDeadLetter failed")
	}

	return nil
}
