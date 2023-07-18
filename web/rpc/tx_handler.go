package rpc

import (
	"encoding/json"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	listenerDb "github/Connector-Gamefi/ConnectorGoSDK/core/chainlistener/common"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/core/treasurealter"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

func NewFtDeposit(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}
	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)
	var p = FtDepositReq{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NewFtDeposit validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
	logger.Logrus.WithFields(logrus.Fields{"params": p}).Info("NewFtDeposit params")

	contract, err := comminfo.GetFTContractByDepositTreasureAndChain(p.Treasure, p.ChainId)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NewFtDeposit get contract by treasure failed")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	order := &model.TFtDepositRecord{
		ChainID:         p.ChainId,
		ContractAddress: strings.ToLower(contract.ContractAddress),
		TargetAddress:   p.Treasure,
		DepositAddress:  p.From,
		Amount:          p.Amount,
		Nonce:           p.Nonce,
		TxHash:          p.Tx,
		OrderStatus:     const_def.CodeDepositOnChain,
		Height:          p.Height,
	}
	err = listenerDb.InsertFTDeposit(order)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("insert new ft deposit failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
}

func NewFtWithdraw(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}
	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = FtWithdrawReq{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NewFtWithdraw validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
	logger.Logrus.WithFields(logrus.Fields{"params": p}).Info("NewFtWithdraw params")

	contract, err := comminfo.GetFTContractByTreasureAndChain(p.Treasure, p.ChainId)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NewFtWithdraw get contract by treasure failed")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	order := &model.TFtWithdrawRecord{
		ChainID:         p.ChainId,
		ContractAddress: strings.ToLower(contract.ContractAddress),
		WithdrawAddress: p.From,
		Amount:          p.Amount,
		Nonce:           p.Nonce,
		TxHash:          p.Tx,
		OrderStatus:     const_def.CodeWithdrawOnChain,
		Height:          p.Height,
	}
	err = listenerDb.UpdateFTWithdraw([]*model.TFtWithdrawRecord{order})
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("UpdateFTWithdraw ft withdraw failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
}

func ConfirmedFtDeposit(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}
	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = TxConfirmedReq{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ConfirmedFtDeposit validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
	logger.Logrus.WithFields(logrus.Fields{"params": p}).Info("ConfirmedFtDeposit params")

	var ftRecord model.TFtDepositRecord
	err = mysql.GetDB().Where("app_id = ? and order_status= ? and tx_hash = ?", p.GameId, const_def.CodeDepositOnChain, p.Tx).First(&ftRecord).Error
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ConfirmedFtDeposit tx not exist")
		r.Code = common.NotFoundTx
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	ftInfo, err := comminfo.GetFTContractByAppAndChain(p.GameId, ftRecord.ChainID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ConfirmedFtDeposit get ft contract info failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	alterAmounts := make(map[string][]string, 0)
	alterAmounts[ftRecord.ContractAddress] = []string{ftRecord.Amount}
	err = treasurealter.HandDepositTreasureAlterLark(ftRecord.AppID, ftRecord.ContractAddress, alterAmounts)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "ContractAddress": ftRecord.ContractAddress, "AlterAmounts": alterAmounts}).Error("ConfirmedFtDeposit HandDepositTreasureAlterLark failed")
	} else {
		logger.Logrus.WithFields(logrus.Fields{"AppID": ftRecord.AppID, "ContractAddress": ftRecord.ContractAddress, "AlterAmounts": alterAmounts}).Info("ConfirmedFtDeposit HandDepositTreasureAlterLark success")
	}

	if p.Status == STATUS_TX_SUCCESS {
		err := listenerDb.FTGameDeposit(&ftRecord, ftInfo)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("FTGameDeposit notify game failed")
			r.Code = common.RequestGameServerFailed
			r.Message = common.ErrorMap[int(r.Code)]
			return
		}
	} else {
		ftRecord.OrderStatus = const_def.CodeDepositTxFailed
	}

	con := map[string]interface{}{
		"app_id":       p.GameId,
		"order_status": const_def.CodeDepositOnChain,
		"tx_hash":      p.Tx}
	value := map[string]interface{}{
		"app_order_id": ftRecord.AppOrderID,
		"height":       p.Height,
		"order_status": ftRecord.OrderStatus,
		"timestamp":    time.Now().UnixMilli()}

	err = mysql.GetDB().Model(&model.TFtDepositRecord{}).Where(con).Updates(value).Error
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("insert new ft deposit failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
}

