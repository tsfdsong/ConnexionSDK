package dashcontroller

import (
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/contracts"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"github/Connector-Gamefi/ConnectorGoSDK/web/dashboard/dashdata"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ChainFTAssets(c *gin.Context) {
	r := &common.Response{
		Code: common.SuccessCode,
	}

	var p = common.PChainFTAssets{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("QueryFtTreasure validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	chainIDStr := c.Request.Header.Get("Chainid")
	chainID, err := strconv.ParseInt(chainIDStr, 0, 64)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ChainID": chainIDStr}).Error("ChainFTAssets invalid chainID")

		r.Code = common.IncorrectParams
		r.Message = "invalid input chainID"
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"params": p, "ChainID": chainID}).Info("ChainFTAssets params")

	//step1 get all ft contracts
	ftContractCache, err := comminfo.GetFTContractAddressMap(int(p.GameID), int(chainID))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ChainFTAssets Find All Ft Contracts Failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"contract_info": ftContractCache}).Info("ChainFTAssets contract info")

	retData, err := contracts.GetGameAllErc20TokenBalance(p.Address, ftContractCache)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ChainFTAssets GetAllErc20Balance Failed")
		r.Code = common.InnerError
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	r.Data = retData
	r.Message = common.ErrorMap[int(r.Code)]
	c.JSON(http.StatusOK, r)
}

func DepositGameERC20Token(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	retData := dashdata.RDepositGameERC20Token{}
	var p = dashdata.PDepositGameERC20Token{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositGameERC20Token validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	chainIDStr := c.Request.Header.Get("Chainid")
	chainID, err := strconv.ParseInt(chainIDStr, 0, 64)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ChainID": chainIDStr}).Error("DepositGameERC20Token invalid chainID")

		r.Code = common.IncorrectParams
		r.Message = "invalid input chainID"
		return
	}
	p.ChainID = chainID

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
	retData.Amount = p.Amount
	retData.Nonce = p.Nonce
	retData.Treasure = p.TargetAddress
	retData.Contract = p.Contract
	r.Data = retData
}

func WithdrawGameERC20Token(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = dashdata.PWithdrawGameERC20Token{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("WithdrawGameERC20Token validator reject")
		r.Code = common.IncorrectParams
		r.Message = common.ErrorMap[int(r.Code)]
		return
	}

	chainIDStr := c.Request.Header.Get("Chainid")
	chainID, err := strconv.ParseInt(chainIDStr, 0, 64)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ChainID": chainIDStr}).Error("WithdrawGameERC20Token invalid chainID")

		r.Code = common.IncorrectParams
		r.Message = "invalid input chainID"
		return
	}
	p.ChainID = chainID

	p.SkywalkingSetting = dashdata.SkywalkingSetting{
		Seq: 0,
		Ctx: c.Request.Context(),
	}

	s := new(dashdata.MultiAssetStrategy)
	s.SetWithdrawImpl(&p)
	err = s.Withdraw()
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

	r.Message = "ft prewithdraw success"
}

func ClaimGameERC20Token(c *gin.Context) {
	r := &common.Response{
		Code:    common.SuccessCode,
		Message: common.ErrorMap[common.SuccessCode],
	}

	defer func(r *common.Response) {
		c.JSON(http.StatusOK, r)
	}(r)

	var p = dashdata.PClaimGameERC20Token{}
	err := c.ShouldBind(&p)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ClaimGameERC20Token validator reject")
		r.Code = common.IncorrectParams
		return
	}

	chainIDStr := c.Request.Header.Get("Chainid")
	chainID, err := strconv.ParseInt(chainIDStr, 0, 64)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ChainID": chainIDStr}).Error("ClaimGameERC20Token invalid chainID")

		r.Code = common.IncorrectParams
		r.Message = "invalid input chainID"
		return
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p, "ChainID": chainID}).Error("ClaimGameERC20Token input data")

	//check order
	withdrawRecord := model.TFtWithdrawRecord{}
	condition := map[string]interface{}{"app_id": p.GameID, "id": p.ID}
	err, found := mysql.WrapFindFirst(model.TableFtWithdrawRecord, &withdrawRecord, condition)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ClaimGameERC20Token WrapFindFirst Error")
		r.Code = common.InnerError
		return
	}

	if !found {
		r.Code = common.RecordNotExist
		return
	}

	if withdrawRecord.Signature == "" || withdrawRecord.OrderStatus != const_def.CodeWithdrawWaitClaim {
		r.Code = common.CantClaimNow
		return
	}

	//check contract switch
	ftContractCacheData, err := comminfo.GetFTContractByAddressAndChainID(int(p.GameID), strings.ToLower(withdrawRecord.ContractAddress), int(chainID))
	if err != nil {

		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("ClaimGameERC20Token GetFTContractByAddressAndChainID Error")

		r.Code = common.InnerError
		return
	}

	if ftContractCacheData.WithdrawSwitch != const_def.SDK_TABLE_SWITCH_OPEN {
		r.Code = common.WithdrawSwitchClose
		return
	}

	retData := dashdata.RClaimGameERC20Token{
		Amount:    withdrawRecord.Amount,
		Nonce:     withdrawRecord.Nonce,
		Treasure:  ftContractCacheData.Treasure,
		Signature: withdrawRecord.Signature,
	}

	r.Data = retData
}
