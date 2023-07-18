package state

import (
	"context"
	"encoding/json"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/core/riskcontrol"
	"github/Connector-Gamefi/ConnectorGoSDK/core/sign"
	"github/Connector-Gamefi/ConnectorGoSDK/distmiddleware/skywalking"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"strings"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

//PrewithdrawHandler prewithdraw
func PrewithdrawHandler(order *model.TNftWithdrawRecord, chainAttrs []commdata.EquipmentAttr, gameAttrs []commdata.EquipmentAttr, imageURI string, chainID int, ctx context.Context) error {
	//construct game equipment data
	gameattrs := make([]byte, 0)
	var err error
	if gameAttrs != nil {
		gameattrs, err = json.Marshal(gameAttrs)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID}).Error("PrewithdrawHandler prewithdraw  game attrs failed")

			return err
		}
	}

	equipData := &model.TGameEquipment{
		AppID:           order.AppID,
		Account:         order.Account,
		ContractAddress: order.ContractAddress,
		EquipmentID:     order.EquipmentID,
		GameAssetName:   order.GameAssetName,
		ImageURI:        imageURI,
		EquipmentAttr:   datatypes.JSON(gameattrs),
		TokenID:         order.TokenID,
		ChainID:         chainID,
		Status:          const_def.SDK_EQUIPMENT_STATUS_DEFAULT,
		WithdrawSwitch:  const_def.CodeEquipmentWithdrawOpen,
	}

	//update or insert game equipment data
	err = UpdateOrInsertGameEquipment(equipData)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID, "GameEquipment": equipData}).Error("PrewithdrawHandler UpdateOrInsertGameEquipment failed")

		//update retry
		err = UpdateOrInsertGameEquipment(equipData)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID, "GameEquipment": equipData}).Error("PrewithdrawHandler UpdateOrInsertGameEquipment retry failed")
			return err
		}
	}

	logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID, "GameEquipment": equipData}).Info("PrewithdrawHandler UpdateOrInsertGameEquipment success")

	//hand nft withdraw order
	defer func(orderdata *model.TNftWithdrawRecord) {
		if orderdata != nil {
			nerr := UpdateNFTWithdrawByID(orderdata)
			if nerr != nil {
				logger.Logrus.WithFields(logrus.Fields{"ErrMsg": nerr.Error(), "OrderID": orderdata.ID, "OrderStatus": orderdata.OrderStatus}).Error("PrewithdrawHandler update withdraw order status failed")
			} else {
				logger.Logrus.WithFields(logrus.Fields{"OrderID": orderdata.ID, "OrderStatus": orderdata.OrderStatus}).Info("PrewithdrawHandler update withdraw order status success")
			}
		}
	}(order)

	//create risk control sub span
	span, subctx, err := skywalking.CreateLocalSpan(ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID}).Error("PrewithdrawHandler CreateLocalSpan failed")
		return err
	}
	span.SetOperationName("NFTWithdraw")
	ctx = subctx
	defer span.End()

	reportSpan, ok := span.(go2sky.ReportedSpan)
	if !ok {
		logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID}).Error("PrewithdrawHandler span type is warong")
		return fmt.Errorf("span type is wrong for risk")
	}

	traceID := reportSpan.Context().TraceID
	go2sky.PutCorrelation(ctx, "trace_id", traceID)
	go2sky.PutCorrelation(ctx, "type", fmt.Sprintf("%d", skywalking.NFTRiskAutoCode))
	seq := int(1)
	err = skywalking.SkyPutCorrelation(config.GetSkyWalkingConfig().KeyWithdraw, seq, config.GetSkyWalkingConfig().ValueWithdraw, ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID, "ErrMsg": err.Error()}).Error("PrewithdrawHandler SkyPutCorrelation failed")

		return err
	}

	//1. check withdraw risk control
	scode, err := riskcontrol.NFTRiskReview(order.AppID, order.Account, order.ContractAddress, order.EquipmentID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID}).Error("PrewithdrawHandler NFTRiskReview failed")
		order.OrderStatus = const_def.CodeWithdrawRiskFailed
		return err
	}

	order.RiskStatus = scode
	if scode == riskcontrol.CodeRiskWaiting {
		logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID}).Info("PrewithdrawHandler order join into risk waiting")

		//skywalking
		span.Log(time.Now(), "[NFTWithdraw]", fmt.Sprintf(" Seq: %d\n orderID: %d\n traceID: %s\n waiting", seq, order.ID, traceID))

		cacheID := fmt.Sprintf("nft:%d", order.ID)
		rerr := redis.SetString(cacheID, traceID, config.GetKeyNoExpireTime())
		if rerr != nil {
			logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID, "ErrMsg": rerr.Error()}).Error("PrewithdrawHandler cache trace id failed")

			return rerr
		}

		//wait for riskcontrol review success, save raw attributes table to db
		order.OrderStatus = const_def.CodeWithdrawRisking
		return nil
	}

	span.Log(time.Now(), "[NFTWithdraw]", fmt.Sprintf(" Seq: %d\n orderID: %d\n traceID: %s\n success", seq, order.ID, traceID))

	//risk control passed and construct sigmature
	err = gamePrewithdraw(order, chainAttrs, ctx, traceID, seq+1)
	if err != nil {
		return err
	}

	logger.Logrus.WithFields(logrus.Fields{"OrderInfo": order}).Info("PrewithdrawHandler success")

	return nil
}

