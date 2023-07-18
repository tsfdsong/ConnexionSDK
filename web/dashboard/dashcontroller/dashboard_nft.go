package dashcontroller

import (
	"encoding/json"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/contracts"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/core/state"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/nonce_gen"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/pool"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"github/Connector-Gamefi/ConnectorGoSDK/web/dashboard/dashdata"
	"github/Connector-Gamefi/ConnectorGoSDK/web/dashboard/dashdata/request"
	"github/Connector-Gamefi/ConnectorGoSDK/web/dashboard/dashdata/response"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

func CheckWalletSign(timestampStr, rawData, signData, walletAddr string) error {
	apiDelayTime := config.GetAPIDelayTime()
	timestmap, err := strconv.ParseInt(timestampStr, 10, 0)
	if err != nil {
		return fmt.Errorf("parse timestamp %s failed, %v", timestampStr, err)
	}

	if timestmap+apiDelayTime < time.Now().Unix() {
		return fmt.Errorf("timestamp %s timeout", timestampStr)
	}

	if !tools.CheckPersonalSign(rawData, signData, walletAddr) {
		return fmt.Errorf("sign not match")
	}

	return nil
}

// NFTPreWithdrawController
func NFTPreWithdrawController(c *gin.Context) {
	r := &common.Response{
		Code: common.SuccessCode,
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	p := dashdata.PreWithdrawParam{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NFTPreWithdrawController prewithdraw bind input parameter failed")
		r.Code = common.IncorrectParams
		r.Message = "invalid parameter"
		return
	}

	rawSignString := fmt.Sprintf("app_id=%d&address=%s&equipment_id=%s&timestamp=%s", p.AppID, strings.ToLower(p.Address), p.EquipmentID, p.Timestamp)
	err = CheckWalletSign(p.Timestamp, rawSignString, p.SignString, p.Address)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "input": p}).Error("NFTPreWithdrawController check wallet signature failed")
		r.Code = common.IncorrectParams
		r.Message = "signature not match"
		return
	}

	if p.GameAssetName == "" {
		logger.Logrus.WithFields(logrus.Fields{"input": p}).Error("NFTPreWithdrawController input game asset name is empty")
		r.Code = common.IncorrectParams
		r.Message = "game asset name cannot empty"
		return
	}

	if p.ImageURI == "" {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": "image uri cannot empty", "input": p}).Error("NFTPreWithdrawController prewithdraw image uri is empty")
		r.Code = common.IncorrectParams
		r.Message = "image uri cannot empty"
		return
	}

	if p.EquipmentID == "" {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": "equipment id cannot empty", "input": p}).Error("NFTPreWithdrawController prewithdraw equipment id is empty")
		r.Code = common.IncorrectParams
		r.Message = "equipment id cannot empty"
		return
	}

	//check prewithdraw order existed or not
	processingorder, err := state.GetNFTPreWithdrawOrderRecord(p.AppID, p.EquipmentID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "input": p}).Error("NFTPreWithdrawController get processing prewithdraw order failed")
		r.Code = common.IncorrectParams
		r.Message = "get processing prewithdraw order failed"
		return
	}

	if processingorder != nil {
		logger.Logrus.WithFields(logrus.Fields{"ProcessedOrder": processingorder}).Info("NFTPreWithdrawController prewithdraw order is being processed")

		r.Data = dashdata.PrewithdrawResponse{
			PrimaryID: processingorder.ID,
		}
		r.Message = "prewithdraw order is being processed"
		return
	}

	//check nft contract withdraw switch
	nftinfo, err := comminfo.GetNftContractByAssetName(int(p.AppID), p.GameAssetName)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "input": p}).Error("NFTPreWithdrawController GetNftContractByAssetName failed")

		r.Code = common.IncorrectParams
		r.Message = "nft contract not found"
		return
	}

	if nftinfo.WithdrawSwitch == const_def.SDK_TABLE_SWITCH_CLOSE {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": "NFT Withdraw switch is closed", "input": p}).Error("NFTPreWithdrawController prewithdraw get nft contract withdraw switch is closed")

		r.Code = common.InnerError
		r.Message = "Withdrawal paused"
		return
	}

	//get mail bind
	bindInfo, err := comminfo.GetBindInfo(int(p.AppID), strings.ToLower(p.Address))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "input": p}).Error("NFTPreWithdrawController GetBindInfo failed")

		r.Code = common.IncorrectParams
		r.Message = "get mail bind failed"
		return
	}

	if bindInfo.UID == 0 {
		logger.Logrus.WithFields(logrus.Fields{"input": p}).Error("NFTPreWithdrawController account is not game account")

		r.Code = common.IncorrectParams
		r.Message = "account must be a game account"
		return
	}

	//check user withdraw switch
	withdrawSwitch, err := comminfo.GetUserWithdrawSwitch(int(p.AppID), bindInfo.Account)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "input": p, "Account": bindInfo.Account}).Error("PrewithdrawHandler GetUserWithdrawSwitch failed")
		r.Code = common.IncorrectParams
		r.Message = "get user withdraw switch failed"
		return
	}

	if withdrawSwitch == const_def.SDK_TABLE_SWITCH_CLOSE {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "input": p, "Account": bindInfo.Account}).Error("PrewithdrawHandler user is not allowed to withdraw")
		r.Code = common.IncorrectParams
		r.Message = "User Withdrawal paused"
		return
	}

	if p.ChainID != nftinfo.ChainID {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "input": p, "ChainID": nftinfo.ChainID}).Error("PrewithdrawHandler chainid is not match")
		r.Code = common.IncorrectParams
		r.Message = "chainid is not match"
		return
	}

	//check equipment withdraw switch
	gameEquip, err := state.GetNFTGameEquipment(int(p.AppID), p.EquipmentID, nftinfo.ContractAddress)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "input": p}).Error("PrewithdrawHandler GetNFTGameEquipment failed")
		r.Code = common.IncorrectParams
		r.Message = "get nft equipment failed"
		return
	}

	if gameEquip.WithdrawSwitch == int8(const_def.CodeEquipmentWithdrawClosed) {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "EquipmentData": gameEquip}).Error("PrewithdrawHandler equipment is not allowed to withdraw")
		r.Code = common.IncorrectParams
		r.Message = "nft equipment is not allowed to withdraw"
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("NFTPreWithdrawController input info")

	attrs := make([]byte, 0)
	chainattr := make([]commdata.EquipmentAttr, 0)
	if p.GameAssetName == config.GetGameEquipmentName() {
		//check attrs match the game equipment attrs
		nftDetail, err := ingame.RequestGameNFTAssetDetail(int(p.AppID), p.EquipmentID, false)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "input": p}).Error("NFTPreWithdrawController RequestGameNFTAssetDetail failed")
			r.Code = common.InnerError
			r.Message = "attrs not found"
			return
		}

		if !reflect.DeepEqual(p.Attrs, nftDetail.Attrs) {
			logger.Logrus.WithFields(logrus.Fields{"input": p, "GameAttrs": nftDetail.Attrs}).Error("NFTPreWithdrawController attrs not match")
			r.Code = common.InnerError
			r.Message = "attrs not match input"
			return
		}

		//construct withdraw order
		chainattr, err = contracts.ConvertAttrsWithDecimal(nftinfo.Decimal, p.Attrs)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "input": p}).Error("NFTPreWithdrawController ConvertAttrsWithDecimal failed")

			r.Code = common.PrewithdrawCode
			r.Message = "convert nft attrs failed"
			return
		}

		attrs, err = json.Marshal(chainattr)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NFTPreWithdrawController prewithdraw marshal chain attrs failed")
			r.Code = common.IncorrectParams
			r.Message = "marshal chain attrs failed"
			return
		}
	}

	//generate nonce
	s := fmt.Sprintf("%d%s%d%d", p.AppID, nftinfo.ContractAddress, bindInfo.UID, time.Now().Unix())
	nonce := nonce_gen.GenNonce(s)

	order := &model.TNftWithdrawRecord{
		AppID:             int(p.AppID),
		UID:               bindInfo.UID,
		AppOrderID:        "",
		EquipmentID:       p.EquipmentID,
		GameAssetName:     p.GameAssetName,
		Account:           bindInfo.Account,
		ContractAddress:   nftinfo.ContractAddress,
		WithdrawAddress:   p.Address,
		TreaseAddress:     nftinfo.Treasure,
		GameMinterAddress: nftinfo.MinterAddress,
		TokenID:           gameEquip.TokenID,
		TxHash:            "",
		OrderStatus:       int(const_def.CodeWithdrawNone),
		Signature:         "",
		SignatureHash:     "",
		SignatureSrc:      datatypes.JSON(attrs),
		Nonce:             nonce,
		RiskStatus:        0,
		RiskReviewer:      "",
		Height:            0,
	}

	err = state.InsertNFTWithdrawRecord(order)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Order": order}).Error("NFTPreWithdrawController InsertNFTWithdrawRecord failed")
		r.Code = common.IncorrectParams
		r.Message = "insert withdraw record failed"
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Order": order}).Info("NFTPreWithdrawController insert withdraw order success")

	err = state.PrewithdrawHandler(order, chainattr, p.Attrs, p.ImageURI, p.ChainID, c.Request.Context())
	if err != nil {
		r.Code = common.PrewithdrawCode
		r.Message = "nft prewithdraw failed"
		return
	}

	if r.Code == common.SuccessCode {
		r.Data = dashdata.PrewithdrawResponse{
			PrimaryID: order.ID,
		}
		r.Message = "nft prewithdraw success"
	}
}

