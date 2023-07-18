package marketplace

import (
	"encoding/json"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/pool"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func getOrderList(pools []RGraphOrderListPool, contractInfo *commdata.NFTContractCacheData) ([]SingleOrderList, error) {
	tokenIDList := make([]string, 0)
	for _, val := range pools {
		tokenIDList = append(tokenIDList, val.TokenId)
	}

	tokenIDEquipMap := make(map[string]model.TGameEquipment, 0)

	//check if pass asset or not
	passAssetMap := config.GetPassAssetConfig()
	passEID, isok := passAssetMap[contractInfo.GameAssetName]
	if isok {
		for _, tokenid := range tokenIDList {
			tokenIDEquipMap[tokenid] = model.TGameEquipment{
				GameAssetName:   contractInfo.GameAssetName,
				ImageURI:        fmt.Sprintf("%s/%s.png", contractInfo.BaseURL, passEID),
				EquipmentID:     passEID,
				ContractAddress: contractInfo.ContractAddress,
			}
		}
	} else {
		equipmentList := []model.TGameEquipment{}
		err := mysql.GetDB().Model(&model.TGameEquipment{}).Where("contract_address = ? AND token_id in ?", contractInfo.ContractAddress, tokenIDList).Find(&equipmentList).Error
		if err != nil {
			return nil, fmt.Errorf("query db failed, %v", err)
		}

		for _, e := range equipmentList {
			tokenIDEquipMap[e.TokenID] = e
		}
	}

	chainToken := config.GetChainToken()
	chainID := contractInfo.ChainID
	purchaseTokenInfo, ok := chainToken[fmt.Sprintf("%d", chainID)]
	if !ok || purchaseTokenInfo.TokenName == "" || purchaseTokenInfo.Decimal == 0 {
		return nil, fmt.Errorf("purches info not found fot %v", chainID)
	}

	//get chain uri ,only for ArchLootPart
	chainURIMap := map[string]string{}
	if contractInfo.GameAssetName == config.GetGameEquipmentName() {
		graphChainURIFormat := `{
					lootAssetsEntities(where:{contract:"%s"%s}) {
					id
					tokenId
					tokenURI
					contract
					}
				}
			`

		tokenIDFilter := strings.Join(tokenIDList, ",")
		finalTokenIDFilter := fmt.Sprintf(",tokenId_in:[%s]", tokenIDFilter)

		graphChainURIString := fmt.Sprintf(graphChainURIFormat, contractInfo.ContractAddress, finalTokenIDFilter)
		graphChainURIRet := RGraphTokenURI{}
		err := pool.GraphRequest(config.GetGameLootEquipmentGraph(), graphChainURIString, &graphChainURIRet)
		if err != nil {
			return nil, fmt.Errorf("graph request: %v", err)
		}

		for _, e := range graphChainURIRet.LootAssetsEntities {
			finalChainURI, err := comminfo.GetChainSVG(e.TokenURI, e.Contract, e.TokenId)
			if err != nil {
				return nil, fmt.Errorf("get chain svg failed, %v %v", e.TokenId, err)
			}

			chainURIMap[e.TokenId] = finalChainURI
		}
	}

	assetSortMap := config.GetAssetSortConfig()

	result := make([]SingleOrderList, 0)
	for _, e := range pools {
		index, err := strconv.Atoi(e.Index)
		if err != nil {
			return nil, fmt.Errorf("parse pool index failed, %v %v", e.Index, err)
		}

		originBalance, isSuccess := big.NewInt(0).SetString(e.AmountTotal1, 0)
		if !isSuccess {
			return nil, fmt.Errorf("format wrong, %v", e.AmountTotal1)
		}

		balance, err := tools.GetTokenAmount(originBalance, int32(purchaseTokenInfo.Decimal), 8)
		if err != nil {
			return nil, fmt.Errorf("convert decimal ,%v", err)
		}

		chainURI, ok := chainURIMap[e.TokenId]
		if !ok {
			chainURI = ""
		}

		equip, ok := tokenIDEquipMap[e.TokenId]
		if !ok {
			return nil, fmt.Errorf("get game image from db failed, %v, %v", e.TokenId, contractInfo.ContractAddress)
		}

		attrs := make([]commdata.EquipmentAttr, 0)
		if equip.EquipmentAttr != nil {
			err := json.Unmarshal([]byte(equip.EquipmentAttr), &attrs)
			if err != nil {
				return nil, fmt.Errorf("unmarshal attrs failed, %v", err)
			}
		} else {
			if contractInfo.GameAssetName == config.GetGameEquipmentName() {
				nftDetail, err := ingame.RequestGameNFTAssetDetail(int(equip.AppID), equip.EquipmentID, false)
				if err != nil {
					return nil, fmt.Errorf("get equipment attrs failed, %s", equip.EquipmentID)
				}

				attrs = nftDetail.Attrs
			}
		}

		sortIndex, ok := assetSortMap[contractInfo.GameAssetName]
		if !ok {
			sortIndex = 0
		}

		//TODO need has a map for other sold nft  contract to show tokenName (saved such as config file dont save with nftContract table)
		item := SingleOrderList{
			Index:         uint64(index),
			ChainImage:    chainURI,
			Image:         equip.ImageURI,
			TokenName:     contractInfo.TokenName,
			TokenID:       e.TokenId,
			Price:         balance,
			PurchaseToken: purchaseTokenInfo.TokenName,
			AssetName:     contractInfo.GameAssetName,
			Attrs:         attrs,
			SortIndex:     sortIndex,
		}
		result = append(result, item)
	}

	return result, nil
}

