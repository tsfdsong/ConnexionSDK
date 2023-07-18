package dashdata

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/contracts"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/core/state"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/nonce_gen"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// PreWithdrawParam
type PreWithdrawParam struct {
	GameAssetName string                   `json:"game_asset_name"`
	AppID         int64                    `json:"app_id"`
	Address       string                   `json:"address"`
	ChainID       int                      `json:"chain_id"`
	EquipmentID   string                   `json:"equipment_id"`
	Attrs         []commdata.EquipmentAttr `json:"attrs"`
	ImageURI      string                   `json:"image_uri"`
	Timestamp     string                   `json:"timestamp" binding:"gt=0"`
	SignString    string                   `json:"signature"  binding:"gt=0"`
}

// PrewithdrawResponse
type PrewithdrawResponse struct {
	PrimaryID uint64 `json:"id"`
}

// WithClaimParam
type WithClaimParam struct {
	AppID      int64  `json:"app_id"`
	Address    string `json:"address"`
	AppOrderID string `json:"app_order_id"`
}

// WithdrawClaimResponse
type WithdrawClaimResponse struct {
	AppOrderID        string `json:"app_order_id"`
	Nonce             string `json:"nonce"`
	ContractAddress   string `json:"contract_address"`
	TokenID           string `json:"token_id"`
	EquipmentID       string `json:"equipment_id"`
	WithdrawAddress   string `json:"withdraw_address"`
	TreaseAddress     string `json:"trease_address"`
	GameMinterAddress string `json:"minter_address"`
	Signature         string `json:"signature"`
	List              commdata.NftSignatureSrcData
}

// DepositParam
type DepositParam struct {
	AppID           int64  `json:"app_id"`
	Address         string `json:"address"`
	ChainID         int    `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
	UID             uint64 `binding:"-"`
	Account         string `binding:"-"`
	GameAssetName   string `binding:"-"`
	TargetAddress   string `binding:"-"`
	Nonce           string `binding:"-"`
	EquipmentID     string `binding:"-"`
}

// DepositResponse
type DepositResponse struct {
	Trease string `json:"treasure"`
	Nonce  string `json:"nonce"`
}

func (p *DepositParam) BeforeDeposit() (bool, error) {
	logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("NFTDepositController input info")

	//check deposit order existed or not
	processedRecord, err := state.GetNFTDepositRecordByCons(p.AppID, p.TokenID, strings.ToLower(p.Address), strings.ToLower(p.ContractAddress))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("NFTDepositController get processed deposit order failed")
		return false, common.NewHpError(err, int(common.IncorrectParams), "get processed deposit order failed")
	}

	if processedRecord != nil {
		logger.Logrus.WithFields(logrus.Fields{"Data": p, "Order": processedRecord}).Error("NFTDepositController deposit order is being processed")

		p.TargetAddress = processedRecord.TargetAddress
		p.Nonce = processedRecord.Nonce

		return true, nil
	}

	//check nft deposit switch
	nftContract, err := comminfo.GetNftContractByContract(int(p.AppID), p.ContractAddress)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("NftDepositHandler GetNftContractByContract failed")
		return false, common.NewHpError(err, int(common.IncorrectParams), "get nft contract info failed")
	}

	if nftContract.DepositSwitch == const_def.SDK_TABLE_SWITCH_CLOSE {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("NftDepositHandler nft deposit switch is closed")
		return false, common.NewHpError(err, int(common.IncorrectParams), "Deposit paused")
	}

	//check owner of the tokenid
	ownerof, err := contracts.GetOwnerOf(p.ContractAddress, p.TokenID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("NftDepositHandler get owner of tokenid failed")
		return false, common.NewHpError(err, int(common.IncorrectParams), "get owner of tokenId failed")
	}

	userAddress := strings.ToLower(p.Address)
	if strings.ToLower(ownerof) != userAddress {
		logger.Logrus.WithFields(logrus.Fields{"Data": p, "Ownerof": strings.ToLower(ownerof)}).Error("NftDepositHandler OwnerOf is not match")
		return false, common.NewHpError(err, int(common.IncorrectParams), "owner of tokenId is not match")
	}

	//get equipment id from game equipment table
	eid := ""
	assetName := nftContract.GameAssetName

	passAssetMap := config.GetPassAssetConfig()
	passID, isOk := passAssetMap[assetName]
	if isOk {
		eid = passID
		logger.Logrus.WithFields(logrus.Fields{"EID": eid, "PassAssetName": assetName}).Info("NftDepositHandler equipment id and pass asset name")
	} else {
		eid, err = state.GetNFTEquipmentIDByTokenID(int(p.AppID), p.TokenID, p.ContractAddress)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("NftDepositHandler nft deposit get equipment id failed")
			return false, common.NewHpError(err, int(common.IncorrectParams), "nft deposit get equipment id failed")
		}
	}

	if eid == "" {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("NftDepositHandler equipment id not found for the tokenId")
		return false, common.NewHpError(err, int(common.IncorrectParams), "nft deposit equipment id not found")
	}

	//mail bind info
	bindinfo, err := comminfo.GetBindInfo(int(p.AppID), userAddress)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("NftDepositHandler GetBindInfo failed")
		return false, common.NewHpError(err, int(common.IncorrectParams), "get mail bind info failed")
	}

	if bindinfo.UID == 0 {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("NftDepositHandler account is not game account")
		return false, common.NewHpError(err, int(common.IncorrectParams), "account must be game account")
	}

	//generate new nft deposit order
	timenow := time.Now()

	s := fmt.Sprintf("%d%s%s%s%s%d", p.AppID, bindinfo.Account, p.ContractAddress, p.TokenID, eid, timenow.Unix())
	nonce := nonce_gen.GenNonce(s)
	p.UID = bindinfo.UID
	p.Account = bindinfo.Account
	p.GameAssetName = nftContract.GameAssetName
	p.EquipmentID = eid
	p.TargetAddress = nftContract.Treasure
	p.Nonce = nonce
	return false, nil
}

func (p *DepositParam) Deposit() error {
	exist, err := p.BeforeDeposit()
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	newRecord := &model.TNftDepositRecord{
		AppID:           int(p.AppID),
		UID:             p.UID,
		Account:         p.Account,
		GameAssetName:   p.GameAssetName,
		EquipmentID:     p.EquipmentID,
		ContractAddress: strings.ToLower(p.ContractAddress),
		DepositAddress:  strings.ToLower(p.Address),
		TargetAddress:   strings.ToLower(p.TargetAddress),
		TokenID:         p.TokenID,
		OrderStatus:     const_def.CodeDepositSign,
		Nonce:           p.Nonce,
		TxHash:          "",
		Height:          0,
	}

	err = state.InsertNFTDepositRecord(newRecord)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("NftDepositHandler nft deposit order insert failed")
		return common.NewHpError(err, int(common.InnerError), "nft deposit order insert failed")
	}

	logger.Logrus.WithFields(logrus.Fields{"DepositOrder": newRecord}).Info("NftDepositHandler nft deposit success")
	return nil
}
