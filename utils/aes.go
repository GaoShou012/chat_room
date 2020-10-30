package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

/*
	注意加密key，长度是16,24,32
*/

const (
	AesModeCBCPk5   = "CBCPk5"
	AesModeCBCPk5Iv = "CBCPk5Iv"
)

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(originData []byte) []byte {
	length := len(originData)

	// 去掉最后一个字节 unPadding 次
	unPadding := int(originData[length-1])
	return originData[:(length - unPadding)]
}

func AesEncrypt(origData []byte, key []byte, iv []byte, mode string) ([]byte, error) {
	if len(origData) < 1 {
		return nil, errors.New("orig data is empty")
	}

	switch mode {
	case AesModeCBCPk5:
		if len(key) < 1 {
			return nil, errors.New("key is empty")
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		blockSize := block.BlockSize()

		// 如果你不指定填充及加密模式的话，将会采用 CBC 模式和 PKCS5Padding 填充进行处理, 这里采用 pkcs5Padding
		origData = pkcs5Padding(origData, blockSize)

		// origData = ZeroPadding(origData, block.BlockSize())
		blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
		crypted := make([]byte, len(origData))

		// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
		// crypted := origData
		blockMode.CryptBlocks(crypted, origData)

		return crypted, nil
	case AesModeCBCPk5Iv:
		if len(key) < 1 || len(iv) < 1 {
			return nil, errors.New("key or iv is empty")
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		blockSize := block.BlockSize()

		// 如果你不指定填充及加密模式的话，将会采用 CBC 模式和 PKCS5Padding 填充进行处理, 这里采用 pkcs5Padding
		origData = pkcs5Padding(origData, blockSize)

		// origData = ZeroPadding(origData, block.BlockSize())
		blockMode := cipher.NewCBCEncrypter(block, iv)
		crypted := make([]byte, len(origData))

		// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
		// crypted := origData
		blockMode.CryptBlocks(crypted, origData)

		return crypted, nil
	default:
		return nil, errors.New("unknown mode")
	}
}

func AesDecrypt(crypted []byte, key []byte, iv []byte, mode string) ([]byte, error) {
	if len(crypted) < 1 {
		return []byte(""), errors.New("crypted is empty")
	}
	if len(key) < 1 {
		return []byte(""), errors.New("key is empty")
	}

	switch mode {
	case AesModeCBCPk5:
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}

		blockSize := block.BlockSize()
		if len(key) < blockSize {
			return nil, errors.New("key too short")
		}

		if len(crypted)%blockSize != 0 {
			return nil, errors.New("crypto/cipher: input not full blocks")
		}

		blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
		origData := make([]byte, len(crypted))

		// origData := crypted
		blockMode.CryptBlocks(origData, crypted)
		origData = pkcs5UnPadding(origData)

		return origData, nil
	case AesModeCBCPk5Iv:
		if len(key) < 1 || len(iv) < 1 {
			return []byte(""), errors.New("key or iv is empty")
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}

		blockMode := cipher.NewCBCDecrypter(block, iv)
		origData := make([]byte, len(crypted))

		// origData := crypted
		blockMode.CryptBlocks(origData, crypted)
		origData = pkcs5UnPadding(origData)

		return origData, nil
	default:
		return nil, errors.New("unknown mode")
	}
}
