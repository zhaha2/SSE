package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

//aes加密用
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//加密
func Fen(key string, origData string) (rst string) {
	text := origData
	AesKey := key //秘钥长度为16的倍数
	encrypted, err := AesEncrypt([]byte(text), []byte(AesKey))
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(encrypted)
}

//AES解密用
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

//解密
func Fde(key, enData string) (rst string) {
	bytEn, err := base64.StdEncoding.DecodeString(enData)
	if err != nil {
		panic(err)
	}

	AesKey := key //秘钥长度为16的倍数
	origin, err := AesDecrypt(bytEn, []byte(AesKey))
	if err != nil {
		panic(err)
	}
	return string(origin)
}

func H(k string, val string) (rst string) {
	t := sha256.New()
	io.WriteString(t, k+val)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func main() {
	var K = "randstr123456"
	kk := H(K, "new1")
	k1 := kk[:len(kk)/2]
	k2 := kk[len(kk)/2:]
	l := H(k1, "1")
	fl := H(k1, "0")
	fst := Fen(k2, l+"|||"+k2)
	println(fl)
	println(string(fst))

	kk = H(K, "zf9p1dz")
	k1 = kk[:len(kk)/2]
	k2 = kk[len(kk)/2:]
	l = H(k1, "1")
	fl = H(k1, "0")
	fst = Fen(k2, l+"|||"+k2)
	println(fl)
	println(string(fst))
}
