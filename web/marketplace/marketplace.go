package marketplace

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
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/pool"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"math/big"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// !!!  if has another nft can sold on markeplace need edit this
func MarketOrderList(c *gin.Context) {
	r := common.Response{
		Code: common.SuccessCode,
	}

	defer func() {
		c.JSON(http.StatusOK, r)
	}()

	var p = POrderList{}
	err := c.ShouldBindQuery(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("MarketOrderList validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("MarketOrderList info")

	if p.PriceSort != 0 && p.PriceSort != const_def.GRAPH_ORDER_ASC && p.PriceSort != const_def.GRAPH_ORDER_DESC {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("MarketOrderList invalid price sort")

		r.Code = common.InnerError
		r.Message = "invalid price sort"
		return
	}

	if p.Creator != "" && !tools.CheckEthAddr(p.Creator) {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("MarketOrderList invalid creator")

		r.Code = common.InnerError
		r.Message = "invalid creator"
		return
	}

	url := config.GetMarketplaceGraph()
	if url == "" {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("MarketOrderList marketplace graph is emtpy")
		r.Code = common.InnerError
		r.Message = "invalid marketplace graph url"
		return
	}

	now := time.Now().Unix()

	creatorfilter := ""
	closedFilter := ""

	if p.Creator != "" {
		creatorfilter = fmt.Sprintf(`creator: "%s",`, strings.ToLower(p.Creator))
	} else {
		closedFilter = fmt.Sprintf(", closeAt_gt: %d", now)
	}

	//query all marketplace order list
	if p.AssetName == "" {
		orderList, lastKey, err := getAllMarketOrder(int(p.GameID), int(p.PageSize), p.LastKey, p.PriceSort, creatorfilter, closedFilter)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "Data": p}).Error("MarketOrderList getAllMarketOrder failed")
			r.Code = common.InnerError
			r.Message = "get all nft market order list failed"
			return
		}

		sort.SliceStable(orderList, func(i, j int) bool {
			return orderList[i].SortIndex > orderList[j].SortIndex
		})

		r.Data = ROrderList{
			LastKey:   lastKey,
			OrderList: orderList,
		}
		r.Message = "marketplace order list"
		return
	}

	//query single nft matketplact order list
	nftContractMap, err := comminfo.GetNftContractByAppID(int(p.GameID))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "Data": p}).Error("MarketOrderList GetNftContractByAppID failed")
		r.Code = common.InnerError
		r.Message = "get nft contract info failed"
		return
	}

	v, ok := nftContractMap[p.AssetName]
	if !ok {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("MarketOrderList contract info not found")

		r.Code = common.InnerError
		r.Message = "chain id not config"
		return
	}

	chainToken := config.GetChainToken()
	purchaseTokenInfo, ok := chainToken[fmt.Sprintf("%d", v.ChainID)]
	if !ok || purchaseTokenInfo.TokenName == "" || purchaseTokenInfo.Decimal == 0 {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "Data": p, "ContractInfo": v}).Error("MarketOrderList chain id not found")

		r.Code = common.InnerError
		r.Message = "chain id not config"
		return
	}

	orderItem, lastKey, err := getMarketOrder(int(p.PageSize), p.LastKey, p.PriceSort, creatorfilter, v.ContractAddress, closedFilter, int32(purchaseTokenInfo.Decimal), v.TokenName, purchaseTokenInfo.TokenName, v.GameAssetName, v.BaseURL)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "Data": p, "ContractInfo": v}).Error("MarketOrderList get order list failed")

		r.Code = common.InnerError
		r.Message = "get market order list failed"
		return
	}

	if r.Code == common.SuccessCode {
		r.Data = ROrderList{
			LastKey:   lastKey,
			OrderList: orderItem,
		}
		r.Message = "get marketplace order list success"
	}
}

