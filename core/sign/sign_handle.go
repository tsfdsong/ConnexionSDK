package sign

import (
	"context"
	"fmt"
	"strings"

	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/contracts"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/distmiddleware/skywalking"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"

	"gorm.io/datatypes"
)

//NFTWithdrawSign send signature request to sign machine
func NFTWithdrawSign(p *model.TNftWithdrawRecord, attrList []commdata.EquipmentAttr, ctx context.Context) (*model.TNftWithdrawRecord, error) {
	if p.TokenID == "" {
		return nftWithdrawSignWithMint(p, attrList, ctx)
	}

	return nftWithdrawSignWithoutMint(p, attrList, ctx)
}

//nftWithdrawSignWithMint sign on gameMint
func nftWithdrawSignWithMint(p *model.TNftWithdrawRecord, attrList []commdata.EquipmentAttr, ctx context.Context) (*model.TNftWithdrawRecord, error) {
	res := p.Copy()

	srcSignData := make([]byte, 0)
	var err error

	if p.GameAssetName == config.GetGameEquipmentName() {
		if len(attrList) == 0 {
			return nil, fmt.Errorf("attrs list is null")
		}

		newIDs := make([]uint64, len(attrList))
		newValues := make([]string, len(attrList))
		for i, v := range attrList {
			newIDs[i] = v.AttributeID
			newValues[i] = v.AttributeValue
		}

		logger.Logrus.WithFields(logrus.Fields{"OrderID": p.ID, "WithdrawAddress": p.WithdrawAddress, "GameMinterAddress": p.GameMinterAddress, "Nonce": p.Nonce, "NewAttrIDs": newIDs, "NewValues": newValues}).Info("NFT Withdraw GameMint Loot sign input")

		// verify(msg.sender, _token, _tokenID, _nonce, _attrIDs, _attrValues, _signature)
		srcSignData, err = GetNftGameMintSignLootSrouceData(p.WithdrawAddress, p.GameMinterAddress, p.Nonce, p.EquipmentID, newIDs, newValues)
		if err != nil {
			return nil, fmt.Errorf("get loot sign failed, %v", err)
		}

		srcData := &commdata.NftSignatureSrcData{
			AttrIDs:          newIDs,
			AttrValues:       newValues,
			AttrIndexsUpdate: make([]int, 0),
			AttrValuesUpdate: make([]string, 0),
			AttrIndexsRM:     make([]int, 0),
		}
		objSrcData, err := srcData.MarshalJson()
		if err != nil {
			return nil, fmt.Errorf("marshal attrs, %v", err)
		}
		res.SignatureSrc = datatypes.JSON(objSrcData)
	} else {
		if len(attrList) != 0 {
			return nil, fmt.Errorf("wrong attrs list")
		}

		logger.Logrus.WithFields(logrus.Fields{"OrderID": p.ID, "WithdrawAddress": p.WithdrawAddress, "GameMinterAddress": p.GameMinterAddress, "Nonce": p.Nonce, "EquipmentID": p.EquipmentID}).Info("NFT Withdraw GameMint sign input")

		// verify(msg.sender, _token, _tokenID,_eid, _nonce)
		srcSignData, err = GetNftGameMintSignSrouceData(p.WithdrawAddress, p.ContractAddress, p.EquipmentID, p.Nonce)
		if err != nil {
			return nil, fmt.Errorf("get sign data failed, %v", err)
		}
	}

	if srcSignData == nil {
		return nil, fmt.Errorf("source sign data is null")
	}

	//sign hash
	hash := crypto.Keccak256Hash(srcSignData)
	paramsHash := hash.Hex()

	// {"contract": [str], "type": ["ERC20", "ERC721"], "tokenId": [int], "amount": [int]}
	signData := fmt.Sprintf(`{"contract":"%s", "type":"%s", "tokenId":%s, "amount":0}`, p.ContractAddress, "ERC721", p.TokenID)
	if p.TokenID == "" {
		signData = fmt.Sprintf(`{"contract":"%s", "type":"%s", "amount":0}`, p.ContractAddress, "ERC721")
	}

	signParam, err := GenSignRequest(paramsHash, signData)
	if err != nil {
		return nil, err
	}

	url := config.GetSignConfig().SignServerURL + const_def.SDK_WITHDRAW_SIGN_URL

	signRet := &RRequestSignature{}
	err = skywalking.SkyPostRequest(url, signParam, &signRet, ctx)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": p.ID}).Error("NFT Withdraw GameMint request signature machine hash failed")

		return nil, fmt.Errorf("post sign request, %v", err)
	}

	if signRet.Code != "success" {
		return nil, fmt.Errorf("sign request return failed status code")
	}

	res.SignatureHash = signRet.ReqHash

	logger.Logrus.WithFields(logrus.Fields{"OrderID": p.ID, "SignInputHash": paramsHash, "SignatureHash": signRet.ReqHash}).Info("NFT Withdraw GameMint sign result hash")

	return res, nil
}

