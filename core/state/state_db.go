package state

import (
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func FindAllFTRecords() ([]model.TFtWithdrawRecord, error) {
	var first model.TFtWithdrawRecord
	result := make([]model.TFtWithdrawRecord, 0)

	err := mysql.WrapFindAllByCondition(first.TableName(), map[string]interface{}{"order_status": int(const_def.CodeWithdrawSign)}, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateFTOrderStatus update order status
func UpdateFTOrderAndRiskStatus(id, orderCode, riskCode int) error {
	res := mysql.GetDB().Model(&model.TFtWithdrawRecord{}).Where("id = ?", id)
	updateColumns := make(map[string]interface{})
	now := time.Now().UnixMilli()
	if orderCode > 0 {
		updateColumns["order_status"] = orderCode
		updateColumns["timestamp"] = now
	}

	if riskCode > 0 {
		updateColumns["risk_status"] = riskCode
		updateColumns["timestamp"] = now
	}
	res = res.Updates(updateColumns)

	if res == nil || res.Error != nil {

		if res == nil {
			return fmt.Errorf("res is nil")
		} else {
			return res.Error
		}

		return res.Error
	}

	return nil
}

// GetNFTGameEquipment get game equipment
func GetNFTGameEquipment(appID int, equipID, contractAddr string) (*model.TGameEquipment, error) {
	var data model.TGameEquipment

	con := map[string]interface{}{"app_id": appID, "equipment_id": equipID, "contract_address": strings.ToLower(contractAddr)}
	err := mysql.GetDB().Model(&model.TGameEquipment{}).Where(con).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.TGameEquipment{
				AppID:       appID,
				EquipmentID: equipID,
				TokenID:     "",
			}, nil
		}

		return nil, err
	}

	return &data, nil
}

// GetNFTGameEquipmentList tokenid => equipment
func GetNFTGameEquipmentList(appID int, contractAddr string) (map[string]model.TGameEquipment, error) {
	var datas []model.TGameEquipment

	con := map[string]interface{}{"app_id": appID, "contract_address": strings.ToLower(contractAddr)}
	err := mysql.GetDB().Model(&model.TGameEquipment{}).Where(con).Find(&datas).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]model.TGameEquipment, 0)
	for _, data := range datas {
		result[data.TokenID] = data
	}

	return result, nil
}

// GetNFTGameEquipmentByEquipmenID get game equipment
func GetNFTGameEquipmentByEquipmenID(appID int, contractAddr, equipmenid string) (*model.TGameEquipment, error) {
	var data model.TGameEquipment

	con := map[string]interface{}{"app_id": appID, "contract_address": strings.ToLower(contractAddr), "equipment_id": equipmenid}
	err := mysql.GetDB().Model(&model.TGameEquipment{}).Where(con).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetNFTGameEquipmentByTokenID get game equipment
func GetNFTEquipmentIDByTokenID(appID int, tokenid string, contractAddr string) (string, error) {
	var data model.TGameEquipment

	con := map[string]interface{}{"app_id": appID, "token_id": tokenid, "contract_address": strings.ToLower(contractAddr)}
	err := mysql.GetDB().Model(&model.TGameEquipment{}).Where(con).First(&data).Error
	if err != nil && err != gorm.ErrRecordNotFound {

		return "", err
	}

	return data.EquipmentID, nil
}

// UpdateGameEquipment insert or update
func UpdateOrInsertGameEquipment(para *model.TGameEquipment) error {
	con := map[string]interface{}{"app_id": para.AppID, "equipment_id": para.EquipmentID, "contract_address": strings.ToLower(para.ContractAddress)}

	if para.TokenID == "" {
		//insert
		if err := mysql.GetDB().Model(&model.TGameEquipment{}).Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&para).Error; err != nil {

			return fmt.Errorf("insert game equipment %v", err)
		}
	} else {
		//update all fields
		values := map[string]interface{}{"token_id": para.TokenID, "account": para.Account, "chain_id": para.ChainID, "withdraw_switch": para.WithdrawSwitch, "status": para.Status, "equipment_attr": para.EquipmentAttr, "image_uri": para.ImageURI}

		if err := mysql.GetDB().Model(&model.TGameEquipment{}).Where(con).Updates(values).Error; err != nil {

			return fmt.Errorf("update game equipment %v", err)
		}
	}
	return nil
}

// FindAllNFTRecordSign get all records of waiting for sign
func FindAllNFTRecordSign() ([]model.TNftWithdrawRecord, error) {
	data := make([]model.TNftWithdrawRecord, 0)

	err := mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(map[string]interface{}{"order_status": int(const_def.CodeWithdrawSign)}).Limit(10).Find(&data).Error
	if err != nil {

		return nil, err
	}

	return data, nil
}

// InsertNFTWithdrawRecord insert order
func InsertNFTWithdrawRecord(para *model.TNftWithdrawRecord) error {
	para.Timestamp = time.Now().UnixMilli()
	if err := mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Create(&para).Error; err != nil {

		return err
	}

	return nil
}

