package dashcontroller

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/state"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func QueryGameFTAssets(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
		Data:    []common.RDashboardGameERC20Asset{},
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	retData := []common.RDashboardGameERC20Asset{}

	var p = common.PQueryGameFTAssets{}
	err := c.ShouldBindQuery(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("QueryGameAssets validator reject")
		return
	}

	chainIDStr := c.Request.Header.Get("Chainid")
	chainID, err := strconv.ParseInt(chainIDStr, 0, 64)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ChainID": chainIDStr}).Error("QueryGameAssets invalid chainID")

		r.Code = common.IncorrectParams
		r.Message = "invalid input chainID"
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p, "ChainID": chainID}).Info("QueryGameAssets input info")

	pUID, err := comminfo.GetUIDByAccount(uint(p.GameID), p.Email)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("GetUIDByAccount Error")
		return
	}

	ftAssets, err := ingame.RequestGameFTAssets(int(p.GameID), p.Email, pUID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("QueryGameAssets RequestGameFTAssets failed")
		return
	}

	ftContractMap, err := comminfo.GetFTContractAssetNameMap(int(p.GameID), int(chainID))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "GameID": p.GameID}).Error("QueryGameAssets GetFTContractAssetNameMap failed")
		return
	}

	//step4 build response
	for _, e := range ftAssets {
		v, ok := ftContractMap[e.AppCoinName]
		if !ok {
			logger.Logrus.WithFields(logrus.Fields{"GameID": p.GameID, "AppCoinName": e.AppCoinName}).Error("QueryGameAssets get ft contract by coin name failed")
			return
		}

		balance, err := tools.AsDashboardDisplayAmount(e.Balance, 8)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("QueryGameAssets AsDashboardDisplayAmount failed")
			return
		}

		erc20Item := common.RDashboardGameERC20Asset{
			Name:        v.TokenName,
			Symbol:      v.TokenSymbol,
			Contract:    v.ContractAddress,
			Decimal:     uint64(v.TokenDecimal),
			Balance:     balance,
			GameDecimal: uint64(v.GameDecimal),
		}

		retData = append(retData, erc20Item)
	}

	logger.Logrus.WithFields(logrus.Fields{"ft": retData}).Info("QueryGameAssets Assets info")

	if r.Code == common.SuccessCode {
		r.Data = retData
		r.Message = "get game assets success"
	}
}

func QueryGameNFTAssets(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
		Data:    nil,
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	retData := common.RQueryGameNFTAssets{
		GameERC721Assets: []common.RDashboardGameERC721Asset{},
		Total:            0,
	}

	var p = common.PQueryGameNFTAssets{}
	err := c.ShouldBindQuery(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("QueryGameAssets validator reject")
		return
	}

	if p.AssetName != "" {
		asName, err := url.QueryUnescape(p.AssetName)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"Data": p, "ErrMsg": err.Error()}).Error("QueryGameAssets parse url failed")
			return
		}
		p.AssetName = asName
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("QueryGameAssets input info")

	pUID, err := comminfo.GetUIDByAccount(uint(p.GameID), p.Email)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("GetUIDByAccount Error")
		return
	}

	nftAssets, total, err := ingame.RequestGameNFTAssets(int(p.GameID), p.Email, p.AssetName, int64(p.Page), int64(p.PageSize), pUID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("QueryGameAssets RequestGameNFTAssets failed")
		return
	}

	nftcontractMap, err := comminfo.GetNftContractByAppID(int(p.GameID))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "GameID": p.GameID}).Error("QueryGameAssets GetNftContractByAppID failed")

		return
	}

	logger.Logrus.WithFields(logrus.Fields{"AppID": p.GameID, "ContractInfo": nftcontractMap}).Info("QueryGameAssets GetNftContractByAppID contract info")

	for _, e := range nftAssets {
		if e.Image == "" {
			logger.Logrus.WithFields(logrus.Fields{"Data": e}).Error("QueryGameAssets input image url is empty")
			return
		}

		nftContract, ok := nftcontractMap[e.GameAssetName]
		if !ok {
			logger.Logrus.WithFields(logrus.Fields{"Data": e}).Error("QueryGameAssets nft contract info not found")
			return
		}

		imageurl := fmt.Sprintf("%s/%s.png", nftContract.BaseURL, e.Image)

		erc721Item := common.RDashboardGameERC721Asset{
			GameAssetName: nftContract.GameAssetName,
			Name:          nftContract.TokenName,
			Symbol:        nftContract.TokenSymbol,
			Contract:      nftContract.ContractAddress,
			TokenID:       e.TokenID,
			EquipmentID:   e.EquipmentID,
			Image:         imageurl,
			Attrs:         e.Attrs,
		}
		retData.GameERC721Assets = append(retData.GameERC721Assets, erc721Item)
	}

	retData.Total = total

	logger.Logrus.WithFields(logrus.Fields{"NFTAssets": retData}).Info("QueryGameAssets NFT Assets info")

	if r.Code == common.SuccessCode {
		r.Data = retData
		r.Message = "get game assets success"
	}
}

