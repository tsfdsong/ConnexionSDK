package eth

import (
	"context"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

// GetLatestHeight get blockchain latest height
func GetLatestHeight(nodeURL []string) (uint64, error) {
	for _, url := range nodeURL {
		client, err := ethclient.Dial(url)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrorMsg": err}).Error("GetLatestHeight dial")
			continue
		}
		height, err := client.BlockNumber(context.Background())
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrorMsg": err}).Error("GetLatestHeight BlockNumber")
			continue
		}

		return height, nil
	}

	return 0, fmt.Errorf("all node url is not used")
}

type EfficientBlockRangeOption struct {
	BlockHeightInterval uint64
	EfficientBlockNum   uint64
}

func GetEfficientBlockRange(nodeURL []string, start uint64, latestHeight uint64, option EfficientBlockRangeOption) (uint64, uint64, error) {

	numInter := option.BlockHeightInterval

	maxHeight := start + numInter
	end := latestHeight - option.EfficientBlockNum

	if end > maxHeight {
		end = maxHeight
	}

	if start < 1 || start > end {
		return 0, 0, fmt.Errorf("input invalided block number,%d %d", start, end)
	}

	first := start
	second := start
	for height := start; height < end; height++ {

		preblock, err := RetryEthGetBlockByNumber(nodeURL, height-1)
		if err != nil {
			return 0, 0, err
		}

		curblock, err := RetryEthGetBlockByNumber(nodeURL, height)
		if err != nil {
			return 0, 0, err
		}

		logger.Logrus.WithFields(logrus.Fields{"parentNumber": preblock.Number().Int64(), "currentNumber": curblock.Number().Int64()}).Debug("check efficient block number")

		if preblock.Hash().String() == curblock.ParentHash().String() {
			second = height - 1
		} else {
			break
		}
	}

	if first > second {
		return 0, 0, fmt.Errorf("no efficient block number, start:%d end:%d", first, second)
	}

	return first, second, nil
}

func RetryEthGetBlockByNumber(nodeURL []string, height uint64) (block *types.Block, err error) {
	for _, url := range nodeURL {
		client, dialErr := ethclient.Dial(url)
		if dialErr != nil {
			continue
		}
		block, err = client.BlockByNumber(context.Background(), big.NewInt(int64(height)))
		if err != nil {
			continue
		}
		return
	}
	return
}

func GetBlockHash(nodeURL []string, height uint64) (string, error) {

	for _, url := range nodeURL {
		client, err := ethclient.Dial(url)
		if err != nil {
			continue
		}

		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(height)))
		if err != nil {
			return "", err
		}

		if block == nil {
			return "", fmt.Errorf("block %v is null", height)
		}

		blockHash := block.Hash().String()
		return blockHash, nil
	}

	return "", fmt.Errorf("%v is not found", height)
}

func GetTxReceipt(nodeURL []string, txhash string) (*types.Receipt, error) {

	for _, url := range nodeURL {
		client, err := ethclient.Dial(url)
		if err != nil {
			continue
		}

		receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(txhash))
		if err != nil {
			return nil, fmt.Errorf("TransactionReceipt %s %v", txhash, err)
		}

		if receipt == nil {
			return nil, fmt.Errorf("TransactionReceipt %s is null", txhash)
		}

		return receipt, nil
	}

	return nil, fmt.Errorf("%v is not found", txhash)
}

func GetLogs(nodeURL []string, start, end *big.Int, contract, tops []string) ([]types.Log, error) {

	contractAddrs := make([]common.Address, 0)
	for _, v := range contract {
		addr := common.HexToAddress(v)
		contractAddrs = append(contractAddrs, addr)
	}

	topics := make([]common.Hash, 0)
	for _, v := range tops {
		topic := common.HexToHash(v)
		topics = append(topics, topic)
	}

	for _, url := range nodeURL {
		client, err := ethclient.Dial(url)
		if err != nil {
			continue
		}

		query := ethereum.FilterQuery{
			FromBlock: start,
			ToBlock:   end,
			Addresses: contractAddrs,
			Topics:    [][]common.Hash{topics},
		}

		return client.FilterLogs(context.Background(), query)
	}

	return nil, fmt.Errorf("filterlogs all node url is not used")
}
