package ingame

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"time"

	"github/Connector-Gamefi/ConnectorGoSDK/utils/http_client"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"

	"github.com/sirupsen/logrus"
)

// RequestPreWithdrawToGame send request to game for prewithdraw
func RequestPreWithdrawToGame(appId int, uid uint64, equipID, name, nonce string) (*model.TNftWithdrawRecord, error) {
	c := config.GetGameConfig()
	if c[appId].GameServerMock {
		t := time.Now()
		mockId := t.Format("20060102150405")
		return &model.TNftWithdrawRecord{
			AppID:         appId,
			UID:           uid,
			AppOrderID:    mockId,
			EquipmentID:   equipID,
			GameAssetName: name,
			Nonce:         nonce,
			OrderStatus:   0,
		}, nil
	}

	body := &WithdrawData{
		AppID: appId,
		Params: []PrewithdrawDataRes{
			{
				UID:           uid,
				EquipmentID:   equipID,
				GameAssetName: name,
				Nonce:         nonce,
			},
		},
	}

	in := make(map[string]interface{}, 0)
	in["appId"] = body.AppID
	in["params"] = body.Params
	sidata, err := InGameSign(body.AppID, in)
	if err != nil {
		return nil, err
	}
	body.SignHash = sidata

	type Resp struct {
		common.GameResponse
		Data []PrewithdrawDataRes `json:"data"`
	}
	res := &Resp{}

	baseUrl, err := comminfo.GetBaseUrl(appId)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s%s", baseUrl, const_def.GameServer_ERC721_PreWithdraw)

	err = http_client.HttpClientReq(url, body, res)
	if err != nil {

		return nil, err
	}

	if res.Code != const_def.GAME_SERVER_SUCCESS_CODE {

		return nil, fmt.Errorf("prewithdraw get error code %d", res.Code)
	}

	if len(res.Data) != 1 {
		return nil, fmt.Errorf("prewithdraw length not match, %v", equipID)
	}

	value := res.Data[0]
	if value.GameAssetName == name && value.EquipmentID == equipID && value.Status == const_def.GAMESERVER_PASS {

		return &model.TNftWithdrawRecord{
			AppID:         appId,
			UID:           uid,
			AppOrderID:    value.AppOrderID,
			EquipmentID:   value.EquipmentID,
			GameAssetName: value.GameAssetName,
			Nonce:         nonce,
			OrderStatus:   0,
		}, nil
	}

	return nil, fmt.Errorf("request nft prewithdraw to game server failed, equipmenid is {%s} and status is {%d}", value.EquipmentID, value.Status)
}

// RequestCommitWithdrawToGame send request to game for commit withdraw
func RequestCommitWithdrawToGame(appId int, uid uint64, assetName, nonce, orderId string, op int) error {
	c := config.GetGameConfig()
	if c[appId].GameServerMock {
		return nil
	}
	baseUrl, err := comminfo.GetBaseUrl(appId)
	if err != nil {
		return fmt.Errorf("get game url %v info %v", appId, err)
	}
	url := fmt.Sprintf("%s%s", baseUrl, const_def.GameServer_ERC721_CommitWithdraw)

	body := &WithdrawData{
		AppID: appId,
		Params: []CommitWithdrawData{
			{
				GameAssetName: assetName,
				UID:           uid,
				Nonce:         nonce,
				AppOrderID:    orderId,
				Operate:       op,
			},
		},
	}

	in := make(map[string]interface{}, 0)
	in["appId"] = body.AppID
	in["params"] = body.Params
	sidata, err := InGameSign(body.AppID, in)
	if err != nil {
		return err
	}
	body.SignHash = sidata

	type Resp struct {
		common.GameResponse
		Data []CommitwithdrawDataRes `json:"data"`
	}
	res := &Resp{}
	err = http_client.HttpClientReq(url, body, res)
	if err != nil {

		return err
	}

	if res.Code != const_def.GAME_SERVER_SUCCESS_CODE {

		return fmt.Errorf("commit withdraw get error code %d", res.Code)
	}

	for _, value := range res.Data {

		if value.GameAssetName == assetName && value.UID == uid && value.AppOrderID == orderId && value.Status == const_def.GAMESERVER_PASS {
			return nil
		}
	}

	return fmt.Errorf("commit order failed, %v %v %d %s %v", assetName, uid, appId, orderId, res.Data)
}

