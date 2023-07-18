package contracts

import (
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"math/big"
	"sort"
	"strings"

	webcommon "github/Connector-Gamefi/ConnectorGoSDK/web/common"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func GetEquipmentTokenURI(contract string, tokenId string) (error, string) {
	nodeURL := config.GetNodeURL()
	if len(nodeURL) == 0 {
		return errors.New("empty node url"), ""
	}
	for _, x := range nodeURL {
		client, err := ethclient.Dial(x)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"url": x, "err": err.Error()}).Error("GetEquipmentTokenURI::dial failed")
			continue
		}

		contractAddress := common.HexToAddress(contract)
		instance, err := NewGameLootEquipment(contractAddress, client)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"url": x, "err": err.Error()}).Error("GetEquipmentTokenURI::NewGameLootEquipment failed")
			continue
		}

		tokenID, flag := big.NewInt(0).SetString(tokenId, 0)
		if !flag {
			return errors.New("invalid tokenID"), ""
		}

		r, err := instance.TokenURI(nil, tokenID)
		return err, r
	}
	return errors.New("ALL Query Failed"), ""
}

func GetERC721ApprovedForAll(contract, owner, operator string) (error, bool) {
	nodeURL := config.GetNodeURL()
	if len(nodeURL) == 0 {
		return errors.New("empty node url"), false
	}
	for _, x := range nodeURL {
		client, err := ethclient.Dial(x)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"url": x, "err": err.Error()}).Error("GetERC721TokenInfo::dial failed")
			continue
		}

		contractAddress := common.HexToAddress(contract)
		instance, err := NewGameLootEquipment(contractAddress, client)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"url": x, "err": err.Error()}).Error("GetERC721TokenInfo::NewGameLootEquipment failed")
			continue
		}

		ownerAddress := common.HexToAddress(owner)
		operatorAddress := common.HexToAddress(operator)
		r, err := instance.IsApprovedForAll(nil, ownerAddress, operatorAddress)
		return err, r
	}
	return errors.New("ALL Query Failed"), false
}

func GetERC721TokenInfo(addr string) (string, string, string, uint8, error) {
	nodeURL := config.GetNodeURL()
	count := 0
	for _, x := range nodeURL {
		client, err := ethclient.Dial(x)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"url": x, "err": err.Error()}).Error("GetERC721TokenInfo::dial failed")
			continue
		}

		contractAddress := common.HexToAddress(addr)
		instance, err := NewGameLootEquipment(contractAddress, client)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"url": x, "err": err.Error()}).Error("GetERC721TokenInfo::NewGameLootEquipment failed")
			continue
		}

		r1, err1 := instance.Name(nil)
		r2, err2 := instance.Symbol(nil)
		r3, err3 := instance.TotalSupply(nil)
		r4, err4 := instance.Decimals(nil)
		if err1 == nil && err2 == nil && err3 == nil && r3 != nil && err4 == nil {
			return r1, r2, r3.String(), r4, nil
		} else {
			logger.Logrus.WithFields(logrus.Fields{"err1": err1, "err2": err2, "err3": err3, "err4": err4}).Error("GetERC721TokenInfo::get info failed")
			count++
			if count >= 2 {
				return "", "", "", 0, fmt.Errorf("incorrect contract %s", addr)
			}
		}
	}
	return "", "", "", 0, fmt.Errorf("contract info %s not found", addr)

}

func GetERC20TokenInfo(addr string) (string, string, string, int, error) {
	nodeURL := config.GetNodeURL()

	count := 0
	for _, x := range nodeURL {
		client, err := ethclient.Dial(x)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"url": x, "err": err.Error()}).Error("GetERC20TokenInfo::dial failed")
			continue
		}

		contractAddress := common.HexToAddress(addr)
		instance, err := NewGameErc20Token(contractAddress, client)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"url": x, "err": err.Error()}).Error("GetERC20TokenInfo::NewGameErc20Token failed")
			continue
		}

		r1, err1 := instance.Name(nil)
		r2, err2 := instance.Symbol(nil)
		r3, err3 := instance.Decimals(nil)
		r4, err4 := instance.Cap(nil)

		if err1 == nil && err2 == nil && err3 == nil && err4 == nil && r4 != nil {

			if r3 != 0 {
				decimal, flag := big.NewInt(0).SetString(fmt.Sprintf("1%+v", strings.Repeat("0", int(r3))), 0)
				if !flag {
					return "", "", "", 0, errors.New("Set String Failed")
				}
				cap := big.NewInt(0).Div(r4, decimal)
				return r1, r2, cap.String(), int(r3), nil
			} else {
				return r1, r2, r4.String(), int(r3), nil
			}

		} else {
			logger.Logrus.WithFields(logrus.Fields{"err1": err1, "err2": err2, "err3": err3, "err4": err4}).Error("GetERC20TokenInfo::get info failed")
			count++
			if count >= 2 {
				return "", "", "", 0, errors.New("Incorrect Contract")
			}
		}
	}
	return "", "", "", 0, errors.New("ALL Query Failed")
}