func QueryGameAssetDetail(c *gin.Context) {
	r := &common.Response{
		Code: common.SuccessCode,
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = common.PQueryGameAssetDetail{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("QueryGameAssetDetail validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("QueryGameAssetDetail input info")

	if p.EquipmentID == "" && p.TokenID == "" {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("QueryGameAssetDetail input equipment_id or token_id is empty")
		r.Code = common.IncorrectParams
		r.Message = "token id or equip id cannot empty"
		return
	}

	if p.GameAssetName == "" {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("QueryGameAssetDetail input game asset name is empty")
		r.Code = common.IncorrectParams
		r.Message = "game asset name cannot empty"
		return
	}

	nftContractInfo, err := comminfo.GetNftContractByAssetName(int(p.GameID), p.GameAssetName)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("QueryGameAssetDetail GetNftContractByAssetName failed")

		r.Code = common.InnerError
		r.Message = "nft contract info not found"
		return
	}

	network, ok := config.GetChainNetwork()[int64(nftContractInfo.ChainID)]
	if !ok {
		logger.Logrus.Error("not found network config")
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("QueryGameAssetDetail get network from chain id")
		r.Code = common.InnerError
		r.Message = "network not found"
		return
	}

	//get pass asset map
	passAssetMap := config.GetPassAssetConfig()
	passEID, isok := passAssetMap[p.GameAssetName]
	if isok {
		imageurl := fmt.Sprintf("%s/%s.png", nftContractInfo.BaseURL, passEID)

		retData := common.RDashboardGameQueryAssetDetail{
			GameAssetName: p.GameAssetName,
			Name:          nftContractInfo.TokenName,
			Symbol:        nftContractInfo.TokenSymbol,
			Contract:      nftContractInfo.ContractAddress,
			Trease:        nftContractInfo.Treasure,
			TokenID:       p.TokenID,
			EquipmentID:   passEID,
			Image:         imageurl,
			ChainImage:    "",
			Attrs:         []commdata.EquipmentAttr{},
			Network:       network,
			Standard:      "ERC721",
		}

		logger.Logrus.WithFields(logrus.Fields{"nftDetail": retData}).Info("QueryGameAssetDetail nftDetail for pass token info")

		r.Data = retData
		r.Message = "nft asset detail for pass success"

		return
	}

	finalChainURI := ""
	if p.TokenID != "" {
		finalChainURI, err = comminfo.GetChainSVG("", nftContractInfo.ContractAddress, p.TokenID)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "TokenID": p.TokenID}).Error("QueryGameAssetDetail GetTokenURI parse svg failed")

			r.Code = common.InnerError
			r.Message = "parse svg failed"
			return
		}
	}

	equipmentID := p.EquipmentID
	if p.EquipmentID == "" && p.TokenID != "" {
		eid, err := state.GetNFTEquipmentIDByTokenID(int(p.GameID), p.TokenID, nftContractInfo.ContractAddress)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("QueryGameAssetDetail get equipmenID by tokenID failed")
			r.Code = common.InnerError
			r.Message = "get equipment id failed"
			return
		}

		if eid == "" {
			retData := common.RDashboardGameQueryAssetDetail{
				GameAssetName: nftContractInfo.GameAssetName,
				Name:          nftContractInfo.TokenName,
				Symbol:        nftContractInfo.TokenSymbol,
				Contract:      nftContractInfo.ContractAddress,
				Trease:        nftContractInfo.Treasure,
				TokenID:       p.TokenID,
				EquipmentID:   eid,
				Image:         "",
				ChainImage:    finalChainURI,
				Attrs:         []commdata.EquipmentAttr{},
				Network:       network,
				Standard:      "ERC721",
			}

			logger.Logrus.WithFields(logrus.Fields{"nftDetail": retData}).Info("QueryGameAssetDetail nftDetail info")

			r.Data = retData
			r.Message = "nft asset detail success"
			return
		}

		logger.Logrus.WithFields(logrus.Fields{"Data": p, "EquipmentID": eid}).Info("QueryGameAssetDetail GetNFTEquipmentIDByTokenID info")

		equipmentID = eid
	}

	nftDetail, err := ingame.RequestGameNFTAssetDetail(int(p.GameID), equipmentID, false)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("QueryGameAssetDetail RequestGameNFTAssetDetail failed")
		r.Code = common.InnerError
		r.Message = "attrs not found"
		return
	}

	if nftDetail.Image == "" {
		logger.Logrus.WithFields(logrus.Fields{"Data": nftDetail}).Error("QueryGameAssetDetail input image url is empty")

		r.Code = common.InnerError
		r.Message = "image url is null"
		return
	}

	imageurl := fmt.Sprintf("%s/%s.png", nftContractInfo.BaseURL, nftDetail.Image)

	retData := common.RDashboardGameQueryAssetDetail{
		GameAssetName: nftDetail.GameAssetName,
		Name:          nftContractInfo.TokenName,
		Symbol:        nftContractInfo.TokenSymbol,
		Contract:      nftContractInfo.ContractAddress,
		Trease:        nftContractInfo.Treasure,
		TokenID:       p.TokenID,
		EquipmentID:   equipmentID,
		Image:         imageurl,
		ChainImage:    finalChainURI,
		Attrs:         nftDetail.Attrs,
		Network:       network,
		Standard:      "ERC721",
	}

	logger.Logrus.WithFields(logrus.Fields{"nftDetail": retData}).Info("QueryGameAssetDetail nftDetail info")

	r.Data = retData
	r.Message = "nft asset detail success"
}