// RequestFTDepositToGame send request to game for prewithdraw
func RequestFTDepositToGame(appId int, para *NotifyFTDepositData, url string) (*model.TFtDepositRecord, error) {
	c := config.GetGameConfig()
	if c[appId].GameServerMock {
		t := time.Now()
		mockId := t.Format("20060102150405")
		return &model.TFtDepositRecord{
			AppID:        appId,
			GameCoinName: para.GameCoinName,
			TxHash:       para.TxHash,
			UID:          uint64(para.Uid),
			OrderStatus:  const_def.CodeDepositSuccess,
			AppOrderID:   mockId,
		}, nil
	}

	item := &model.TFtDepositRecord{
		AppID:        appId,
		GameCoinName: para.GameCoinName,
		TxHash:       para.TxHash,
		UID:          uint64(para.Uid),
		OrderStatus:  const_def.CodeDepositGameFailed,
	}

	body := &NotifyFTDeposits{
		AppID:    appId,
		Params:   []NotifyFTDepositData{*para},
		SignHash: "",
	}

	in := make(map[string]interface{}, 0)
	in["appId"] = body.AppID
	in["params"] = body.Params
	sidata, err := InGameSign(body.AppID, in)
	if err != nil {
		return item, err
	}
	body.SignHash = sidata

	type Resp struct {
		common.GameResponse
		Data []GameFTDepositRes `json:"data"`
	}
	res := &Resp{}

	baseurl := fmt.Sprintf("%s%s", url, const_def.GameServer_ERC20_Deposit)

	err = http_client.HttpClientReq(baseurl, body, res)
	logger.Logrus.WithFields(logrus.Fields{"res": *res, "body": *body, "url": baseurl}).Info("RequestFTDepositToGame result")
	if err != nil {

		return item, fmt.Errorf("ft deposit send http request failed, %v", err)
	}
	if res.Code != const_def.GAME_SERVER_SUCCESS_CODE {

		return item, fmt.Errorf("ft deposit return error code %d", res.Code)
	}

	for _, v := range res.Data {
		if v.TxHash == para.TxHash && v.GameCoinName == para.GameCoinName && v.Status == const_def.GAMESERVER_PASS && v.AppOrderID != "" {
			item.AppOrderID = v.AppOrderID
			item.OrderStatus = const_def.CodeDepositSuccess
			return item, nil
		}
	}

	return item, fmt.Errorf("record not match, appid={%d} para={%v} resData={%v}", appId, para, res.Data)
}

// RequestNFTDepositToGame send request to game for nft deposit
func RequestNFTDepositToGame(appId int, para *GameNFTDepositData, url string) (*model.TNftDepositRecord, error) {
	c := config.GetGameConfig()
	if c[appId].GameServerMock {
		t := time.Now()
		mockId := t.Format("20060102150405")
		return &model.TNftDepositRecord{
			AppID:         appId,
			TxHash:        para.TxHash,
			TokenID:       para.TokenID,
			GameAssetName: para.GameAssetName,
			EquipmentID:   para.EquipmentID,
			UID:           uint64(para.Uid),
			OrderStatus:   const_def.CodeDepositSuccess,
			AppOrderID:    mockId,
		}, nil
	}

	item := &model.TNftDepositRecord{
		AppID:         appId,
		TxHash:        para.TxHash,
		TokenID:       para.TokenID,
		GameAssetName: para.GameAssetName,
		EquipmentID:   para.EquipmentID,
		UID:           uint64(para.Uid),
		OrderStatus:   const_def.CodeDepositGameFailed,
	}

	if para.Attrs == nil {
		para.Attrs = make([]commdata.EquipmentAttr, 0)
	}

	body := &NotifyNFTDeposits{
		AppID:    appId,
		Params:   []GameNFTDepositData{*para},
		SignHash: "",
	}

	if para.TxHash == "" {
		return item, fmt.Errorf("txHash is null")
	}

	if para.EquipmentID == "" {
		return item, fmt.Errorf("the equipment id of %s is null,txHash is %s", para.TokenID, para.TxHash)
	}

	if para.Uid == 0 {
		return item, fmt.Errorf("the uid of %s is null,txHash is %s", para.TokenID, para.TxHash)
	}

	in := make(map[string]interface{}, 0)
	in["appId"] = body.AppID
	in["params"] = body.Params
	sidata, err := InGameSign(body.AppID, in)
	if err != nil {
		return item, err
	}
	body.SignHash = sidata

	type Resp struct {
		common.GameResponse
		Data []GameNFTDepositRes `json:"data"`
	}
	res := &Resp{}

	baseurl := fmt.Sprintf("%s%s", url, const_def.GameServer_ERC721_Deposit)

	err = http_client.HttpClientReq(baseurl, body, res)
	logger.Logrus.WithFields(logrus.Fields{"res": *res, "body": *body, "url": baseurl}).Info("RequestNFTDepositToGame result")
	if err != nil {

		return item, fmt.Errorf("nft deposit send request failed,%v", err)
	}

	if res.Code != const_def.GAME_SERVER_SUCCESS_CODE {

		return item, fmt.Errorf("nft deposit return error code %d", res.Code)
	}

	for _, v := range res.Data {
		if v.GameAssetName == para.GameAssetName && v.TxHash == para.TxHash && v.Status == const_def.GAMESERVER_PASS && v.AppOrderID != "" {
			item.GameAssetName = v.GameAssetName
			item.AppOrderID = v.AppOrderID
			item.OrderStatus = const_def.CodeDepositSuccess
			item.EquipmentID = v.EquipmentId
			item.TokenID = v.TokenId
			return item, err
		}
	}

	return item, fmt.Errorf("record not match, appid={%d} para={%v} resData={%v}", appId, para, res.Data)
}

