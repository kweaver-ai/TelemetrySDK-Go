package cipters

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/pkg/errors"
)

var (
	rsaCipter = NewRSACipher(RSAPRIVATEKEY, RSAPUBLICKEY)
)

type rsaCipher struct {
	privateKey string
	publicKey  string
}

func NewRSACipher(privateKey string, publicKey string) *rsaCipher {
	ci := &rsaCipher{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
	return ci
}

// RsaEncryptBase64 使用 RSA 公钥加密数据, 返回加密后并编码为 base64 的数据
func RsaEncryptBase64(originalData string) (string, error) {
	block, _ := pem.Decode([]byte(rsaCipter.publicKey))
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", errors.Wrap(err, "解析公钥失败")
	}
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(originalData))
	if err != nil {
		return "", errors.Wrap(err, "加密失败")
	}
	return base64.StdEncoding.EncodeToString(encryptedData), err
}

// RsaDecryptBase64 使用 RSA 私钥解密数据
func RsaDecryptBase64(encryptedData string) (string, error) {
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", errors.Wrap(err, "base64解析失败")
	}
	block, _ := pem.Decode([]byte(rsaCipter.privateKey))
	priKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", errors.Wrap(err, "解析私钥失败")
	}

	originalData, err := rsa.DecryptPKCS1v15(rand.Reader, priKey.(*rsa.PrivateKey), encryptedDecodeBytes)
	return string(originalData), err
}