// !!! if has another nft contract. need edit this api
func OrderDetail(c *gin.Context) {
	r := common.Response{
		Code: common.SuccessCode,
	}

	defer func() {
		c.JSON(http.StatusOK, r)
	}()

	ret := ROrderDetail{}

	var p = POrderDetail{}
	err := c.ShouldBindQuery(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("OrderDetail validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("OrderDetail input info")

	if p.Owned {
		if p.OwnedContractAddress == "" || !tools.CheckEthAddr(p.OwnedContractAddress) {
			logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("OrderDetail owned need contract address")
			r.Code = common.IncorrectParams
			r.Message = "owned contract address is invalid"
			return
		}
	}
	url := config.GetMarketplaceGraph()
	if url == "" {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("OrderDetail marketplace graph url is null")
		r.Code = common.InnerError
		r.Message = "marketplace graph url is null"
		return
	}

	marketplaceContract := strings.ToLower(config.GetMarketplaceContract())
	if marketplaceContract == "" || !tools.CheckEthAddr(marketplaceContract) {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("OrderDetail marketplace address is null")
		r.Code = common.InnerError
		r.Message = "marketplace address is null"
		return
	}

	nftContractInfo, err := comminfo.GetNftContractByAssetName(int(p.GameID), p.AssetName)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"Err": err, "Data": p}).Error("OrderDetail GetNftContractByAssetName failed")
		r.Code = common.InnerError
		r.Message = "contract info get failed"
		return
	}

	chainToken := config.GetChainToken()
	purchaseTokenInfo, ok := chainToken[fmt.Sprintf("%d", nftContractInfo.ChainID)]
	if !ok || purchaseTokenInfo.TokenName == "" || purchaseTokenInfo.Decimal == 0 {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("OrderDetail token config not found")
		r.Code = common.InnerError
		r.Message = "token info not found"
		return
	}

	network, ok := config.GetChainNetwork()[int64(nftContractInfo.ChainID)]
	if !ok {
		logger.Logrus.Error("OrderDetail not found network config")
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("OrderDetail network config not found")
		r.Code = common.InnerError
		r.Message = "network info not found"
		return
	}

	//get attr&image from database

	nattrs := make([]commdata.EquipmentAttr, 0)
	equipImage := ""

	finalChainURI := ""

	//check if pass asset or not
	passAssetMap := config.GetPassAssetConfig()
	passEID, isok := passAssetMap[p.AssetName]
	if isok {
		equipImage = fmt.Sprintf("%s/%s.png", nftContractInfo.BaseURL, passEID)
	} else {
		equipmentInfo := model.TGameEquipment{}
		condition := map[string]interface{}{"app_id": p.GameID, "token_id": p.TokenID, "contract_address": strings.ToLower(nftContractInfo.ContractAddress)}
		err, found := mysql.WrapFindFirst(model.TableGameEquipment, &equipmentInfo, condition)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "Data": p}).Error("OrderDetail QueryGameAssetDetail get equipmenID by tokenID failed")
			r.Code = common.InnerError
			r.Message = "equipment id get failed"
			return
		}

		if found {
			equipImage = equipmentInfo.ImageURI
			if equipmentInfo.EquipmentAttr != nil {
				err = json.Unmarshal([]byte(equipmentInfo.EquipmentAttr), &nattrs)
				if err != nil {
					logger.Logrus.WithFields(logrus.Fields{
						"ErrMsg":        err.Error(),
						"tokenID":       p.TokenID,
						"EquipmentAttr": equipmentInfo.EquipmentAttr,
					}).Error("OrderDetail get attrs by tokenid from db failed")

					r.Code = common.InnerError
					r.Message = "unmarshal attrs from db failed"
					return
				}
			} else if equipmentInfo.EquipmentID != "" {
				details, err := ingame.RequestGameNFTAssetDetail(equipmentInfo.AppID, equipmentInfo.EquipmentID, false)
				if err != nil {
					logger.Logrus.WithFields(logrus.Fields{
						"ErrMsg":      err.Error(),
						"tokenID":     p.TokenID,
						"EquipmentID": equipmentInfo.EquipmentID,
					}).Error("OrderDetail RequestGameNFTAssetDetail failed")

					r.Code = common.InnerError
					r.Message = "game attrs get failed"
					return
				}

				nattrs = details.Attrs
			} else {
				logger.Logrus.WithFields(logrus.Fields{
					"tokenID":     p.TokenID,
					"EquipmentID": equipmentInfo.EquipmentID,
				}).Error("OrderDetail game equipment record not found")
				r.Code = common.InnerError
				r.Message = "game equipment record not found"
				return
			}
		}

		//get uri from chain
		finalChainURI, err = comminfo.GetChainSVG("", nftContractInfo.ContractAddress, p.TokenID)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"tokenID": p.TokenID, "ErrMsg": err}).Error("OrderDetail get chain uri from contract failed")
			r.Code = common.InnerError
			r.Message = "uri from contract get failed"
			return
		}
	}

	if p.Owned {
		err, approved := contracts.GetERC721ApprovedForAll(p.OwnedContractAddress, p.UserAddress, marketplaceContract)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"err": err}).Error("OrderDetail QueryApproved Failed")
			r.Code = common.InnerError
			r.Message = "query approved failed"
			return
		}
		retData := ROrderDetail{
			Index:               0,
			TokenID:             p.TokenID,
			Price:               "0",
			OriginPrice:         "0",
			PurchaseToken:       purchaseTokenInfo.TokenName,
			Status:              const_def.MARKETPLACE_ORDER_SELL,
			ContractAddress:     strings.ToLower(p.OwnedContractAddress),
			Attrs:               nattrs,
			ChainImage:          finalChainURI,
			Image:               equipImage,
			Name:                nftContractInfo.TokenName,
			Symbol:              nftContractInfo.TokenSymbol,
			Standard:            "ERC721",
			Network:             network,
			Approved:            approved,
			MarketplaceContract: marketplaceContract,
			Creator:             "",
		}
		ret = retData
	} else {
		index := fmt.Sprintf("%x", p.Index)

		format := `{
				pool(id: "0x%s") {
				  id
				  index
				  token0
				  tokenId
				  amountTotal1
				  closeAt
				  swapped
				  redeemed
				  closed
				  creator
				}
			  }
			  `
		graphString := fmt.Sprintf(format, index)

		graphRet := RGraphOrderDetailPool{}
		err = pool.GraphRequest(config.GetMarketplaceGraph(), graphString, &graphRet)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"err": err, "GraphQueryString": graphString}).Error("OrderDetail graph request failed")
			r.Code = common.InnerError
			r.Message = "pool graph query failed"
			return
		}

		if graphRet.Pool.TokenId != p.TokenID {
			logger.Logrus.WithFields(logrus.Fields{"tokenID": p.TokenID, "poolTokenID": graphRet.Pool.TokenId}).Error("OrderDetail tokenid not match")
			r.Code = common.InnerError
			r.Message = "token id not match"
			return
		}
		logger.Logrus.WithFields(logrus.Fields{"graphRet": graphRet}).Info("OrderDetail RGraphOrderDetailData")
		//build returns

		originBalance, _ := big.NewInt(0).SetString(graphRet.Pool.AmountTotal1, 0)
		balance, err := tools.GetTokenAmount(originBalance, int32(purchaseTokenInfo.Decimal), 8)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"err": err}).Error("OrderDetail GetTokenAmount failed")
			r.Code = common.InnerError
			r.Message = "token amount get failed"
			return
		}

		status := 0
		now := time.Now().Unix()

		//buy
		//cacel
		//redeem
		//notlist
		if !graphRet.Pool.Swapped && !graphRet.Pool.Closed && fmt.Sprintf("%+v", now) < graphRet.Pool.CloseAt {
			if strings.EqualFold(p.UserAddress, graphRet.Pool.Creator) {
				status = const_def.MARKETPLACE_ORDER_CANCEL
			} else {
				status = const_def.MARKETPLACE_ORDER_BUY
			}
		} else {
			if strings.EqualFold(p.UserAddress, graphRet.Pool.Creator) && !graphRet.Pool.Closed && !graphRet.Pool.Swapped && !graphRet.Pool.Redeemed && fmt.Sprintf("%+v", now) >= graphRet.Pool.CloseAt {
				status = const_def.MARKETPLACE_ORDER_REDEEM
			} else {
				status = const_def.MARKETPLACE_ORDER_NOT_LIST
			}
		}

		//TODO need has a map for other sold nft contract to show tokenName/tokenSymbol
		retData := ROrderDetail{
			Index:               p.Index,
			TokenID:             graphRet.Pool.TokenId,
			Price:               balance,
			OriginPrice:         graphRet.Pool.AmountTotal1,
			PurchaseToken:       purchaseTokenInfo.TokenName,
			Status:              status,
			ContractAddress:     graphRet.Pool.Token0,
			Attrs:               nattrs,
			ChainImage:          finalChainURI,
			Image:               equipImage,
			Name:                nftContractInfo.TokenName,
			Symbol:              nftContractInfo.TokenSymbol,
			Standard:            "ERC721",
			Network:             network,
			Approved:            false,
			MarketplaceContract: marketplaceContract,
			Creator:             graphRet.Pool.Creator,
		}
		ret = retData
	}

	if r.Code == common.SuccessCode {
		r.Data = ret
	}
}

