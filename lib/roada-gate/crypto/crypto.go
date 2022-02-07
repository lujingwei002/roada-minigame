package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/roada-go/gat/log"
)

//CBC加密
func DesCBCEncrypt(src, key string) string {
	data := []byte(src)
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		panic(err)
	}
	data = PKCS5Padding(data, block.BlockSize())
	//获取CBC加密模式
	iv := keyByte //用密钥作为向量(不建议这样使用)
	mode := cipher.NewCBCEncrypter(block, iv)
	out := make([]byte, len(data))
	mode.CryptBlocks(out, data)
	return fmt.Sprintf("%X", out)
}

//CBC解密
func DesCBCDecrypt(src, key string) string {
	keyByte := []byte(key)
	data, err := hex.DecodeString(src)
	if err != nil {
		panic(err)
	}
	block, err := des.NewCipher(keyByte)
	if err != nil {
		panic(err)
	}
	iv := keyByte //用密钥作为向量(不建议这样使用)
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(data))
	mode.CryptBlocks(plaintext, data)
	plaintext = PKCS5UnPadding(plaintext)
	return string(plaintext)
}

//ECB加密
func DesEncrypt(data []byte, key string) ([]byte, error) {
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	//对明文数据进行补码
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//对明文按照blocksize进行分块加密
		//必要时可以使用go关键字进行并行加密
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

//ECB解密
func DesDecrypt(data []byte, key string) ([]byte, error) {
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = PKCS5UnPadding(out)
	return out, nil
}

//明文补码算法
func PKCS5Padding(encryptedData []byte, blockSize int) []byte {
	padding := blockSize - len(encryptedData)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(encryptedData, padtext...)
}

//明文减码算法
func PKCS5UnPadding(plainData []byte) []byte {
	length := len(plainData)
	unpadding := int(plainData[length-1])
	return plainData[:(length - unpadding)]
}

//解密：对采用sha1算法加密后转base64格式的数据进行解密（私钥PKCS1格式）
func RsaDecryptWithSha1Base64(encryptedData, privateKey string) (string, error) {
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	log.Println("ggggg", len(encryptedDecodeBytes), base64.StdEncoding.DecodedLen(len(encryptedData)))
	key, _ := base64.StdEncoding.DecodeString(privateKey)
	prvKey, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		log.Println(err)
		return "", err
	}
	originalData, err := rsa.DecryptPKCS1v15(rand.Reader, prvKey, encryptedDecodeBytes)
	return string(originalData), err
}

func RsaEncryptWithSha1Base64(plainData, publicKey string) (string, error) {
	derBytes, _ := base64.StdEncoding.DecodeString(publicKey)
	pub, err := x509.ParsePKCS1PublicKey(derBytes)
	if err != nil {
		log.Println(err)
		return "", err
	}
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(plainData))
	return base64.StdEncoding.EncodeToString(encryptedData), err
}