func ConfirmedFtWithdraw(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}
	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = TxConfirmedReq{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ConfirmedFtWithdraw validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
	logger.Logrus.WithFields(logrus.Fields{"params": p}).Info("ConfirmedFtWithdraw params")

	var ftRecord model.TFtWithdrawRecord
	err = mysql.GetDB().Where("app_id = ? and order_status= ? and tx_hash = ?", p.GameId, const_def.CodeWithdrawOnChain, p.Tx).First(&ftRecord).Error
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ConfirmedFtWithdraw tx not exist")
		r.Code = common.NotFoundTx
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	alterAmounts := make(map[string][]string, 0)
	alterAmounts[ftRecord.ContractAddress] = []string{ftRecord.Amount}
	err = treasurealter.HandTreasureAlterLark(ftRecord.AppID, ftRecord.ContractAddress, alterAmounts)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "ContractAddress": ftRecord.ContractAddress, "AlterAmounts": alterAmounts}).Error("ConfirmedFtDeposit HandDepositTreasureAlterLark failed")
	} else {
		logger.Logrus.WithFields(logrus.Fields{"AppID": ftRecord.AppID, "ContractAddress": ftRecord.ContractAddress, "AlterAmounts": alterAmounts}).Info("ConfirmedFtDeposit HandDepositTreasureAlterLark success")
	}

	orderStatus := const_def.CodeWithdrawTxFailed
	if p.Status == STATUS_TX_SUCCESS {
		orderStatus = const_def.CodeWithdrawSuccess
	}
	con := map[string]interface{}{
		"app_id":       p.GameId,
		"order_status": const_def.CodeWithdrawOnChain,
		"tx_hash":      p.Tx}
	value := map[string]interface{}{
		"order_status": orderStatus,
		"height":       p.Height,
		"timestamp":    time.Now().UnixMilli()}

	err = mysql.GetDB().Model(&model.TFtWithdrawRecord{}).Where(con).Updates(value).Error
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("update FT withdraw record failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
}