//NftCommitWithdrawHandler query sign status and update db status
func NftCommitWithdrawHandler(data *model.TNftWithdrawRecord) error {
	signdata, err := sign.GetNftSignStatus(data.SignatureHash)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": data.ID, "SignHash": data.SignatureHash}).Error("NftCommitWithdrawHandler get sign status handle failed")

		data.Signature = ""
		return err
	}

	defer func(order *model.TNftWithdrawRecord) {
		err := UpdateNFTWithdrawOrderStatus(order)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": data.ID, "OrderStatus": data.OrderStatus}).Error("NftCommitWithdrawHandler update nft withdraw order status failed")
		} else {
			logger.Logrus.WithFields(logrus.Fields{"OrderID": data.ID, "OrderStatus": data.OrderStatus}).Info("NftCommitWithdrawHandler update nft withdraw order status success")
		}
	}(data)

	opStatus := int(0)

	//sign failed
	if signdata == "" {
		logger.Logrus.WithFields(logrus.Fields{"OrderID": data.ID}).Error("NftCommitWithdrawHandler get sign failed")

		//you can add data.OrderStatus = const_def.CodeWithdrawSignFailed when no need sign send retry
		data.OrderStatus = const_def.CodeWithdrawSignFailed
		opStatus = const_def.NOTI_GAMESERVER_RECOVER
	} else {
		//sign sucess
		data.Signature = signdata
		opStatus = const_def.NOTI_GAMESERVER_DELETE
	}

	//request game to delete asset
	err = ingame.RequestCommitWithdrawToGame(data.AppID, data.UID, data.GameAssetName, data.Nonce, data.AppOrderID, opStatus)
	if err != nil {
		data.OrderStatus = const_def.CodeWithdrawCommitFailed

		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": data.ID, "OpStatus": opStatus, "OrderStatus": data.OrderStatus}).Error("NftCommitWithdrawHandler request game to delete nft asset failed")
	} else {
		if opStatus == const_def.NOTI_GAMESERVER_DELETE {
			data.OrderStatus = const_def.CodeWithdrawWaitClaim
		} else {
			data.OrderStatus = const_def.CodeNotiRecorverSuccess
		}

		logger.Logrus.WithFields(logrus.Fields{"OrderID": data.ID, "OpStatus": opStatus, "OrderStatus": data.OrderStatus}).Info("NftCommitWithdrawHandler request game to delete or recover nft asset success")
	}

	if data.OrderStatus == const_def.CodeWithdrawWaitClaim {
		equipCon := map[string]interface{}{
			"app_id":           data.AppID,
			"contract_address": strings.ToLower(data.ContractAddress),
			"account":          data.Account,
			"equipment_id":     data.EquipmentID,
		}

		equipValue := map[string]interface{}{
			"status": const_def.SDK_EQUIPMENT_STATUS_WITHDRAW}

		err = mysql.GetDB().Model(&model.TGameEquipment{}).Where(equipCon).Updates(equipValue).Error
		if err != nil {

			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": data.ID}).Error("NftCommitWithdrawHandler update game equipment status failed")

			return fmt.Errorf("update equipment withdraw status %v", err)
		}
	}

	logger.Logrus.WithFields(logrus.Fields{"OrderID": data.ID}).Info("NftCommitWithdrawHandler prewithdraw success")

	return nil
}

