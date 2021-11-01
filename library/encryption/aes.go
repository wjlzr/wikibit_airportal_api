package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

//指定Aes加密
func SpecifyAesEncrypt(orig string, key string, ivKey string) string {
	origData := []byte(orig)
	k := []byte(key)

	block, err := aes.NewCipher(k)
	if err == nil {
		blockSize := block.BlockSize()
		origData = PKCS7Padding(origData, blockSize)
		blockMode := cipher.NewCBCEncrypter(block, []byte(ivKey))
		cryted := make([]byte, len(origData))
		blockMode.CryptBlocks(cryted, origData)
		return base64.StdEncoding.EncodeToString(cryted)
	}
	return ""
}

//指定Aes解密
func SpecifyAesDecrypt(cryted string, key string, ivKey string) string {
	crytedByte, err := base64.StdEncoding.DecodeString(cryted)
	if err == nil {
		k := []byte(key)
		block, err := aes.NewCipher(k)
		if err == nil {
			blockMode := cipher.NewCBCDecrypter(block, []byte(ivKey))
			orig := make([]byte, len(crytedByte))
			blockMode.CryptBlocks(orig, crytedByte)

			orig = PKCS7UnPadding(orig)
			return string(orig)
		}
	}
	return ""
}

//aes加密
func AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err == nil {
		blockSize := block.BlockSize()                            // 获取秘钥块的长度
		origData = PKCS7Padding(origData, blockSize)              // 补全码
		blockMode := cipher.NewCBCEncrypter(block, k[:blockSize]) // 加密模式
		cryted := make([]byte, len(origData))                     // 创建数组
		blockMode.CryptBlocks(cryted, origData)                   // 加密
		return base64.StdEncoding.EncodeToString(cryted)
	}
	return ""
}

//aes解密
func AesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, err := base64.StdEncoding.DecodeString(cryted)
	if err == nil {
		k := []byte(key)
		block, err := aes.NewCipher(k) // 分组秘钥
		if err == nil {
			blockSize := block.BlockSize()                            // 获取秘钥块的长度
			blockMode := cipher.NewCBCDecrypter(block, k[:blockSize]) // 加密模式
			orig := make([]byte, len(crytedByte))                     // 创建数组
			blockMode.CryptBlocks(orig, crytedByte)                   // 解密
			orig = PKCS7UnPadding(orig)                               // 去补全码
			return string(orig)
		}
	}
	return ""
}

//PKCS7填充
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS7解除
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length-unpadding < 0 {
		return nil
	}
	return origData[:(length - unpadding)]
}

//PKCS5填充
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS5解除
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//zero填充
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

//zero解除
func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