func ConfirmedNftWithdraw(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}
	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = TxConfirmedReq{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ConfirmedNftWithdraw validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
	logger.Logrus.WithFields(logrus.Fields{"params": p}).Info("ConfirmedNftWithdraw params")

	orderStatus := const_def.CodeWithdrawTxFailed
	if p.Status == STATUS_TX_SUCCESS {
		orderStatus = const_def.CodeWithdrawSuccess
	}
	//update db
	con := map[string]interface{}{
		"app_id":       p.GameId,
		"order_status": const_def.CodeWithdrawOnChain,
		"tx_hash":      p.Tx}
	value := map[string]interface{}{
		"order_status": orderStatus,
		"height":       p.Height,
		"timestamp":    time.Now().UnixMilli()}

	err = mysql.GetDB().Model(&model.TNftWithdrawRecord{}).Where(con).Updates(value).Error
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("update NFT withdraw record failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
}

func NewNftMintWithdraw(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}
	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = NftWithdrawReq{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NewNftLootWithdraw validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
	logger.Logrus.WithFields(logrus.Fields{"params": p}).Info("NewNftLootWithdraw params")

	order := &model.TNftWithdrawRecord{
		GameMinterAddress: p.MinterAddr,
		WithdrawAddress:   p.From,
		ContractAddress:   p.ContractAddr,
		Nonce:             p.Nonce,
		TokenID:           p.TokenID,
		EquipmentID:       p.EquipID,
		TxHash:            p.Tx,
		OrderStatus:       const_def.CodeWithdrawOnChain,
		Height:            p.Height,
	}
	err = listenerDb.UpdateNFTWithdrawTxHashForMint([]*model.TNftWithdrawRecord{order})
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("UpdateNFTWithdrawTxHashForMint nft withdraw failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
}

func NewNftUpdateWithdraw(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}
	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = NftWithdrawReq{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NewNftUpdateWithdraw validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
	logger.Logrus.WithFields(logrus.Fields{"params": p}).Info("NewNftUpdateWithdraw params")

	order := &model.TNftWithdrawRecord{
		TreaseAddress:   p.Treasure,
		ContractAddress: p.ContractAddr,
		WithdrawAddress: p.From,
		Nonce:           p.Nonce,
		TokenID:         p.TokenID,
		TxHash:          p.Tx,
		OrderStatus:     const_def.CodeWithdrawOnChain,
		Height:          p.Height,
	}
	err = listenerDb.UpdateNFTWithdrawTxHashForUpdate([]*model.TNftWithdrawRecord{order})
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("UpdateNFTWithdrawTxHashForUpdate nft withdraw failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
}

func NewNftDeposit(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}
	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = NftDepositReq{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NewNftDeposit validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
	logger.Logrus.WithFields(logrus.Fields{"params": p}).Info("NewNftDeposit params")

	order := &model.TNftDepositRecord{
		ContractAddress: p.ContractAddr,
		TargetAddress:   p.TargetAddr,
		DepositAddress:  p.From,
		Nonce:           p.Nonce,
		TokenID:         p.TokenID,
		TxHash:          p.Tx,
		OrderStatus:     const_def.CodeDepositOnChain,
		Height:          p.Height,
	}
	err = listenerDb.UpdateNFTDeposit([]*model.TNftDepositRecord{order})
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("UpdateNFTDeposit nft deposit failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
}

func ConfirmedNftDeposit(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}
	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = TxConfirmedReq{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ConfirmedNftDeposit validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
	logger.Logrus.WithFields(logrus.Fields{"params": p}).Info("ConfirmedNftDeposit params")

	var nftRecord model.TNftDepositRecord
	err = mysql.GetDB().Where("app_id = ? and order_status= ? and tx_hash = ?", p.GameId, const_def.CodeDepositOnChain, p.Tx).First(&nftRecord).Error
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ConfirmedNftDeposit tx not exist")
		r.Code = common.NotFoundTx
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	if p.Status == STATUS_TX_SUCCESS {
		err := listenerDb.NFTGameDeposit(&nftRecord)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NFTGameDeposit notify game failed")
			r.Code = common.RequestGameServerFailed
			r.Message = common.ErrorMap[int(r.Code)]
			return
		}
	} else {
		nftRecord.OrderStatus = const_def.CodeDepositTxFailed
	}
	con := map[string]interface{}{
		"app_id":       p.GameId,
		"order_status": const_def.CodeDepositOnChain,
		"tx_hash":      p.Tx}
	value := map[string]interface{}{
		"app_order_id":    nftRecord.AppOrderID,
		"equipment_id":    nftRecord.EquipmentID,
		"uid":             nftRecord.UID,
		"game_asset_name": nftRecord.GameAssetName,
		"height":          nftRecord.Height,
		"order_status":    nftRecord.OrderStatus,
		"timestamp":       time.Now().UnixMilli()}

	err = mysql.GetDB().Model(&model.TNftDepositRecord{}).Where(con).Updates(value).Error
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("update nft deposit failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}
	passAssetMap := config.GetPassAssetConfig()
	//check pass asset
	_, isOK := passAssetMap[nftRecord.GameAssetName]
	if !isOK && nftRecord.OrderStatus == const_def.CodeDepositSuccess {
		// get nft details by equipment id
		nftDetails, err := ingame.RequestGameNFTAssetDetail(nftRecord.AppID, nftRecord.EquipmentID, true)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "NFTDepositData": nftRecord}).Error("RequestGameNFTAssetDetail failed")
			return
		}

		attrs, err := json.Marshal(nftDetails.Attrs)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "NFTDepositData": nftRecord, "Attrs": nftDetails.Attrs}).Error("ConfirmedNftDeposit NFT marshal attrs failed")
			return
		}

		// update game equipment table
		equipCon := map[string]interface{}{
			"app_id":           nftRecord.AppID,
			"contract_address": strings.ToLower(nftRecord.ContractAddress),
			"token_id":         nftRecord.TokenID}

		equipValue := map[string]interface{}{
			"equipment_id":   nftRecord.EquipmentID,
			"account":        nftRecord.Account,
			"equipment_attr": datatypes.JSON(attrs),
			"status":         const_def.SDK_EQUIPMENT_STATUS_DEPOSIT}

		err = mysql.GetDB().Model(&model.TGameEquipment{}).Where(equipCon).Updates(equipValue).Error
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "NFTDepositData": nftRecord}).Error("ConfirmedNftDeposit NFT update equipment record failed")
			return
		}
	}
}