func getAllMarketOrder(appID, pageSize int, lastKey string, priceSort int, creatFilter, closedFilter string) ([]SingleOrderList, string, error) {
	var orderDir string
	orderByFilter := "amountTotal1Key"
	sortedKeyField := "AmountTotal1Key"
	var sortedCond string
	if priceSort == const_def.GRAPH_ORDER_ASC {
		orderDir = "asc"
		sortedCond = fmt.Sprintf(",amountTotal1Key_gt:\"%s\"", lastKey)
	} else if priceSort == const_def.GRAPH_ORDER_DESC {
		orderDir = "desc"
		sortedCond = fmt.Sprintf(",amountTotal1Key_lt:\"%s\"", lastKey)
	} else {
		orderDir = "desc"
		orderByFilter = "createdTimeKey"
		sortedCond = fmt.Sprintf(",createdTimeKey_lt:\"%s\"", lastKey)
		sortedKeyField = "CreatedTimeKey"
	}
	if lastKey == "" {
		sortedCond = ""
	}
	format := `{
		pools(first: %d,  orderBy: %s, orderDirection: %s, where: {%s closed: false, redeemed: false, swapped: false %s %s %s}) {
		  id
		  index
		  token0
		  tokenId
		  amountTotal1
		  amountTotal1Key
		  createdTimeKey
		}
	  }
	  `

	conMap, err := comminfo.GetNftContractByAppID(appID)
	if err != nil {
		return nil, "", fmt.Errorf("contract list: %v", err)
	}

	contractList := make([]string, 0)
	for _, v := range conMap {
		contractList = append(contractList, strconv.Quote(v.ContractAddress))
	}

	contractStr := strings.Join(contractList, ",")
	contractFilter := fmt.Sprintf(",token0_in:[%s]", contractStr)

	graphString := fmt.Sprintf(format, pageSize, orderByFilter, orderDir, creatFilter, closedFilter, contractFilter, sortedCond)

	graphRet := RGraphOrderListPools{}
	err = pool.GraphRequest(config.GetMarketplaceGraph(), graphString, &graphRet)
	if err != nil {
		return nil, "", fmt.Errorf("graph request: %v", err)
	}

	contractPoolMap := make(map[string][]RGraphOrderListPool, 0)
	keyField := ""
	for k, v := range graphRet.Pools {
		contractAddr := v.NFTContract
		item, ok := contractPoolMap[contractAddr]
		if !ok {
			item = make([]RGraphOrderListPool, 0)
		}

		item = append(item, v)
		contractPoolMap[contractAddr] = item
		if k == (len(graphRet.Pools) - 1) {
			r := reflect.ValueOf(v)
			keyField = reflect.Indirect(r).FieldByName(sortedKeyField).String()
		}
	}

	result := make([]SingleOrderList, 0)
	for addr, list := range contractPoolMap {
		tokenIDList := make([]string, 0)
		for _, val := range list {
			tokenIDList = append(tokenIDList, val.TokenId)
		}

		if len(tokenIDList) < 1 {
			continue
		}

		contractInfo, err := comminfo.GetNftContractByContract(appID, addr)
		if err != nil {
			return nil, "", err
		}

		item, err := getOrderList(list, contractInfo)
		if err != nil {
			return nil, "", err
		}

		result = append(result, item...)
	}

	return result, keyField, nil
}

