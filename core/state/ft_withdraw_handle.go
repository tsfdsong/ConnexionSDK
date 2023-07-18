package state

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/core/riskcontrol"
	"github/Connector-Gamefi/ConnectorGoSDK/core/sign"
	"github/Connector-Gamefi/ConnectorGoSDK/distmiddleware/skywalking"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/http_client"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	webcommon "github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"math/big"
	"strings"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

// FTRiskControlHandler,you must handle the error if it is not nil.if success, you can update risk code into the risk_status
func FTRiskControlHandler(appID int, account, amount, coinName string) (int, error) {
	//check risk control logic
	return riskcontrol.FTRiskReview(appID, account, amount, coinName)
}

// FTPreWithdrawHandler CodeWithdrawRisking -> CodeWithdrawPreWithdraw
func FTPreWithdrawHandler(param *webcommon.PGameERC20PreWithdraw, url string) (webcommon.RERC20PreWithdraw, error, error) {
	//send prewithdraw request to game

	//update db state to CodeWithdrawPreWithdraw
	r := webcommon.RERC20PreWithdraw{}
	err1, err2 := http_client.HttpClientReqWithExactError(url, *param, &r)
	if err1 != nil || err2 != nil {
		return r, err1, err2
	}

	return r, nil, nil
}

func getFTSignHash(user, treasure, token, amount, nonce string, chainId int64) (string, error) {
	uint256Ty, err := abi.NewType("uint256", "uint256", nil)
	if err != nil {
		return "", errors.New("nettype failed")
	}

	addressTy, err := abi.NewType("address", "address", nil)
	if err != nil {
		return "", errors.New("nettype failed")
	}

	bytesTy, err := abi.NewType("bytes4", "bytes4", nil)
	if err != nil {
		return "", errors.New("nettype failed")
	}

	arguments := abi.Arguments{
		{
			Type: addressTy,
		},
		{
			Type: addressTy,
		},
		{
			Type: addressTy,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: bytesTy,
		},
	}

	if chainId != 97 && chainId != 56 {
		arguments = append(arguments, abi.Argument{
			Type: uint256Ty,
		})
	}

	bAmount, ok := new(big.Int).SetString(amount, 0)
	if !ok {
		return "", errors.New("bigint set string failed")
	}
	bNonce, ok := new(big.Int).SetString(nonce, 0)
	if !ok {
		return "", errors.New("bigint set string failed")
	}

	bChainID := big.NewInt(chainId)

	b, _ := hex.DecodeString(const_def.SDK_ERC20_WITHDRAWSELECTOR)
	b4 := [4]byte{}
	copy(b4[:], b[:4])

	var bytes []byte
	if chainId != 97 && chainId != 56 {
		bytes, err = arguments.Pack(
			common.HexToAddress(user),
			common.HexToAddress(treasure),
			common.HexToAddress(token),
			bAmount,
			bNonce,
			b4,
			bChainID,
		)
		if err != nil {
			return "", err
		}
	} else {
		bytes, err = arguments.Pack(
			common.HexToAddress(user),
			common.HexToAddress(treasure),
			common.HexToAddress(token),
			bAmount,
			bNonce,
			b4,
		)
		if err != nil {
			return "", err
		}
	}

	hash := crypto.Keccak256Hash(bytes)
	paramsHash := hash.Hex()

	return paramsHash, nil
}