func getNodeURLByChainID(chainID int64) []string {
	if chainID != 97 && chainID != 56 {
		return config.GetZKNodeURL()
	} else {
		return config.GetNodeURL()
	}
}

func GetERC20TokenBalance(userAddr, contractAddr string, chainID int64) (string, error) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("GetERC20TokenBalance failed")
		}
	}()

	nodeURL := getNodeURLByChainID(chainID)
	contractAddress := common.HexToAddress(contractAddr)
	userAddress := common.HexToAddress(userAddr)
	index := 0
	nodeLen := len(nodeURL)
	errCount := 0
	zeroCount := 0
	for {
		client, err := ethclient.Dial(nodeURL[index])
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"url": nodeURL[index], "contract": contractAddr, "err": err.Error()}).Error("GetERC20TokenBalance::dial failed")
			continue
		}

		instance, err := NewGameErc20Token(contractAddress, client)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"url": nodeURL[index], "contract": contractAddr, "err": err.Error()}).Error("GetERC20TokenBalance::NewGameErc20Token failed")
			continue
		}

		r, err := instance.BalanceOf(nil, userAddress)

		if err == nil && r != nil {
			if r.Cmp(big.NewInt(0)) == 0 {
				zeroCount++
				if zeroCount >= 2 {
					return "0", nil
				} else {
					index = (index + 1) % nodeLen
					continue
				}

			} else {
				return r.String(), nil
			}

		} else {
			ferr := err
			if err == nil && r == nil {
				ferr = errors.New("get balance nil response")
			}
			logger.Logrus.WithFields(logrus.Fields{"err": err, "contract": contractAddr}).Error("GetERC20TokenBalance::get balance failed")
			errCount++
			if errCount >= 2 {
				return "0", ferr
			} else {
				index = (index + 1) % nodeLen
				continue
			}
		}
	}

}

func GetGameAllErc20TokenBalance(addr string, contracts map[string]commdata.FTContractCacheData) ([]webcommon.RChainFTAssets, error) {
	if len(contracts) == 0 {
		return []webcommon.RChainFTAssets{}, nil
	}

	retData := []webcommon.RChainFTAssets{}

	appCoinNameList := []string{}
	newContractMap := map[string]commdata.FTContractCacheData{}
	for _, e := range contracts {
		appCoinNameList = append(appCoinNameList, e.GameCoinName)
		newContractMap[e.GameCoinName] = e
	}
	sort.Strings(appCoinNameList)

	for _, e := range appCoinNameList {
		v, ok := newContractMap[e]
		if !ok {
			logger.Logrus.Error("what happend!")
			continue
		}
		balance, err := GetERC20TokenBalance(addr, v.ContractAddress, int64(v.ChainID))
		if err != nil {
			return []webcommon.RChainFTAssets{}, err
		}
		displayBalance, err := tools.GetTokenAmount2(balance, int32(v.TokenDecimal), 8)
		if err != nil {
			return []webcommon.RChainFTAssets{}, err
		}

		item := webcommon.RChainFTAssets{
			Contract:        strings.ToLower(v.ContractAddress),
			Symbol:          v.TokenSymbol,
			Treasure:        strings.ToLower(v.Treasure),
			DepositTreasure: strings.ToLower(v.DepositTreasure),
			GameDecimal:     v.GameDecimal,
			Decimal:         v.TokenDecimal,
			Balance:         displayBalance,
			OriginBalance:   balance,
			ChainID:         v.ChainID,
		}
		retData = append(retData, item)
	}
	return retData, nil
}