func getMarketOrder(pageSize int, lastKey string, priceSort int, creatFilter, contractAddr, closedFilter string, decimal int32, nftTokenName, purchaseTokenName, assetName, baseURL string) ([]SingleOrderList, string, error) {
	var orderDir string
	orderByFilter := "amountTotal1Key"
	sortedKeyField := "AmountTotal1Key"
	var sortedCond string
	if priceSort == const_def.GRAPH_ORDER_ASC {
		orderDir = "asc"
		sortedCond = fmt.Sprintf(",amountTotal1Key_gt:\"%s\"", lastKey)
	} else if priceSort == const_def.GRAPH_ORDER_DESC {
		orderDir = "desc"
		sortedCond = fmt.Sprintf(",amountTotal1Key_lt:\"%s\"", lastKey)
	} else {
		orderDir = "desc"
		orderByFilter = "createdTimeKey"
		sortedCond = fmt.Sprintf(",createdTimeKey_lt:\"%s\"", lastKey)
		sortedKeyField = "CreatedTimeKey"
	}
	if lastKey == "" {
		sortedCond = ""
	}
	format := `{
		pools(first: %d, orderBy: %s, orderDirection: %s, where: {%s closed: false,token0: "%s", token1:"0x0000000000000000000000000000000000000000", redeemed: false, swapped: false %s %s}) {
		  id
		  index
		  token0
		  tokenId
		  amountTotal1
		  amountTotal1Key
		  createdTimeKey
		}
	  }
	  `

	graphString := fmt.Sprintf(format, pageSize, orderByFilter, orderDir, creatFilter, contractAddr, closedFilter, sortedCond)

	graphRet := RGraphOrderListPools{}
	err := pool.GraphRequest(config.GetMarketplaceGraph(), graphString, &graphRet)
	if err != nil {
		return nil, "", fmt.Errorf("graph request: %v", err)
	}

	tokenIDFilter := ""
	tokenIDs := make([]string, 0)
	result := make([]SingleOrderList, 0)
	for _, v := range graphRet.Pools {
		if v.TokenId != "" {
			tokenIDFilter += fmt.Sprintf("%s,", v.TokenId)
			tokenIDs = append(tokenIDs, v.TokenId)
		}
	}

	if len(tokenIDFilter) == 0 {
		return result, "", nil
	}

	//get chain uri
	chainURIMap := map[string]string{}
	if assetName == config.GetGameEquipmentName() {
		graphChainURIFormat := `{
			lootAssetsEntities(where:{contract:"%s"%s}) {
			id
			tokenId
			tokenURI
			contract
			}
		}
		`

		finalTokenIDFilter := ""
		if len(tokenIDFilter) > 0 {
			tokenIDFilter = tokenIDFilter[:len(tokenIDFilter)-1]
			finalTokenIDFilter = fmt.Sprintf(",tokenId_in:[%s]", tokenIDFilter)
		}

		graphChainURIString := fmt.Sprintf(graphChainURIFormat, contractAddr, finalTokenIDFilter)
		graphChainURIRet := RGraphTokenURI{}
		err = pool.GraphRequest(config.GetGameLootEquipmentGraph(), graphChainURIString, &graphChainURIRet)
		if err != nil {
			return nil, "", fmt.Errorf("graph request: %v", err)
		}

		for _, e := range graphChainURIRet.LootAssetsEntities {
			finalChainURI, err := comminfo.GetChainSVG(e.TokenURI, e.Contract, e.TokenId)
			if err != nil {
				return nil, "", fmt.Errorf("get chain svg failed, %v %v", e.TokenId, err)
			}

			chainURIMap[e.TokenId] = finalChainURI
		}
	}

	tokenIDEquipMap := make(map[string]model.TGameEquipment, 0)

	//check if pass asset or not
	passAssetMap := config.GetPassAssetConfig()
	passEID, isok := passAssetMap[assetName]
	if isok {
		for _, tokenid := range tokenIDs {
			tokenIDEquipMap[tokenid] = model.TGameEquipment{
				GameAssetName:   assetName,
				ImageURI:        fmt.Sprintf("%s/%s.png", baseURL, passEID),
				EquipmentID:     passEID,
				ContractAddress: contractAddr,
			}
		}
	} else {
		equipmentList := []model.TGameEquipment{}
		err = mysql.GetDB().Model(&model.TGameEquipment{}).Where("contract_address = ? AND token_id in ?", contractAddr, tokenIDs).Find(&equipmentList).Error
		if err != nil {
			return nil, "", fmt.Errorf("query db failed, %v", err)
		}

		for _, e := range equipmentList {
			tokenIDEquipMap[e.TokenID] = e
		}
	}

	keyField := ""
	for k, e := range graphRet.Pools {
		index, err := strconv.Atoi(e.Index)
		if err != nil {
			return nil, "", fmt.Errorf("parse pool index failed, %v %v", e.Index, err)
		}

		originBalance, isSuccess := big.NewInt(0).SetString(e.AmountTotal1, 0)
		if !isSuccess {
			return nil, "", fmt.Errorf("format wrong, %v", e.AmountTotal1)
		}

		balance, err := tools.GetTokenAmount(originBalance, decimal, 8)
		if err != nil {
			return nil, "", fmt.Errorf("convert decimal ,%v", err)
		}

		chainURI, ok := chainURIMap[e.TokenId]
		if !ok {
			chainURI = ""
		}

		equip, ok := tokenIDEquipMap[e.TokenId]
		if !ok {
			return nil, "", fmt.Errorf("get game image from db failed, %v, %v", e.TokenId, contractAddr)
		}

		attrs := make([]commdata.EquipmentAttr, 0)
		if equip.EquipmentAttr != nil {
			err := json.Unmarshal([]byte(equip.EquipmentAttr), &attrs)
			if err != nil {
				return nil, "", fmt.Errorf("unmarshal attrs failed, %v", err)
			}
		} else {
			if assetName == config.GetGameEquipmentName() {
				nftDetail, err := ingame.RequestGameNFTAssetDetail(int(equip.AppID), equip.EquipmentID, false)
				if err != nil {
					return nil, "", fmt.Errorf("get equipment attrs failed, %s", equip.EquipmentID)
				}

				attrs = nftDetail.Attrs
			}
		}

		//TODO need has a map for other sold nft  contract to show tokenName (saved such as config file dont save with nftContract table)
		item := SingleOrderList{
			Index:         uint64(index),
			ChainImage:    chainURI,
			Image:         equip.ImageURI,
			TokenName:     nftTokenName,
			TokenID:       e.TokenId,
			Price:         balance,
			PurchaseToken: purchaseTokenName,
			AssetName:     assetName,
			Attrs:         attrs,
		}
		result = append(result, item)
		if k == (len(graphRet.Pools) - 1) {
			r := reflect.ValueOf(e)
			keyField = reflect.Indirect(r).FieldByName(sortedKeyField).String()
		}
	}

	return result, keyField, nil
}