// NFTClaimController
func NFTClaimController(c *gin.Context) {
	r := &common.Response{
		Code: common.SuccessCode,
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	primaryID := c.Query("id")

	logger.Logrus.WithFields(logrus.Fields{"PrimaryID": primaryID}).Info("NFTClaimController input parameters")

	pid, err := strconv.ParseInt(primaryID, 0, 64)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "PrimaryID": primaryID}).Error("NFTClaimController parse primary key failed")
		r.Code = common.IncorrectParams
		r.Message = "parse id failed"

		return
	}

	record, attlist, err := state.NftWithdrawClaimHandler(uint64(pid))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "PrimaryID": primaryID}).Error("NFTClaimController NftWithdrawClaimHandler failed")

		r.Code = common.WithdrawCode
		r.Message = "claim failed"
		return
	}

	if r.Code == common.SuccessCode {
		claimData := &dashdata.WithdrawClaimResponse{
			AppOrderID:        record.AppOrderID,
			Nonce:             record.Nonce,
			ContractAddress:   strings.ToLower(record.ContractAddress),
			TokenID:           record.TokenID,
			EquipmentID:       record.EquipmentID,
			TreaseAddress:     record.TreaseAddress,
			GameMinterAddress: record.GameMinterAddress,
			WithdrawAddress:   record.WithdrawAddress,
			Signature:         record.Signature,
			List:              attlist,
		}

		r.Data = claimData
		r.Message = "nft commit withdraw success"
	}
}

