package sdkcontroller

import (
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/riskcontrol"
	"github/Connector-Gamefi/ConnectorGoSDK/core/state"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"github/Connector-Gamefi/ConnectorGoSDK/web/sdkbackend/sdkdata"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//cache:2 DB:1
func GetUserAssets(c *gin.Context) {
	r := &common.Response{
		Code: common.SuccessCode,
	}

	retData := sdkdata.RGetUserAssets{
		Erc20:      []common.RBackendERC20Asset{},
		UpdateTime: "",
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = sdkdata.PGetUserAssets{}
	err := c.ShouldBindQuery(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("GetUserAssets validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	chainIDStr := c.Request.Header.Get("Chainid")
	chainID, err := strconv.ParseInt(chainIDStr, 0, 64)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ChainID": chainIDStr}).Error("GetUserAssets invalid chainID")

		r.Code = common.IncorrectParams
		r.Message = "invalid input chainID"
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p, "ChainID": chainID}).Info("GetUserAssets input info")

	gameErc20Map, err := comminfo.GetFTContractAssetNameMap(int(p.GameID), int(chainID))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("GetUserAssets GetFTContractAssetNameMap failed")
		r.Code = common.InnerError
		r.Message = "get ft contract info failed"
		return
	}

	//4 request game server
	ftData, err := ingame.RequestGameFTAssets(int(p.GameID), p.Email, p.UID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("GetUserAssets RequestGameFTAssets failed")
		r.Code = common.InnerError
		r.Message = "get ft gage assets failed"
		return
	}

	for _, e := range ftData {
		v, ok := gameErc20Map[e.AppCoinName]
		if !ok {
			logger.Logrus.WithFields(logrus.Fields{"game_coin_name": e.AppCoinName}).Error("GetUserAssets find contract info by game_coin_name falied")
			continue
		}
		balance, err := tools.AsDashboardDisplayAmount(e.Balance, 8)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"balance": e.Balance}).Error("GetUserAssets error balance")
			continue
		}

		frozenBalance, err := tools.AsDashboardDisplayAmount(e.FrozenBalance, 8)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"FrozenBalance": e.FrozenBalance}).Error("GetUserAssets error FrozenBalance")
			continue
		}

		erc20Item := common.RBackendERC20Asset{
			Symbol:        v.TokenSymbol,
			Contract:      v.ContractAddress,
			Balance:       balance,
			FrozenBalance: frozenBalance,
		}
		retData.Erc20 = append(retData.Erc20, erc20Item)
	}

	retData.UpdateTime = time.Now().UTC().Format("2006-01-02 15:04:05")

	r.Data = retData
	r.Message = "get user ft game assets success"
}

func GetUserNFTAssets(c *gin.Context) {
	r := common.U8ListResponse{
		Code: common.SuccessCode,
	}

	defer func() {
		r.Message = common.ErrorMap[int(r.Code)]
		c.JSON(http.StatusOK, r)
	}()

	retData := sdkdata.RGetUserNFTAssets{
		Erc721:     []common.RBackendERC721Asset{},
		UpdateTime: "",
	}

	var p = sdkdata.PGetUserNFTAssets{}
	err := c.ShouldBindQuery(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("GetUserNFTAssets validator reject")
		r.Code = common.IncorrectParams
		return
	}

	chainIDStr := c.Request.Header.Get("Chainid")
	chainID, err := strconv.ParseInt(chainIDStr, 0, 64)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ChainID": chainIDStr}).Error("GetUserAssets invalid chainID")

		r.Code = common.IncorrectParams
		r.Message = "invalid input chainID"
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"input": p, "ChainID": chainID}).Info("GetUserNFTAssets input info")

	nftContractMap, err := comminfo.GetNftContractByAppIDAndChainID(int(p.GameID), int(chainID))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "input": p}).Error("GetUserNFTAssets GetNftContractByAppIDAndChainID failed")
		r.Code = common.InnerError
		return
	}

	//TODO NOTICE!!!  if change nft contract. must associated with game server
	nftData, total, err := ingame.RequestGameNFTAssets(int(p.GameID), p.Email, p.AssetName, int64(p.Page), int64(p.PageSize), p.UID)
	logger.Logrus.WithFields(logrus.Fields{"nftData": nftData}).Info("QueryGameAssets nftData")
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("GetUserNFTAssets RequestGameNFTAssets failed")
		r.Code = common.InnerError
		return
	}

	for _, e := range nftData {
		balance := 1
		frozenBalance := 0
		if e.Frozen {
			balance = 0
			frozenBalance = 1
		}

		gameErc721Info, ok := nftContractMap[e.GameAssetName]
		if !ok {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err, "AssetName": e.GameAssetName}).Error("GetUserNFTAssets get nft contract from asset name failed")
			r.Code = common.InnerError
			return
		}

		erc721Item := common.RBackendERC721Asset{
			Symbol:        gameErc721Info.TokenSymbol,
			TokenContract: gameErc721Info.ContractAddress,
			TokenID:       e.TokenID,
			EquipmentID:   e.EquipmentID,
			Balance:       balance,
			FrozenBalance: frozenBalance,
		}
		retData.Erc721 = append(retData.Erc721, erc721Item)
	}

	retData.UpdateTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	r.Total = total

	if r.Code == common.SuccessCode {
		r.Data = retData
	}
}

//cache:3 db:1
func WithdrawExamineSet(c *gin.Context) {
	r := common.Response{
		Code: common.SuccessCode,
	}

	defer func() {
		r.Message = common.ErrorMap[int(r.Code)]
		c.JSON(http.StatusOK, r)
	}()

	var p = sdkdata.PWithdrawExamineSet{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("validator reject")
		r.Code = common.IncorrectParams
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("WithdrawExamineSet input data")

	if p.Status != riskcontrol.CodeRiskSuccess && p.Status != riskcontrol.CodeRiskRejected {
		logger.Logrus.WithFields(logrus.Fields{"Status": p.Status}).Error("WithdrawExamineSet incorrect status")
		r.Code = common.IncorrectParams
		return
	}

	if p.ContractType == const_def.U8_CONTRACT_ERC20 {
		err := state.FTRiskReviewHandle(p.ID, p.Status, p.Reviewer, c.Request.Context())
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "ID": p.ID,
				"riskstatus": p.Status, "reviewer": p.Reviewer}).Error("ft risk control review failed")
			r.Code = common.InnerError
			return
		}
	} else if p.ContractType == const_def.U8_CONTRACT_ERC721 {
		err := state.RiskReviewHandle(p.ID, p.Status, p.Reviewer, c.Request.Context())
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "ID": p.ID,
				"riskstatus": p.Status, "reviewer": p.Reviewer}).Error("nft risk control review failed")
			r.Code = common.InnerError
			return
		}
	} else {
		logger.Logrus.Info("incorrect contract type")
		r.Code = common.InnerError
		return
	}
}
