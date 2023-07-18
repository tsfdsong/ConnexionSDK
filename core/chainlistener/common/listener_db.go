package common

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	BatchCount int = 10
)

func InsertHeight(appId int, height uint64) error {
	data := model.TBlockHeight{
		LatestParsedHeight: height,
		AppID:              appId,
	}
	return mysql.GetDB().Model(&model.TBlockHeight{}).Create(&data).Error
}

func UpdateHeight(appId int, height uint64) error {
	blockHeight := model.TBlockHeight{}
	err := mysql.GetDB().Where("app_id = ?", appId).First(&blockHeight).Error
	if err != nil {
		return fmt.Errorf("get height %v", err)
	}

	heigthCondition := map[string]interface{}{"id": blockHeight.ID}
	heightUpdateMap := map[string]interface{}{"latest_parsed_height": height}
	return mysql.GetDB().Model(&model.TBlockHeight{}).Where(heigthCondition).Updates(heightUpdateMap).Error
}

func GetHeight(appId int) (uint64, error) {
	height := model.TBlockHeight{}
	err := mysql.GetDB().Where("app_id = ?", appId).First(&height).Error
	if err != nil {
		return 0, err
	}

	return height.LatestParsedHeight, nil
}

// if program down. start height from db. not async
func InitLatestHeight(appId int, initHeight int64) (uint64, error) {
	height, err := GetHeight(appId)
	if err != nil {
		err := InsertHeight(appId, uint64(initHeight))
		if err != nil {
			return 0, err
		}

		return uint64(initHeight), nil
	}
	return height, nil
}

// CheckSwitch checkout switch of parse log
func CheckSwitch(appId int) (bool, error) {
	lswitch, err := comminfo.GetFilterLogSwitchByApp(appId)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error(fmt.Sprintf("APPID:%d get filter log switch failed", appId))
		return false, err
	}

	dbSwitch := lswitch == const_def.SDK_TABLE_SWITCH_OPEN

	return dbSwitch, nil
}

func InsertTxHash(txs []model.TTxHash) error {
	err := mysql.GetDB().Model(&model.TTxHash{}).CreateInBatches(txs, len(txs)).Error
	return err
}

func InsertFTDeposit(v *model.TFtDepositRecord) error {
	err, found := mysql.WrapFindFirst(model.TableFtDepositRecord, &model.TFtDepositRecord{}, map[string]interface{}{"contract_address": v.ContractAddress, "nonce": v.Nonce, "chain_id": v.ChainID})
	if err != nil {
		return err
	}
	if found {
		return nil
	}
	ftContractConf, err := comminfo.GetFTContractInfoByDepositTreasureAndChain(v.ContractAddress, v.TargetAddress, v.ChainID)
	if err != nil {
		return err
	}

	bindInfo, err := comminfo.GetEmailBindInfo(uint(ftContractConf.AppID), uint(ftContractConf.ChainID), v.DepositAddress)
	if err != nil {
		return err
	}

	v.AppID = ftContractConf.AppID
	v.GameCoinName = ftContractConf.GameCoinName
	v.UID = bindInfo.UID
	v.Account = bindInfo.Account
	v.Timestamp = time.Now().UnixMilli()

	err = mysql.WrapInsertSingle(model.TableFtDepositRecord, v)
	if err != nil {
		return err
	}
	return nil
}

func UpdateFTWithdraw(rs []*model.TFtWithdrawRecord) error {
	for _, v := range rs {
		con := map[string]interface{}{
			"chain_id":         v.ChainID,
			"contract_address": strings.ToLower(v.ContractAddress),
			"withdraw_address": strings.ToLower(v.WithdrawAddress),
			"amount":           v.Amount,
			"nonce":            v.Nonce}
		value := map[string]interface{}{
			"tx_hash":      v.TxHash,
			"height":       v.Height,
			"order_status": v.OrderStatus,
			"timestamp":    time.Now().UnixMilli()}

		err := mysql.GetDB().Model(&model.TFtWithdrawRecord{}).Where(con).Updates(value).Error
		if err != nil {

			return err
		}
	}

	return nil
}

func UpdateNFTDeposit(rs []*model.TNftDepositRecord) error {
	for _, v := range rs {
		con := map[string]interface{}{
			"contract_address": strings.ToLower(v.ContractAddress),
			"deposit_address":  strings.ToLower(v.DepositAddress),
			"trease_address":   strings.ToLower(v.TargetAddress),
			"token_id":         v.TokenID,
			"nonce":            v.Nonce}
		value := map[string]interface{}{
			"tx_hash":      v.TxHash,
			"height":       v.Height,
			"order_status": v.OrderStatus,
			"timestamp":    time.Now().UnixMilli()}

		err := mysql.GetDB().Model(&model.TNftDepositRecord{}).Where(con).Updates(value).Error
		if err != nil {

			return err
		}
	}

	return nil
}

