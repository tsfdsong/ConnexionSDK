package treasurealter

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/contracts"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/http_client"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/math"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"strings"

	"github.com/sirupsen/logrus"
)

type LarkText struct {
	Text string `json:"text"`
}

type LarkAlterMsg struct {
	MsgType string   `json:"msg_type"`
	Content LarkText `json:"content"`
}

// RequestTreasureAlterLark send request to game for prewithdraw
func RequestTreasureAlterLark(text string) error {
	body := &LarkAlterMsg{
		MsgType: "text",
		Content: LarkText{
			Text: text,
		},
	}

	type Resp struct {
		common.GameResponse
		Data interface{} `json:"data"`
	}
	res := &Resp{}

	url := config.GetAdminNodeURL()

	logger.Logrus.WithFields(logrus.Fields{"AlterInfo": text, "URL": url}).Debugf("RequestTreasureAlterLark info")

	err := http_client.HttpClientReq(url, body, res)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if res.Code != const_def.GAME_SERVER_SUCCESS_CODE {
		return fmt.Errorf("%d, %s", res.Code, res.Message)
	}

	return nil
}

func HandTreasureAlterLark(appId int, conAddr string, amountMap map[string][]string) error {
	ftContractInfo, err := comminfo.GetFTContractByAddress(appId, conAddr)
	if err != nil {
		return err
	}

	amounts, ok := amountMap[conAddr]
	if !ok {
		return fmt.Errorf("%s not found amounts", conAddr)
	}

	bal, err := contracts.GetERC20TokenBalance(ftContractInfo.Treasure, conAddr, int64(ftContractInfo.ChainID))
	if err != nil {
		return fmt.Errorf("get balance failed, %v", err)
	}

	amount := math.ToDecimal(bal, ftContractInfo.TokenDecimal)

	amountFormat := make([]string, 0)
	for _, v := range amounts {
		item := math.ToDecimal(v, ftContractInfo.TokenDecimal)

		amountFormat = append(amountFormat, item.String())
	}

	disAmt := strings.Join(amountFormat, " ,")

	text := fmt.Sprintf("%s withdrawal treasury balance:{ %s }, Amount: { %s }", ftContractInfo.TokenSymbol, amount.String(), disAmt)

	err = RequestTreasureAlterLark(text)
	if err != nil {
		return fmt.Errorf("admin alter failed, %v", err)
	}

	return nil
}

func HandDepositTreasureAlterLark(appId int, conAddr string, amountMap map[string][]string) error {
	ftContractInfo, err := comminfo.GetFTContractByAddress(appId, conAddr)
	if err != nil {
		return err
	}

	amounts, ok := amountMap[conAddr]
	if !ok {
		return fmt.Errorf("%s not found amounts", conAddr)
	}

	bal, err := contracts.GetERC20TokenBalance(ftContractInfo.DepositTreasure, conAddr, int64(ftContractInfo.ChainID))
	if err != nil {
		return fmt.Errorf("get deposit treasure balance failed, %v", err)
	}

	amount := math.ToDecimal(bal, ftContractInfo.TokenDecimal)

	amountFormat := make([]string, 0)
	for _, v := range amounts {
		item := math.ToDecimal(v, ftContractInfo.TokenDecimal)

		amountFormat = append(amountFormat, item.String())
	}

	disAmt := strings.Join(amountFormat, " ,")

	text := fmt.Sprintf("%s deposit treasury balance:{ %s }, Amount: { %s }", ftContractInfo.TokenSymbol, amount.String(), disAmt)

	err = RequestTreasureAlterLark(text)
	if err != nil {
		return fmt.Errorf("deposit alter failed, %v", err)
	}

	return nil
}
