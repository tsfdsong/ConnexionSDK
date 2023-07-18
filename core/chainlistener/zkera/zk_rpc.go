package zkera

import (
	"context"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
	"github.com/zksync-sdk/zksync2-go"
)

const BATCH_STATUS_SEALED = "sealed"
const BATCH_STATUS_VERIFIED = "verified"

func GetLatestHeight(nodeURL []string) (uint64, error) {
	for _, url := range nodeURL {
		client, err := zksync2.NewDefaultProvider(url)
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

		logger.Logrus.WithFields(logrus.Fields{"parentNumber": preblock.Number.ToInt().Int64(), "currentNumber": curblock.Number.ToInt().Int64(), "parentHash": preblock.Hash.String(), "currentHash": curblock.ParentHash.String()}).Debug("check efficient block number")

		if preblock.Hash == curblock.ParentHash {
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

type Block struct {
	Number     hexutil.Big `json:"number"`
	Hash       common.Hash `json:"hash"`
	ParentHash common.Hash
}

func RetryEthGetBlockByNumber(nodeURL []string, height uint64) (block *Block, err error) {

	for _, url := range nodeURL {
		client, dialErr := rpc.Dial(url)
		if dialErr != nil {
			continue
		}
		//GetBlockByNumber() will return error. So use HeaderByNumber() instead
		//https://github.com/zksync-sdk/zksync2-go/issues/9
		//block, err = client.GetBlockByNumber(zksync2.BlockNumber(hexutil.EncodeUint64(height)))
		err = client.Call(&block, "eth_getBlockByNumber", hexutil.EncodeUint64(height), false)

		if err != nil {
			continue
		}
		return
	}
	return
}

func GetTxReceipt(nodeURL []string, txhash string) (*zksync2.TransactionReceipt, error) {

	for _, url := range nodeURL {
		client, err := zksync2.NewDefaultProvider(url)
		if err != nil {
			continue
		}

		receipt, err := client.GetTransactionReceipt(common.HexToHash(txhash))
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

func GetLogs(nodeURL []string, start, end *big.Int, contract, tops []string) ([]zksync2.Log, error) {

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
		provider, err := zksync2.NewDefaultProvider(url)
		if err != nil {
			continue
		}

		fromBlock := zksync2.BlockNumber(hexutil.EncodeBig(start))
		toBlock := zksync2.BlockNumber(hexutil.EncodeBig(end))
		query := zksync2.FilterQuery{
			FromBlock: &fromBlock,
			ToBlock:   &toBlock,
			Addresses: contractAddrs,
			Topics:    [][]common.Hash{topics},
		}

		return provider.GetLogs(query)
	}

	return nil, fmt.Errorf("filterlogs all node url is not used")
}

func GetZkBlockDetails(nodeURL []string, block uint32) (*zksync2.BlockDetails, error) {

	for _, url := range nodeURL {
		client, err := zksync2.NewDefaultProvider(url)
		if err != nil {
			continue
		}

		block, err := client.ZksGetBlockDetails(block)
		if err != nil {
			return nil, fmt.Errorf("ZksGetBlockDetails %v", err)
		}

		return block, nil
	}

	return nil, fmt.Errorf("block %v is not found", block)
}

func ConfirmTxVerified(nodeURL []string, txhash string) (*zksync2.TransactionReceipt, bool, error) {
	receipt, err := GetTxReceipt(nodeURL, txhash)
	if err != nil {
		return nil, false, err
	}
	if receipt.Status != 1 {
		return receipt, false, nil
	}
	blockNum := uint32(receipt.BlockNumber.Uint64())
	block, err := GetZkBlockDetails(nodeURL, blockNum)
	if err != nil || block.Status == BATCH_STATUS_SEALED {
		return receipt, false, fmt.Errorf("%s in block %d is sealing", txhash, blockNum)
	}
	return receipt, block.Status == BATCH_STATUS_VERIFIED, nil
}
