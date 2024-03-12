package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
)

// 本文都是处理加密解密的方法
const (
	encryptionKey  = "WHFanxingwangluo" // 密钥
	encryptionSalt = "fx-"              // 盐
)

// FxEncrypt 内部加密函数
func FxEncrypt(encryptStr string) string {
	key := []byte(encryptionKey)
	salt := []byte(encryptionSalt)
	plainText := []byte(encryptStr)

	// 将盐与明文数据进行拼接
	data := append(plainText, salt...)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], data)

	return base64.URLEncoding.EncodeToString(cipherText)
}

// FxDecrypt 内部解密函数
func FxDecrypt(cipherText string) ([]byte, error) {
	key := []byte(encryptionKey)
	salt := []byte(encryptionSalt)

	encryptedData, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(encryptedData) < aes.BlockSize {
		return nil, err
	}

	iv := encryptedData[:aes.BlockSize]
	encryptedData = encryptedData[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encryptedData, encryptedData)

	// 去除盐部分
	decryptedData := encryptedData[:len(encryptedData)-len(salt)]

	return decryptedData, nil
}
