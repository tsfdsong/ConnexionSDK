package sdkcontroller

import (
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"github/Connector-Gamefi/ConnectorGoSDK/web/sdkbackend/sdkdata"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func DepositRepairOrder(c *gin.Context) {
	r := common.Response{
		Code: common.SuccessCode,
	}

	var p = sdkdata.PRepairOrder{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositRepairOrder validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		c.JSON(http.StatusOK, r)
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"input": p}).Info("DepositRepairOrder input info")

	var condition = map[string]interface{}{"id": p.ID}
	if p.ContractType == const_def.U8_CONTRACT_ERC20 {
		var erc20DepositRecord = model.TFtDepositRecord{}
		err, found := mysql.WrapFindFirst(model.TableFtDepositRecord, &erc20DepositRecord, condition)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositRepairOrder WrapFindFirst Error")
			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		if !found {
			r.Code = common.RecordNotExist
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		if erc20DepositRecord.OrderStatus != const_def.CodeDepositGameFailed {
			r.Code = common.CantRepairOrder
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}
		//repair it
		//TODO ft contract record must be lower
		ftContractCacheData, err := comminfo.GetFTContractByAddress(erc20DepositRecord.AppID, strings.ToLower(erc20DepositRecord.ContractAddress))
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositRepairOrder GetFTContractByAddress failed")

			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		if erc20DepositRecord.UID == 0 {
			r.Code = common.MustBindGameAccount
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		baseURL, err := comminfo.GetBaseUrlOriginError(erc20DepositRecord.AppID)
		if err == gorm.ErrRecordNotFound || baseURL == "" {
			r.Code = common.GameInfoNotExist
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositRepairOrder WrapFindFirst Error")
			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		balance, err := tools.GetTokenExactAmount(erc20DepositRecord.Amount, int32(ftContractCacheData.TokenDecimal))
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositRepairOrder GetTokenExactAmount")
			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		para := &ingame.NotifyFTDepositData{
			GameCoinName: ftContractCacheData.GameCoinName,
			Amount:       balance,
			TxHash:       erc20DepositRecord.TxHash,
			Uid:          int64(erc20DepositRecord.UID),
		}

		logger.Logrus.WithFields(logrus.Fields{"data": para}).Info("DepositRepairOrder RequestFTDepositToGame input")

		order, err := ingame.RequestFTDepositToGame(erc20DepositRecord.AppID, para, baseURL)
		if err != nil || order == nil {
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": erc20DepositRecord}).Error("DepositRepairOrder RequestFTDepositToGame failed")
			}

			r.Code = common.DepositNotiGameReject
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		} else {

			err := mysql.WrapUpdateByCondition(model.TableFtDepositRecord, condition, map[string]interface{}{"order_status": const_def.CodeDepositNotiSuccess, "app_order_id": order.AppOrderID})
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositRepairOrder WrapUpdateByCondition Error")
				r.Code = common.InnerError
				r.Message = common.ErrorMap[int(r.Code)]
				c.JSON(http.StatusOK, r)
				return
			}
		}
	} else if p.ContractType == const_def.U8_CONTRACT_ERC721 {
		var erc721DepositRecord = model.TNftDepositRecord{}
		err, found := mysql.WrapFindFirst(model.TableNftDepositRecord, &erc721DepositRecord, condition)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositRepairOrder WrapFindFirst nft failed")

			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		if !found {
			r.Code = common.RecordNotExist
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		if erc721DepositRecord.OrderStatus != const_def.CodeDepositGameFailed {
			r.Code = common.CantRepairOrder
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		gameEquipment := model.TGameEquipment{}
		err, found = mysql.WrapFindFirst(model.TableGameEquipment, &gameEquipment, map[string]interface{}{"app_id": erc721DepositRecord.AppID, "contract_address": strings.ToLower(erc721DepositRecord.ContractAddress), "token_id": erc721DepositRecord.TokenID})
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositRepairOrder WrapFindFirst equipment failed")
			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		passAssetMap := config.GetPassAssetConfig()
		_, isOk := passAssetMap[erc721DepositRecord.GameAssetName]
		if !isOk {
			if !found {
				r.Code = common.EquipmentNotExist
				r.Message = common.ErrorMap[int(r.Code)]
				c.JSON(http.StatusOK, r)
				return
			}
		}

		if erc721DepositRecord.UID == 0 {
			r.Code = common.MustBindGameAccount
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		baseURL, err := comminfo.GetBaseUrlOriginError(erc721DepositRecord.AppID)
		if err == gorm.ErrRecordNotFound || baseURL == "" {
			r.Code = common.GameInfoNotExist
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositRepairOrder GetBaseUrlOriginError failed")

			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		data := &ingame.GameNFTDepositData{
			GameAssetName: erc721DepositRecord.GameAssetName,
			TxHash:        erc721DepositRecord.TxHash,
			Uid:           int64(erc721DepositRecord.UID),
			TokenID:       erc721DepositRecord.TokenID,
			EquipmentID:   erc721DepositRecord.EquipmentID,
		}

		order, err := ingame.RequestNFTDepositToGame(erc721DepositRecord.AppID, data, baseURL)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": erc721DepositRecord}).Error("DepositRepairOrder RequestNFTDepositToGame failed")

			r.Code = common.DepositNotiGameReject
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		} else {
			tx := mysql.WrapBeginTx()
			result := tx.Table(model.TableNftDepositRecord).Where(condition).Updates(map[string]interface{}{"order_status": const_def.CodeDepositNotiSuccess, "app_order_id": order.AppOrderID})
			if result == nil || result.Error != nil {
				tx.Rollback()
				if result != nil && result.Error != nil {
					logger.Logrus.WithFields(logrus.Fields{"ErrMsg": result.Error.Error()}).Error("DepositRepairOrder update deposit record failed")
				}
				r.Code = common.InnerError
				r.Message = common.ErrorMap[int(r.Code)]
				c.JSON(http.StatusOK, r)
				return
			}

			gameResult := tx.Table(model.TableGameEquipment).Where(map[string]interface{}{"app_id": erc721DepositRecord.AppID, "contract_address": strings.ToLower(erc721DepositRecord.ContractAddress), "token_id": erc721DepositRecord.TokenID}).Updates(map[string]interface{}{"equipment_id": order.EquipmentID})
			if gameResult == nil || gameResult.Error != nil {
				tx.Rollback()
				if gameResult != nil && gameResult.Error != nil {
					logger.Logrus.WithFields(logrus.Fields{"ErrMsg": gameResult.Error.Error()}).Error("DepositRepairOrder update equipment record failed")
				}

				r.Code = common.InnerError
				r.Message = common.ErrorMap[int(r.Code)]
				c.JSON(http.StatusOK, r)
				return
			}
			err = tx.Commit().Error
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositRepairOrder tx commit failed")
				r.Code = common.InnerError
				r.Message = common.ErrorMap[int(r.Code)]
				c.JSON(http.StatusOK, r)
				return
			}
		}
	} else {
		logger.Logrus.WithFields(logrus.Fields{"ContractType": p.ContractType}).Info("DepositRepairOrder incorrect contract type")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		c.JSON(http.StatusOK, r)
		return
	}

	r.Message = common.ErrorMap[int(r.Code)]
	c.JSON(http.StatusOK, r)
}

func WithdrawRepairOrder(c *gin.Context) {
	r := common.Response{
		Code: common.SuccessCode,
	}

	var p = sdkdata.PRepairOrder{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("WithdrawRepairOrder validator reject")

		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		c.JSON(http.StatusOK, r)
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("Withdraw repair order info")

	var erc20WithdrawRecord = model.TFtWithdrawRecord{}
	var erc721WithdrawRecord = model.TNftWithdrawRecord{}
	var condition = map[string]interface{}{"id": p.ID}
	if p.ContractType == const_def.U8_CONTRACT_ERC20 {
		err, found := mysql.WrapFindFirst(model.TableFtWithdrawRecord, &erc20WithdrawRecord, condition)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("WrapFindFirst Error")
			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		if !found {
			r.Code = common.RecordNotExist
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		if erc20WithdrawRecord.OrderStatus != int(const_def.CodeWithdrawSignFailed) &&
			erc20WithdrawRecord.OrderStatus != int(const_def.CodeWithdrawCommitFailed) && erc20WithdrawRecord.OrderStatus != int(const_def.CodeInnerError) {
			r.Code = common.CantRepairOrder
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		//repair it
		//TODO ft contract record must be lower

		if erc20WithdrawRecord.UID == 0 {
			r.Code = common.MustBindGameAccount
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		baseURL, err := comminfo.GetBaseUrlOriginError(erc20WithdrawRecord.AppID)
		if err == gorm.ErrRecordNotFound || baseURL == "" {
			r.Code = common.GameInfoNotExist
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("GetBaseUrlOriginError Error")
			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		//noti gameserver recover assets
		secondResult, err := ingame.RequestFTRecoverToGame(erc20WithdrawRecord.AppID, erc20WithdrawRecord.UID, erc20WithdrawRecord.AppOrderID, erc20WithdrawRecord.GameCoinName, erc20WithdrawRecord.Nonce)

		logger.Logrus.WithFields(logrus.Fields{"Ret": secondResult, "p": p, "err": err}).Info("WithdrawRepairOrder")

		if err == nil && secondResult.IsSuccess() {
			code := const_def.CodeNotiRecorverSuccess
			err := mysql.WrapUpdateByCondition(model.TableFtWithdrawRecord, condition, map[string]interface{}{"order_status": code})
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("WrapUpdateByCondition Failed")
				r.Code = common.InnerError
				r.Message = common.ErrorMap[int(r.Code)]
				c.JSON(http.StatusOK, r)
				return
			}
		} else {
			r.Code = common.WithdrawRecorverNotiFailed
		}
	} else if p.ContractType == const_def.U8_CONTRACT_ERC721 {
		err, found := mysql.WrapFindFirst(model.TableNftWithdrawRecord, &erc721WithdrawRecord, condition)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("WrapFindFirst Error")
			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}
		if !found {
			r.Code = common.RecordNotExist
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		if erc721WithdrawRecord.UID == 0 {
			r.Code = common.MustBindGameAccount
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}
		if erc721WithdrawRecord.OrderStatus != int(const_def.CodeWithdrawSignFailed) &&
			erc721WithdrawRecord.OrderStatus != int(const_def.CodeWithdrawCommitFailed) {
			r.Code = common.CantRepairOrder
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		//repair it
		baseURL, err := comminfo.GetBaseUrlOriginError(erc721WithdrawRecord.AppID)
		if err == gorm.ErrRecordNotFound || baseURL == "" {
			r.Code = common.GameInfoNotExist
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("GetBaseUrlOriginError Error")
			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		err = ingame.RequestCommitWithdrawToGame(erc721WithdrawRecord.AppID, erc721WithdrawRecord.UID, erc721WithdrawRecord.GameAssetName, erc721WithdrawRecord.Nonce, erc721WithdrawRecord.AppOrderID, const_def.NOTI_GAMESERVER_RECOVER)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NFT withdraw Repair Error")
			r.Code = common.WithdrawRecorverNotiFailed
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}

		err = mysql.WrapUpdateByCondition(model.TableNftWithdrawRecord, condition, map[string]interface{}{"order_status": const_def.CodeNotiRecorverSuccess, "timestamp": time.Now().UnixMilli()})
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NFT withdraw Repair update order status Failed")
			r.Code = common.InnerError
			r.Message = common.ErrorMap[int(r.Code)]
			c.JSON(http.StatusOK, r)
			return
		}
	} else {
		logger.Logrus.Info("incorrect contract type")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		c.JSON(http.StatusOK, r)
		return
	}

	r.Message = common.ErrorMap[int(r.Code)]
	c.JSON(http.StatusOK, r)
}

func ParseLogSwitch(c *gin.Context) {
	r := common.Response{
		Code: common.SuccessCode,
	}

	var p = sdkdata.PParseLogSwitch{}
	err := c.ShouldBind(&p)
	logger.Logrus.WithFields(logrus.Fields{"p": p}).Info("ParseLogSwitch info")
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		c.JSON(http.StatusOK, r)
		return
	}

	if p.Status != const_def.SDK_TABLE_SWITCH_OPEN && p.Status != const_def.SDK_TABLE_SWITCH_CLOSE {
		logger.Logrus.Info("incorrect status")
		r.Code = common.IncorrectParams
		r.Message = "invalid switch status"
		c.JSON(http.StatusOK, r)
		return
	}
	height := model.TBlockHeight{}
	err = mysql.GetDB().Where("app_id = ?", p.GameID).First(&height).Error
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Info("get log switch from db failed")
		r.Code = common.RecordNotExist
		r.Message = "record not exist"
		c.JSON(http.StatusOK, r)
		return
	}

	condition := map[string]interface{}{"id": height.ID}
	udpateMap := map[string]interface{}{"switch": p.Status}

	err = mysql.WrapUpdateByCondition(model.TableBlockHeight, condition, udpateMap)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Info("WrapUpdateByCondition Error")
		r.Code = common.InnerError
		r.Message = "update log switch failed"
		c.JSON(http.StatusOK, r)
		return
	}

	err = comminfo.DeleteFilterLogSwitchByApp(p.GameID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Info("Delete filter log switch cache failed")
		r.Code = common.InnerError
		r.Message = "delete log switch cache failed"
		c.JSON(http.StatusOK, r)
		return
	}

	r.Message = common.ErrorMap[int(r.Code)]
	c.JSON(http.StatusOK, r)
}