func Activity(c *gin.Context) {
	r := common.Response{
		Code: common.SuccessCode,
		Data: nil,
	}

	defer func() {
		c.JSON(http.StatusOK, r)
	}()

	var p = PActivity{}
	err := c.ShouldBindQuery(&p)

	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("market activity validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("market activity input info")

	url := config.GetMarketplaceGraph()
	if url == "" {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("market activity graph url is null")
		r.Code = common.InnerError
		r.Message = "get graph url failed"
		return
	}

	typeFilter := ""
	if !(p.Canceled == p.List && p.List == p.Purchase && p.Purchase == p.Sale && p.Sale == p.Redeem) {
		isList := ""
		if p.List {
			isList = fmt.Sprintf("%d,", const_def.MARKETPLACE_LIST)
		}
		isCancel := ""
		if p.Canceled {
			isCancel = fmt.Sprintf("%d,", const_def.MARKETPLACE_CANCEL)
		}
		isPurchase := ""
		if p.Purchase {
			isPurchase = fmt.Sprintf("%d,", const_def.MARKETPLACE_PURCHASE)
		}
		isSale := ""
		if p.Sale {
			isSale = fmt.Sprintf("%d,", const_def.MARKETPLACE_SALE)
		}

		isRedeem := ""
		if p.Redeem {
			isRedeem = fmt.Sprintf("%d,", const_def.MARKETPLACE_REDEEM)
		}
		final := isList + isCancel + isPurchase + isSale + isRedeem

		final = final[:len(final)-1]

		typeFilter = fmt.Sprintf("activityType_in:[%s],", final)
	}

	if p.AssetName == "" {
		list, lastKey, err := getAllActivitys(int(p.GameID), int(p.PageSize), p.LastKey, typeFilter, strings.ToLower(p.Address))

		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"Data": p, "ErrMsg": err.Error()}).Error("market activity getAllActivitys failed")
			r.Code = common.InnerError
			r.Message = "get all activitys failed"
			return
		}

		r.Data = RActivity{
			LastKey:      lastKey,
			ActivityList: list,
		}

		r.Message = "activitys success"
		return
	}

	contractInfo, err := comminfo.GetNftContractByAssetName(int(p.GameID), p.AssetName)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"Data": p, "ErrMsg": err.Error()}).Error("market activity GetNftContractByAssetName failed")
		r.Code = common.InnerError
		r.Message = "get nft contract info failed"
		return
	}

	chainToken := config.GetChainToken()
	purchaseTokenInfo, ok := chainToken[fmt.Sprintf("%d", contractInfo.ChainID)]
	if !ok || purchaseTokenInfo.TokenName == "" || purchaseTokenInfo.Decimal == 0 {
		logger.Logrus.WithFields(logrus.Fields{"Data": p, "ContractInfo": contractInfo}).Error("market activity wrong chain id config")
		r.Code = common.InnerError
		r.Message = "error chain id config"
		return
	}

	list, lastKey, err := getActivitySignalAssetName(int(p.PageSize), p.LastKey, typeFilter, strings.ToLower(p.Address), contractInfo.ContractAddress, contractInfo.TokenName, contractInfo.TokenSymbol, int(purchaseTokenInfo.Decimal), contractInfo.GameAssetName, contractInfo.BaseURL)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"Data": p, "ContractInfo": contractInfo, "ErrMsg": err.Error()}).Error("market activity getActivitySignalAssetName failed")
		r.Code = common.InnerError
		r.Message = "activity for asset name failed"
		return
	}

	r.Data = RActivity{
		LastKey:      lastKey,
		ActivityList: list,
	}

	r.Message = "activity for asset name success"
}
