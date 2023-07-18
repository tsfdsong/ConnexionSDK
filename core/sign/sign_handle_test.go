package sign

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/contracts"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/http_client"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestNFTSign(t *testing.T) {
	//1. [1 2 3 4] => [5 8 6 7]
	//2. [1 2 5] => [5 5 9]
	//result: [5] => [9], [0,1] => [5, 5] , [3,2]
	withdrawAddr := "0x6cd5ff3387995d6EC2a8760851fd016C3Bd7C290"
	trease := "0x8383DB04AC31DD2061c6Fa7CD17eFA5969b829c5"
	contractAddr := "0x86F302797Cac919979736B2C177cE8Cc1723B61F"
	tokenID := "2"
	nonce := "103295893377079487058644722598565023988939674223239183021321897898802326893804"

	logger.Init("./sdk.log")

	err := config.LoadConf("../../")
	if err != nil {
		t.Errorf("LoadConf: %v", err)
		return
	}

	// verify(msg.sender, _token, _tokenID, _nonce, _attrIDs, _attrValues, _attrIndexesUpdate, _attrValuesUpdate, _attrIndexesRM, _signature)
	srcSignData, err := GetNftSignSrouceData(withdrawAddr, trease, contractAddr, tokenID, nonce)
	if err != nil {
		t.Errorf("get sign data: %v", err)
		return
	}

	//sign hash
	hash := crypto.Keccak256Hash(srcSignData)
	paramsHash := hash.Hex()

	fmt.Printf("hash: %s\n", paramsHash)

	// {"contract": [str], "type": ["ERC20", "ERC721"], "tokenId": [int], "amount": [int]}
	signData := fmt.Sprintf(`{"contract":"%s", "type":"%s", "tokenId":%s, "amount":0}`, contractAddr, "ERC721", tokenID)
	signParam, err := GenSignRequest(paramsHash, signData)
	if err != nil {
		t.Errorf("get sign request: %v", err)
		return
	}

	url := config.GetSignConfig().SignServerURL + const_def.SDK_WITHDRAW_SIGN_URL

	signRet := &RRequestSignature{}
	err = http_client.HttpClientReq(url, signParam, &signRet)
	if err != nil {
		t.Errorf("send http request: %v", err)
		return
	}

	if signRet.Code != "success" {
		t.Errorf("http resp status code: %v", signRet.Code)
		return
	}

	fmt.Printf("sign machaine return hash: %v\n", signRet.ReqHash)
	// t.Logf("object src data: %v\n", objSrcData)

	signdata, err := GetNftSignStatus(signRet.ReqHash)
	if err != nil {
		t.Errorf("GetNftSignStatus: %v", err)
		return
	}

	fmt.Printf("sign machaine return sig data: %v\n", signdata)
}

func TestNFTDepositSign(t *testing.T) {
	withdrawAddr := "0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2"
	trease := "0xa131AD247055FD2e2aA8b156A11bdEc81b9eAD95"
	contractAddr := "0x652c9ACcC53e765e1d96e2455E618dAaB79bA595"
	tokenID := "1"
	nonce := "1"

	logger.Init("./sdk.log")

	err := config.LoadConf("../../")
	if err != nil {
		t.Errorf("LoadConf: %v", err)
		return
	}

	tokenid, err := math.NewFromString(tokenID)
	if err != nil {
		t.Errorf("tokenid: %v", err)
		return
	}

	non, err := math.NewFromString(nonce)
	if err != nil {
		t.Errorf("nonce: %v", err)
		return
	}

	// verify(msg.sender, _token, _tokenID, _nonce, _attrIDs, _attrValues, _attrIndexesUpdate, _attrValuesUpdate, _attrIndexesRM, _signature)

	output, err := contracts.GetLootTreaseSignBytes("topUp", common.HexToAddress(withdrawAddr), common.HexToAddress(trease), common.HexToAddress(contractAddr), tokenid, non)
	if err != nil {
		t.Errorf("GetLootTreaseSignBytes: %v", err)
		return
	}

	//sign hash
	hash := crypto.Keccak256Hash(output)
	paramsHash := hash.Hex()

	fmt.Printf("hash: %s\n", paramsHash)

	// {"contract": [str], "type": ["ERC20", "ERC721"], "tokenId": [int], "amount": [int]}
	signData := fmt.Sprintf(`{"contract":"%s", "type":"%s", "tokenId":%s, "amount":0}`, contractAddr, "ERC721", tokenID)
	signParam, err := GenSignRequest(paramsHash, signData)
	if err != nil {
		t.Errorf("get sign request: %v", err)
		return
	}

	url := config.GetSignConfig().SignServerURL + const_def.SDK_DEPOSIT_SIGN_URL
	signRet := &RRequestSignature{}
	err = http_client.HttpClientReq(url, signParam, &signRet)
	if err != nil {
		t.Errorf("send http request: %v", err)
		return
	}

	if signRet.Code != "success" {
		t.Errorf("http resp status code: %v", signRet.Code)
		return
	}

	fmt.Printf("sign machaine return hashid: %v\n", signRet.ReqHash)

	signdata, err := GetNftSignStatus(signRet.ReqHash)
	if err != nil {
		t.Errorf("GetNftSignStatus: %v", err)
		return
	}

	fmt.Printf("GetNftSignStatus return sig data: %v\n", signdata)
}
