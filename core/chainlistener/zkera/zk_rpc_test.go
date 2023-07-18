package zkera_test

import (
	"encoding/json"
	"github/Connector-Gamefi/ConnectorGoSDK/core/chainlistener/zkera"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"math/big"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGetLatestHeight(t *testing.T) {
	node := []string{"https://testnet.era.zksync.dev"}
	height, err := zkera.GetLatestHeight(node)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(height)
}

func TestGetEfficientBlockRange(t *testing.T) {
	logger.Logrus = logrus.New()
	logger.Logrus.SetOutput(os.Stdout)
	node := []string{"https://testnet.era.zksync.dev"}
	start, end, err := zkera.GetEfficientBlockRange(node, 5145518, 5145528, zkera.EfficientBlockRangeOption{BlockHeightInterval: 1000, EfficientBlockNum: 0})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Start Block:", start, " End block:", end)
}

func TestRetryEthGetBlockByNumber(t *testing.T) {
	logger.Logrus = logrus.New()
	logger.Logrus.SetOutput(os.Stdout)
	node := []string{"https://testnet.era.zksync.dev"}
	block, err := zkera.RetryEthGetBlockByNumber(node, 4547014)
	if err != nil {
		t.Error(err)
		return
	}

	b, _ := json.Marshal(block)
	t.Log(string(b))
}
func TestGetTxReceipt(t *testing.T) {
	node := []string{"https://testnet.era.zksync.dev"}
	log, err := zkera.GetTxReceipt(node, "0x609bd10f2464a2e227291abcaad1567ac8a782b63f6430e3aab36b5f03f1a116")
	if err != nil {
		t.Error(err)
		return
	}
	b, _ := json.Marshal(log)
	t.Log(string(b))
}

func TestGetLogs(t *testing.T) {
	node := []string{"https://testnet.era.zksync.dev"}
	contract := []string{"0xEF9E47AB5c9570a345C5069eb72B3F3D3283B33F"}
	//Transfer(address,address,uint256)
	topics := []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}
	log, err := zkera.GetLogs(node, big.NewInt(4636753), big.NewInt(4636755), contract, topics)
	if err != nil {
		t.Error(err)
		return
	}
	b, _ := json.Marshal(log)
	t.Log(string(b))
}

func TestGetZkBlockDetails(t *testing.T) {
	node := []string{"https://testnet.era.zksync.dev"}
	block, err := zkera.GetZkBlockDetails(node, 5189999)
	if err != nil {
		t.Error(err)
		return
	}
	b, _ := json.Marshal(block)
	t.Log(string(b))
}
