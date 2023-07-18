package sign

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/contracts"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/math"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
)

func getIndexAndValues(oldIds, newIds []uint64, newValues []string) ([]uint64, []string, []int, []string, []int) {
	msrc := make(map[uint64]int)
	mall := make(map[uint64]int)
	mdest := make(map[uint64]int)
	set := make(map[uint64]int)

	for i, v := range oldIds {
		msrc[v] = i
		mall[v] = i
	}

	for j, v := range newIds {
		mdest[v] = j

		l := len(mall)
		mall[v] = j
		if l == len(mall) {
			set[v] = j
		}
	}

	for k := range set {
		delete(mall, k)
	}

	added := make([]uint64, 0)
	deleted := make([]int, 0)
	for v := range mall {
		j, exist := msrc[v]
		if exist {
			deleted = append(deleted, j)
		} else {
			added = append(added, v)
		}
	}

	updatedIndex := make([]int, 0)
	updateValue := make([]string, 0)
	for k := range set {
		j, exist := msrc[k]
		if exist {
			updatedIndex = append(updatedIndex, j)
		}

		i, ok := mdest[k]
		if ok {
			updateValue = append(updateValue, newValues[i])
		}
	}

	newValue := make([]string, 0)
	for _, v := range added {
		i, ok := mdest[v]
		if ok {
			newValue = append(newValue, newValues[i])
		}
	}

	return added, newValue, updatedIndex, updateValue, deleted
}

//GetNftLootSignSrouceData get sign source data
func GetNftLootSignSrouceData(sender, trease, nftcontract, tokenid, nonce string, oldIds, newIds []uint64, newValues []string) ([]byte, []byte, error) {
	if len(newIds) != len(newValues) {
		return nil, nil, fmt.Errorf("input length is not match")
	}
	attrIDs, attrValues, attrIndexsUpdate, attrValuesUpdate, attrIndexsDelete := getIndexAndValues(oldIds, newIds, newValues)

	//delete index sort desc
	sort.Sort(sort.Reverse(sort.IntSlice(attrIndexsDelete)))

	//seralize object
	srcData := &commdata.NftSignatureSrcData{
		AttrIDs:          attrIDs,
		AttrValues:       attrValues,
		AttrIndexsUpdate: attrIndexsUpdate,
		AttrValuesUpdate: attrValuesUpdate,
		AttrIndexsRM:     attrIndexsDelete,
	}

	srcBytes, err := srcData.MarshalJson()
	if err != nil {
		return nil, nil, err
	}

	if len(attrIDs) != len(attrValues) {
		return nil, nil, fmt.Errorf("id and value length is not match")
	}

	if len(attrIndexsUpdate) != len(attrValuesUpdate) {
		return nil, nil, fmt.Errorf("id and value length is not match")
	}

	ids := make([]*big.Int, 0)
	for _, v := range attrIDs {
		ids = append(ids, big.NewInt(int64(v)))
	}

	values := make([]*big.Int, 0)
	for _, v := range attrValues {
		b, err := math.NewFromString(v)
		if err != nil {
			return nil, nil, err
		}

		values = append(values, b)
	}

	inUpdates := make([]*big.Int, 0)
	for _, v := range attrIndexsUpdate {
		inUpdates = append(inUpdates, big.NewInt(int64(v)))
	}

	valUpdate := make([]*big.Int, 0)
	for _, v := range attrValuesUpdate {
		b, err := math.NewFromString(v)
		if err != nil {
			return nil, nil, err
		}

		valUpdate = append(valUpdate, b)
	}

	inDelete := make([]*big.Int, 0)
	for _, v := range attrIndexsDelete {
		b := big.NewInt(int64(v))

		inDelete = append(inDelete, b)
	}

	tokenID, err := math.NewFromString(tokenid)
	if err != nil {
		return nil, nil, err
	}

	non, err := math.NewFromString(nonce)
	if err != nil {
		return nil, nil, err
	}

	output, err := contracts.GetLootTreaseSignBytes("upChainLoot", common.HexToAddress(sender),
		common.HexToAddress(trease), common.HexToAddress(nftcontract), tokenID, non, ids, values, inUpdates, valUpdate, inDelete)
	if err != nil {
		return nil, nil, err
	}
	return output, srcBytes, nil
}

//GetNftSignSrouceData get sign source data
func GetNftSignSrouceData(sender, trease, nftcontract, tokenid, nonce string) ([]byte, error) {

	tokenID, err := math.NewFromString(tokenid)
	if err != nil {
		return nil, err
	}

	non, err := math.NewFromString(nonce)
	if err != nil {
		return nil, err
	}

	output, err := contracts.GetLootTreaseSignBytes("upChain", common.HexToAddress(sender),
		common.HexToAddress(trease), common.HexToAddress(nftcontract), tokenID, non)
	if err != nil {
		return nil, err
	}
	return output, nil
}

//GetNftGameMintSignLootSrouceData get sign source data for gameMint of equipment contract
func GetNftGameMintSignLootSrouceData(sender, nftcontract, nonce, equipID string, newIds []uint64, newValues []string) ([]byte, error) {
	if len(newIds) != len(newValues) {
		return nil, fmt.Errorf("input length is not match")
	}

	ids := make([]*big.Int, 0)
	for _, v := range newIds {
		ids = append(ids, big.NewInt(int64(v)))
	}

	values := make([]*big.Int, 0)
	for _, v := range newValues {
		b, err := math.NewFromString(v)
		if err != nil {
			return nil, err
		}

		values = append(values, b)
	}

	non, err := math.NewFromString(nonce)
	if err != nil {
		return nil, err
	}

	eid, err := math.NewFromString(equipID)
	if err != nil {
		return nil, err
	}

	output, err := contracts.GetLootTreaseSignBytes("gameMintLoot", common.HexToAddress(sender), common.HexToAddress(nftcontract), non, eid, ids, values)
	if err != nil {
		return nil, err
	}
	return output, nil
}

//GetNftGameMintSignSrouceData get sign source data for gameMint of erc721 contract
func GetNftGameMintSignSrouceData(sender, nftcontract, equipID, nonce string) ([]byte, error) {
	eid, err := math.NewFromString(equipID)
	if err != nil {
		return nil, err
	}

	non, err := math.NewFromString(nonce)
	if err != nil {
		return nil, err
	}

	output, err := contracts.GetLootTreaseSignBytes("gameMint", common.HexToAddress(sender), common.HexToAddress(nftcontract), eid, non)
	if err != nil {
		return nil, err
	}
	return output, nil
}
