package aes_rsa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESEncryptDecrypt(t *testing.T) {
	key := "TestAESEncrypt"
	textString := "HelloConnexion"
	privateKey := GenBytesPrivateKey(key)
	assert.Equal(t, true, len(privateKey) > 0)
	encData, err := AESEncrypt(privateKey, []byte(textString))
	assert.Nil(t, err)
	assert.Equal(t, true, len(encData) > 0)
	decData, err := AESDecrypt(privateKey, encData)
	assert.Nil(t, err)
	assert.EqualValues(t, textString, decData)
}

func TestRSAEncDecWithPKCS1V15(t *testing.T) {
	text := "TestRSAEncDecWithPKCS1V15"
	publickey, err := ParsePKIXPublicKey("./rsa_public_key.pem")
	assert.Nil(t, err)
	encData, err := RSAPKCS1V15Encrypt(publickey, []byte(text))
	assert.Nil(t, err)
	assert.Equal(t, true, len(encData) > 0)
	privatekey, err := ParsePrivateKey("./rsa_private_key.pem")
	assert.Nil(t, err)
	decData, err := RSAPKCS1V15Decrypt(privatekey, encData)
	assert.Nil(t, err)
	assert.Equal(t, true, string(decData) == text)

}