//nftWithdrawSignWithoutMint get sign data without mint,but upchain
func nftWithdrawSignWithoutMint(p *model.TNftWithdrawRecord, attrList []commdata.EquipmentAttr, ctx context.Context) (*model.TNftWithdrawRecord, error) {
	res := p.Copy()

	srcSignData := make([]byte, 0)
	objSrcData := make([]byte, 0)
	var err error
	if p.GameAssetName == config.GetGameEquipmentName() {
		if len(attrList) == 0 {
			return nil, fmt.Errorf("attrs list is null")
		}

		//generate input parameter
		oldAttrs, err := contracts.GetAttrsFromChain(p.ContractAddress, p.TokenID)
		if err != nil {
			return nil, fmt.Errorf("get loot attrs, %v", err)
		}

		//old attribute array
		oldIDs := make([]uint64, len(oldAttrs))
		for i, val := range oldAttrs {
			oldIDs[i] = val.AttributeID
		}

		newIDs := make([]uint64, len(attrList))
		newValues := make([]string, len(attrList))
		for i, v := range attrList {
			key := v.AttributeID
			newIDs[i] = key
			newValues[i] = v.AttributeValue
		}

		// verify(msg.sender, _token, _tokenID, _nonce, _attrIDs, _attrValues, _attrIndexesUpdate, _attrValuesUpdate, _attrIndexesRM, _signature)
		srcSignData, objSrcData, err = GetNftLootSignSrouceData(strings.ToLower(p.WithdrawAddress), p.TreaseAddress, strings.ToLower(p.ContractAddress), p.TokenID, p.Nonce, oldIDs, newIDs, newValues)
		if err != nil {
			return nil, fmt.Errorf("get loot sign source data, %v", err)
		}

		res.SignatureSrc = datatypes.JSON(objSrcData)
	} else {
		if len(attrList) != 0 {
			return nil, fmt.Errorf("wrong attrs list")
		}

		// verify(msg.sender, _token, _tokenID, _nonce)
		srcSignData, err = GetNftSignSrouceData(strings.ToLower(p.WithdrawAddress), p.TreaseAddress, strings.ToLower(p.ContractAddress), p.TokenID, p.Nonce)
		if err != nil {
			return nil, fmt.Errorf("get sign source data, %v", err)
		}
	}

	if srcSignData == nil {
		return nil, fmt.Errorf("source sign data is null")
	}

	//sign hash
	hash := crypto.Keccak256Hash(srcSignData)
	paramsHash := hash.Hex()

	// {"contract": [str], "type": ["ERC20", "ERC721"], "tokenId": [int], "amount": [int]}
	signData := fmt.Sprintf(`{"contract":"%s", "type":"%s", "tokenId":%s, "amount":0}`, p.ContractAddress, "ERC721", p.TokenID)

	signParam, err := GenSignRequest(paramsHash, signData)
	if err != nil {
		return nil, fmt.Errorf("gen sign hash {%v} request,%v", paramsHash, err)
	}

	url := config.GetSignConfig().SignServerURL + const_def.SDK_WITHDRAW_SIGN_URL

	signRet := &RRequestSignature{}
	err = skywalking.SkyPostRequest(url, signParam, &signRet, ctx)
	if err != nil {

		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "OrderID": p.ID}).Error("NFT Withdraw upChain request signature machine hash failed")

		return nil, fmt.Errorf("post sign hash {%v} request,%v", paramsHash, err)
	}

	if signRet.Code != "success" {

		return nil, fmt.Errorf("send sign hash {%v} request failed", paramsHash)
	}

	res.SignatureHash = signRet.ReqHash

	logger.Logrus.WithFields(logrus.Fields{"OrderID": p.ID, "SignInputHash": paramsHash, "SignatureHash": signRet.ReqHash}).Info("NFT Withdraw upChain sign result hash")

	return res, nil
}
