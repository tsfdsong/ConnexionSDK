package sign

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/aes_rsa"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/http_client"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"time"

	"github.com/sirupsen/logrus"
)

func GenSignRequest(paramHash, signData string) (PRequestSignature, error) {
	r := PRequestSignature{}
	key := tools.GenCode(8)

	privateKey := aes_rsa.GenBytesPrivateKey(key)
	encodeData, err := aes_rsa.AESEncrypt(privateKey, []byte(paramHash))
	if err != nil {
		return r, err
	}

	r.EncodeData = encodeData
	assetData, err := aes_rsa.AESEncrypt(privateKey, []byte(signData))
	if err != nil {
		return r, err
	}

	r.AssetsData = assetData
	r.Timestamp = fmt.Sprintf("%v", time.Now().Unix())

	//oaep := aes_rsa.NewRSAOaep(sha256.New())

	publickey, err := aes_rsa.ParsePKIXPublicKey(config.GetSignConfig().SignaturePublicFile)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("parse public key failed")
		return r, err
	}

	k, err := aes_rsa.RSAPKCS1V15Encrypt(publickey, []byte(key))
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("rsa encrypt key failed")
		return r, err
	}

	r.Key = k
	return r, nil
}

//GetNftSignStatus get sign status
func GetNftSignStatus(signHash string) (string, error) {
	r := &RQuerySignature{}

	signURL := config.GetSignConfig().SignResultServerURL + const_def.SDK_QUERY_SIGN_URL + "/" + signHash

	err := http_client.HttpClientReqWithGet(signURL, &r)
	if err != nil {
		return "", err
	}

	if r.Code != "success" {
		logger.Logrus.WithFields(logrus.Fields{"Res": r, "SignHash": signHash}).Error("sign return status code is failed")
		return "", nil
	}

	if r.Sign == const_def.SIGN_PENDING {
		return "", fmt.Errorf("%s is pending", signHash)
	}

	sign, err := tools.EthSignFix(r.Sign)
	if err != nil {
		return "", err
	}

	return sign, nil
}
