package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

func GenerateRandomPassword() []byte {
	randomPassword := make([]byte, 32)
	_, err := rand.Read(randomPassword)
	if err != nil {
		fmt.Errorf("failed to generate random password: %v", err)
		return []byte("")
	}
	return randomPassword
}

func EncryptQueryParams(params map[string]string, key []byte) string {
	paramsString, err := json.Marshal(params)
	if err != nil {
		fmt.Errorf("failed to convert query params to string: %v", err)
		return ""
	}

	encryptedParams, err := CustomEncrypt(string(paramsString), key)
	if err != nil {
		fmt.Errorf("failed to encrypt query params: %v", err)
		return ""
	}

	return encryptedParams
}

func EncryptString(plainText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainTextBytes := []byte(plainText)
	ciphertext := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainTextBytes)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptString(cipherText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherTextBytes, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(cipherTextBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherTextBytes, cipherTextBytes)

	return string(cipherTextBytes), nil
}

func SwapCase(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r - 'a' + 'A'
		} else if r >= 'A' && r <= 'Z' {
			return r - 'A' + 'a'
		} else {
			return r
		}
	}, str)
}

func ReverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func CustomEncrypt(str string, key []byte) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(str))

	reversed := ReverseString(encoded)

	swapped := SwapCase(reversed)

	encrypted, err := EncryptString(swapped, key)
	if err != nil {
		return "", err
	}

	final := ReverseString(encrypted)

	return final, nil
}

func CustomDecrypt(cipherText string, key []byte) (string, error) {
	reversed := ReverseString(cipherText)

	decrypted, err := DecryptString(reversed, key)
	if err != nil {
		return "", err
	}

	swapped := SwapCase(decrypted)

	reversedAgain := ReverseString(swapped)

	decoded, err := base64.StdEncoding.DecodeString(reversedAgain)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
