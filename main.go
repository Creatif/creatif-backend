package main

import (
	"creatif/pkg/lib/sdk"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"time"
)

func mdHashing(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}

func encryptIt(value []byte, keyPhrase string) ([]byte, error) {

	aesBlock, err := aes.NewCipher([]byte(mdHashing(keyPhrase)))
	if err != nil {
		return []byte{}, err
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return []byte{}, err
	}

	nonce := make([]byte, gcmInstance.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return []byte{}, err
	}

	cipheredText := gcmInstance.Seal(nonce, nonce, value, nil)

	return cipheredText, nil
}

func decryptIt(ciphered []byte, keyPhrase string) ([]byte, error) {
	hashedPhrase := mdHashing(keyPhrase)
	aesBlock, err := aes.NewCipher([]byte(hashedPhrase))
	if err != nil {
		return []byte{}, err
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return []byte{}, err
	}

	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]

	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return []byte{}, err
	}

	return originalText, err
}

func main() {
	uid, _ := sdk.NewULID()
	cookie := fmt.Sprintf("%s:%s", uid, time.Now().String())

	for i := 0; i < 10; i++ {
		encrypted, err := encryptIt([]byte(cookie), "digital1986")
		if err != nil {
			log.Fatal(err)
		}

		decrypted, err := decryptIt(encrypted, "digital1986")

		fmt.Println(string(decrypted))
	}
}
