package common

import (
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"time"

	"github.com/sirupsen/logrus"
)

func FTGameDeposit(v *model.TFtDepositRecord, ftInfo []commdata.FTContractCacheData) error {
	info, err := comminfo.GetFtContractFromCacheData(v.AppID, v.ContractAddress, ftInfo)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "FTDepositOrder": v}).Error("FTGameDeposit GetFTContractByAddress failed")
		return err
	}

	balance, err := tools.GetTokenExactAmount(v.Amount, int32(info.TokenDecimal))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "FTDepositOrder": v, "decimal": info.TokenDecimal}).Error("FTGameDeposit Convert Balance failed")
		return err
	}

	data := &ingame.NotifyFTDepositData{
		GameCoinName: v.GameCoinName,
		Amount:       balance,
		TxHash:       v.TxHash,
		Uid:          int64(v.UID),
	}

	baseurl, err := comminfo.GetBaseUrl(v.AppID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "FTDepositOrder": v}).Error("FTGameDeposit GetBaseUrl failed")
		return err
	}

	order, err := ingame.RequestFTDepositToGame(v.AppID, data, baseurl)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "FTDepositOrder": order}).Error("FTGameDeposit RequestFTDepositToGame failed")
	}

	v.OrderStatus = order.OrderStatus
	v.AppOrderID = order.AppOrderID

	return nil
}

func NFTGameDeposit(v *model.TNftDepositRecord) error {
	if v.UID == 0 || v.Account == "" {
		bindinfo, err := comminfo.GetBindInfo(v.AppID, v.DepositAddress)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "NFTDepositData": v}).Error("NFTGameDeposit GetBindInfo failed")
			return err
		}

		v.UID = bindinfo.UID
		v.Account = bindinfo.Account
	}

	data := &ingame.GameNFTDepositData{
		GameAssetName: v.GameAssetName,
		TxHash:        v.TxHash,
		Uid:           int64(v.UID),
		TokenID:       v.TokenID,
		EquipmentID:   v.EquipmentID,
	}

	baseurl, err := comminfo.GetBaseUrl(v.AppID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "NFTDepositData": v}).Error("NFTGameDeposit GetBaseUrl failed")
		return err
	}

	order, err := ingame.RequestNFTDepositToGame(v.AppID, data, baseurl)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "NFTDepositData": order}).Error("NFTGameDeposit RequestNFTDepositToGame failed")
	}
	v.OrderStatus = order.OrderStatus
	v.AppOrderID = order.AppOrderID
	v.EquipmentID = order.EquipmentID

	return nil
}

func NFTGameWithdraw(rs []*model.TNftWithdrawRecord) {
	for _, v := range rs {
		con := map[string]interface{}{
			"app_id":  v.AppID,
			"tx_hash": v.TxHash}
		value := map[string]interface{}{
			"order_status": v.OrderStatus,
			"timestamp":    time.Now().UnixMilli()}

		err := mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(con).Updates(value).Error
		if err != nil {

			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": v}).Error("NFTGameWithdraw update nft withdraw order status failed")
			continue
		}

		logger.Logrus.WithFields(logrus.Fields{"Data": v}).Info("NFTGameWithdraw success")
	}
}