//NftWithdrawClaimHandler claim order
func NftWithdrawClaimHandler(id uint64) (*model.TNftWithdrawRecord, commdata.NftSignatureSrcData, error) {
	var attrList commdata.NftSignatureSrcData

	obj, err := GetWithdrawOrderByPrimKey(id)
	if err != nil {
		return nil, attrList, err
	}

	if obj.SignatureSrc == nil {
		return obj, attrList, nil
	}

	err = json.Unmarshal([]byte(obj.SignatureSrc), &attrList)
	if err != nil {
		return nil, attrList, err
	}

	if obj.OrderStatus != const_def.CodeWithdrawWaitClaim {
		return nil, attrList, fmt.Errorf("primary id %v is waiting for sign success, please waiting", id)
	}

	return obj, attrList, nil
}

//RiskReviewHandle risk control review
func RiskReviewHandle(id, riskStatus int, reviewer string, ctx context.Context) error {
	order, err := GetWithdrawOrderByPrimKey(uint64(id))
	if err != nil {
		return err
	}

	order.RiskStatus = riskStatus
	order.RiskReviewer = reviewer

	//create risk control sub span
	span, subctx, err := skywalking.CreateLocalSpan(ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID}).Error("RiskReviewHandle CreateLocalSpan failed")
		return err
	}
	span.SetOperationName("NFTRiskReview")
	ctx = subctx
	defer span.End()

	//get trace id from cache
	go2sky.PutCorrelation(ctx, "type", fmt.Sprintf("%d", skywalking.NFTRiskReviewCode))

	traceID, err := redis.GetString(fmt.Sprintf("nft:%d", order.ID))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID}).Error("RiskReviewHandle get trace id from cache failed")
		return err
	}
	go2sky.PutCorrelation(ctx, "pre_trace_id", traceID)

	seq := int(2)
	err = skywalking.SkyPutCorrelation(config.GetSkyWalkingConfig().KeyRiskReview, seq, config.GetSkyWalkingConfig().ValueRiskReview, ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID, "ErrMsg": err.Error()}).Error("RiskReviewHandle SkyPutCorrelation failed")

		return err
	}

	defer func() {
		err := redis.DeleteString(fmt.Sprintf("nft:%d", order.ID))
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID}).Error("RiskReviewHandle delete trace id from cache failed")
		} else {
			logger.Logrus.WithFields(logrus.Fields{"TraceID": traceID, "OrderID": order.ID}).Info("RiskReviewHandle delete trace id from cache success")
		}
	}()

	if order.RiskStatus == riskcontrol.CodeRiskSuccess {
		if order.SignatureSrc == nil {
			logger.Logrus.WithFields(logrus.Fields{"Order": order}).Error("RiskReviewHandle source of signature is empty")
			return fmt.Errorf("source of signature is empty for %d", order.ID)
		}

		var attrList []commdata.EquipmentAttr
		err = json.Unmarshal([]byte(order.SignatureSrc), &attrList)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Order": order}).Error("RiskReviewHandle unmarshal attrs failed")
			return err
		}

		span.Log(time.Now(), "[NFTRiskReview]", fmt.Sprintf("OrderID: %d\n traceID:%s\n success", order.ID, traceID))

		err = gamePrewithdraw(order, attrList, ctx, traceID, seq+1)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Order": order}).Error("RiskReviewHandle game prewithdraw failed")
		}
	} else {
		span.Log(time.Now(), "[NFTRiskReview]", fmt.Sprintf("OrderID: %d\n traceID:%s\n failed", order.ID, traceID))
		order.OrderStatus = const_def.CodeWithdrawRiskFailed
	}

	err = UpdateNFTWithdrawOrderStatus(order)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Order": order}).Error("RiskReviewHandle update order status failed")

		return err
	}

	logger.Logrus.WithFields(logrus.Fields{"OrderInfo": order}).Info("RiskReviewHandle update order status success")

	return nil
}