func RequestFTPreWithdrawToGame(appId int, uid uint64, balance string, coinName string, nonce string) (*GameFTPreWithdraResp, error) {
	c := config.GetGameConfig()
	if c[appId].GameServerMock {
		t := time.Now()
		mockId := t.Format("20060102150405")

		return &GameFTPreWithdraResp{
			GameCommonResponse: GameCommonResponse{
				Code:    const_def.GAME_SERVER_SUCCESS_CODE,
				Message: "",
			},
			Data: GameFTPreWithdrawData{
				AppOrderID: mockId,
				Status:     const_def.GAMESERVER_PASS,
			},
		}, nil
	}

	baseUrl, err := comminfo.GetBaseUrl(appId)
	if err != nil {
		return nil, err
	}
	url := baseUrl + const_def.GameServer_FT_Pre_Withdraw

	innerParams := common.PERC20PreWithdraw{
		Uid:         int64(uid),
		AppCoinName: coinName,
		Amount:      balance,
		Nonce:       nonce,
	}
	params := common.PGameERC20PreWithdraw{
		Params: innerParams,
		AppID:  appId,
		//SignHash: ,
	}

	in := make(map[string]interface{}, 0)
	in["appId"] = appId
	in["params"] = innerParams
	sidata, err := InGameSign(appId, in)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("WithdrawGameERC20Token InGame Sign Failed")
		return nil, common.NewHpError(err, int(common.InnerError), "")
	}
	params.SignHash = sidata
	logger.Logrus.WithFields(logrus.Fields{"params": params, "url": url}).Info("WithdrawGameERC20Token FTPreWithdrawHandler request")

	res := &GameFTPreWithdraResp{}
	err = http_client.HttpClientReq(url, &params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func RequestFTRecoverToGame(appId int, uid uint64, appOrderId string, coinName string, nonce string) (*GameWithdrawComfirmResp, error) {
	return RequestFTOperationToGame(appId, uid, appOrderId, coinName, nonce, const_def.NOTI_GAMESERVER_RECOVER)
}

func RequestFTOperationToGame(appId int, uid uint64, appOrderId string, coinName string, nonce string, op int) (*GameWithdrawComfirmResp, error) {
	c := config.GetGameConfig()
	if c[appId].GameServerMock {
		return &GameWithdrawComfirmResp{
			GameCommonResponse: GameCommonResponse{
				Code:    const_def.GAME_SERVER_SUCCESS_CODE,
				Message: "",
			},
			Data: []GameWithdrawComfirmData{
				{
					AppOrderID:   appOrderId,
					Nonce:        nonce,
					GameCoinName: coinName,
					Status:       const_def.GAMESERVER_PASS,
				},
			},
		}, nil
	}
	baseUrl, err := comminfo.GetBaseUrl(appId)
	if err != nil {
		return nil, err
	}
	url := baseUrl + const_def.GameServer_FT_Withdraw
	param := []common.PERC20WithdrawComfirm{{
		AppOrderID:   appOrderId,
		Nonce:        nonce,
		Uid:          int64(uid),
		GameCoinName: coinName,
		Operation:    op,
	}}

	params := common.PGameERC20WithdrawComfirm{
		Params: param,
		AppID:  int(appId),
		//SignHash: ,
	}

	secondIn := make(map[string]interface{}, 0)
	secondIn["appId"] = appId
	secondIn["params"] = param
	secondSigdata, err := InGameSign(int(appId), secondIn)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"err": err.Error()}).Error("WithdrawGameERC20Token InGameSign failed")
		return nil, err
	}

	params.SignHash = secondSigdata
	logger.Logrus.WithFields(logrus.Fields{"params": params, "url": url}).Info("WithdrawGameERC20Token FTCommitWithdrawHandler request")
	res := &GameWithdrawComfirmResp{}
	err = http_client.HttpClientReq(url, &params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