func UpdateNFTDepositByTokenID(rs []*model.TNftDepositRecord) error {
	for _, v := range rs {
		con := map[string]interface{}{
			"app_id":           v.AppID,
			"contract_address": strings.ToLower(v.ContractAddress),
			"nonce":            v.Nonce}
		value := map[string]interface{}{
			"app_order_id": v.AppOrderID,
			"equipment_id": v.EquipmentID,
			"height":       v.Height,
			"order_status": v.OrderStatus,
			"timestamp":    time.Now().UnixMilli()}

		err := mysql.GetDB().Model(&model.TNftDepositRecord{}).Where(con).Updates(value).Error
		if err != nil {

			return fmt.Errorf("update nft deposit record %v", err)
		}

		if v.OrderStatus == const_def.CodeDepositSuccess {
			equipCon := map[string]interface{}{
				"app_id":           v.AppID,
				"contract_address": strings.ToLower(v.ContractAddress),
				"token_id":         v.TokenID}

			equipValue := map[string]interface{}{
				"equipmen_id": v.EquipmentID,
				"status":      const_def.SDK_EQUIPMENT_STATUS_DEPOSIT}

			err = mysql.GetDB().Model(&model.TGameEquipment{}).Where(equipCon).Updates(equipValue).Error
			if err != nil {

				return fmt.Errorf("update equipment id %v", err)
			}
		}
	}

	return nil
}

func UpdateNFTWithdrawTxHashForMint(rs []*model.TNftWithdrawRecord) error {
	for _, v := range rs {
		con := map[string]interface{}{
			"withdraw_address": strings.ToLower(v.WithdrawAddress),
			"nonce":            v.Nonce}

		if v.GameMinterAddress != "" {
			con["minter_address"] = strings.ToLower(v.GameMinterAddress)
		}

		if v.ContractAddress != "" {
			con["contract_address"] = strings.ToLower(v.ContractAddress)
		}

		value := map[string]interface{}{
			"tx_hash":      v.TxHash,
			"height":       v.Height,
			"token_id":     v.TokenID,
			"order_status": v.OrderStatus,
			"timestamp":    time.Now().UnixMilli()}

		err := mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(con).Updates(value).Error
		if err != nil {
			return err
		}

		//update token id of game equipment
		var order model.TNftWithdrawRecord
		err = mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(con).Find(&order).Error
		if err != nil {
			return err
		}

		if order.EquipmentID == "" {
			return fmt.Errorf("the equipment id of txhash %s is null", v.TxHash)
		}

		equipCon := map[string]interface{}{
			"app_id":       order.AppID,
			"equipment_id": order.EquipmentID,
		}

		equipValue := map[string]interface{}{
			"token_id": order.TokenID,
			"status":   const_def.SDK_EQUIPMENT_STATUS_WITHDRAW}

		err = mysql.GetDB().Model(&model.TGameEquipment{}).Where(equipCon).Updates(equipValue).Error
		if err != nil {
			return err
		}

		logger.Logrus.WithFields(logrus.Fields{"Data": order}).Info("update nft withdraw mint record and game equipment record success")
	}

	return nil
}

func UpdateNFTWithdrawTxHashForUpdate(rs []*model.TNftWithdrawRecord) error {
	for _, v := range rs {
		con := map[string]interface{}{
			"contract_address": strings.ToLower(v.ContractAddress),
			"withdraw_address": strings.ToLower(v.WithdrawAddress),
			"nonce":            v.Nonce}
		value := map[string]interface{}{
			"tx_hash":      v.TxHash,
			"height":       v.Height,
			"token_id":     v.TokenID,
			"order_status": v.OrderStatus,
			"timestamp":    time.Now().UnixMilli()}

		err := mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(con).Updates(value).Error
		if err != nil {

			return err
		}
	}

	return nil
}

func UpdateNFTGameEquipment(rs []*model.TNftWithdrawRecord) error {
	for _, v := range rs {
		con := map[string]interface{}{
			"withdraw_address": strings.ToLower(v.WithdrawAddress),
			"tx_hash":          v.TxHash,
			"nonce":            v.Nonce}

		if v.GameMinterAddress != "" {
			con["minter_address"] = strings.ToLower(v.GameMinterAddress)
		}

		if v.ContractAddress != "" {
			con["contract_address"] = strings.ToLower(v.ContractAddress)
		}

		var order model.TNftWithdrawRecord
		err := mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(con).Find(&order).Error
		if err != nil {

			return err
		}

		if order.EquipmentID == "" {
			return fmt.Errorf("the equipment if of txhash %s is null", v.TxHash)
		}

		equipCon := map[string]interface{}{
			"app_id":       order.AppID,
			"equipment_id": order.EquipmentID,
		}

		equipValue := map[string]interface{}{
			"token_id": order.TokenID,
			"status":   const_def.SDK_EQUIPMENT_STATUS_WITHDRAW}

		err = mysql.GetDB().Model(&model.TGameEquipment{}).Where(equipCon).Updates(equipValue).Error
		if err != nil {

			return err
		}

		logger.Logrus.WithFields(logrus.Fields{"Data": order}).Info("UpdateNFTGameEquipment update nft withdraw game equipment success")
	}

	return nil
}