// FTSignHandler CodeWithdrawPreWithdraw -> CodeWithdrawSign
// !!!FUNCTION DEPEND ON SKYWORKING TRACE. DOES NOT WORK ALONE
func FTSignHandler(signAmount, withdrawAddr, contractAddr, treasure, nonce string, chainID int64, ctx context.Context) (sign.RRequestSignature, error) {
	r := sign.RRequestSignature{}

	paramsHash, err := getFTSignHash(withdrawAddr, treasure, contractAddr, signAmount, nonce, chainID)
	if err != nil {
		return r, err
	}

	signData := fmt.Sprintf(`{"contract":"%s", "type":"%s","tokenId":0, "amount":%s}`, contractAddr, "ERC20", signAmount)
	signParam, err := sign.GenSignRequest(paramsHash, signData)
	if err != nil {
		return r, err
	}

	logger.Logrus.WithFields(logrus.Fields{"signparam": signParam, "paramsHash": paramsHash, "signData:": signData}).Info("FT sign result")

	secondUrl := config.GetSignConfig().SignServerURL + const_def.SDK_WITHDRAW_SIGN_URL

	err = skywalking.SkyPostRequest(secondUrl, signParam, &r, ctx)
	if err != nil {

		return r, err
	}

	if r.Code != const_def.SIGN_SUCCESS || r.ReqHash == "" {
		return r, fmt.Errorf("send sign request failed")
	}

	return r, nil
}

func FTCommitWithdrawHandler(param webcommon.PGameERC20WithdrawComfirm, url string) (webcommon.RERC20WithdrawComfirm, error, error) {
	//send commit withdraw request to game

	r := webcommon.RERC20WithdrawComfirm{}
	err1, err2 := http_client.HttpClientReqWithExactError(url, param, &r)
	if err1 != nil || err2 != nil {
		return r, err1, err2
	}
	return r, nil, nil
}

func FTWithdrawSignResultRequest(param sign.PQuerySignature) (sign.RQuerySignature, error) {
	r := sign.RQuerySignature{}
	signURL := config.GetSignConfig().SignResultServerURL + const_def.SDK_QUERY_SIGN_URL + "/" + param.ReqHash
	err := http_client.HttpClientReqWithGet(signURL, &r)
	if err != nil {

		return r, err
	}

	return r, nil
}