// NFTDepositController
func NFTDepositController(c *gin.Context) {
	r := &common.Response{
		Code: common.SuccessCode,
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	p := dashdata.DepositParam{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("NFTDepositController decode input request failed")
		r.Code = common.IncorrectParams
		r.Message = "invalid parameter"
		return
	}

	s := new(dashdata.MultiAssetStrategy)
	s.SetDepositImpl(&p)
	err = s.Deposit()

	if err != nil {
		if hpErr, ok := err.(*common.HpError); ok {
			r.Code = int64(hpErr.Code)
			r.Message = hpErr.CodeMsg()
		} else {
			r.Code = common.InnerError
			r.Message = err.Error()
		}
		return

	}
	claimData := &dashdata.DepositResponse{
		Trease: p.TargetAddress,
		Nonce:  p.Nonce,
	}

	r.Data = claimData
	r.Message = "nft deposit success"
}

func getChainAllNFTAssets(appID, page int, lastTokenId string, owner string, cacheData map[string]commdata.NFTContractCacheData) ([]response.NFTChainAssets, string, error) {
	graphFormat := `{
		lootAssetsEntities(first: %d, orderBy: tokenId, orderDirection: desc,where: {%s}) {
			tokenId
			tokenURI
			contract
		}
	}`

	whereCondition := "owner: \"" + owner + "\""
	if lastTokenId != "" {
		whereCondition += ", tokenId_lt:\"" + lastTokenId + "\""
	}

	graphString := fmt.Sprintf(graphFormat, page, whereCondition)
	graphResp := dashdata.RGraphLootAssets{}
	err := pool.GraphRequest(config.GetGameLootEquipmentGraph(), graphString, &graphResp)
	if err != nil {
		return nil, "", fmt.Errorf("request graph, %v", err)
	}

	result := make([]response.NFTChainAssets, 0)
	if len(graphResp.Assets) < 1 {
		return result, "", nil
	}

	tokenIDList := make(map[string][]string, 0)
	for _, v := range graphResp.Assets {
		list, ok := tokenIDList[v.Contract]
		if !ok {
			list = make([]string, 0)
		}

		list = append(list, v.TokenId)

		tokenIDList[v.Contract] = list
	}

	//contract address => contract cache data
	contractAddrMap := make(map[string]commdata.NFTContractCacheData, 0)
	for _, v := range cacheData {
		contractAddrMap[v.ContractAddress] = v
	}

	//get pass asset map
	passAssetMap := config.GetPassAssetConfig()

	//map: contractAddr => (tokenID => GameEquipment)
	doubleImageMap := make(map[string]map[string]model.TGameEquipment, 0)
	for contractAddr, tokeIDs := range tokenIDList {
		tokenIDMap, ok := doubleImageMap[contractAddr]
		if !ok {
			tokenIDMap = make(map[string]model.TGameEquipment, 0)
		}

		//check if pass asset or not
		info, ok := contractAddrMap[contractAddr]
		if !ok {
			// return nil, "", fmt.Errorf("contract not found, %v", err)
			continue
		}

		assetName := info.GameAssetName
		eid, isok := passAssetMap[assetName]
		if isok {
			item := model.TGameEquipment{
				GameAssetName:   assetName,
				EquipmentID:     eid,
				ImageURI:        fmt.Sprintf("%s/%s.png", info.BaseURL, eid),
				ContractAddress: contractAddr,
			}

			//all token id has the same game equipment data for pass asset
			for _, tokenid := range tokeIDs {
				tokenIDMap[tokenid] = item
			}

			doubleImageMap[contractAddr] = tokenIDMap

			continue
		}

		//other assets
		equipmentList := []model.TGameEquipment{}

		err := mysql.GetDB().Model(&model.TGameEquipment{}).Where("app_id = ? and contract_address = ? and token_id in ?", appID, contractAddr, tokeIDs).Find(&equipmentList).Error
		if err != nil {
			return nil, "", fmt.Errorf("query equipment, %v", err)
		}

		for _, val := range equipmentList {
			if val.GameAssetName == "" {
				val.GameAssetName = info.GameAssetName
			}
			tokenIDMap[val.TokenID] = val
		}

		doubleImageMap[contractAddr] = tokenIDMap
	}

	assetSortMap := config.GetAssetSortConfig()

	for _, v := range graphResp.Assets {
		data, ok := doubleImageMap[v.Contract][v.TokenId]
		if !ok {
			logger.Logrus.WithFields(logrus.Fields{"Data": v}).Error("equipment info of token id of contract not found")
			continue
		}

		equipImage := data.ImageURI

		finalChainURI := ""
		if data.GameAssetName == config.GetGameEquipmentName() {
			finalChainURI, err = comminfo.GetChainSVG(v.TokenURI, v.Contract, v.TokenId)
			if err != nil {
				return nil, "", fmt.Errorf("get chain svg failed, %v %v", v.TokenId, err)
			}
		}

		attrs := make([]commdata.EquipmentAttr, 0)

		_, isok := passAssetMap[data.GameAssetName]
		if !isok {
			if data.EquipmentAttr == nil {
				details, err := ingame.RequestGameNFTAssetDetail(appID, data.EquipmentID, false)
				if err != nil {
					return nil, "", fmt.Errorf("get asset details from game,%v", err)
				}

				attrs = details.Attrs
			} else {
				err := json.Unmarshal([]byte(data.EquipmentAttr), &attrs)
				if err != nil {
					return nil, "", fmt.Errorf("parse db attrs,%v", err)
				}
			}
		}

		info, ok := cacheData[data.GameAssetName]
		if !ok {
			return nil, "", fmt.Errorf("contract info not found for %v", data.GameAssetName)
		}

		sortIndex, ok := assetSortMap[data.GameAssetName]
		if !ok {
			sortIndex = 0
		}

		asset := response.NFTChainAssets{
			ChainImage:    finalChainURI,
			Image:         equipImage,
			Name:          info.TokenName,
			Contract:      v.Contract,
			Trease:        info.Treasure,
			TokenID:       v.TokenId,
			GameAssetName: info.GameAssetName,
			Attrs:         attrs,
			SortIndex:     sortIndex,
		}

		result = append(result, asset)
	}

	nextTokenId := ""
	if len(result) > 0 {
		nextTokenId = result[len(result)-1].TokenID
	}

	return result, nextTokenId, nil
}

func getChainNFTAssets(appID, page int, lastTokenId string, owner, contractAddr, treasure, assetName, tokenName, baseURL string) ([]response.NFTChainAssets, string, error) {
	graphFormat := `{
		lootAssetsEntities(first: %d, orderBy: tokenId, orderDirection: desc,where: {%s}) {
			tokenId
			tokenURI
			contract
		}
	}`

	whereCondition := fmt.Sprintf("contract:\"%s\",owner: \"%s\"", contractAddr, owner)
	if lastTokenId != "" {
		whereCondition += ", tokenId_lt:\"" + lastTokenId + "\""
	}

	graphString := fmt.Sprintf(graphFormat, page, whereCondition)
	graphResp := dashdata.RGraphLootAssets{}
	err := pool.GraphRequest(config.GetGameLootEquipmentGraph(), graphString, &graphResp)
	if err != nil {
		return nil, "", fmt.Errorf("request graph, %v", err)
	}

	result := make([]response.NFTChainAssets, 0)

	if len(graphResp.Assets) < 1 {
		return result, "", nil
	}

	tokenIDs := make([]string, 0)
	for _, e := range graphResp.Assets {
		tokenIDs = append(tokenIDs, e.TokenId)
	}

	tokenIDImageMap := make(map[string]model.TGameEquipment, 0)

	//check if pass asset or not
	passAssetMap := config.GetPassAssetConfig()

	eid, isok := passAssetMap[assetName]
	if isok {
		item := model.TGameEquipment{
			GameAssetName:   assetName,
			EquipmentID:     eid,
			ImageURI:        fmt.Sprintf("%s/%s.png", baseURL, eid),
			ContractAddress: contractAddr,
		}

		//all token id has the same game equipment data for pass asset
		for _, tokenid := range tokenIDs {
			tokenIDImageMap[tokenid] = item
		}
	} else {
		//other assets
		equipmentList := []model.TGameEquipment{}

		err := mysql.GetDB().Model(&model.TGameEquipment{}).Where("app_id = ? and contract_address = ? and token_id in ?", appID, contractAddr, tokenIDs).Find(&equipmentList).Error
		if err != nil {
			return nil, "", fmt.Errorf("query equipment, %v", err)
		}

		for _, val := range equipmentList {
			tokenIDImageMap[val.TokenID] = val
		}
	}

	for _, v := range graphResp.Assets {
		equip, ok := tokenIDImageMap[v.TokenId]
		equipImage := ""
		if !ok {
			return nil, "", fmt.Errorf("game image not found for token id, %v", v.TokenId)
		} else {
			equipImage = equip.ImageURI
		}

		finalChainURI := ""
		if assetName == config.GetGameEquipmentName() {
			finalChainURI, err = comminfo.GetChainSVG(v.TokenURI, v.Contract, v.TokenId)
			if err != nil {
				return nil, "", fmt.Errorf("get chain svg failed, %v %v", v.TokenId, err)
			}
		}

		attrs := make([]commdata.EquipmentAttr, 0)

		//pass has not asset detail
		_, isok := passAssetMap[assetName]
		if !isok {
			if equip.EquipmentAttr == nil {
				details, err := ingame.RequestGameNFTAssetDetail(appID, equip.EquipmentID, false)
				if err != nil {
					return nil, "", fmt.Errorf("get asset details from game,%v", err)
				}

				attrs = details.Attrs
			} else {
				err := json.Unmarshal([]byte(equip.EquipmentAttr), &attrs)
				if err != nil {
					return nil, "", fmt.Errorf("parse db attrs,%v", err)
				}
			}
		}

		asset := response.NFTChainAssets{
			ChainImage:    finalChainURI,
			Image:         equipImage,
			Name:          tokenName,
			Contract:      v.Contract,
			Trease:        treasure,
			TokenID:       v.TokenId,
			GameAssetName: assetName,
			Attrs:         attrs,
		}

		result = append(result, asset)
	}
	nextTokenId := ""
	if len(result) > 0 {
		nextTokenId = result[len(result)-1].TokenID
	}

	return result, nextTokenId, nil
}

/**
 * @Description: query equipment assets on chain
 * @author: shalom
 * @date: 2021/12/24 10:40 AM
 * @version: V1.0
 */
func QueryEquipment(c *gin.Context) {

	r := &common.Response{
		Code: common.SuccessCode,
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = request.QueryEquipment{}
	err := c.ShouldBindQuery(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("validator reject")

		r.Code = common.IncorrectParams
		r.Message = "invalid input parameters"
		return
	}

	chainIDStr := c.Request.Header.Get("Chainid")
	chainID, err := strconv.ParseInt(chainIDStr, 0, 64)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ChainID": chainIDStr}).Error("invalid chainID")

		r.Code = common.IncorrectParams
		r.Message = "invalid input chainID"
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p, "ChainID": chainID}).Info("QueryEquipment input info")

	nftContractMap, err := comminfo.GetNftContractByAppIDAndChainID(int(p.GameID), int(chainID))
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Logrus.WithFields(logrus.Fields{
			"ErrMsg": err.Error(),
			"GameID": p.GameID,
		}).Error("QueryEquipment GetNftContractByAppIDAndChainID error")

		r.Code = common.InnerError
		r.Message = "nft contract lists get failed"
		return
	}

	if len(nftContractMap) == 0 {
		logger.Logrus.WithFields(logrus.Fields{
			"Data": p,
		}).Info("QueryEquipment GetNftContractByAppIDAndChainID contract not config")

		r.Code = common.SuccessCode
		r.Data = response.QueryEquipmentResp{
			LastKey:   "",
			NFTAssets: make([]response.NFTChainAssets, 0),
		}
		r.Message = "nft contract not config"
		return
	}

	page := int(p.PageSize)
	lastTokenId := p.LastKey
	owner := strings.ToLower(p.Owner)

	result := make([]response.NFTChainAssets, 0)

	//query all nft assets
	if p.AssetName == "" {
		list, lastKey, err := getChainAllNFTAssets(int(p.GameID), page, lastTokenId, owner, nftContractMap)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{
				"Data":   p,
				"ErrMsg": err.Error(),
			}).Error("QueryEquipment get all nft assets on chain error")

			r.Code = common.InnerError
			r.Message = "get all nft assets on chain failed"
			return
		}

		result = append(result, list...)

		sort.SliceStable(result, func(i, j int) bool {
			return result[i].SortIndex > result[j].SortIndex
		})

		r.Data = response.QueryEquipmentResp{
			LastKey:   lastKey,
			NFTAssets: result,
		}
		r.Message = "get equipment nft assets on chain success"
		return
	}

	//query single nft asset
	v, ok := nftContractMap[p.AssetName]
	if !ok {
		logger.Logrus.WithFields(logrus.Fields{
			"Data": p,
		}).Error("QueryEquipment get contract info error")

		r.Code = common.InnerError
		r.Message = "get contract info failed"
		return
	}

	list, lastKey, err := getChainNFTAssets(int(p.GameID), page, lastTokenId, owner, v.ContractAddress, v.Treasure, v.GameAssetName, v.TokenName, v.BaseURL)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{
			"Data":         p,
			"ContractInfo": v,
			"ErrMsg":       err.Error(),
		}).Error("QueryEquipment get chain nft assets error")

		r.Code = common.InnerError
		r.Message = "get nft assets on chain failed"
		return
	}

	result = append(result, list...)

	r.Data = response.QueryEquipmentResp{
		LastKey:   lastKey,
		NFTAssets: result,
	}
	r.Message = "get nft assets on chain success"
}
