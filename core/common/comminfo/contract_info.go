package comminfo

import (
	"encoding/json"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"strings"
)

const (
	CachePrefixNFTContractInfo = "nft.contract.info"
	CachePrefixNFTContractList = "nft.contract.list"

	CachePrefixFTContractList = "ft.contract.list"
	CachePrefixFTContractInfo = "ft.contract.info"
)

type NFTCacheData struct {
	NFTInfo []commdata.NFTContractCacheData `json:"nft_info"`
}

func GetNFTContractCache() ([]commdata.NFTContractCacheData, error) {
	nftObj, err := redis.GetString(CachePrefixNFTContractInfo)
	if err == nil {
		var res NFTCacheData
		errn := json.Unmarshal([]byte(nftObj), &res)
		if errn == nil && len(res.NFTInfo) > 0 {
			return res.NFTInfo, nil
		}
	}

	var v []model.TNftContract
	err = mysql.GetDB().Model(&model.TNftContract{}).Find(&v).Error
	if err != nil {
		return nil, fmt.Errorf("nft contract info get all record failed, %v", err)
	}

	nftTable := make([]commdata.NFTContractCacheData, 0)
	for _, val := range v {
		if val.ContractAddress == "" || val.Treasure == "" {
			continue
		}

		item := commdata.NFTContractCacheData{
			ID:              val.ID,
			AppID:           val.AppID,
			ChainID:         val.ChainID,
			Treasure:        strings.ToLower(val.Treasure),
			ContractAddress: strings.ToLower(val.ContractAddress),
			MinterAddress:   strings.ToLower(val.MinterAddress),
			TokenName:       val.TokenName,
			TokenSymbol:     val.TokenSymbol,
			TokenSupply:     val.TokenSupply,
			GameAssetName:   val.GameAssetName,
			Decimal:         val.Decimal,
			DepositSwitch:   val.DepositSwitch,
			WithdrawSwitch:  val.WithdrawSwitch,
			BaseURL:         val.BaseURL,
		}

		nftTable = append(nftTable, item)
	}

	rawRes := &NFTCacheData{
		NFTInfo: nftTable,
	}

	bytes, err := json.Marshal(rawRes)
	if err != nil {
		return nil, fmt.Errorf("nft contract info marshal failed, %v", err)
	}

	err = redis.SetString(CachePrefixNFTContractInfo, string(bytes), config.GetKeyNoExpireTime())
	if err != nil {
		return nil, fmt.Errorf("nft contract info set cache failed, {%v}", err)
	}

	return nftTable, nil
}

