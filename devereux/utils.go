package devereux

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/howeyc/gopass"
)

var DEBUG = true

func Debugf(format string, args ...interface{}) {
	if DEBUG {
		log.Printf("DEBUG " + format, args...)
	}
}

func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// encrypt string to base64 crypto using AES
func encrypt(keyText string, rawData []byte) ([]byte, error) {
	x := sha256.Sum256([]byte(keyText))
	key := x[:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize + len(rawData))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte(""), err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], rawData)

	// convert to base64
	return []byte(base64.URLEncoding.EncodeToString(ciphertext)), nil
}

// decrypt from base64 to decrypted string
func decrypt(keyText string, cryptoBytes []byte) ([]byte, error) {
	x := sha256.Sum256([]byte(keyText))
	key := x[:]

	ciphertext, _ := base64.URLEncoding.DecodeString(string(cryptoBytes))

	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return []byte(""), errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return []byte(fmt.Sprintf("%s", ciphertext)), nil
}

func concatByteSlices(byteSlices ...[]byte) []byte {
	newByteSlice := make([]byte, 0)
	for i := 0; i < len(byteSlices); i++ {
		newByteSlice = append(newByteSlice, byteSlices[i]...)
	}
	return newByteSlice
}

func promptUserForInput(prompt string) (string, error) {
	var err error
	var keyBytes []byte

	fmt.Print(prompt)
	keyBytes, err = gopass.GetPasswd()
	if err != nil {
		return "", err
	}

	return string(keyBytes), nil
}