// UpdateFTOrderStatus update order status
func UpdateNFTWithdrawOrderStatus(para *model.TNftWithdrawRecord) error {
	con := map[string]interface{}{"id": para.ID}

	values := map[string]interface{}{"app_order_id": para.AppOrderID, "order_status": para.OrderStatus, "signature": para.Signature, "signature_hash": para.SignatureHash, "signature_source": para.SignatureSrc, "trease_address": para.TreaseAddress, "risk_status": para.RiskStatus, "timestamp": time.Now().UnixMilli()}
	err := mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(con).Updates(values).Error

	if err != nil {

		return err
	}

	return nil
}

// UpdateNFTWithdrawByID update order status
func UpdateNFTWithdrawByID(para *model.TNftWithdrawRecord) error {
	con := map[string]interface{}{"id": para.ID}

	values := map[string]interface{}{"token_id": para.TokenID, "app_order_id": para.AppOrderID, "order_status": para.OrderStatus, "signature": para.Signature, "signature_hash": para.SignatureHash, "signature_source": para.SignatureSrc, "timestamp": time.Now().UnixMilli()}
	err := mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(con).Updates(values).Error

	if err != nil {

		return err
	}

	return nil
}

// GetNFTPreWithdrawOrderRecord get record by appid , uid and order id
func GetNFTPreWithdrawOrderRecord(appId int64, equipmentId string) (*model.TNftWithdrawRecord, error) {
	var data model.TNftWithdrawRecord

	err := mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where("app_id=? and equipment_id=? and order_status in ?", appId, equipmentId, []int{const_def.CodeWithdrawRisking, const_def.CodeWithdrawSign, const_def.CodeWithdrawWaitClaim}).First(&data).Error
	if err != nil && err != gorm.ErrRecordNotFound {

		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &data, nil
}

// GetWithdrawOrderByPrimKey get record by primary key
func GetWithdrawOrderByPrimKey(id uint64) (*model.TNftWithdrawRecord, error) {
	var data model.TNftWithdrawRecord

	con := map[string]interface{}{"id": id}
	err := mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(con).First(&data).Error
	if err != nil {

		return nil, err
	}

	return &data, nil
}

// InsertNFTDepositRecord insert order
func InsertNFTDepositRecord(para *model.TNftDepositRecord) error {
	para.Timestamp = time.Now().UnixMilli()
	if err := mysql.GetDB().Model(&model.TNftDepositRecord{}).Create(&para).Error; err != nil {

		return err
	}

	return nil
}

// UpdateNFTDepositRecord update
func UpdateNFTDepositRecord(para *model.TNftDepositRecord) error {
	con := map[string]interface{}{"app_id": para.AppID, "account": para.Account, "equipment_id": para.EquipmentID, "contract_address": strings.ToLower(para.ContractAddress), "deposit_address": strings.ToLower(para.DepositAddress), "trease_address": strings.ToLower(para.TargetAddress), "token_id": para.TokenID, "order_status": para.OrderStatus}

	values := map[string]interface{}{
		"app_order_id": para.AppOrderID, "uid": para.UID, "tx_hash": para.TxHash, "order_status": para.OrderStatus, "nonce": para.Nonce, "timestamp": time.Now().UnixMilli(),
	}

	err := mysql.GetDB().Model(&model.TNftDepositRecord{}).Where(con).Updates(values).Error
	if err != nil {

		return err
	}

	return nil
}

// GetNFTDepositRecordByCons
func GetNFTDepositRecordByCons(appId int64, tokenId, depositAddr, contractAddr string) (*model.TNftDepositRecord, error) {
	var res model.TNftDepositRecord
	err := mysql.GetDB().Model(&model.TNftDepositRecord{}).Where("app_id=? and token_id=? and deposit_address=? and contract_address=? and order_status in ?", appId, tokenId, depositAddr, contractAddr, []int{const_def.CodeDepositSign, const_def.CodeDepositOnChain}).First(&res).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &res, nil
}

// find sign
func FTFindAllSendSignRecords() ([]model.TFtWithdrawRecord, error) {
	var first model.TFtWithdrawRecord
	result := []model.TFtWithdrawRecord{}
	count := int64(0)
	err := mysql.WrapFindAllByConditionCheckOrderLimit(first.TableName(), map[string]interface{}{"order_status": int(const_def.CodeWithdrawSign)}, &result, "created_at", int(config.GetSignConfig().SignRequestNum), &count)
	if err != nil {

		return nil, err
	}

	return result, nil
}

// InsertFTWithdrawRecord insert order
func InsertFTWithdrawRecord(para *model.TFtWithdrawRecord) error {
	para.Timestamp = time.Now().UnixMilli()
	if err := mysql.GetDB().Model(&model.TFtWithdrawRecord{}).Create(&para).Error; err != nil {

		return err
	}

	return nil
}

// UpdateFTOrderStatus update order status
func UpdateFTWithdrawOrderStatus(para *model.TFtWithdrawRecord) error {
	con := map[string]interface{}{"id": para.ID}

	values := map[string]interface{}{"app_order_id": para.AppOrderID, "order_status": para.OrderStatus, "signature": para.Signature, "signature_hash": para.SignatureHash, "risk_status": para.RiskStatus, "timestamp": time.Now().UnixMilli()}
	err := mysql.GetDB().Model(&model.TFtWithdrawRecord{}).Where(con).Updates(values).Error

	if err != nil {

		return err
	}

	return nil
}