// GetNftContractByAssetName get nft contract info by asset name
func GetNftContractByAssetName(appid int, assetName string) (*commdata.NFTContractCacheData, error) {
	cacheList, err := GetNFTContractCache()
	if err != nil {
		return nil, fmt.Errorf("GetNftContractByAssetName GetNFTContractCache failed, %v", err)
	}

	for _, v := range cacheList {
		if v.AppID == appid && v.GameAssetName == assetName {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("GetNftContractByAssetName get {%d} not found", appid)
}

// GetNftContractByContract get nft contract info by contract address
func GetNftContractByContract(appid int, contractAddr string) (*commdata.NFTContractCacheData, error) {
	cacheList, err := GetNFTContractCache()
	if err != nil {
		return nil, fmt.Errorf("GetNftContractByContract GetNFTContractCache failed, %v", err)
	}

	for _, v := range cacheList {
		if v.AppID == appid && v.ContractAddress == strings.ToLower(contractAddr) {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("GetNftContractByContract get {%s} not found", contractAddr)
}

// GetNftContractByAppID get nft contract info by appid
func GetNftContractByAppID(appid int) (map[string]commdata.NFTContractCacheData, error) {
	cacheList, err := GetNFTContractCache()
	if err != nil {
		return nil, fmt.Errorf("GetNftContractByAppID GetNFTContractCache failed, %v", err)
	}

	res := make(map[string]commdata.NFTContractCacheData, 0)
	for _, v := range cacheList {
		if v.AppID == appid {
			res[v.GameAssetName] = v
		}
	}

	return res, nil
}

// GetNftContractByAppIDAndChainID get nft contract info by appid
func GetNftContractByAppIDAndChainID(appid, chainID int) (map[string]commdata.NFTContractCacheData, error) {
	cacheList, err := GetNFTContractCache()
	if err != nil {
		return nil, fmt.Errorf("GetNftContractByAppID GetNFTContractCache failed, %v", err)
	}

	res := make(map[string]commdata.NFTContractCacheData, 0)
	for _, v := range cacheList {
		if v.AppID == appid && v.ChainID == chainID {
			res[v.GameAssetName] = v
		}
	}

	return res, nil
}

func getNFTContractList() ([]string, error) {
	key := CachePrefixNFTContractList
	result, err := redis.LRange(key)
	if err == nil && len(result) > 0 {
		return result, nil
	}

	nftInfo, err := GetNFTContractCache()
	if err != nil {
		return nil, err
	}

	listMap := make(map[string]bool)
	for _, v := range nftInfo {
		if v.Treasure != "" {
			listMap[v.Treasure] = true
		}

		if v.MinterAddress != "" {
			listMap[v.MinterAddress] = true
		}

		if v.ContractAddress != "" {
			listMap[v.ContractAddress] = true
		}
	}

	for name := range listMap {
		result = append(result, name)
	}

	if len(result) > 0 {
		err = redis.LPush(key, config.GetKeyNoExpireTime(), result)
		if err != nil {
			return nil, fmt.Errorf("nft push key={%s} failed, {%v}", key, err)
		}
	}

	return result, nil
}

func getFTContractList() ([]string, error) {
	key := CachePrefixFTContractList
	result, err := redis.LRange(key)
	if err == nil && len(result) > 0 {
		return result, nil
	}

	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	for _, v := range ftInfo {
		result = append(result, v.Treasure, v.ContractAddress, v.DepositTreasure)
	}

	if len(result) > 0 {
		err = redis.LPush(key, config.GetKeyNoExpireTime(), result)
		if err != nil {
			return nil, fmt.Errorf("ft push key={%s} failed, {%v}", key, err)
		}
	}

	return result, nil
}

func GetAllContractList() ([]string, error) {
	res1, err := getFTContractList()
	if err != nil {
		return nil, fmt.Errorf("ft contract list %v", err)
	}

	res2, err := getNFTContractList()
	if err != nil {
		return nil, fmt.Errorf("nft contract list %v", err)
	}

	result := make([]string, 0)
	result = append(result, res1...)
	result = append(result, res2...)

	return result, nil
}

func getNFTContractListByAppId(appId int) ([]string, error) {

	nftInfo, err := GetNftContractByAppID(appId)
	if err != nil {
		return nil, err
	}

	listMap := make(map[string]bool)
	for _, v := range nftInfo {
		if v.Treasure != "" {
			listMap[v.Treasure] = true
		}

		if v.MinterAddress != "" {
			listMap[v.MinterAddress] = true
		}

		if v.ContractAddress != "" {
			listMap[v.ContractAddress] = true
		}
	}

	result := make([]string, 0)
	for name := range listMap {
		result = append(result, name)
	}

	return result, nil
}

type FTCacheData struct {
	FTInfo []commdata.FTContractCacheData `json:"ft_info"`
}

func GetFTContractCache() ([]commdata.FTContractCacheData, error) {
	ftObj, err := redis.GetString(CachePrefixFTContractInfo)
	if err == nil {
		var res FTCacheData
		errn := json.Unmarshal([]byte(ftObj), &res)
		if errn == nil && len(res.FTInfo) > 0 {
			return res.FTInfo, nil
		}
	}

	var v []model.TFtContract
	err = mysql.GetDB().Model(&model.TFtContract{}).Find(&v).Error
	if err != nil {
		return nil, fmt.Errorf("ft contract info get all record failed, %v", err)
	}

	ftTable := make([]commdata.FTContractCacheData, 0)
	for _, val := range v {
		if val.ContractAddress == "" || val.Treasure == "" {
			continue
		}

		item := commdata.FTContractCacheData{
			ID:              val.ID,
			AppID:           val.AppID,
			ChainID:         val.ChainID,
			Treasure:        strings.ToLower(val.Treasure),
			DepositTreasure: strings.ToLower(val.DepositTreasure),
			ContractAddress: strings.ToLower(val.ContractAddress),
			TokenName:       val.TokenName,
			TokenSymbol:     val.TokenSymbol,
			TokenSupply:     val.TokenSupply,
			TokenDecimal:    val.TokenDecimal,
			GameCoinName:    val.GameCoinName,
			GameDecimal:     val.GameDecimal,
			DepositSwitch:   val.DepositSwitch,
			WithdrawSwitch:  val.WithdrawSwitch,
		}

		ftTable = append(ftTable, item)
	}

	rawRes := &FTCacheData{
		FTInfo: ftTable,
	}

	bytes, err := json.Marshal(rawRes)
	if err != nil {
		return nil, fmt.Errorf("ft contract info marshal failed, %v", err)
	}

	err = redis.SetString(CachePrefixFTContractInfo, string(bytes), config.GetKeyNoExpireTime())
	if err != nil {
		return nil, fmt.Errorf("ft contract info set cache failed, {%v}", err)
	}

	return ftTable, nil
}

func GetFTContractByAppAndChain(appid int, chainId uint64) ([]commdata.FTContractCacheData, error) {
	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	res := make([]commdata.FTContractCacheData, 0)
	for _, v := range ftInfo {
		if v.AppID == appid && uint64(v.ChainID) == chainId {
			res = append(res, v)
		}
	}

	return res, nil
}

func GetFTContractByAddress(appid int, address string) (*commdata.FTContractCacheData, error) {
	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	for _, v := range ftInfo {
		if v.AppID == appid && v.ContractAddress == strings.ToLower(address) {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("ft contract info for {%s} not found", address)
}

func GetFTContractByAddressAndChainID(appid int, address string, chainID int) (*commdata.FTContractCacheData, error) {
	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	for _, v := range ftInfo {
		if v.AppID == appid && v.ContractAddress == strings.ToLower(address) && v.ChainID == chainID {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("ft contract info for {%s} not found", address)
}

func GetFirstFTContractByCoinName(appid int, coinName string) (*commdata.FTContractCacheData, error) {
	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	for _, v := range ftInfo {
		if v.AppID == appid && v.GameCoinName == coinName {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("ft contract info for coinName {%s} not found", coinName)
}

func GetFtContractFromCacheData(appid int, address string, src []commdata.FTContractCacheData) (*commdata.FTContractCacheData, error) {
	for _, v := range src {
		if v.AppID == appid && v.ContractAddress == strings.ToLower(address) {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("ft contract info for {%s} not found", address)
}

func GetFTContractListByAppIDAndChainID(appid, chainid int) ([]commdata.FTContractCacheData, error) {
	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	result := make([]commdata.FTContractCacheData, 0)
	for _, v := range ftInfo {
		if v.AppID == appid && v.ChainID == chainid {
			result = append(result, v)
		}
	}

	return result, nil
}

// get ft contract by deposit treasury address
func GetFTContractByDepositTreasureAndChain(treasure string, chainId uint64) (*commdata.FTContractCacheData, error) {
	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	for _, v := range ftInfo {
		if v.DepositTreasure == strings.ToLower(treasure) && uint64(v.ChainID) == chainId {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("ft deposit treasure info for {%s} not found", treasure)
}

func GetFTContractByTreasureAndChain(treasure string, chainId uint64) (*commdata.FTContractCacheData, error) {
	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	for _, v := range ftInfo {
		if v.Treasure == strings.ToLower(treasure) && uint64(v.ChainID) == chainId {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("ft treasure info for {%s} not found", treasure)
}

func GetFTContractAssetNameMap(appid, chainID int) (map[string]commdata.FTContractCacheData, error) {
	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	result := make(map[string]commdata.FTContractCacheData, 0)
	for _, v := range ftInfo {
		if v.AppID == appid && v.ChainID == chainID {
			result[v.GameCoinName] = v
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("GetFTContractAssetNameMap for {%d} not found", appid)
	}

	return result, nil
}

func GetFTContractAddressMap(appid, chainID int) (map[string]commdata.FTContractCacheData, error) {
	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	result := make(map[string]commdata.FTContractCacheData, 0)
	for _, v := range ftInfo {
		if v.AppID == appid && v.ChainID == chainID {
			result[v.ContractAddress] = v
		}
	}

	return result, nil
}

func DeleteFTContractList() error {
	return redis.LTrim(CachePrefixFTContractList)
}

func DeleteFTInfoCache() error {
	err := redis.DeleteString(CachePrefixFTContractInfo)
	if err != nil {
		return fmt.Errorf("DeleteFTInfoCache %v", err)
	}
	return nil
}

func DeleteNFTContractList() error {
	return redis.LTrim(CachePrefixNFTContractList)
}

func DeleteNFTInfoCache() error {
	err := redis.DeleteString(CachePrefixNFTContractInfo)
	if err != nil {
		return fmt.Errorf("DeleteNFTInfoCache %v", err)
	}
	return nil
}

func DeleteAllContractCache() error {
	//FT
	err := DeleteFTContractList()
	if err != nil {
		return err
	}

	err = DeleteFTInfoCache()
	if err != nil {
		return err
	}

	//NFT
	err = DeleteNFTContractList()
	if err != nil {
		return err
	}

	err = DeleteNFTInfoCache()
	if err != nil {
		return err
	}

	return nil
}

func GetFTContractInfoByAddress(contractAddr string, treasureAddr string) (*commdata.FTContractCacheData, error) {
	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	for _, v := range ftInfo {
		if v.ContractAddress == strings.ToLower(contractAddr) && v.Treasure == strings.ToLower(treasureAddr) {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("GetFTContractInfoByAddress contract_address:%s,treasure_address:%s not exist", contractAddr, treasureAddr)
}

func GetFTContractInfoByDepositTreasureAndChain(contractAddr string, treasureAddr string, chainId uint64) (*commdata.FTContractCacheData, error) {
	ftInfo, err := GetFTContractCache()
	if err != nil {
		return nil, err
	}

	for _, v := range ftInfo {
		if v.ContractAddress == strings.ToLower(contractAddr) && v.DepositTreasure == strings.ToLower(treasureAddr) && uint64(v.ChainID) == chainId {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("GetFTContractInfoByDepositTreasureAndChain contract_address:%s,treasure_address:%s,chain_id:%d not exist", contractAddr, treasureAddr, chainId)
}