// FTRiskReviewHandle risk control review
func FTRiskReviewHandle(id, riskStatus int, reviewer string, ctx context.Context) error {
	var data model.TFtWithdrawRecord

	con := map[string]interface{}{"id": id}
	err := mysql.GetDB().Model(&model.TFtWithdrawRecord{}).Where(con).First(&data).Error
	if err != nil {
		return err
	}

	data.RiskStatus = riskStatus
	data.RiskReviewer = reviewer

	//create risk control sub span
	span, subctx, err := skywalking.CreateLocalSpan(ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": data.ID}).Error("FTRiskReviewHandle CreateLocalSpan failed")

		return err
	}
	span.SetOperationName("FTRiskReview")
	ctx = subctx
	defer span.End()

	//get trace id from cache
	go2sky.PutCorrelation(ctx, "type", fmt.Sprintf("%d", skywalking.FTRiskReviewCode))

	traceID, err := redis.GetString(fmt.Sprintf("ft:%s", data.Nonce))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": data.ID}).Error("FTRiskReviewHandle get trace id from cache failed")
		return err
	}
	go2sky.PutCorrelation(ctx, "pre_trace_id", traceID)

	seq := int(2)
	err = skywalking.SkyPutCorrelation(config.GetSkyWalkingConfig().KeyRiskReview, seq, config.GetSkyWalkingConfig().ValueRiskReview, ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"OrderID": data.ID, "ErrMsg": err.Error()}).Error("FTRiskReviewHandle SkyPutCorrelation failed")

		return err
	}

	defer func() {
		err := redis.DeleteString(fmt.Sprintf("ft:%s", data.Nonce))
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": data.ID}).Error("FTRiskReviewHandle delete trace id from cache failed")
		} else {
			logger.Logrus.WithFields(logrus.Fields{"TraceID": traceID, "OrderID": data.ID}).Info("FTRiskReviewHandle delete trace id from cache success")
		}
	}()

	if data.RiskStatus == riskcontrol.CodeRiskSuccess {
		//create risk control sub span
		gamespan, gamectx, err := skywalking.CreateLocalSpan(ctx)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": data.ID}).Error("FTRiskReviewHandle GamePrewithdraw CreateLocalSpan failed")
			return err
		}
		gamespan.SetOperationName("FTGameWithdraw")
		ctx = gamectx
		defer gamespan.End()

		seq = seq + 1
		err = skywalking.SkyPutCorrelation(config.GetSkyWalkingConfig().KeyGameWithdraw, seq, config.GetSkyWalkingConfig().ValueGameWithdraw, ctx)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"OrderID": data.ID, "ErrMsg": err.Error()}).Error("FTRiskReviewHandle GamePrewithdraw SkyPutCorrelation failed")

			return err
		}

		cacheData, err := comminfo.GetFTContractByAddress(data.AppID, data.ContractAddress)
		if err != nil {
			return err
		}

		balance, err := tools.GetTokenExactAmount(data.Amount, int32(cacheData.TokenDecimal))
		if err != nil {
			return err
		}

		rsp, err := ingame.RequestFTPreWithdrawToGame(int(data.AppID), data.UID, balance, data.GameCoinName, data.Nonce)
		if err != nil || rsp.IsFailure() || rsp.Data.AppOrderID == "" {
			retryData := &commdata.FTPrewithdrawRetryData{
				ContractAddress: strings.ToLower(data.ContractAddress),
				UserAddress:     strings.ToLower(data.WithdrawAddress),
				Nonce:           data.Nonce,
				Account:         data.Account,
				AppID:           int(data.AppID),
				Uid:             int64(data.UID),
				AppCoinName:     data.GameCoinName,
				Amount:          data.Amount,
			}
			jobid, nerr := FTPublishPrewithdraw(retryData)
			if nerr != nil {
				logger.Logrus.WithFields(logrus.Fields{"ErrMsg": nerr.Error(), "RetryData": retryData}).Error("FTRiskReviewHandle FTPublishPrewithdraw failed")
				return nerr
			} else {
				logger.Logrus.WithFields(logrus.Fields{"RetryData": retryData, "JobID": jobid}).Info("FTRiskReviewHandle FTPublishPrewithdraw success")
			}

			data.OrderStatus = const_def.CodeWithdrawTimeout

		} else {
			data.AppOrderID = rsp.Data.AppOrderID
		}

		signResult, signErr := FTSignHandler(data.Amount, data.WithdrawAddress, data.ContractAddress, cacheData.Treasure, data.Nonce, int64(data.ChainID), ctx)
		if signErr != nil {
			logger.Logrus.WithFields(logrus.Fields{"OrderID": data.ID, "ErrMsg": signErr, "SignResult": signResult}).Error("FTRiskReviewHandle FTSignHandler failed")

			//noti gameserver recover assets
			secondResult, err := ingame.RequestFTRecoverToGame(data.AppID, data.UID, rsp.Data.AppOrderID, data.GameCoinName, data.Nonce)

			if err != nil || secondResult.IsFailure() || secondResult.IsReject() {

				logger.Logrus.WithFields(logrus.Fields{"Order": data, "ErrMsg": err, "GaemResponse": secondResult}).Error("FTRiskReviewHandle RequestFTRecoverToGame failed")

				//noti failed. save db and let timer to handle it
				data.OrderStatus = const_def.CodeWithdrawCommitFailed
			} else if secondResult.IsSuccess() {
				data.OrderStatus = const_def.CodeNotiRecorverSuccess
			} else {
				logger.Logrus.WithFields(logrus.Fields{"secondresult": secondResult}).Error("FTRiskReviewHandle RequestFTRecoverToGame Unknow Status")
			}
		} else {
			data.OrderStatus = const_def.CodeWithdrawSign
			data.SignatureHash = signResult.ReqHash
		}
	} else {
		data.OrderStatus = const_def.CodeWithdrawRiskFailed
	}

	logger.Logrus.WithFields(logrus.Fields{"Order": data}).Error("FTRiskReviewHandle order status info")

	values := map[string]interface{}{"app_order_id": data.AppOrderID, "order_status": data.OrderStatus, "signature": data.Signature, "signature_hash": data.SignatureHash, "risk_reviewer": data.RiskReviewer, "risk_status": data.RiskStatus, "timestamp": time.Now().UnixMilli()}
	err = mysql.GetDB().Model(&model.TFtWithdrawRecord{}).Where(con).Updates(values).Error
	if err != nil {
		return err
	}

	return nil
}
