package aes_rsa

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

//string->byte->hash
func GenBytesPrivateKey(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func AESEncrypt(key []byte, text []byte) (string, error) {
	//cipher.Block
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("aes new cipher failed")
		return "", err
	}
	blockSize := block.BlockSize()
	originData := pad(text, blockSize)

	iv := []byte(tools.GenCode(blockSize))
	blockMode := cipher.NewCBCEncrypter(block, iv)

	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)

	var buffer bytes.Buffer
	buffer.Write(iv)
	buffer.Write(crypted)
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//has base64
func AESDecrypt(key []byte, text string) (string, error) {
	decode_data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("base 64 decode failed")
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("aes new cipher failed")
		return "", err
	}

	iv := decode_data[:16]
	blockMode := cipher.NewCBCDecrypter(block, iv)

	origin_data := make([]byte, len(decode_data[16:]))
	blockMode.CryptBlocks(origin_data, decode_data[16:])

	return string(unpad(origin_data)), nil
}

func unpad(ciphertext []byte) []byte {
	length := len(ciphertext)
	unpadding := int(ciphertext[length-1])
	return ciphertext[:(length - unpadding)]
}

func RSAPKCS1V15Encrypt(publicKey *rsa.PublicKey, message []byte) (string, error) {
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("rsa encrypt with oaep failed")
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
func RSAPKCS1V15Decrypt(privateKey *rsa.PrivateKey, base64string string) ([]byte, error) {
	decode_data, err := base64.StdEncoding.DecodeString(base64string)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("base64 decode failed")
		return nil, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decode_data)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("rsa decrypt with oaep failed")
		return nil, err
	}
	return plainText, nil
}

func ParsePKIXPublicKey(file string) (*rsa.PublicKey, error) {
	publicKey, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("read file failed")
		return nil, err
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		logger.Logrus.Error("block is nil")
		return nil, errors.New("block is nil")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("parse publickey failed")
		return nil, err
	}
	return pubInterface.(*rsa.PublicKey), nil
}

func ParsePrivateKey(file string) (*rsa.PrivateKey, error) {
	prv, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("read  privatekey  file failed")
		return nil, err
	}
	block, _ := pem.Decode(prv)
	if block == nil {
		logger.Logrus.Error("block is nil")
		return nil, errors.New("block is nil")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("parse private key failed")
		return nil, err
	}
	return priv, nil
}
