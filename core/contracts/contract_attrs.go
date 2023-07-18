package contracts

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

// GameTemplateMetaData contains all meta data concerning the GameTemplate contract.
var GameTemplateMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wallet_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"eqID_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce_\",\"type\":\"uint256\"}],\"name\":\"gameMint\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wallet_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"eqID_\",\"type\":\"uint256\"},{\"internalType\":\"uint128[]\",\"name\":\"attrIDs_\",\"type\":\"uint128[]\"},{\"internalType\":\"uint128[]\",\"name\":\"attrValues_\",\"type\":\"uint128[]\"}],\"name\":\"gameMintLoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_this\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"upChain\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_this\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint128[]\",\"name\":\"_attrIDs\",\"type\":\"uint128[]\"},{\"internalType\":\"uint128[]\",\"name\":\"_attrValues\",\"type\":\"uint128[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_attrIndexesUpdate\",\"type\":\"uint256[]\"},{\"internalType\":\"uint128[]\",\"name\":\"_attrValuesUpdate\",\"type\":\"uint128[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_attrIndexesRMs\",\"type\":\"uint256[]\"}],\"name\":\"upChainLoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// GameTemplateABI is the input ABI used to generate the binding from.
// Deprecated: Use GameTemplateMetaData.ABI instead.
var GameTemplateABI = GameTemplateMetaData.ABI

//GetAttrsFromChain get attributes from contract by tokenid
func GetAttrsFromChain(contractAddr, tokenID string) ([]commdata.EquipmentAttr, error) {
	// Connect to a geth node (when using Infura, you need to use your own API key)
	url := config.GetNodeURL()
	if len(url) < 1 {
		return nil, fmt.Errorf("node url is not config")
	}
	conn, err := ethclient.Dial(url[0])
	if err != nil {
		return nil, err
	}

	// Instantiate the contract and display its name
	address := common.HexToAddress(contractAddr)
	inst, err := NewGameLootEquipment(address, conn)
	if err != nil {
		return nil, err
	}

	id, err := math.NewFromString(tokenID)
	if err != nil {
		return nil, err
	}

	newAttrs, err := inst.Attributes(&bind.CallOpts{}, id)
	if err != nil {
		return nil, err
	}

	attrs := make([]commdata.EquipmentAttr, 0)
	for _, e := range newAttrs {
		item := commdata.EquipmentAttr{
			AttributeID:    e.AttrID.Uint64(),
			AttributeValue: e.AttrValue.String(),
		}

		attrs = append(attrs, item)
	}

	return attrs, err
}

//GetAttrs get attributes from contract by tokenid
func GetDecimal(contractAddr string) (uint8, error) {
	// Connect to a geth node (when using Infura, you need to use your own API key)
	url := config.GetNodeURL()
	if len(url) < 1 {
		return 0, fmt.Errorf("node url is not config")
	}
	conn, err := ethclient.Dial(url[0])
	if err != nil {
		return 0, err
	}

	// Instantiate the contract and display its name
	address := common.HexToAddress(contractAddr)
	inst, err := NewGameLootEquipment(address, conn)
	if err != nil {
		return 0, err
	}

	dec, err := inst.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, err
	}

	return dec, nil
}

//ConvertAttrsWithDecimal (attr value * decimal)
func ConvertAttrsWithDecimal(decimal int, attrs []commdata.EquipmentAttr) ([]commdata.EquipmentAttr, error) {
	//the first element of attrs is special handed
	if len(attrs) < 1 {
		return nil, fmt.Errorf("attrs is invalied")
	}

	first := attrs[0]

	//special attribute id is 0
	if first.AttributeID != 0 {
		return nil, fmt.Errorf("attr id has error format, %d", first.AttributeID)
	}

	distinctMap := make(map[uint64]*big.Int, 0)
	for _, e := range attrs[1:] {
		value := math.ToWei(e.AttributeValue, decimal)
		plusValue, ok := distinctMap[e.AttributeID]
		if !ok {
			distinctMap[e.AttributeID] = value
		} else {
			distinctMap[e.AttributeID] = new(big.Int).Add(value, plusValue)
		}
	}

	resattrs := make([]commdata.EquipmentAttr, 0)
	resattrs = append(resattrs, first)
	for k, v := range distinctMap {
		item := commdata.EquipmentAttr{
			AttributeID:    k,
			AttributeValue: v.String(),
		}

		resattrs = append(resattrs, item)
	}

	return resattrs, nil
}

//GetLootTreaseSignBytes get input bytes
func GetLootTreaseSignBytes(method string, params ...interface{}) ([]byte, error) {
	parsed, err := abi.JSON(strings.NewReader(GameTemplateABI))
	if err != nil {
		return nil, err
	}

	// Otherwise pack up the parameters and invoke the contract
	input, err := parsed.Pack(method, params...)
	if err != nil {
		return nil, err
	}

	//remove first 4 bytes of method id
	return input[4:], nil
}

//GetOwnerOf get owner of tokenid
func GetOwnerOf(contractAddr, tokenID string) (string, error) {
	// Connect to a geth node (when using Infura, you need to use your own API key)
	urllist := config.GetNodeURL()

	for _, url := range urllist {
		conn, err := ethclient.Dial(url)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ContractAddress": contractAddr, "TokenID": tokenID, "err": err}).Error("GetOwnerOf dial failed")
			continue
		}

		// Instantiate the contract and display its name
		address := common.HexToAddress(contractAddr)
		inst, err := NewGameLootEquipment(address, conn)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ContractAddress": contractAddr, "TokenID": tokenID, "err": err}).Error("GetOwnerOf NewGameLootEquipment failed")
			continue
		}

		id, err := math.NewFromString(tokenID)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ContractAddress": contractAddr, "TokenID": tokenID, "err": err}).Error("GetOwnerOf NewFromString failed")
			continue
		}

		tokenAddress, err := inst.OwnerOf(&bind.CallOpts{}, id)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ContractAddress": contractAddr, "TokenID": tokenID, "err": err}).Error("GetOwnerOf OwnerOf failed")
			continue
		}

		return tokenAddress.String(), nil
	}

	return "", fmt.Errorf("all node is invalied, %s, %s", contractAddr, tokenID)
}

//GetTokenURI get  token uri of tokenid
func GetTokenURI(contractAddr, tokenID string) (string, error) {
	// Connect to a geth node (when using Infura, you need to use your own API key)
	urllist := config.GetNodeURL()

	for _, url := range urllist {
		conn, err := ethclient.Dial(url)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ContractAddress": contractAddr, "TokenID": tokenID, "err": err}).Error("GetTokenURI dial failed")
			continue
		}

		// Instantiate the contract and display its name
		address := common.HexToAddress(contractAddr)
		inst, err := NewGameLootEquipment(address, conn)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ContractAddress": contractAddr, "TokenID": tokenID, "err": err}).Error("GetTokenURI NewGameLootEquipment failed")
			continue
		}

		id, err := math.NewFromString(tokenID)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ContractAddress": contractAddr, "TokenID": tokenID, "err": err}).Error("GetTokenURI NewFromString failed")
			continue
		}

		tokenURI, err := inst.TokenURI(&bind.CallOpts{}, id)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ContractAddress": contractAddr, "TokenID": tokenID, "err": err}).Error("GetTokenURI OwnerOf failed")
			continue
		}

		return tokenURI, nil
	}

	return "", fmt.Errorf("GetTokenURI all node is invalied, %s, %s", contractAddr, tokenID)
}
