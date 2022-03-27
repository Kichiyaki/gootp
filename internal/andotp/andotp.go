package andotp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math/big"

	"golang.org/x/crypto/pbkdf2"
)

const (
	ivLen         = 12
	keyLen        = 32
	iterationLen  = 4
	saltLen       = 12
	maxIterations = 160000
	minIterations = 140000
)

func Encrypt(password, plaintext []byte) ([]byte, error) {
	iter := make([]byte, iterationLen)
	iv := make([]byte, ivLen)
	salt := make([]byte, saltLen)

	maxMinIterationsSubtracted, err := rand.Int(rand.Reader, big.NewInt(int64(maxIterations-minIterations)))
	if err != nil {
		return nil, fmt.Errorf("rand.Int: %w", err)
	}

	iterations := int(maxMinIterationsSubtracted.Int64() + minIterations)
	binary.BigEndian.PutUint32(iter, uint32(iterations))

	_, err = rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("rand.Read: %w", err)
	}

	_, err = rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("rand.Read: %w", err)
	}

	secretKey := pbkdf2.Key(password, salt, iterations, keyLen, sha1.New)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}

	cipherText := aesGCM.Seal(nil, iv, plaintext, nil)

	finalCipher := make([]byte, 0, len(iter)+len(salt)+len(iv)+len(cipherText))
	finalCipher = append(finalCipher, iter...)
	finalCipher = append(finalCipher, salt...)
	finalCipher = append(finalCipher, iv...)
	finalCipher = append(finalCipher, cipherText...)

	return finalCipher, nil

}

func Decrypt(password, text []byte) ([]byte, error) {
	iterations := text[:iterationLen]
	salt := text[iterationLen : iterationLen+saltLen]
	iv := text[iterationLen+saltLen : iterationLen+saltLen+ivLen]
	cipherText := text[iterationLen+saltLen+ivLen:]
	iter := int(binary.BigEndian.Uint32(iterations))
	secretKey := pbkdf2.Key(password, salt, iter, keyLen, sha1.New)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}

	plaintext, err := aesGCM.Open(nil, iv, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("aesGCM.Open: %w", err)
	}

	return plaintext, nil
}