func getActivitySignalAssetName(pageSize int, lastKey string, typeFilter, user, contractAddr, name, symbol string, decimal int, assetName, baseURL string) ([]SingleActivity, string, error) {
	format := `{
		activities(first: %d, orderBy: timeKey, orderDirection: desc,where:{%sactivityer: "%s",token0: "%s",token1:"0x0000000000000000000000000000000000000000" %s}) {
			id
			activityType
			token0
			amountTotal1
			tokenId
			activityer
			creator
			buyer
			time
			hash
			timeKey
		}
	  }
	  `
	sortedCond := ""
	if lastKey != "" {
		sortedCond = fmt.Sprintf(",timeKey_lt:\"%s\"", lastKey)
	}

	graphString := fmt.Sprintf(format, pageSize, typeFilter, user, contractAddr, sortedCond)

	graphRet := RGraphActivities{}
	err := pool.GraphRequest(config.GetMarketplaceGraph(), graphString, &graphRet)
	if err != nil {
		return nil, "", err
	}

	result := make([]SingleActivity, 0)
	if len(graphRet.Activities) < 1 {
		return result, "", nil
	}

	tokenIDList := make([]string, 0)
	keyField := ""
	for k, v := range graphRet.Activities {
		tokenIDList = append(tokenIDList, v.TokenId)
		if k == (len(graphRet.Activities) - 1) {
			keyField = v.TimeKey
		}
	}

	tokenIDEquipMap := make(map[string]model.TGameEquipment, 0)

	//check if pass asset or not
	passAssetMap := config.GetPassAssetConfig()
	passEID, isok := passAssetMap[assetName]
	if isok {
		for _, tokenid := range tokenIDList {
			tokenIDEquipMap[tokenid] = model.TGameEquipment{
				GameAssetName:   assetName,
				ImageURI:        fmt.Sprintf("%s/%s.png", baseURL, passEID),
				EquipmentID:     passEID,
				ContractAddress: contractAddr,
			}
		}
	} else {
		equipmentList := []model.TGameEquipment{}
		err = mysql.GetDB().Model(&model.TGameEquipment{}).Where("contract_address = ? AND token_id in ?", contractAddr, tokenIDList).Find(&equipmentList).Error
		if err != nil {
			return nil, "", fmt.Errorf("query db failed, %v", err)
		}

		for _, e := range equipmentList {
			tokenIDEquipMap[e.TokenID] = e
		}
	}

	for _, e := range graphRet.Activities {
		activityType, err := strconv.Atoi(e.ActivityType)
		if err != nil {
			return nil, "", err
		}

		originBalance, _ := big.NewInt(0).SetString(e.AmountTotal1, 0)
		balance, err := tools.GetTokenAmount(originBalance, int32(decimal), 8)
		if err != nil {
			return nil, "", err
		}

		to := ""
		if activityType == const_def.MARKETPLACE_PURCHASE || activityType == const_def.MARKETPLACE_SALE {
			to = e.Buyer
		}

		gameEquip, ok := tokenIDEquipMap[e.TokenId]
		if !ok {
			logger.Logrus.WithFields(logrus.Fields{"Data": e}).Error("market activity for asset name and get image url failed")
			continue
		}

		item := SingleActivity{
			ActivityType: activityType,
			TokenName:    name,
			TokenSymbol:  symbol,
			Price:        balance,
			TokenID:      e.TokenId,
			From:         e.Creator,
			To:           to,
			Time:         e.Time,
			Hash:         e.Hash,
			GameImageURL: gameEquip.ImageURI,
		}

		result = append(result, item)
	}

	return result, keyField, nil
}

