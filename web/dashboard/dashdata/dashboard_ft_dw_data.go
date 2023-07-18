package dashdata

import (
	"context"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/core/riskcontrol"
	"github/Connector-Gamefi/ConnectorGoSDK/core/state"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/math"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/nonce_gen"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PDepositGameERC20Token struct {
	GameID        uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	Address       string `json:"address" binding:"eth_addr,ne=0x0000000000000000000000000000000000000000"`
	Contract      string `form:"contract" json:"contract" uri:"contract" binding:"eth_addr,ne=0x0000000000000000000000000000000000000000"`
	Amount        string `json:"amount" binding:"gt=0"`
	UID           uint64 `binding:"-"`
	Account       string `binding:"-"`
	GameCoinName  string `binding:"-"`
	TargetAddress string `binding:"-"`
	Nonce         string `binding:"-"`
	ChainID       int64  `binding:"-"`
}

type RDepositGameERC20Token struct {
	Amount   string `json:"amount"`
	Nonce    string `json:"nonce"`
	Treasure string `json:"treasure"`
	//Signature string `json:"signature"`
	Contract string `json:"contract"`
}

type PWithdrawGameERC20Token struct {
	GameID            uint64                        `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	Address           string                        `json:"address" binding:"eth_addr,ne=0x0000000000000000000000000000000000000000"`
	Contract          string                        `form:"contract" json:"contract" uri:"contract" binding:"eth_addr,ne=0x0000000000000000000000000000000000000000"`
	Amount            string                        `json:"amount" binding:"gt=0"`
	Timestamp         string                        `json:"timestamp" binding:"gt=0"`
	SignString        string                        `json:"signature"  binding:"gt=0"`
	UID               uint64                        `binding:"-"`
	Account           string                        `binding:"-"`
	GameCoinName      string                        `binding:"-"`
	TargetAddress     string                        `binding:"-"`
	Nonce             string                        `binding:"-"`
	OrderStatus       int                           `binding:"-"`
	SignatureHash     string                        `binding:"-"`
	AppOrderID        string                        `binding:"-"`
	RiskStatus        int                           `binding:"-"`
	SkywalkingSetting SkywalkingSetting             `binding:"-"`
	GameInfo          string                        `binding:"-"`
	ContractInfo      *commdata.FTContractCacheData `binding:"-"`
	SignDeadline      int64                         `binding:"-"`
	ChainID           int64                         `binding:"-"`
}

type SkywalkingSetting struct {
	Seq     int
	TraceID string
	Ctx     context.Context
}

type PClaimGameERC20Token struct {
	GameID  uint64 `form:"app_id" json:"app_id" uri:"app_id" binding:"gt=0"`
	Address string `json:"address" binding:"eth_addr,ne=0x0000000000000000000000000000000000000000"`
	ID      uint64 `form:"id" json:"id" uri:"id" binding:"gt=0"`
}

type RClaimGameERC20Token struct {
	Amount    string `json:"amount"`
	Nonce     string `json:"nonce"`
	Treasure  string `json:"treasure"`
	Signature string `json:"signature"`
}

func (p *PDepositGameERC20Token) BeforeDeposit() error {
	logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("DepositGameERC20Token input info")

	bindEmail, err := comminfo.GetEmailBindInfo(uint(p.GameID), uint(p.ChainID), p.Address)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("DepositGameERC20Token GetEmailBindInfo Error")
		return common.NewHpError(err, int(common.InnerError), "")
	}

	ftContractCacheData, err := comminfo.GetFTContractByAddressAndChainID(int(p.GameID), strings.ToLower(p.Contract), int(p.ChainID))

	if err == gorm.ErrRecordNotFound {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("DepositGameERC20Token ft contract not found")
		return common.NewHpError(err, int(common.ContractInfoNotExist), "")
	}
	if err != nil || ftContractCacheData == nil {
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("DepositGameERC20Token GetFTContractByAddress Error")
		}
		return common.NewHpError(err, int(common.InnerError), "")
	}

	if ftContractCacheData.Treasure == "" || ftContractCacheData.TokenDecimal == 0 {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("DepositGameERC20Token treasure or decimal info not found")
		return common.NewHpError(err, int(common.ContractInfoNotExist), "")
	}

	if ftContractCacheData.DepositSwitch != const_def.SDK_TABLE_SWITCH_OPEN {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("DepositGameERC20Token contract deposit switch is closed")
		return common.NewHpError(err, int(common.DepositSwitchClose), "")
	}

	err = math.CheckAmountValidator(p.Amount, ftContractCacheData.GameDecimal, ftContractCacheData.TokenDecimal)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("DepositGameERC20Token CheckAmountValidator failed")
		return common.NewHpError(err, int(common.InnerError), err.Error())
	}

	//TODO nonce is a resource. avoid to cheap invoke this

	//step1  gen nonce
	nonce := nonce_gen.GenNonce(fmt.Sprintf("%s%s%s%d%d%s", p.Address, p.Amount, p.Contract, p.GameID, time.Now().UnixNano(), tools.GenCode(8)))

	logger.Logrus.WithFields(logrus.Fields{"Data": p, "Nonce": nonce, "Treasure": ftContractCacheData.Treasure}).Info("DepositGameERC20Token generate nonce success")
	p.UID = bindEmail.UID
	p.Account = bindEmail.Account
	p.GameCoinName = ftContractCacheData.GameCoinName
	p.TargetAddress = ftContractCacheData.DepositTreasure
	p.Nonce = nonce
	return nil
}

func (p *PDepositGameERC20Token) Deposit() error {
	err := p.BeforeDeposit()
	if err != nil {
		return err
	}
	return nil
}

func (p *PWithdrawGameERC20Token) BeforeWithdraw() error {
	logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("FT withdraw input info")

	apiDelayTime := config.GetAPIDelayTime()
	timestmap, err := strconv.ParseInt(p.Timestamp, 10, 0)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("FT withdraw parse timestamp failed")
		return common.NewHpError(err, int(common.IncorrectParams), "")
	}
	if timestmap+apiDelayTime < time.Now().Unix() {
		logger.Logrus.WithFields(logrus.Fields{"Now": time.Now().Unix(), "timestmap": timestmap, "apiDelayTime": apiDelayTime}).Error("FT withdraw check timestamp timeout")
		return common.NewHpError(err, int(common.InvalidTimestamp), "")
	}

	waitSignString := fmt.Sprintf("app_id=%d&address=%s&contract=%s&amount=%s&timestamp=%s&chainId=%d", p.GameID, strings.ToLower(p.Address), strings.ToLower(p.Contract), p.Amount, p.Timestamp, p.ChainID)
	if !tools.CheckPersonalSign(waitSignString, p.SignString, p.Address) {
		logger.Logrus.WithFields(logrus.Fields{"Data": p, "RawString": waitSignString}).Error("FT withdraw check sign failed")
		return common.NewHpError(err, int(common.InvalidSignature), "")
	}

	//ensure sign only use once
	signKey := const_def.GetFTSignKey(p.Address, p.SignString)
	signValue, err := redis.GetStringAcceptable(signKey)
	logger.Logrus.WithFields(logrus.Fields{"signValue": signValue, "err": err, "p": p}).Info("FT withdraw ft prewithdraw get sign from cache")
	if err != nil {
		return common.NewHpError(err, int(common.InnerError), "")
	} else if signValue != "" {
		return common.NewHpError(err, int(common.InvalidSignature), "")
	} else {
		p.SignDeadline = timestmap + apiDelayTime
	}

	//step1 bind email info
	bindEmail, err := comminfo.GetBindInfoWithChainID(int(p.GameID), strings.ToLower(p.Address), p.ChainID)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("FT withdraw GetBindInfoWithChainID failed")
		return common.NewHpError(err, int(common.InnerError), "")
	}

	if bindEmail.UID == 0 {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("FT withdraw game account uid not found")
		return common.NewHpError(err, int(common.MustBindGameAccount), "")
	}

	//step2 user info
	var userInfo = model.TUser{}
	usercondition := map[string]interface{}{"app_id": p.GameID, "account": bindEmail.Account}
	err, found := comminfo.GetUserInfo(usercondition, &userInfo)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("FT withdraw GetUserInfo failed")
		return common.NewHpError(err, int(common.InnerError), "")
	}

	if !found {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("FT withdraw user not exist")
		return common.NewHpError(err, int(common.UserNotExist), "")
	}

	if userInfo.WithdrawSwitch != const_def.SDK_TABLE_SWITCH_OPEN {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("FT withdraw user withdraw switch is closed")
		return common.NewHpError(err, int(common.UserCantWithdraw), "")
	}

	//step3 check contract info
	p.ContractInfo, err = comminfo.GetFTContractByAddressAndChainID(int(p.GameID), strings.ToLower(p.Contract), int(p.ChainID))
	if err == gorm.ErrRecordNotFound {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("FT withdraw ft contract info not found")
		return common.NewHpError(err, int(common.ContractInfoNotExist), "")
	}
	if err != nil || p.ContractInfo == nil {
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("FT withdraw GetFTContractByAddressAndChainID Error")
		}
		return common.NewHpError(err, int(common.InnerError), "")
	}

	if p.ContractInfo.Treasure == "" {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("FT withdraw treasure is empty")
		return common.NewHpError(err, int(common.ContractInfoNotExist), "")
	}

	if p.ContractInfo.WithdrawSwitch != const_def.SDK_TABLE_SWITCH_OPEN {
		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("FT withdraw contract withdraw switch is closed")
		return common.NewHpError(err, int(common.WithdrawSwitchClose), "")
	}

	//step4 check decimal info
	err = math.CheckAmountValidator(p.Amount, p.ContractInfo.GameDecimal, p.ContractInfo.TokenDecimal)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("FT withdraw CheckAmountValidator failed")
		return common.NewHpError(err, int(common.IncorrectParams), "")
	}

	//step6 check game info
	p.GameInfo, err = comminfo.GetBaseUrlOriginError(int(p.GameID))
	if p.GameInfo == "" || err == gorm.ErrRecordNotFound {
		return common.NewHpError(err, int(common.GameInfoNotExist), "")
	}
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "Data": p}).Error("FT withdraw GetBaseUrlOriginError Error")
		return common.NewHpError(err, int(common.InnerError), "")
	}

	//step7 begin withdraw flow

	//step1  gen nonce
	nonce := nonce_gen.GenNonce(fmt.Sprintf("%d%s%s%s%d%d%s", p.ChainID, p.Address, p.Amount, p.Contract, p.GameID, time.Now().UnixNano(), tools.GenCode(8)))

	//need insert this

	p.UID = bindEmail.UID
	p.GameCoinName = p.ContractInfo.GameCoinName

	p.Account = bindEmail.Account
	p.Nonce = nonce

	return nil
}

func (p *PWithdrawGameERC20Token) RiskControl() (bool, error) {
	seq := int(1)
	riskControlSubSpan := new(RiskControlSubSpan)
	err := riskControlSubSpan.Handle(p.SkywalkingSetting.Ctx, seq)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"OrderNonce": p.Nonce}).Error(err.Error())
	}
	defer riskControlSubSpan.End()

	//risk control
	withdrawCode, err := state.FTRiskControlHandler(int(p.GameID), p.Account, p.Amount, p.GameCoinName)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "code": withdrawCode}).Error("FT withdraw erc20 withdraw risk failed")
		return false, common.NewHpError(err, int(common.InnerError), "")
	}

	logger.Logrus.WithFields(logrus.Fields{"Data": p, "code": withdrawCode}).Info("FT withdraw erc20 withdraw risk result")

	p.RiskStatus = withdrawCode
	p.SkywalkingSetting = SkywalkingSetting{
		Seq:     seq,
		TraceID: riskControlSubSpan.TraceID,
		Ctx:     riskControlSubSpan.GetContext(),
	}
	if withdrawCode == riskcontrol.CodeRiskWaiting {
		logger.Logrus.WithFields(logrus.Fields{"OrderNonce": p.Nonce}).Info("FT withdraw order join into risk waiting")

		//skywalking
		riskControlSubSpan.GetSpan().Log(time.Now(), "[FTWithdraw]", fmt.Sprintf(" Seq: %d\n OrderNonce: %s\n traceID: %s\n waiting", seq, p.Nonce, riskControlSubSpan.TraceID))

		cacheID := fmt.Sprintf("ft:%s", p.Nonce)
		rerr := redis.SetString(cacheID, riskControlSubSpan.TraceID, config.GetKeyNoExpireTime())
		if rerr != nil {
			logger.Logrus.WithFields(logrus.Fields{"OrderNonce": p.Nonce, "ErrMsg": rerr.Error()}).Error("FT withdraw cache trace id failed")
		}

		p.OrderStatus = const_def.CodeWithdrawRisking

		return true, nil
	}

	return false, nil
}

func (p *PWithdrawGameERC20Token) Withdraw() error {
	err := p.BeforeWithdraw()
	if err != nil {
		return err
	}

	// create init record
	record, err := p.CreateWithdrawRecord()
	if err != nil {
		return common.NewHpError(err, int(common.InnerError), "")
	}

	isRisking, err := p.RiskControl()
	if err != nil {
		return err
	}

	if isRisking {
		record.OrderStatus = p.OrderStatus
		record.RiskStatus = p.RiskStatus
		err := state.UpdateFTWithdrawOrderStatus(record)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"err": err.Error(), "data": p}).Error("FT withdraw update mysql risk control status failed")
			return common.NewHpError(err, int(common.InnerError), "risking control")
		}
		return nil
	}

	//create game server withdraw sub span
	gameServerSubSpan := new(GameServerSubSpan)
	p.SkywalkingSetting.Seq = p.SkywalkingSetting.Seq + 1
	err = gameServerSubSpan.Handle(p.SkywalkingSetting.Ctx, p.SkywalkingSetting.Seq)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"OrderNonce": p.Nonce, "ErrMsg": err.Error()}).Error("FT withdraw handle skywalking failed")
		return common.NewHpError(err, int(common.InnerError), "")
	}
	p.SkywalkingSetting.Ctx = gameServerSubSpan.GetContext()
	defer gameServerSubSpan.End()

	//continue flow

	//step1 pre withdraw
	//1 --send pre withdraw failed. return frontend failed
	//2 --send success but game server reject. return frontend failed.
	//3 --send success and game server pass. gen nonce

	balance, err := tools.GetTokenExactAmount(p.Amount, int32(p.ContractInfo.TokenDecimal))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("FT withdraw GetTokenExactAmount failed")
		return common.NewHpError(err, int(common.InnerError), "")
	}
	logger.Logrus.WithFields(logrus.Fields{"balance": balance, "amount": p.Amount, "decimal": p.ContractInfo.TokenDecimal}).Info("FT withdraw Convert Balance success")

	rsp, err := ingame.RequestFTPreWithdrawToGame(int(p.GameID), p.UID, balance, p.ContractInfo.GameCoinName, p.Nonce)

	//must save db
	if err != nil {
		retryData := &commdata.FTPrewithdrawRetryData{
			ContractAddress: strings.ToLower(p.Contract),
			UserAddress:     strings.ToLower(p.Address),
			Nonce:           p.Nonce,
			Account:         p.Account,
			AppID:           int(p.GameID),
			Uid:             int64(p.UID),
			AppCoinName:     p.ContractInfo.GameCoinName,
			Amount:          p.Amount,
			TraceID:         p.SkywalkingSetting.TraceID,
			SeqNumber:       p.SkywalkingSetting.Seq + 1,
			ChainID:         p.ChainID,
		}
		jobid, nerr := state.FTPublishPrewithdraw(retryData)
		if nerr != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": nerr.Error(), "RetryData": retryData, "JobID": jobid}).Error(const_def.NEED_MANUAL_REPAIR_THIS + " FT withdraw PrewithdrawHandler FTPublishPrewithdraw failed")

		} else {
			logger.Logrus.WithFields(logrus.Fields{"RetryData": retryData, "JobID": jobid}).Info("FT withdraw PrewithdrawHandler FTPublishPrewithdraw success")
		}
		//r.code success
		p.WithdrawTimeout(record)
		return nil

	}

	if rsp.IsFailure() {
		p.OrderStatus = const_def.CodeWithdrawPreWithdrawFailed

		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Error("FT withdraw prewithdraw to game server failed")

		go p.UpdateRecordStatus(record)
		return common.NewHpError(err, int(common.WithdrawFailed), "")
	}

	//get apporder id
	p.AppOrderID = rsp.Data.AppOrderID
	//step3 sign
	//1 --send sign failed.
	//2 --send sign success but code is failure
	//3 --send sign success and code is success and save db

	signResult, signErr := state.FTSignHandler(p.Amount, p.Address, p.Contract, p.ContractInfo.Treasure, p.Nonce, p.ChainID, p.SkywalkingSetting.Ctx)
	if signErr == nil {
		p.OrderStatus = const_def.CodeWithdrawSign
		p.SignatureHash = signResult.ReqHash

		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("FT withdraw send request to signmachine success")

		p.WithdrawSuccess(record)
		return nil
	} else {
		p.OrderStatus = const_def.CodeWithdrawSignFailed

		logger.Logrus.WithFields(logrus.Fields{"SignResult": signResult, "ErMsg": signErr, "Data": p}).Error("FT withdraw send request to signmachine failed")
	}

	//noti gameserver recover assets
	recoverRsp, err := ingame.RequestFTRecoverToGame(int(p.GameID), p.UID, p.AppOrderID, p.ContractInfo.GameCoinName, p.Nonce)
	if err != nil || recoverRsp.IsFailure() || recoverRsp.IsReject() {
		p.OrderStatus = const_def.CodeWithdrawCommitFailed
		go p.UpdateRecordStatus(record)

		logger.Logrus.WithFields(logrus.Fields{"Data": p, "err": err, "RecoverFTResult": recoverRsp}).Error("FT withdraw recorver failed")

		return common.NewHpError(err, int(common.WithdrawFailed), "")

	} else if recoverRsp.IsSuccess() {
		//must have
		p.OrderStatus = const_def.CodeNotiRecorverSuccess

		logger.Logrus.WithFields(logrus.Fields{"Data": p}).Info("FT withdraw recorver success")

		go p.UpdateRecordStatus(record)
		return common.NewHpError(err, int(common.WithdrawFailed), "")
	} else {
		logger.Logrus.WithFields(logrus.Fields{"Data": p, "recoverRsp": recoverRsp}).Error("FT withdraw recover Unknow Status")

		go p.UpdateRecordStatus(record)
		return common.NewHpError(err, int(common.WithdrawFailed), "")
	}
}

func (p *PWithdrawGameERC20Token) UpdateRecordStatus(r *model.TFtWithdrawRecord) error {
	r.OrderStatus = p.OrderStatus
	r.AppOrderID = p.AppOrderID
	r.SignatureHash = p.SignatureHash
	err := state.UpdateFTWithdrawOrderStatus(r)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"err": err.Error(), "data": p}).Error(const_def.NEED_MANUAL_REPAIR_THIS + "FT withdraw update mysql failed")
	}
	return err
}

func (p *PWithdrawGameERC20Token) CreateWithdrawRecord() (*model.TFtWithdrawRecord, error) {
	erc20W := model.TFtWithdrawRecord{}
	erc20W.UID = p.UID
	erc20W.GameCoinName = p.GameCoinName
	erc20W.AppID = int(p.GameID)
	erc20W.ChainID = uint64(p.ChainID)

	erc20W.Account = p.Account
	erc20W.ContractAddress = strings.ToLower(p.Contract)
	erc20W.WithdrawAddress = strings.ToLower(p.Address)
	erc20W.Amount = p.Amount
	erc20W.Nonce = p.Nonce
	erc20W.RiskStatus = p.RiskStatus
	erc20W.OrderStatus = const_def.CodeWithdrawNone

	err := state.InsertFTWithdrawRecord(&erc20W)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"erc20WithdrawInfo": erc20W, "err": err.Error()}).Error(const_def.NEED_MANUAL_REPAIR_THIS + "FT withdraw insert mysql failed")
	}
	return &erc20W, err

}

func (p *PWithdrawGameERC20Token) SetFTSignKeyUsed() {
	t := time.Now().Unix()
	if (p.SignDeadline - t) > 0 {
		redis.SetString(const_def.GetFTSignKey(p.Address, p.SignString), "sign", p.SignDeadline-t)
	}
}

func (p *PWithdrawGameERC20Token) WithdrawSuccess(r *model.TFtWithdrawRecord) {
	p.UpdateRecordStatus(r)
	p.SetFTSignKeyUsed()
}

func (p *PWithdrawGameERC20Token) WithdrawTimeout(r *model.TFtWithdrawRecord) {
	p.OrderStatus = const_def.CodeWithdrawTimeout
	p.UpdateRecordStatus(r)
	p.SetFTSignKeyUsed()
}