func gamePrewithdraw(order *model.TNftWithdrawRecord, attrList []commdata.EquipmentAttr, ctx context.Context, traceID string, seq int) error {
	//create risk control sub span
	span, subctx, err := skywalking.CreateLocalSpan(ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID}).Error("GamePrewithdraw CreateLocalSpan failed")
		return err
	}
	span.SetOperationName("NFTGameWithdraw")
	ctx = subctx
	defer span.End()

	err = skywalking.SkyPutCorrelation(config.GetSkyWalkingConfig().KeyGameWithdraw, seq, config.GetSkyWalkingConfig().ValueGameWithdraw, ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID, "ErrMsg": err.Error()}).Error("GamePrewithdraw SkyPutCorrelation failed")

		return err
	}

	//game withdraw
	//get withdraw app order id
	apporder, err := ingame.RequestPreWithdrawToGame(order.AppID, order.UID, order.EquipmentID, order.GameAssetName, order.Nonce)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID}).Error("GamePrewithdraw request nft prewithdraw to game server failed")

		span.Log(time.Now(), "[GameWithdraw]", fmt.Sprintf(" Seq: %d\n OrderID: %d\n game prewithdraw failed", seq, order.ID))

		retryData := &commdata.PrewithdrawRetryData{
			ID:            order.ID,
			AppID:         order.AppID,
			UID:           order.UID,
			EquipmentID:   order.EquipmentID,
			GameAssetName: order.GameAssetName,
			Nonce:         order.Nonce,
			Attrs:         attrList,
			TraceID:       traceID,
			SeqNumber:     seq + 1,
		}

		jobid, nerr := NFTPublishPrewithdraw(retryData)
		if nerr != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": nerr.Error(), "OrderID": order.ID, "JobID": jobid}).Error("GamePrewithdraw add nft prewithdraw task queue failed")
		} else {
			logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID, "JobID": jobid}).Info("GamePrewithdraw add nft prewithdraw task queue success")
		}

		order.OrderStatus = const_def.CodeWithdrawPreWithdrawFailed
		return err
	}

	if apporder.AppOrderID != "" {
		order.AppOrderID = apporder.AppOrderID
	}

	span.Log(time.Now(), "[GameWithdraw]", fmt.Sprintf(" Seq: %d\n OrderID: %d\n GameOrderID: %s\n success", seq, order.ID, order.AppOrderID))

	//get signature from sign machine
	neworder, err := sign.NFTWithdrawSign(order, attrList, ctx)
	if err == nil {
		order.SignatureHash = neworder.SignatureHash
		order.SignatureSrc = neworder.SignatureSrc

		order.OrderStatus = const_def.CodeWithdrawSign

		logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID}).Info("GamePrewithdraw get signature hash and source data success")

		return nil
	}

	logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": order.ID}).Error("GamePrewithdraw get signature hash and source data failed")

	//auto recover asset
	gameerr := ingame.RequestCommitWithdrawToGame(order.AppID, order.UID, order.GameAssetName, order.Nonce, order.AppOrderID, const_def.NOTI_GAMESERVER_RECOVER)
	if gameerr != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": gameerr.Error(), "OrderID": order.ID}).Error("GamePrewithdraw auto recover asset failed")

		order.OrderStatus = const_def.CodeWithdrawSignFailed
	} else {
		logger.Logrus.WithFields(logrus.Fields{"OrderID": order.ID}).Info("GamePrewithdraw auto recover asset success")

		order.OrderStatus = const_def.CodeNotiRecorverSuccess
	}

	return err
}