func getAllActivitys(appID, pageSize int, lastKey string, typeFilter, user string) ([]SingleActivity, string, error) {
	format := `{
		activities(first: %d, orderBy: timeKey, orderDirection: desc,where:{%sactivityer: "%s",token1:"0x0000000000000000000000000000000000000000" %s %s}) {
			id
			activityType
			token0
			amountTotal1
			tokenId
			activityer
			creator
			buyer
			time
			hash
			timeKey
		}
	  }
	  `
	sortedCond := ""
	if lastKey != "" {
		sortedCond = fmt.Sprintf(",timeKey_lt:\"%s\"", lastKey)
	}

	conMap, err := comminfo.GetNftContractByAppID(appID)
	if err != nil {
		return nil, "", fmt.Errorf("contract list: %v", err)
	}

	contractList := make([]string, 0)
	for _, v := range conMap {
		contractList = append(contractList, strconv.Quote(v.ContractAddress))
	}

	contractStrs := strings.Join(contractList, ",")
	contractFilter := fmt.Sprintf(",token0_in:[%s]", contractStrs)

	graphString := fmt.Sprintf(format, pageSize, typeFilter, user, contractFilter, sortedCond)

	graphRet := RGraphActivities{}
	err = pool.GraphRequest(config.GetMarketplaceGraph(), graphString, &graphRet)
	if err != nil {
		return nil, "", err
	}

	result := make([]SingleActivity, 0)
	if len(graphRet.Activities) < 1 {
		return result, "", nil
	}

	contractActsMap := make(map[string][]RGraphSingleActivity, 0)
	keyField := ""
	for k, v := range graphRet.Activities {
		contractAddr := v.Token0
		list, ok := contractActsMap[contractAddr]
		if !ok {
			list = make([]RGraphSingleActivity, 0)
		}

		list = append(list, v)
		contractActsMap[contractAddr] = list
		if k == (len(graphRet.Activities) - 1) {
			keyField = v.TimeKey
		}
	}

	for contractAddr, list := range contractActsMap {
		if len(list) < 1 {
			continue
		}

		tokenIDList := make([]string, 0)
		for _, v := range list {
			tokenIDList = append(tokenIDList, v.TokenId)
		}

		contractInfo, err := comminfo.GetNftContractByContract(appID, contractAddr)
		if err != nil {
			return nil, "", err
		}

		tokenIDEquipMap := make(map[string]model.TGameEquipment, 0)

		//check if pass asset or not
		passAssetMap := config.GetPassAssetConfig()
		passEID, isok := passAssetMap[contractInfo.GameAssetName]
		if isok {
			for _, tokenid := range tokenIDList {
				tokenIDEquipMap[tokenid] = model.TGameEquipment{
					GameAssetName:   contractInfo.GameAssetName,
					ImageURI:        fmt.Sprintf("%s/%s.png", contractInfo.BaseURL, passEID),
					EquipmentID:     passEID,
					ContractAddress: contractAddr,
				}
			}
		} else {
			equipmentList := []model.TGameEquipment{}
			err = mysql.GetDB().Model(&model.TGameEquipment{}).Where("contract_address = ? AND token_id in ?", contractAddr, tokenIDList).Find(&equipmentList).Error
			if err != nil {
				return nil, "", fmt.Errorf("query db failed, %v", err)
			}

			for _, e := range equipmentList {
				tokenIDEquipMap[e.TokenID] = e
			}
		}

		chainToken := config.GetChainToken()
		purchaseTokenInfo, ok := chainToken[fmt.Sprintf("%d", contractInfo.ChainID)]
		if !ok || purchaseTokenInfo.TokenName == "" || purchaseTokenInfo.Decimal == 0 {
			return nil, "", fmt.Errorf("purchase token chain id not found for %s", contractInfo.ContractAddress)
		}

		for _, e := range list {
			activityType, err := strconv.Atoi(e.ActivityType)
			if err != nil {
				return nil, "", err
			}

			originBalance, _ := big.NewInt(0).SetString(e.AmountTotal1, 0)
			balance, err := tools.GetTokenAmount(originBalance, int32(purchaseTokenInfo.Decimal), 8)
			if err != nil {
				return nil, "", err
			}

			to := ""
			if activityType == const_def.MARKETPLACE_PURCHASE || activityType == const_def.MARKETPLACE_SALE {
				to = e.Buyer
			}

			gameEquip, ok := tokenIDEquipMap[e.TokenId]
			if !ok {
				logger.Logrus.WithFields(logrus.Fields{"Data": e}).Error("market activity get image url failed")
				continue
			}

			item := SingleActivity{
				ActivityType: activityType,
				TokenName:    contractInfo.TokenName,
				TokenSymbol:  contractInfo.TokenSymbol,
				Price:        balance,
				TokenID:      e.TokenId,
				From:         e.Creator,
				To:           to,
				Time:         e.Time,
				Hash:         e.Hash,
				GameImageURL: gameEquip.ImageURI,
			}

			result = append(result, item)
		}
	}

	return result, keyField, nil
}
